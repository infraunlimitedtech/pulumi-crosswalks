package github

import (
	"github.com/pulumi/pulumi-github/sdk/v4/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func SetServerPublicKey(ctx *pulumi.Context, serverKey string) error {

	_, err := github.NewUserSshKey(ctx, "serverKey", &github.UserSshKeyArgs{
		Title: pulumi.String("server-key"),
		Key:   pulumi.String(serverKey),
	})
	if err != nil {
		return err
	}

	return nil
}
