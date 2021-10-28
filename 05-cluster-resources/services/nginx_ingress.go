package services

import (
	helmv3 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (infra *Infra) RunNginxIngress() error {
	name := infra.LB.NginxIngress.Name

	_, err := helmv3.NewChart(infra.ctx, name, helmv3.ChartArgs{
		Chart:            pulumi.String(name),
		Namespace:        pulumi.String(infra.Namespace),
		Version:          pulumi.String(infra.LB.NginxIngress.Helm.Version),
		SkipCRDRendering: pulumi.Bool(true),
		FetchArgs: helmv3.FetchArgs{
			Repo: pulumi.String("https://helm.nginx.com/stable"),
		},
		Values: pulumi.Map{
			"controller": pulumi.Map{
				"name":         pulumi.String(name),
				"replicaCount": pulumi.Int(1),
				"nodeSelector": pulumi.Map{
					"node-role.kubernetes.io/master": pulumi.String("true"),
				},
				"service": pulumi.Map{
					"annotations": pulumi.Map{
						"consul.hashicorp.com/service-name": pulumi.String("ingress"),
					},
					"externalTrafficPolicy": pulumi.String("Cluster"),
					"loadBalancerSourceRanges": pulumi.Array{
						pulumi.String("192.168.0.0/16"),
						pulumi.String("10.0.0.0/8"),
					},
				},
				"globalConfiguration": pulumi.Map{
					"create": pulumi.Bool(true),
					"spec": pulumi.Map{
						"listeners": pulumi.Array{
							pulumi.Map{
								"name":     infra.LB.NginxIngress.UDPDNSListener.Name,
								"port":     infra.LB.NginxIngress.UDPDNSListener.Port,
								"protocol": infra.LB.NginxIngress.UDPDNSListener.Protocol,
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}
