package spec

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

type ClusterSpec struct {
	InternalDomainZone string
}

func Init(ctx *pulumi.Context) (*ClusterSpec, error) {
	var pulumiSpecCfg ClusterSpec
	cfg := config.New(ctx, "")
	cfg.RequireSecretObject("spec", &pulumiSpecCfg)

	return &pulumiSpecCfg, nil
}
