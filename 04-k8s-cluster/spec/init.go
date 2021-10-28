package spec

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

type ClusterSpec struct {
	InternalDomainZone string
	ExternalIP string
}

func Init(ctx *pulumi.Context) (*ClusterSpec, error) {
	var pulumiSpecCfg ClusterSpec
	cfg := config.New(ctx, "")
	cfg.RequireSecretObject("spec", &pulumiSpecCfg)

	ctx.Export("cluster:spec:externalIP", pulumi.String(pulumiSpecCfg.ExternalIP))

	return &pulumiSpecCfg, nil
}
