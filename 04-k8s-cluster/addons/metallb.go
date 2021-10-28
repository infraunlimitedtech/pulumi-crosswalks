package addons

import (
	helmv3 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (a *Addons) RunMetalLb() error {
	addonName := "metallb-addon"
	externalPoolName := "external"

	_, err := helmv3.NewChart(a.ctx, "metallbAddon", helmv3.ChartArgs{
		Chart:     pulumi.String("metallb"),
		Namespace: pulumi.String(a.Namespace),
		FetchArgs: helmv3.FetchArgs{
			Repo: pulumi.String("https://metallb.github.io/metallb"),
		},
		Values: pulumi.Map{
			"fullnameOverride": pulumi.String(addonName),
			"speaker": pulumi.Map{
				"nodeSelector": pulumi.Map{
					"node-role.kubernetes.io/master": pulumi.String("true"),
				},
			},
			"configInline": pulumi.Map{
				"address-pools": pulumi.Array{
					pulumi.Map{
						"name":     pulumi.String("default"),
						"protocol": pulumi.String("layer2"),
						"addresses": pulumi.Array{
							pulumi.String(a.MetalLb.DefaultNetworkPool),
						},
						"avoid-buggy-ips": pulumi.Bool(true),
					},
					pulumi.Map{
						"name":     pulumi.String("kubeapi"),
						"protocol": pulumi.String("layer2"),
						"addresses": pulumi.Array{
							pulumi.String(a.MetalLb.KubeapiIP + "/32"),
						},
						"auto-assign": pulumi.Bool(false),
					},
					pulumi.Map{
						"name":     pulumi.String(externalPoolName),
						"protocol": pulumi.String("layer2"),
						"addresses": pulumi.Array{
							pulumi.String(a.MetalLb.ExternalIP + "/32"),
						},
						"auto-assign": pulumi.Bool(false),
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
