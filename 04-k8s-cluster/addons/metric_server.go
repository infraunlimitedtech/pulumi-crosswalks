package addons

import (
	helmv3 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (a *Addons) RunMetricServer() error {
	_, err := helmv3.NewChart(a.ctx, "metrics-server", helmv3.ChartArgs{
		Chart:     pulumi.String("metrics-server"),
		Version:   pulumi.String("v5.9.0"),
		Namespace: pulumi.String(a.Namespace),
		FetchArgs: helmv3.FetchArgs{
			Repo: pulumi.String("https://charts.bitnami.com/bitnami"),
		},
	})
	if err != nil {
		return err
	}
	return nil
}
