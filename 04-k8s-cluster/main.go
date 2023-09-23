package main

import (
	"k8s-cluster/addons"
	"k8s-cluster/rbac"
	"k8s-cluster/rbac/sa"
	"k8s-cluster/services"
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

		kubeRBAC, err := rbac.Init(ctx)
		if err != nil {
			return err
		}

		kubeClusterAddons, err := addons.NewRunner(ctx, kubeClusterSpec, infraLayerNodeInfo)
		if err != nil {
			return err
		}

		if err := kubeClusterAddons.Run(); err != nil {
			return err
		}

		kubeClusterServices, err := services.NewRunner(ctx)
		if err != nil {
			return err
		}

		if err := kubeClusterServices.Run(); err != nil {
			return err
		}

		if kubeRBAC.ServiceAccounts.Prometheus {
			sa.PrometheusAccount(kubeClusterAddons.Namespace, ctx)
		}


		return nil
	})
}
