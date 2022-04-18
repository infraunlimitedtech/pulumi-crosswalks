package services

import (
	helmv3 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (infra *Infra) RunPrometheus() error {
	name := "kube-prometheus-stack"

	_, err := helmv3.NewRelease(infra.ctx, "kube-monitoring", &helmv3.ReleaseArgs{
		Chart:         pulumi.String(name),
		Namespace:     pulumi.String(infra.Namespace),
		CleanupOnFail: pulumi.Bool(true),
		WaitForJobs:   pulumi.Bool(false),
		Timeout:       pulumi.Int(600),
		//Version:          pulumi.String(infra.LB.NginxIngress.Helm.Version),
		RepositoryOpts: helmv3.RepositoryOptsArgs{
			Repo: pulumi.String("https://prometheus-community.github.io/helm-charts"),
		},
		Values: pulumi.Map{
			"prometheusOperator": pulumi.Map{
				"clusterDomain": pulumi.String("local.intra.infraunlimited.tech"),
				"resources": pulumi.Map{
					"requests": pulumi.Map{
						"memory": pulumi.String("256Mi"),
					},
					"limits": pulumi.Map{
						"memory": pulumi.String("368Mi"),
					},
				},
			},
			"grafana": pulumi.Map{
				"sidecar": pulumi.Map{
					"dashboards": pulumi.Map{
						"multicluster": pulumi.Map{
					        	"enabled": pulumi.Bool(true),
						},
					},
				},
				"ingress": pulumi.Map{
					"enabled": pulumi.Bool(true),
					"hosts": pulumi.StringArray{
						pulumi.String("grafana.local.intra.infraunlimited.tech"),
					},
					"tls": pulumi.MapArray{
						pulumi.Map{
							"secretName": pulumi.String("nginx-ingress-nginx-ingress-default-server-tls"),
							"hosts": pulumi.StringArray{
								pulumi.String("grafana.local.intra.infraunlimited.tech"),
							},
						},
					},
				},
			},
			"kubelet": pulumi.Map{
				"enabled": pulumi.Bool(true),
			},
			"kubeProxy": pulumi.Map{
				"enabled": pulumi.Bool(false),
			},
			"kubeScheduler": pulumi.Map{
				"enabled": pulumi.Bool(false),
			},
			"kubeEtcd": pulumi.Map{
				"enabled": pulumi.Bool(false),
			},
			"alertmanager": pulumi.Map{
				"alertmanagerSpec": pulumi.Map{
					"resources": pulumi.Map{
						"requests": pulumi.Map{
							"memory": pulumi.String("32Mi"),
						},
						"limits": pulumi.Map{
							"memory": pulumi.String("64Mi"),
						},
					},
				},
			},
			"prometheus-node-exporter": pulumi.Map{
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
						"cpu":    pulumi.String("25m"),
					},
					"limits": pulumi.Map{
						"memory": pulumi.String("32Mi"),
						"cpu":    pulumi.String("50m"),
					},
				},
			},
			"prometheus": pulumi.Map{
				"prometheusSpec": pulumi.Map{
					"externalLabels": pulumi.Map{
						"cluster": pulumi.String("mgmt"),
					},
					"scrapeInterval": pulumi.String("60s"),
					"resources": pulumi.Map{
						"requests": pulumi.Map{
							"memory": pulumi.String("512Mi"),
							"cpu":    pulumi.String("500m"),
						},
						"limits": pulumi.Map{
							"memory": pulumi.String("1024Mi"),
							"cpu":    pulumi.String("1"),
						},
					},
					"storageSpec": pulumi.Map{
						"volumeClaimTemplate": pulumi.Map{
							"spec": pulumi.Map{
								"resources": pulumi.Map{
									"requests": pulumi.Map{
										"storage": pulumi.String("50Gi"),
									},
								},
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}
