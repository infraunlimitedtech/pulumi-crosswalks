package cloudflared

import (
	"fmt"
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

	defaultBackend string
}

type Rule struct {
	Hostname string
	Service  string
}

func New(cfg *Cloudflared, DefaultBackend string) *Cloudflared {
	if cfg == nil {
		cfg = &Cloudflared{
			Enabled: false,
		}
	}

	cfg.defaultBackend = DefaultBackend

	return cfg
}

func (c *Cloudflared) IsEnabled() bool {
	return c.Enabled
}

func (c *Cloudflared) Manage(ctx *pulumi.Context, ns *corev1.Namespace) error {
	repo := "https://cloudflare.github.io/helm-charts"
	name := "cloudflare-tunnel"

	ingress := pulumi.Array{}
	for _, rule := range c.Ingress {

		if rule.Service == "" {
			rule.Service = fmt.Sprintf("https://%s", c.defaultBackend)
		}

		ingress = append(ingress, pulumi.Map{
			"hostname": pulumi.String(rule.Hostname),
			"service":  pulumi.String(rule.Service),
			"originRequest": pulumi.Map{
				"noTLSVerify":    pulumi.Bool(true),
				"httpHostHeader": pulumi.String(rule.Hostname),
			},
		})
	}

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
			"replicaCount": pulumi.Int(1),
			"cloudflare": pulumi.Map{
				"account":    pulumi.String(c.Account),
				"secret":     pulumi.String(c.Secret),
				"tunnelId":   pulumi.String(c.TunnelID),
				"tunnelName": pulumi.String(c.TunnelName),
				"ingress":    ingress,
			},
		},
	})

	if err != nil {
		return nil
	}

	return nil
}
