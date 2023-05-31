package monitoring

import (
	"fmt"

	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	helmv3 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v3"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const vmPort = 8428

func (m *Stack) runVM() error {
	appName := "victoria-metrics"
	m.VictoriaMetrics.Server.Port = vmPort

	deployed, err := helmv3.NewRelease(m.ctx, appName, &helmv3.ReleaseArgs{
		Name:      pulumi.String(appName),
		Chart:     pulumi.String("victoria-metrics-single"),
		Namespace: m.Namespace.Metadata.Name().Elem(),
		RepositoryOpts: helmv3.RepositoryOptsArgs{
			Repo: pulumi.String("https://victoriametrics.github.io/helm-charts"),
		},
		Version: pulumi.String(m.VictoriaMetrics.Helm.Version),
		Values: pulumi.Map{
			"server": pulumi.Map{
				"scrape": pulumi.Map{
					"enabled": pulumi.Bool(true),
					"config": pulumi.Map{
						"global": pulumi.Map{
							"scrape_interval": pulumi.String("60s"),
						},
					},
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
			},
		},
	})
	if err != nil {
		return fmt.Errorf("helm: %w", err)
	}

	// Create a service with clusterIP for VictoriaMetrics.
	// Because some of our services outside of the cluster need to connect to VM somehow without inCluster dns.
	// But with Stateful set the helm chart creates only headless service
	_, err = corev1.NewService(m.ctx, appName, &corev1.ServiceArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(appName),
			Namespace: m.Namespace.Metadata.Name().Elem(),
		},
		Spec: &corev1.ServiceSpecArgs{
			Ports: corev1.ServicePortArray{
				corev1.ServicePortArgs{
					Name:       pulumi.String("http"),
					Port:       pulumi.Int(m.VictoriaMetrics.Server.Port),
					TargetPort: pulumi.Int(m.VictoriaMetrics.Server.Port),
				},
			},
			Selector: pulumi.StringMap{
				"app.kubernetes.io/instance": pulumi.String(appName),
			},
			Type:      pulumi.String("ClusterIP"),
			ClusterIP: pulumi.String(m.VictoriaMetrics.Server.ClusterIP),
		},
	}, pulumi.DependsOn([]pulumi.Resource{deployed}))

	if err != nil {
		return fmt.Errorf("k8s service: %w", err)
	}

	return nil
}
