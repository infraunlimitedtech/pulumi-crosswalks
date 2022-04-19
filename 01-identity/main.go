package main

import (
	"fmt"
	"identity/github"
	"identity/yandex"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

type LocalUsersConfig struct {
	Root RootConfig
}

type RootConfig struct {
	Password string
}

type SSHConfig struct {
	ServerAccess ServerAccessConfig `json:"server_access"`
}
type ServerAccessConfig struct {
	Credentials CredetialsConfig
}

type CredetialsConfig struct {
	User       string
	PrivateKey string
	PublicKey  string
}

type GithubConfig struct {
	Managed bool
}

type YandexConfig struct {
	Managed bool
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		cfg := config.New(ctx, "")

		var sshCfg SSHConfig
		cfg.RequireObject("ssh", &sshCfg)

		var githubCfg GithubConfig
		cfg.RequireObject("github", &githubCfg)

		var yandexCfg YandexConfig
		cfg.RequireObject("yandex", &yandexCfg)

		var localUsers LocalUsersConfig
		cfg.RequireObject("local_users", &localUsers)

		if githubCfg.Managed {
			_ = github.ManageOrganization(ctx)
		}

		creds := make(map[string]interface{})
		creds["user"] = sshCfg.ServerAccess.Credentials.User
		creds["privatekey"] = pulumi.ToSecret(sshCfg.ServerAccess.Credentials.PrivateKey)
		creds["publickey"] = sshCfg.ServerAccess.Credentials.PublicKey

		ctx.Export("identity:ssh:server_access:credentials", pulumi.ToMap(creds))

		ctx.Export("identity:organization", pulumi.String(cfg.Require("organization")))

		ctx.Export("identity:local_users:root:password", pulumi.String(localUsers.Root.Password))

		if yandexCfg.Managed {
			y, err := yandex.ManageS3Helper(ctx)
			if err != nil {
				err = fmt.Errorf("error configure S3 Yandex service accounts: %w", err)
				ctx.Log.Error(err.Error(), nil)
				return err
			}

			ctx.Export("identity:yandex:s3", pulumi.ToMap(y))
		}

		return nil
	})
}
