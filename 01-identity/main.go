package main

import (
	"encoding/base64"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
	"identity/github"
)

type SSHConfig struct {
	ServerAccess ServerAccessConfig `json:"server_access"`
}
type ServerAccessConfig struct {
	UploadToGithub bool `json:"upload_to_github"`
	Credentials    CredetialsConfig
}

type CredetialsConfig struct {
	User       string
	PrivateKey string
	PublicKey  string
}

type GithubConfig struct {
	Managed bool
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		cfg := config.New(ctx, "")

		var sshCfg SSHConfig
		cfg.RequireObject("ssh", &sshCfg)

		var githubCfg GithubConfig
		cfg.RequireObject("github", &githubCfg)

		if githubCfg.Managed {
			_ = github.ManageOrganization(ctx)
		}

		if sshCfg.ServerAccess.UploadToGithub {
			decoded, err := base64.StdEncoding.DecodeString(sshCfg.ServerAccess.Credentials.PublicKey)
			if err != nil {
				return err
			}
			_ = github.SetServerPublicKey(ctx, string(decoded))
		}

		creds := make(map[string]interface{})
		creds["user"] = sshCfg.ServerAccess.Credentials.User
		creds["privatekey"] = pulumi.ToSecret(sshCfg.ServerAccess.Credentials.PrivateKey)

		ctx.Export("identity:ssh:server_access:credentials", pulumi.ToMap(creds))

		return nil
	})
}
