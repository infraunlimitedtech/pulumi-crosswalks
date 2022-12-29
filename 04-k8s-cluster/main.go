package main

import (
	"k8s-cluster/addons"
	"k8s-cluster/infra/firewall"
	"k8s-cluster/rbac"
	"k8s-cluster/spec"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		mainCfg := getMainCfg(ctx)

		infraStack, err := pulumi.NewStackReference(ctx, mainCfg.InfraStack, nil)
		if err != nil {
			return err
		}

		infraLayerNodeInfo := infraStack.GetOutput(pulumi.String("infra:nodes:info"))

		kubeClusterSpec, err := spec.Init(ctx)
		if err != nil {
			return err
		}
		if err := rbac.Init(ctx); err != nil {
			return err
		}

		kubeClusterAddons, err := addons.Init(ctx, kubeClusterSpec)
		if err != nil {
			return err
		}

		if err := kubeClusterAddons.RunMetricServer(); err != nil {
			return err
		}

		if err := kubeClusterAddons.RunNginxIngress(); err != nil {
			return err
		}

		kilo, err := kubeClusterAddons.RunKilo()
		if err != nil {
			return err
		}

		return firewall.Manage(ctx, infraLayerNodeInfo, kilo.GetRequiredFirewallRules())
	})
}
