package nginx_ingress

import (
	"k8s-cluster/config"
	nginxv1 "k8s-cluster/crds/generated/nginxinc/kubernetes-ingress/k8s/v1alpha1"

	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	helmv3 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v3"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type NginxIngress struct {
	Enabled bool
	Name    string
	ClusterIP string
	Domain  string
	KubeAPI NginxKubeAPI
	Replicas int
	Helm    *config.HelmParams
}

type NginxKubeAPI struct {
	ClusterIP string
}

func New(cfg *NginxIngress) *NginxIngress {
	if cfg == nil {
		cfg = &NginxIngress{}
	}
	cfg.Enabled = true
	return cfg
}

func (n *NginxIngress) IsEnabled() bool {
	return n.Enabled
}

func (n *NginxIngress) Manage(ctx *pulumi.Context, ns *corev1.Namespace) error {
	addonName := n.Name
	ingressClassName := n.Name

	deploy, err := helmv3.NewRelease(ctx, addonName, &helmv3.ReleaseArgs{
		Name:      pulumi.String(addonName),
		Chart:     pulumi.String("nginx-ingress"),
		Namespace: ns.Metadata.Name().Elem(),
		Version:   pulumi.String(n.Helm.Version),
		RepositoryOpts: &helmv3.RepositoryOptsArgs{
			Repo: pulumi.String("https://helm.nginx.com/stable"),
		},
		Values: pulumi.Map{
			"controller": pulumi.Map{
				"name":         pulumi.String(addonName),
				"replicaCount": pulumi.Int(n.Replicas),
				"nodeSelector": pulumi.Map{
					"node-role.kubernetes.io/master": pulumi.String("true"),
				},
				"tolerations": pulumi.MapArray{
					pulumi.Map{
						"operator": pulumi.String("Exists"),
						"key":      pulumi.String("CriticalAddonsOnly"),
					},
					pulumi.Map{
						"operator": pulumi.String("Exists"),
						"key":      pulumi.String("node-role.kubernetes.io/control-plane"),
					},
				},
				"enableTLSPassthrough": pulumi.Bool(true),
				"ingressClass":         pulumi.String(ingressClassName),
				"resources": pulumi.Map{
					"requests": pulumi.Map{
						"memory": pulumi.String("64Mi"),
						"cpu":    pulumi.String("25m"),
					},
					"limits": pulumi.Map{
						"memory": pulumi.String("96Mi"),
						"cpu":    pulumi.String("50m"),
					},
				},
				"service": pulumi.Map{
					"create": pulumi.Bool(false),
				},
				"globalConfiguration": pulumi.Map{
					"create": pulumi.Bool(false),
				},
			},
		},
	})
	if err != nil {
		return err
	}

	kubeUpstreamName := "kube-masters"
	kubeAPIAddr := "kubernetes." + n.Domain

	_, err = nginxv1.NewTransportServer(ctx, "nginxIngressAddonTransportServerKubernetes", &nginxv1.TransportServerArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String("kube-masters"),
			Namespace: pulumi.String("default"),
		},
		Spec: &nginxv1.TransportServerSpecArgs{
			Listener: &nginxv1.TransportServerSpecListenerArgs{
				Name:     pulumi.String("tls-passthrough"),
				Protocol: pulumi.String("TLS_PASSTHROUGH"),
			},
			Host:             pulumi.String(kubeAPIAddr),
			IngressClassName: pulumi.String(ingressClassName),
			Upstreams: &nginxv1.TransportServerSpecUpstreamsArray{
				&nginxv1.TransportServerSpecUpstreamsArgs{
					Name:    pulumi.String(kubeUpstreamName),
					Service: pulumi.String("kubernetes"),
					Port:    pulumi.Int(443),
				},
			},
			Action: &nginxv1.TransportServerSpecActionArgs{
				Pass: pulumi.String(kubeUpstreamName),
			},
		},
	})

	if err != nil {
		return err
	}

	_, err = corev1.NewService(ctx, "nginxIngressAddonLoadBalancerKubeApi", &corev1.ServiceArgs{
		ApiVersion: pulumi.String("v1"),
		Kind:       pulumi.String("Service"),
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String("kube-api"),
			Namespace: deploy.Namespace,
		},
		Spec: &corev1.ServiceSpecArgs{
			Selector: pulumi.StringMap{
				"app": pulumi.String(n.Name),
			},
			Type:      pulumi.String("ClusterIP"),
			ClusterIP: pulumi.String(n.KubeAPI.ClusterIP),
			Ports: corev1.ServicePortArray{
				&corev1.ServicePortArgs{
					Protocol: pulumi.String("TCP"),
					Port:     pulumi.Int(443),
				},
			},
		},
	}, pulumi.DeleteBeforeReplace(true), pulumi.DependsOn([]pulumi.Resource{deploy}))

	if err != nil {
		return err
	}

	_, err = corev1.NewService(ctx, "nginxIngressAddonLoadBalancerMain", &corev1.ServiceArgs{
		ApiVersion: pulumi.String("v1"),
		Kind:       pulumi.String("Service"),
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String("ingress"),
			Namespace: deploy.Namespace,
		},
		Spec: &corev1.ServiceSpecArgs{
			Selector: pulumi.StringMap{
				"app": pulumi.String(n.Name),
			},
			Type:      pulumi.String("ClusterIP"),
			ClusterIP: pulumi.String(n.ClusterIP),
			Ports: corev1.ServicePortArray{
				&corev1.ServicePortArgs{
					Name: pulumi.String("http"),
					Protocol: pulumi.String("TCP"),
					Port:     pulumi.Int(80),
				},
				&corev1.ServicePortArgs{
					Name: pulumi.String("https"),
					Protocol: pulumi.String("TCP"),
					Port:     pulumi.Int(443),
				},
			},
		},
	}, pulumi.DeleteBeforeReplace(true), pulumi.DependsOn([]pulumi.Resource{deploy}))

	if err != nil {
		return err
	}

	return nil
}
