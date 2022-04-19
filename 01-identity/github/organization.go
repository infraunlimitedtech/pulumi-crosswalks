package github

import (
	"fmt"

	"github.com/pulumi/pulumi-github/sdk/v4/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var users = []map[string]string{
	{
		"username": "spigell",
		"role":     "admin",
	},
	{
		"username": "infraunlimitedBot",
		"role":     "member",
	},
}

func ManageOrganization(ctx *pulumi.Context) error {
	for _, v := range users {
		_, err := github.NewMembership(ctx, fmt.Sprintf("orgUser_%s_%s", v["username"], v["role"]), &github.MembershipArgs{
			Role:     pulumi.String(v["role"]),
			Username: pulumi.String(v["username"]),
		})
		if err != nil {
			return err
		}
	}
	return nil
}
