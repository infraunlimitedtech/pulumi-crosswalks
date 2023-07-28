package kilo

import (
	"fmt"
	"k8s-cluster/config"
	"path/filepath"

	kilov1 "k8s-cluster/crds/generated/squat/kilo/v1alpha1"

	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	rbacv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/rbac/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/yaml"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type StartedKilo struct {
	Deployed  bool
	Port      int
	Firewalls *config.Firewalls
}

type Kilo struct {
	Enabled     *config.Status
	Name        string
	Port        int
	PrivateKey  string
	CRDS        *config.CRDS
	Version     string
	Peers       []Peer
	Firewalls   *config.Firewalls
	ExternalIPs []string
}

type Peer struct {
	Name       string
	PublicKey  string
	AllowedIPs []string
}

// Using the same wireguard port for all instances because it's not possible to use different.
// Since kilo makes some labels on node it rewrite wg endpoint and port is taken from node.
const port = 31200

func RunKilo(ctx *pulumi.Context, ns *corev1.Namespace, cfg *Kilo) (*StartedKilo, error) {
	serviceAccount, err := corev1.NewServiceAccount(ctx, fmt.Sprintf("%s_KiloSA", cfg.Name), &corev1.ServiceAccountArgs{
		ApiVersion: pulumi.String("v1"),
		Kind:       pulumi.String("ServiceAccount"),
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(cfg.Name),
			Namespace: ns.Metadata.Name(),
		},
	})
	if err != nil {
		return nil, err
	}

	clusterRole, err := rbacv1.NewClusterRole(ctx, fmt.Sprintf("%s_kiloClusterRole", cfg.Name), &rbacv1.ClusterRoleArgs{
		ApiVersion: pulumi.String("rbac.authorization.k8s.io/v1"),
		Kind:       pulumi.String("ClusterRole"),
		Metadata: &metav1.ObjectMetaArgs{
			Name:      serviceAccount.Metadata.Name(),
			Namespace: serviceAccount.Metadata.Namespace(),
		},
		Rules: rbacv1.PolicyRuleArray{
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String(""),
				},
				Resources: pulumi.StringArray{
					pulumi.String("nodes"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("list"),
					pulumi.String("patch"),
					pulumi.String("watch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("kilo.squat.ai"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("peers"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("list"),
					pulumi.String("watch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("apiextensions.k8s.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("customresourcedefinitions"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	_, err = rbacv1.NewClusterRoleBinding(ctx, fmt.Sprintf("%s_kiloClusterRoleBinding", cfg.Name), &rbacv1.ClusterRoleBindingArgs{
		ApiVersion: pulumi.String("rbac.authorization.k8s.io/v1"),
		Kind:       pulumi.String("ClusterRoleBinding"),
		Metadata: &metav1.ObjectMetaArgs{
			Name: serviceAccount.Metadata.Name().Elem(),
		},
		RoleRef: &rbacv1.RoleRefArgs{
			ApiGroup: pulumi.String("rbac.authorization.k8s.io"),
			Kind:     pulumi.String("ClusterRole"),
			Name:     clusterRole.Metadata.Name().Elem(),
		},
		Subjects: rbacv1.SubjectArray{
			&rbacv1.SubjectArgs{
				Kind:      pulumi.String("ServiceAccount"),
				Name:      serviceAccount.Metadata.Name().Elem(),
				Namespace: serviceAccount.Metadata.Namespace().Elem(),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	scripts, err := corev1.NewConfigMap(ctx, fmt.Sprintf("%s_kiloScripts", cfg.Name), &corev1.ConfigMapArgs{
		ApiVersion: pulumi.String("v1"),
		Kind:       pulumi.String("ConfigMap"),
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String("kilo-scripts"),
			Namespace: ns.Metadata.Name().Elem(),
		},
		Data: pulumi.StringMap{
			"init.sh": pulumi.String(`#!/bin/sh
cat > /etc/kubernetes/kubeconfig <<EOF
apiVersion: v1
kind: Config
name: kilo
clusters:
- cluster:
    server: $(sed -n 's/.*server: \(.*\)/\1/p' /var/lib/rancher/k3s/agent/kubelet.kubeconfig)
    certificate-authority: /var/lib/rancher/k3s/agent/server-ca.crt
users:
- name: kilo
  user:
    token: $(cat /var/run/secrets/kubernetes.io/serviceaccount/token)
contexts:
- name: kilo
  context:
    cluster: kilo
    namespace: ${NAMESPACE}
    user: kilo
current-context: kilo
EOF
cp secrets/* /var/lib/kilo/`),
		},
	})
	if err != nil {
		return nil, err
	}

	privateKey, err := corev1.NewSecret(ctx, fmt.Sprintf("%s_kiloPrivateKey", cfg.Name), &corev1.SecretArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String("kilo-private-key"),
			Namespace: ns.Metadata.Name(),
		},
		Type: pulumi.String("Opaque"),
		StringData: pulumi.StringMap{
			"key": pulumi.String(cfg.PrivateKey),
		},
	})
	if err != nil {
		return nil, err
	}

	deploy, err := appsv1.NewDeployment(ctx, fmt.Sprintf("%s_kiloDeployment", cfg.Name), &appsv1.DeploymentArgs{
		ApiVersion: pulumi.String("apps/v1"),
		Kind:       pulumi.String("Deployment"),
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(cfg.Name),
			Namespace: ns.Metadata.Name(),
			Labels: pulumi.StringMap{
				"app.kubernetes.io/name":    pulumi.String("kilo"),
				"app.kubernetes.io/part-of": pulumi.String("kilo"),
			},
		},
		Spec: &appsv1.DeploymentSpecArgs{
			Strategy: &appsv1.DeploymentStrategyArgs{
				Type: pulumi.String("Recreate"),
			},
			Replicas: pulumi.Int(1),
			Selector: &metav1.LabelSelectorArgs{
				MatchLabels: pulumi.StringMap{
					"app.kubernetes.io/name":    pulumi.String("kilo"),
					"app.kubernetes.io/part-of": pulumi.String("kilo"),
				},
			},
			Template: &corev1.PodTemplateSpecArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Labels: pulumi.StringMap{
						"app.kubernetes.io/name":    pulumi.String("kilo"),
						"app.kubernetes.io/part-of": pulumi.String("kilo"),
					},
				},
				Spec: &corev1.PodSpecArgs{
					Tolerations: corev1.TolerationArray{
						&corev1.TolerationArgs{
							Key:      pulumi.String("CriticalAddonsOnly"),
							Operator: pulumi.String("Exists"),
						},
						&corev1.TolerationArgs{
							Key:      pulumi.String("node-role.kubernetes.io/control-plane"),
							Operator: pulumi.String("Exists"),
						},
					},
					Affinity: &corev1.AffinityArgs{
						NodeAffinity: &corev1.NodeAffinityArgs{
							PreferredDuringSchedulingIgnoredDuringExecution: corev1.PreferredSchedulingTermArray{
								&corev1.PreferredSchedulingTermArgs{
									Weight: pulumi.Int(1),
									Preference: &corev1.NodeSelectorTermArgs{
										MatchExpressions: corev1.NodeSelectorRequirementArray{
											&corev1.NodeSelectorRequirementArgs{
												Key:      pulumi.String("infraunlimited.tech/kilo-vpn-node"),
												Operator: pulumi.String("In"),
												Values: pulumi.StringArray{
													pulumi.String("true"),
												},
											},
										},
									},
								},
							},
						},
					},
					ServiceAccountName: serviceAccount.Metadata.Name().Elem(),
					HostNetwork:        pulumi.Bool(false),
					Containers: corev1.ContainerArray{
						&corev1.ContainerArgs{
							Name:  pulumi.String(cfg.Name),
							Image: pulumi.Sprintf("squat/kilo:%v", cfg.Version),
							Args: pulumi.StringArray{
								pulumi.String("--kubeconfig=/etc/kubernetes/kubeconfig"),
								// Without hostname it will not work
								pulumi.Sprintf("%v%v%v", "--hostname=", "$", "(NODE_NAME)"),
								pulumi.String("--cni=false"),
								pulumi.String("--log-level=info"),
								pulumi.Sprintf("--port=%d", port),
								pulumi.Sprintf("--interface=%s", cfg.Name),
								pulumi.String("--compatibility=flannel"),
								pulumi.String("--local=false"),
							},
							Env: corev1.EnvVarArray{
								&corev1.EnvVarArgs{
									Name: pulumi.String("NODE_NAME"),
									ValueFrom: &corev1.EnvVarSourceArgs{
										FieldRef: &corev1.ObjectFieldSelectorArgs{
											FieldPath: pulumi.String("spec.nodeName"),
										},
									},
								},
							},
							Resources: &corev1.ResourceRequirementsArgs{
								Requests: pulumi.StringMap{
									"memory": pulumi.String("32Mi"),
									"cpu":    pulumi.String("50m"),
								},
								Limits: pulumi.StringMap{
									"memory": pulumi.String("64Mi"),
									"cpu":    pulumi.String("50m"),
								},
							},
							Ports: corev1.ContainerPortArray{
								&corev1.ContainerPortArgs{
									ContainerPort: pulumi.Int(1107),
									Name:          pulumi.String("metrics"),
								},
							},
							SecurityContext: &corev1.SecurityContextArgs{
								Privileged: pulumi.Bool(false),
								Capabilities: &corev1.CapabilitiesArgs{
									Add: pulumi.StringArray{
										pulumi.String("NET_ADMIN"),
										pulumi.String("SYS_MODULE"),
									},
								},
							},
							VolumeMounts: corev1.VolumeMountArray{
								&corev1.VolumeMountArgs{
									Name:      pulumi.String("kubeconfig"),
									MountPath: pulumi.String("/etc/kubernetes"),
									ReadOnly:  pulumi.Bool(true),
								},
								&corev1.VolumeMountArgs{
									Name:      pulumi.String("lib-modules"),
									MountPath: pulumi.String("/lib/modules"),
									ReadOnly:  pulumi.Bool(true),
								},
								&corev1.VolumeMountArgs{
									Name:      pulumi.String("kilo-dir"),
									MountPath: pulumi.String("/var/lib/kilo"),
								},
							},
						},
					},
					InitContainers: corev1.ContainerArray{
						&corev1.ContainerArgs{
							Name:  pulumi.String("prepare-configs"),
							Image: pulumi.String("squat/kilo"),
							Command: pulumi.StringArray{
								pulumi.String("/bin/sh"),
							},
							Args: pulumi.StringArray{
								pulumi.String("/scripts/init.sh"),
							},
							VolumeMounts: corev1.VolumeMountArray{
								&corev1.VolumeMountArgs{
									Name:      pulumi.String("kubeconfig"),
									MountPath: pulumi.String("/etc/kubernetes"),
								},
								&corev1.VolumeMountArgs{
									Name:      pulumi.String("scripts"),
									MountPath: pulumi.String("/scripts/"),
									ReadOnly:  pulumi.Bool(true),
								},
								&corev1.VolumeMountArgs{
									Name:      pulumi.String("k3s-agent"),
									MountPath: pulumi.String("/var/lib/rancher/k3s/agent/"),
									ReadOnly:  pulumi.Bool(true),
								},
								&corev1.VolumeMountArgs{
									Name:      pulumi.String("kilo-dir"),
									MountPath: pulumi.String("/var/lib/kilo"),
								},
								&corev1.VolumeMountArgs{
									Name:      pulumi.String("key"),
									MountPath: pulumi.String("secrets"),
								},
							},
							Env: corev1.EnvVarArray{
								&corev1.EnvVarArgs{
									Name: pulumi.String("NAMESPACE"),
									ValueFrom: &corev1.EnvVarSourceArgs{
										FieldRef: &corev1.ObjectFieldSelectorArgs{
											FieldPath: pulumi.String("metadata.namespace"),
										},
									},
								},
							},
						},
					},
					Volumes: corev1.VolumeArray{
						&corev1.VolumeArgs{
							Name:     pulumi.String("kilo-dir"),
							EmptyDir: nil,
						},
						&corev1.VolumeArgs{
							Name: pulumi.String("key"),
							Secret: &corev1.SecretVolumeSourceArgs{
								SecretName: privateKey.Metadata.Name().Elem(),
							},
						},
						&corev1.VolumeArgs{
							Name:     pulumi.String("kubeconfig"),
							EmptyDir: nil,
						},
						&corev1.VolumeArgs{
							Name: pulumi.String("scripts"),
							ConfigMap: &corev1.ConfigMapVolumeSourceArgs{
								Name: scripts.Metadata.Name().Elem(),
							},
						},
						&corev1.VolumeArgs{
							Name: pulumi.String("k3s-agent"),
							HostPath: &corev1.HostPathVolumeSourceArgs{
								Path: pulumi.String("/var/lib/rancher/k3s/agent"),
							},
						},
						&corev1.VolumeArgs{
							Name: pulumi.String("lib-modules"),
							HostPath: &corev1.HostPathVolumeSourceArgs{
								Path: pulumi.String("/lib/modules"),
							},
						},
					},
				},
			},
		},
	}, pulumi.DeleteBeforeReplace(true))
	if err != nil {
		return nil, err
	}

	_, err = corev1.NewService(ctx, fmt.Sprintf("%s_kiloVpnService", cfg.Name), &corev1.ServiceArgs{
		ApiVersion: pulumi.String("v1"),
		Kind:       pulumi.String("Service"),
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String("kilo-vpn"),
			Namespace: deploy.Metadata.Namespace().Elem(),
		},
		Spec: &corev1.ServiceSpecArgs{
			Type: pulumi.String("NodePort"),
			Selector: pulumi.StringMap{
				"app.kubernetes.io/name":    pulumi.String("kilo"),
				"app.kubernetes.io/part-of": pulumi.String("kilo"),
			},
			Ports: corev1.ServicePortArray{
				&corev1.ServicePortArgs{
					Name:     pulumi.String("vpn"),
					Protocol: pulumi.String("UDP"),
					Port:     pulumi.Int(port),
					NodePort: pulumi.Int(cfg.Port),
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	// We need to wait the deployed kilo before configuration
	configurationDeps := []pulumi.Resource{deploy}

	if cfg.CRDS.Install {
		deployCRDS, err := yaml.NewConfigGroup(ctx, "kiloCrds", &yaml.ConfigGroupArgs{
			Files: []string{filepath.Join(cfg.CRDS.Path, "*.yaml")},
		})
		if err != nil {
			return nil, err
		}
		configurationDeps = append(configurationDeps, deployCRDS)
	}

	for _, peer := range cfg.Peers {
		_, err = kilov1.NewPeer(ctx, fmt.Sprintf("%s_kiloPeer_%s", cfg.Name, peer.Name), &kilov1.PeerArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name:      pulumi.String(peer.Name),
				Namespace: ns.Metadata.Name(),
			},
			Spec: &kilov1.PeerSpecArgs{
				AllowedIPs:          pulumi.ToStringArray(peer.AllowedIPs),
				PublicKey:           pulumi.String(peer.PublicKey),
				PersistentKeepalive: pulumi.Int(10),
			},
		}, pulumi.DependsOn(configurationDeps))
		if err != nil {
			return nil, err
		}
	}

	startedKilo := &StartedKilo{
		Port:     cfg.Port,
		Deployed: true,
	}

	if cfg.Firewalls != nil {
		startedKilo.Firewalls = cfg.Firewalls
	}

	return startedKilo, nil
}

func (k *Kilo) IsEnabled() config.Status {
	return k.Enabled.WithDefault(true)
}
