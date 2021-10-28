package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"

	"managed-infrastructure/infra"
	"managed-infrastructure/providers/vagrant"
)

type pulumiConfig struct {
	SSH SSHConfig
}

type SSHConfig struct {
	User       string
	PrivateKey string `json:"private_key"`
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		var pulumiCfg pulumiConfig
		cfg := config.New(ctx, "")
		cfg.RequireSecretObject("main", &pulumiCfg)

		var i infra.Infra
		if ctx.Stack() == "local" {
			ctx.Log.Warn("Local stack is not implemented yet. Controlled via Vagrant", nil)
			i = vagrant.Init(pulumiCfg.SSH.User, pulumi.ToSecret(pulumiCfg.SSH.PrivateKey))
		}

		ctx.Export("infra:nodes:info", pulumi.ToMapMap(i.GetNodes()))
		return nil

	})
}
