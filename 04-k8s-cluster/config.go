package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

type pulumiConfig struct {
	Key string `json:"key"`
}

func parseConfig(ctx *pulumi.Context) *pulumiConfig {
	var pulumiCfg pulumiConfig
	cfg := config.New(ctx, "")
	cfg.RequireSecretObject("main", &pulumiCfg)

	return &pulumiCfg
}
