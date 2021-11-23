package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"

	"managed-infrastructure/infra"
	"managed-infrastructure/providers/libvirt"
	"managed-infrastructure/providers/vagrant"
)

type mainConfig struct {
	Provider      string
	IdentityStack string `json:"identitystack"`
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		var mainCfg mainConfig
		cfg := config.New(ctx, "")
		cfg.RequireObject("main", &mainCfg)

		identityStack, err := pulumi.NewStackReference(ctx, mainCfg.IdentityStack, nil)
		if err != nil {
			return err
		}

		sshCreds := identityStack.GetOutput(pulumi.String("identity:ssh:credentials"))

		var i infra.Infra
		switch mainCfg.Provider {
		case "vagrant":
			ctx.Log.Warn("Vagrant stack is not implemented yet. Controlled via Vagrant", nil)
			i = vagrant.Init(sshCreds)
		case "libvirt":
			var libvirtCfg libvirt.Config
			var err error
			cfg.RequireSecretObject("libvirt", &libvirtCfg)
			i, err = libvirt.Init(ctx, sshCreds, &libvirtCfg)
			if err != nil {
				return err
			}
		default:
			ctx.Log.Error("Unknown provider", nil)
			return nil
		}

		ctx.Export("infra:nodes:info", pulumi.ToMapMap(i.GetNodes()))
		return nil
	})
}
