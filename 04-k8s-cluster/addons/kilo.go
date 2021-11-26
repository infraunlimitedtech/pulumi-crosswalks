package addons

import (
	"fmt"
	"path/filepath"

	kilov1 "k8s-cluster/crds/generated/squat/kilo/v1alpha1"

	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	rbacv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/rbac/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/yaml"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (a *Addons) RunKilo() error {
	serviceName := "kilo"
	vpnPort := 51821

	serviceAccount, err := corev1.NewServiceAccount(a.ctx, "kiloServiceAccount", &corev1.ServiceAccountArgs{
		ApiVersion: pulumi.String("v1"),
		Kind:       pulumi.String("ServiceAccount"),
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(serviceName),
			Namespace: pulumi.String(a.Namespace),
		},
	})
	if err != nil {
		return err
	}

	clusterRole, err := rbacv1.NewClusterRole(a.ctx, "kiloClusterRole", &rbacv1.ClusterRoleArgs{
		ApiVersion: pulumi.String("rbac.authorization.k8s.io/v1"),
		Kind:       pulumi.String("ClusterRole"),
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(serviceName),
			Namespace: pulumi.String(a.Namespace),
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
		return err
	}

	_, err = rbacv1.NewClusterRoleBinding(a.ctx, "kiloClusterRoleBinding", &rbacv1.ClusterRoleBindingArgs{
		ApiVersion: pulumi.String("rbac.authorization.k8s.io/v1"),
		Kind:       pulumi.String("ClusterRoleBinding"),
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(serviceName),
			Namespace: pulumi.String(a.Namespace),
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
				Namespace: pulumi.String(a.Namespace),
			},
		},
	})
	if err != nil {
		return err
	}

	_, err = corev1.NewConfigMap(a.ctx, "kiloScripts", &corev1.ConfigMapArgs{
		ApiVersion: pulumi.String("v1"),
		Kind:       pulumi.String("ConfigMap"),
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String("kilo-scripts"),
			Namespace: pulumi.String(a.Namespace),
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
		return err
	}

	privateKey, err := corev1.NewSecret(a.ctx, "kiloPrivateKey", &corev1.SecretArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String("kilo-private-key"),
			Namespace: pulumi.String(a.Namespace),
		},
		Type: pulumi.String("Opaque"),
		StringData: pulumi.StringMap{
			"key": pulumi.String(a.Kilo.PrivateKey),
		},
	})
	if err != nil {
		return err
	}

	_, err = appsv1.NewDeployment(a.ctx, "kiloDeployment", &appsv1.DeploymentArgs{
		ApiVersion: pulumi.String("apps/v1"),
		Kind:       pulumi.String("Deployment"),
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(serviceName),
			Namespace: pulumi.String(a.Namespace),
			Labels: pulumi.StringMap{
				"app.kubernetes.io/name":    pulumi.String("kilo"),
				"app.kubernetes.io/part-of": pulumi.String("kilo"),
			},
		},
		Spec: &appsv1.DeploymentSpecArgs{
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
					/*											Tolerations: corev1.TolerationArray{
						&corev1.TolerationArgs{
							Key:               pulumi.String("node.kubernetes.io/unreachable"),
							Operator:          pulumi.String("Exists"),
							Effect:            pulumi.String("NoExecute"),
							TolerationSeconds: pulumi.Int(2),
						},
						&corev1.TolerationArgs{
							Key:               pulumi.String("node.kubernetes.io/not-ready"),
							Operator:          pulumi.String("Exists"),
							Effect:            pulumi.String("NoExecute"),
							TolerationSeconds: pulumi.Int(2),
						},
					},
					*/NodeSelector: pulumi.StringMap{
						"node-role.kubernetes.io/master": pulumi.String("true"),
					},
					ServiceAccountName: serviceAccount.Metadata.Name().Elem(),
					HostNetwork:        pulumi.Bool(false),
					Containers: corev1.ContainerArray{
						&corev1.ContainerArgs{
							Name:  pulumi.String("kilo"),
							Image: pulumi.String(fmt.Sprintf("squat/kilo:%v", a.Kilo.Version)),
							Args: pulumi.StringArray{
								pulumi.String("--kubeconfig=/etc/kubernetes/kubeconfig"),
								pulumi.String(fmt.Sprintf("%v%v%v", "--hostname=", "$", "(NODE_NAME)")),
								pulumi.String("--cni=false"),
								pulumi.String("--log-level=warn"),
								pulumi.String(fmt.Sprintf("--port=%d", vpnPort)),
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
									Name:      pulumi.String("xtables-lock"),
									MountPath: pulumi.String("/run/xtables.lock"),
									ReadOnly:  pulumi.Bool(false),
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
								Name: pulumi.String("kilo-scripts"),
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
						&corev1.VolumeArgs{
							Name: pulumi.String("xtables-lock"),
							HostPath: &corev1.HostPathVolumeSourceArgs{
								Path: pulumi.String("/run/xtables.lock"),
								Type: pulumi.String("FileOrCreate"),
							},
						},
					},
				},
			},
		},
	}, pulumi.DeleteBeforeReplace(true))
	if err != nil {
		return err
	}

	_, err = corev1.NewService(a.ctx, "kiloVpnService", &corev1.ServiceArgs{
		ApiVersion: pulumi.String("v1"),
		Kind:       pulumi.String("Service"),
		Metadata: &metav1.ObjectMetaArgs{
			Name: pulumi.String("kilo-vpn"),
			Namespace: pulumi.String(a.Namespace),
		},
		Spec: &corev1.ServiceSpecArgs{
			ExternalTrafficPolicy: pulumi.String("Local"),
			LoadBalancerIP:        pulumi.String(a.MetalLb.ExternalIP),
			Type:                  pulumi.String("LoadBalancer"),
			Selector: pulumi.StringMap{
				"app.kubernetes.io/name":    pulumi.String("kilo"),
				"app.kubernetes.io/part-of": pulumi.String("kilo"),
			},
			Ports: corev1.ServicePortArray{
				&corev1.ServicePortArgs{
					Name:     pulumi.String("vpn"),
					Protocol: pulumi.String("UDP"),
					Port:     pulumi.Int(vpnPort),
				},
			},
		},
	})
	if err != nil {
		return err
	}

	_, err = yaml.NewConfigGroup(a.ctx, "kiloCrds", &yaml.ConfigGroupArgs{
		Files: []string{filepath.Join(a.Kilo.Crds.Path, "*.yaml")},
	})

	if err != nil {
		return err
	}

	for _, peer := range a.Kilo.Peers {
		_, err = kilov1.NewPeer(a.ctx, fmt.Sprintf("kiloPeer_%v", peer.Name), &kilov1.PeerArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name:      pulumi.String(peer.Name),
				Namespace: pulumi.String(a.Namespace),
			},
			Spec: &kilov1.PeerSpecArgs{
				AllowedIPs:          pulumi.ToStringArray(peer.AllowedIPs),
				PublicKey:           pulumi.String(peer.PublicKey),
				PersistentKeepalive: pulumi.Int(10),
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}
