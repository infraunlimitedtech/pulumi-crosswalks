package addons

import (
	helmv3 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (a *Addons) RunMetalLb() error {
	addonName := "metallb-addon"

	_, err := helmv3.NewChart(a.ctx, "metallbAddon", helmv3.ChartArgs{
		Chart:     pulumi.String("metallb"),
		Namespace: a.Namespace.Metadata.Name().Elem(),
		Version:   pulumi.String(a.MetalLb.Helm.Version),
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
							pulumi.String(a.MetalLb.Pools.Default.Network),
						},
						"avoid-buggy-ips": pulumi.Bool(true),
					},
					pulumi.Map{
						"name":     pulumi.String("kubeapi"),
						"protocol": pulumi.String("layer2"),
						"addresses": pulumi.Array{
							pulumi.String(a.MetalLb.Pools.Kubeapi.Network),
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
