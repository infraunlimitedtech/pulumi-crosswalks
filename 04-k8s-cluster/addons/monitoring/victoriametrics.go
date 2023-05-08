package monitoring

import (
	"fmt"

	helmv3 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (m *Stack) runVM() error {
	appName := "victoria-metrics"

	_, err := helmv3.NewRelease(m.ctx, appName, &helmv3.ReleaseArgs{
		Name:      pulumi.String(appName),
		Chart:     pulumi.String("victoria-metrics-single"),
		Namespace: m.Namespace.Metadata.Name().Elem(),
		RepositoryOpts: helmv3.RepositoryOptsArgs{
			Repo: pulumi.String("https://victoriametrics.github.io/helm-charts"),
		},
		Version: pulumi.String("0.8.59"),
		Values: pulumi.Map{
			"server": pulumi.Map{
				"scrape": pulumi.Map{
					"enabled": pulumi.Bool(true),
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

	return nil
}
