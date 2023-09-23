package kilo_vpn

import (
	"k8s-cluster/packages/kilo"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

)

const (
	vpnPort     = 32200
	serviceName = "kilo-vpn"
)

type KiloVPN struct {
	cfg    *kilo.Kilo
}

func New(cfg *kilo.Kilo) *KiloVPN {
	if cfg == nil {
		cfg = &kilo.Kilo{
			Enabled: false,
		}
	}
	
	cfg.Name = serviceName
	cfg.Port = vpnPort

	return &KiloVPN{
		cfg: cfg,
	}
}

func (k *KiloVPN) IsEnabled() bool {
	return k.cfg.Enabled
}

func (k *KiloVPN) Manage(ctx *pulumi.Context) error {
	ns, err := kilo.CreateNS(ctx, "kilo-vpn")
	if err != nil {
		return err
	}

	_, err = kilo.RunKilo(ctx, ns, k.cfg)
	if err != nil {
		return err
	}

	return nil
}
