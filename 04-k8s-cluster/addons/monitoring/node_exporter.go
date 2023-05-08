package monitoring

import (
	"fmt"

	helmv3 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (m *Stack) runNodeExporter() error {
	addonName := "node-exporter"

	_, err := helmv3.NewRelease(m.ctx, addonName, &helmv3.ReleaseArgs{
		Name:      pulumi.String(addonName),
		Chart:     pulumi.String("prometheus-node-exporter"),
		Namespace: m.Namespace.Metadata.Name().Elem(),
		RepositoryOpts: helmv3.RepositoryOptsArgs{
			Repo: pulumi.String("https://prometheus-community.github.io/helm-charts"),
		},
		Version: pulumi.String(m.NodeExporter.Helm.Version),
		Values: pulumi.Map{
			// It is allowed to run anywhere!
			"tolerations": pulumi.Array{
				pulumi.Map{
					"operator": pulumi.String("Exists"),
					"effect":   pulumi.String("NoExecute"),
				},
				pulumi.Map{
					"operator": pulumi.String("Exists"),
					"effect":   pulumi.String("NoSchedule"),
				},
			},
			"annotations": pulumi.Map{
				"prometheus.io/scrape": pulumi.String("true"),
				"prometheus.io/port":   pulumi.String("9100"),
			},
		},
	})
	if err != nil {
		return fmt.Errorf("helm: %w", err)
	}

	return nil
}
