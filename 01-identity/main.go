package main

import (
	"github.com/pulumi/pulumi-github/sdk/v4/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"os"
)

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
		return nil
	})
}
