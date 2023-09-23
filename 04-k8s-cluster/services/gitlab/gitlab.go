package gitlab

import (
	"fmt"
	"k8s-cluster/config"

	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	helmv3 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v3"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)


type Gitlab struct {
	Enabled bool
	Domain string
	Helm *config.HelmParams
}

func New(cfg *Gitlab) *Gitlab {
	if cfg == nil {
		cfg = &Gitlab{
			Enabled: false,
		}
	}
	return cfg
}

func (g *Gitlab) IsEnabled() bool {
	return g.Enabled
}


func (g *Gitlab) Manage(ctx *pulumi.Context) error {
	repo := "https://charts.gitlab.io/"
	name := "gitlab"

	ns, err := corev1.NewNamespace(ctx, name, &corev1.NamespaceArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name: pulumi.String(name),
		},
	})
	if err != nil {
		return fmt.Errorf("%s namespace: %w", name, err)
	}

	_, err = helmv3.NewRelease(ctx, name, &helmv3.ReleaseArgs{
		Chart:     pulumi.String(name),
		Version:   pulumi.String(g.Helm.Version),
		Namespace: ns.Metadata.Name().Elem(),
		RepositoryOpts: &helmv3.RepositoryOptsArgs{
			Repo: pulumi.String(repo),
		},
		Values: pulumi.Map{
			"global": pulumi.Map{
				"edition": pulumi.String("ce"),
				"hosts": pulumi.Map{
					"domain": pulumi.String(g.Domain),
					"https": pulumi.Bool(false),
					"gitlab": pulumi.Map{
						"name": pulumi.Sprintf("gitlab-dev.%s", g.Domain),
					},
				},
				"minio": pulumi.Map{
					"enabled": pulumi.Bool(true),
				},
				"ingress": pulumi.Map{
					"configureCertmanager": pulumi.Bool(false),
					"class": pulumi.String("nginx-ingress-addon"),
				},
				"kas": pulumi.Map{
					"enabled": pulumi.Bool(false),
				},
			},
			"gitlab": pulumi.Map{
				"webservice": pulumi.Map{
					"resources": pulumi.Map{
						"requests": pulumi.Map{
							"cpu": pulumi.String("500m"),
							"memory": pulumi.String("1000Mi"),
						},
					},
					"minReplicas": pulumi.Int(1),
					"maxReplicas": pulumi.Int(1),
				},
				"sidekiq": pulumi.Map{
					"resources": pulumi.Map{
						"requests": pulumi.Map{
							"cpu": pulumi.String("400m"),
							"memory": pulumi.String("500Mi"),
						},
						"limits": pulumi.Map{
							"cpu": pulumi.String("600m"),
						},
					},
					"minReplicas": pulumi.Int(1),
					"maxReplicas": pulumi.Int(1),
				},
				"gitlab-shell": pulumi.Map{
					"enabled": pulumi.Bool(false),
				},
				"migrations": pulumi.Map{
					"resources": pulumi.Map{
						"requests": pulumi.Map{
							"cpu": pulumi.String("50m"),
							"memory": pulumi.String("100Mi"),
						},
					},
				},
			},
			"certmanager": pulumi.Map{
				"install": pulumi.Bool(false),
			},
			"prometheus": pulumi.Map{
				"install": pulumi.Bool(false),
			},
			"nginx-ingress": pulumi.Map{
				"enabled": pulumi.Bool(false),
			},
			"gitlab-runner": pulumi.Map{
				"install": pulumi.Bool(false),
			},
			"registry": pulumi.Map{
				"hpa": pulumi.Map{
					"minReplicas": pulumi.Int(1),
					"maxReplicas": pulumi.Int(1),
				},
			},
			"certmanager-issuer": pulumi.Map{
				"email": pulumi.String("test@example.com"),
			},
		},
	})
	if err != nil {
		return err
	}
	return nil
}
