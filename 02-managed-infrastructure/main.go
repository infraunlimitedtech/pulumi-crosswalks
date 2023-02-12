package main

import (
	"managed-infrastructure/config"
	"managed-infrastructure/infra"
	"managed-infrastructure/providers/external"
	"managed-infrastructure/providers/hetzner"
	"managed-infrastructure/providers/libvirt"
	"managed-infrastructure/providers/yandex"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		cfg := config.ParseConfig(ctx)

		var err error
		identityStack, err := pulumi.NewStackReference(ctx, cfg.Main.IdentityStack, nil)
		if err != nil {
			return err
		}

		sshCreds := identityStack.GetOutput(pulumi.String("identity:ssh:server_access:credentials"))

		var i infra.ComputeInfra
		switch cfg.Main.Providers.Compute {
		case "libvirt":
			i, err = libvirt.ManageCompute(ctx, sshCreds, cfg.Compute.Libvirt)
			if err != nil {
				return err
			}
		case "hetzner":
			i, err = hetzner.ManageCompute(ctx, sshCreds, cfg.Compute.Hetzner)
			if err != nil {
				return err
			}
		case "external":
			i = external.Init(sshCreds, cfg.Compute.External)
		default:
			ctx.Log.Error("Unknown compute provider", nil)
			return nil
		}

		var s infra.S3Infra
		switch cfg.Main.Providers.S3 {
		case "yandex":
			creds := identityStack.GetOutput(pulumi.String("identity:yandex:s3"))
			p, err := yandex.InitProvider(ctx, creds)
			if err != nil {
				return err
			}
			s, err = yandex.ManageS3(ctx, cfg.S3.Yandex, creds, p)
			if err != nil {
				return err
			}
		case "none":
		default:
			ctx.Log.Error("Unknown S3 provider", nil)
			return nil
		}

		ctx.Export("infra:nodes:info", pulumi.ToMapMap(i.GetNodes()))

		if cfg.Main.Providers.S3 != "none" {
			ctx.Export("infra:storage:info", pulumi.ToMapMap(s.GetStorage()))
		}
		return nil
	})
}
