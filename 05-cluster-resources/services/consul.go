package services

import (
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	helmv3 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v3"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (infra *Infra) RunConsulStack() error {
	serviceName := "consul"

	gossipSecret, err := corev1.NewSecret(infra.ctx, "gossip-secret", &corev1.SecretArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Namespace: pulumi.String(infra.Namespace),
		},
		Type: pulumi.String("Opaque"),
		StringData: pulumi.StringMap{
			"key": pulumi.String(infra.Consul.GossipSecret),
		},
	})
	if err != nil {
		return err
	}

	_, err = helmv3.NewChart(infra.ctx, "consul-stack", helmv3.ChartArgs{
		Chart:     pulumi.String("consul"),
		Version:   pulumi.String("v0.32.1"),
		Namespace: pulumi.String(infra.Namespace),
		FetchArgs: helmv3.FetchArgs{
			Repo: pulumi.String("https://helm.releases.hashicorp.com"),
		},
		Values: pulumi.Map{
			"global": pulumi.Map{
				"name":       pulumi.String(serviceName),
				"domain":     pulumi.String(infra.DNS.InternalDomain),
				"datacenter": pulumi.String("CHANGE_ME"),
				"recursors": pulumi.Array{
					pulumi.String("8.8.8.8"),
					pulumi.String("1.1.1.1"),
				},

				"gossipEncryption": pulumi.Map{
					"secretName": gossipSecret.Metadata.Name(),
					"secretKey":  pulumi.String("key"),
				},
			},
			"server": pulumi.Map{
				"replicas": pulumi.Int(3),
				"disruptionBudget": pulumi.Map{
					"enabled": pulumi.Bool(false),
				},
				"storage": pulumi.String("1Gi"),
				"resources": pulumi.Map{
					"requests": pulumi.Map{
						"cpu": pulumi.String("50m"),
					},
					"limits": pulumi.Map{
						"cpu": pulumi.String("50m"),
					},
				},
				"extraConfig": pulumi.String("{    \"log_level\": \"DEBUG\"  }"),
			},
			"client": pulumi.Map{
				"resources": pulumi.Map{
					"requests": pulumi.Map{
						"memory": pulumi.String("64Mi"),
						"cpu":    pulumi.String("50m"),
					},
					"limits": pulumi.Map{
						"memory": pulumi.String("64Mi"),
						"cpu":    pulumi.String("50m"),
					},
				},
			},
			"syncCatalog": pulumi.Map{
				"enabled":               pulumi.Bool(true),
				"addK8SNamespaceSuffix": pulumi.Bool(false),
				"syncClusterIPServices": pulumi.Bool(true),
				"default":               pulumi.Bool(false),
				"k8sDenyNamespaces": pulumi.StringArray{
					pulumi.String("none"),
				},
				"toK8S": pulumi.Bool(false),
				"resources": pulumi.Map{
					"requests": pulumi.Map{
						"memory": pulumi.String("16Mi"),
						"cpu":    pulumi.String("20m"),
					},
					"limits": pulumi.Map{
						"memory": pulumi.String("32Mi"),
						"cpu":    pulumi.String("30m"),
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
