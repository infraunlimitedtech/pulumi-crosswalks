package config

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
	"managed-infrastructure/providers/hetzner"
	"managed-infrastructure/providers/libvirt"
	"managed-infrastructure/providers/yandex"
)

type mainConfig struct {
	Providers     *Providers
	IdentityStack string `json:"identitystack"`
}

type Providers struct {
	S3      string
	Compute string
}

type S3Config struct {
	Yandex *yandex.S3Config
}

type ComputeConfig struct {
	Libvirt *libvirt.ComputeConfig
	Hetzner *hetzner.ComputeConfig
}

type GatheredConfig struct {
	Main    *mainConfig
	Compute *ComputeConfig
	S3      *S3Config
}

func ParseConfig(ctx *pulumi.Context) *GatheredConfig {
	cfg := config.New(ctx, "")
	mainCfg, computeCfg, s3Cfg := &mainConfig{}, &ComputeConfig{}, &S3Config{}
	cfg.RequireObject("main", &mainCfg)
	cfg.RequireSecretObject("compute", &computeCfg)
	cfg.RequireSecretObject("s3", &s3Cfg)

	return &GatheredConfig{
		Main:    mainCfg,
		Compute: computeCfg,
		S3:      s3Cfg,
	}
}
