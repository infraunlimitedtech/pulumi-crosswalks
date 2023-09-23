package services

import (
	"k8s-cluster/packages/kilo"
	"k8s-cluster/services/gitlab"
	kiloVPN "k8s-cluster/services/kilo-vpn"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	pulumiConfig "github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

type service interface {
	IsEnabled() bool
	Manage(*pulumi.Context) error
}

type Services struct {
	KiloVPN *kilo.Kilo
	Gitlab  *gitlab.Gitlab
}

type Runner struct {
	ctx     *pulumi.Context
	services []service
}

func NewRunner(ctx *pulumi.Context) (*Runner, error) {
	// Init vars from stack's config
	var pulumiServicesCfg Services
	cfg := pulumiConfig.New(ctx, "")
	cfg.RequireSecretObject("services", &pulumiServicesCfg)

	services := []service{
		gitlab.New(pulumiServicesCfg.Gitlab),
		kiloVPN.New(pulumiServicesCfg.KiloVPN),
	}

	s := &Runner{
		ctx:      ctx,
		services: services,
	}

	return s, nil
}

func (s *Runner) Run() error {
	for _, service := range s.services {
		if service.IsEnabled() {
			if err := service.Manage(s.ctx); err != nil {
				return err
			}
		}
	}

	return nil
}