package addons

import (
	nginxv1 "k8s-cluster/crds/generated/nginxinc/kubernetes-ingress/k8s/v1alpha1"

	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	helmv3 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v3"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (a *Addons) RunNginxIngress() error {
	addonName := a.NginxIngress.Name
	ingressClassName := a.NginxIngress.Name

	deploy, err := helmv3.NewChart(a.ctx, addonName, helmv3.ChartArgs{
		Chart:     pulumi.String("nginx-ingress"),
		Namespace: pulumi.String(a.Namespace),
		Version:   pulumi.String(a.NginxIngress.Helm.Version),
		FetchArgs: helmv3.FetchArgs{
			Repo: pulumi.String("https://helm.nginx.com/stable"),
		},
		Values: pulumi.Map{
			"controller": pulumi.Map{
				"name":         pulumi.String(addonName),
				"replicaCount": pulumi.Int(2),
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
				"config": pulumi.Map{
					"entries": pulumi.Map{
						"stream-snippets": pulumi.String("map_hash_bucket_size 512;\n map_hash_max_size 2048;"),
					},
				},
			},
		},
	})
	if err != nil {
		return err
	}

	kubeUpstreamName := "kube-masters"
	kubeAPIAddr := "kubernetes." + a.NginxIngress.Domain

	_, err = nginxv1.NewTransportServer(a.ctx, "nginxIngressAddonTransportServerKubernetes", &nginxv1.TransportServerArgs{
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
	a.ctx.Export("clusters:addons:kubeAPIAddress", pulumi.String(kubeAPIAddr))

	_, err = corev1.NewService(a.ctx, "nginxIngressAddonLoadBalancerKubeApi", &corev1.ServiceArgs{
		ApiVersion: pulumi.String("v1"),
		Kind:       pulumi.String("Service"),
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String("kube-api"),
			Namespace: pulumi.String(a.Namespace),
		},
		Spec: &corev1.ServiceSpecArgs{
			Selector: pulumi.StringMap{
				"app": pulumi.String(a.NginxIngress.Name),
			},
			Type:      pulumi.String("ClusterIP"),
			ClusterIP: pulumi.String(a.NginxIngress.KubeAPI.ClusterIP),
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

	return nil
}
