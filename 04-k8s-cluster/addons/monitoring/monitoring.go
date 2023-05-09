package monitoring

import (
	"fmt"
	"k8s-cluster/addons"

	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Stack struct {
	ctx             *pulumi.Context
	Namespace       *corev1.Namespace
	NodeExporter    *addons.NodeExporter
	VictoriaMetrics *addons.VictoriaMetrics
}

func Run(ctx *pulumi.Context, params *addons.Monitoring) error {
	namespace := "monitoring"

	// Setup all monitoring services and deployments to mon namespace
	ns, err := corev1.NewNamespace(ctx, namespace, &corev1.NamespaceArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name: pulumi.String(namespace),
		},
	})
	if err != nil {
		return fmt.Errorf("monitoring namespace: %w", err)
	}

	// Setup node-exporter
	mon := &Stack{
		ctx:             ctx,
		Namespace:       ns,
		VictoriaMetrics: params.VictoriaMetrics,
		NodeExporter:    params.NodeExporter,
	}

	err = mon.runNodeExporter()
	if err != nil {
		return fmt.Errorf("node-exporter: %w", err)
	}

	err = mon.runVM()
	if err != nil {
		return fmt.Errorf("victoria-metrics %w", err)
	}

	return nil
}
