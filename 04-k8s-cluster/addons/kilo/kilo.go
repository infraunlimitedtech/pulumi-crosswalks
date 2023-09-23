package addons

import (
	"fmt"
	"k8s-cluster/infra/firewall"
	"k8s-cluster/packages/kilo"

	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Kilo struct {
	nodeInfo pulumi.AnyOutput
	cfg   *kilo.Kilo
}

const (
	vpnPort     = 31200
	serviceName = "kilo"
)

func New(cfg *kilo.Kilo, infraLayerNodeInfo pulumi.AnyOutput) *Kilo {
	if cfg == nil {
		cfg = &kilo.Kilo{}
	}

	cfg.Enabled = true
	cfg.Name = serviceName
	cfg.Port = vpnPort

	return &Kilo{
		nodeInfo: infraLayerNodeInfo,
		cfg: cfg,
	}
}

func (k *Kilo) IsEnabled() bool {
	return k.cfg.Enabled
}

func (k *Kilo) Manage(ctx *pulumi.Context, ns *corev1.Namespace) error {
	deployed, err := kilo.RunKilo(ctx, ns, k.cfg)
	if err != nil {
		return fmt.Errorf("deploy kilo: %w", err)
	}

	if deployed.Firewalls != nil {
		err = firewall.Manage(ctx, k.nodeInfo, deployed.GetRequiredFirewallRules())
		if err != nil {
			return fmt.Errorf("manage kilo firewall: %w", err)
		}
	}

	return nil
}
