package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

type mainCfg struct {
	InfraStack string
}

func getMainCfg(ctx *pulumi.Context) mainCfg {
	var pulumiCfg mainCfg
	cfg := config.New(ctx, "")
	cfg.RequireSecretObject("main", &pulumiCfg)

	return pulumiCfg
}
