package main

import (
	"k8s-cluster/addons"
	"k8s-cluster/rbac"
	"k8s-cluster/spec"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
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

		if err := kubeClusterAddons.RunMetalLb(); err != nil {
			return err
		}

		if err := kubeClusterAddons.RunNginxIngress(); err != nil {
			return err
		}

		return kubeClusterAddons.RunKilo()
	})
}
