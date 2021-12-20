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
					"memory": pulumi.String("16Mi"),
				},
				"limits": pulumi.Map{
					"memory": pulumi.String("32Mi"),
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
