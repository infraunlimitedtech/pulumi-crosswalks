package services

import (
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	helmv3 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v3"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (infra *Infra) RunNginxIngress() error {
	name := infra.LB.NginxIngress.Name

	deploy, err := helmv3.NewChart(infra.ctx, name, helmv3.ChartArgs{
		Chart:            pulumi.String(name),
		Namespace:        pulumi.String(infra.Namespace),
		Version:          pulumi.String(infra.LB.NginxIngress.Helm.Version),
		SkipCRDRendering: pulumi.Bool(true),
		FetchArgs: helmv3.FetchArgs{
			Repo: pulumi.String("https://helm.nginx.com/stable"),
		},
		Values: pulumi.Map{
			"controller": pulumi.Map{
				"name":                pulumi.String(name),
				"replicaCount":        pulumi.Int(1),
				"setAsDefaultIngress": pulumi.Bool(true),
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

	_, err = corev1.NewService(infra.ctx, "nginx-ingress", &corev1.ServiceArgs{
		ApiVersion: pulumi.String("v1"),
		Kind:       pulumi.String("Service"),
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String("ingress"),
			Namespace: pulumi.String(infra.Namespace),
		},
		Spec: &corev1.ServiceSpecArgs{
			Selector: pulumi.StringMap{
				"app": pulumi.String(infra.LB.NginxIngress.Name),
			},
			Type:      pulumi.String("ClusterIP"),
			ClusterIP: pulumi.String(infra.LB.NginxIngress.ClusterIP),
			Ports: corev1.ServicePortArray{
				&corev1.ServicePortArgs{
					Name:     pulumi.String("http"),
					Protocol: pulumi.String("TCP"),
					Port:     pulumi.Int(80),
				},
				&corev1.ServicePortArgs{
					Name:     pulumi.String("https"),
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
