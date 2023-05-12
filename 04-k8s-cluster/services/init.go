package services

import (
	"k8s-cluster/packages/kilo"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	pulumiConfig "github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

type Services struct {
	ctx     *pulumi.Context
	KiloVPN *kilo.Kilo
}

func Init(ctx *pulumi.Context) (*Services, error) {
	// Init vars from stack's config
	var pulumiServicesCfg Services
	cfg := pulumiConfig.New(ctx, "")
	cfg.RequireSecretObject("services", &pulumiServicesCfg)

	s := &Services{
		ctx:     ctx,
		KiloVPN: pulumiServicesCfg.KiloVPN,
	}

	return s, nil
}
