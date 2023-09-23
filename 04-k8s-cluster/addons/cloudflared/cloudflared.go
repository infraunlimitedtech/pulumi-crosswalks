package cloudflared

import (
	"k8s-cluster/config"

	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	helmv3 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Cloudflared struct {
	Enabled    bool
	Helm       *config.HelmParams
	Account    string
	TunnelID   string
	TunnelName string
	Secret     string
	Ingress    []Rule

	backend string
}

type Rule struct {
	Hostname string
	Service  string
}

func New(cfg *Cloudflared, backend string) *Cloudflared {
	if cfg == nil {
		cfg = &Cloudflared{
			Enabled: false,
		}
	}

	cfg.backend = backend

	return cfg
}

func (c *Cloudflared) IsEnabled() bool {
	return c.Enabled
}

func (c *Cloudflared) Manage(ctx *pulumi.Context, ns *corev1.Namespace) error {
	repo := "https://cloudflare.github.io/helm-charts"
	name := "cloudflare-tunnel"
	_, err := helmv3.NewRelease(ctx, name, &helmv3.ReleaseArgs{
		Chart:     pulumi.String(name),
		Version:   pulumi.String(c.Helm.Version),
		Namespace: ns.Metadata.Name().Elem(),
		RepositoryOpts: &helmv3.RepositoryOptsArgs{
			Repo: pulumi.String(repo),
		},
		Values: pulumi.Map{
			"image": pulumi.Map{
				"tag": pulumi.String("2023.8.2"),
			},
			"cloudflare": pulumi.Map{
				"account":    pulumi.String(c.Account),
				"secret":     pulumi.String(c.Secret),
				"tunnelId":   pulumi.String(c.TunnelID),
				"tunnelName": pulumi.String(c.TunnelName),
				"ingress": pulumi.Array{
					pulumi.Map{
						"hostname": pulumi.String("*.infraunlimited.tech"),
						"service":  pulumi.Sprintf("https://%s", c.backend),
						"originRequest": pulumi.Map{
							"noTLSVerify":    pulumi.Bool(true),
							"httpHostHeader": pulumi.String("gitlab-dev.infraunlimited.tech"),
						},
					},
				},
			},
		},
	})

	if err != nil {
		return nil
	}

	return nil
}
