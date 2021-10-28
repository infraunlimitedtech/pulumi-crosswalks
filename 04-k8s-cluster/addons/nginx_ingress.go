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
				"replicaCount": pulumi.Int(1),
				"nodeSelector": pulumi.Map{
					"node-role.kubernetes.io/master": pulumi.String("true"),
				},
				"enableTLSPassthrough": pulumi.Bool(true),
				"ingressClass":         pulumi.String(ingressClassName),
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
			Type:           pulumi.String("LoadBalancer"),
			LoadBalancerIP: pulumi.String("192.168.74.1"),
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
