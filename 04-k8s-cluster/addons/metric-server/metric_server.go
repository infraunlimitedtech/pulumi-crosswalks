package metrics_server

import (
	helmv3 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v3"
		corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type MetricServer struct {
	Enabled bool
}

func New() *MetricServer {
	return &MetricServer{
		Enabled: true,
	}
}

func (m *MetricServer) IsEnabled() bool {
	return m.Enabled
}

func (m *MetricServer) Manage(ctx *pulumi.Context, ns *corev1.Namespace) error {
	_, err := helmv3.NewRelease(ctx, "metrics-server", &helmv3.ReleaseArgs{
		Chart:     pulumi.String("metrics-server"),
		Version:   pulumi.String("v6.2.17"),
		Namespace: ns.Metadata.Name().Elem(),
		RepositoryOpts: &helmv3.RepositoryOptsArgs{
			Repo: pulumi.String("https://charts.bitnami.com/bitnami"),
		},
		Values: pulumi.Map{
			"nodeSelector": pulumi.Map{
				"node-role.kubernetes.io/control-plane": pulumi.String("true"),
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
			"resources": pulumi.Map{
				"requests": pulumi.Map{
					"memory": pulumi.String("48Mi"),
				},
				"limits": pulumi.Map{
					"memory": pulumi.String("96Mi"),
				},
			},
			"apiService": pulumi.Map{
				"create": pulumi.Bool(true),
			},
		},
	})
	if err != nil {
		return err
	}
	return nil
}
