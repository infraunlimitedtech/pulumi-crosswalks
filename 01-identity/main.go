package main

import (
	"github.com/pulumi/pulumi-github/sdk/v4/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
	"os"
)

type SSHConfig struct {
	Credentials CredetialsConfig
}

type CredetialsConfig struct {
	User       string
	PrivateKey string
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		serverKey, _ := os.ReadFile("ssh/keys/server.pub")
		_, err := github.NewUserSshKey(ctx, "serverKey", &github.UserSshKeyArgs{
			Title: pulumi.String("server-key"),
			Key:   pulumi.String(serverKey),
		})
		if err != nil {
			return err
		}

		var sshCfg SSHConfig
		cfg := config.New(ctx, "")
		cfg.RequireObject("ssh", &sshCfg)

		creds := make(map[string]interface{})
		creds["user"] = sshCfg.Credentials.User
		creds["privatekey"] = pulumi.ToSecret(sshCfg.Credentials.PrivateKey)

		ctx.Export("identity:ssh:credentials", pulumi.ToMap(creds))
		return nil
	})
}
