package config

import (
	"github.com/imdario/mergo"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func ParseConfig(ctx *pulumi.Context) *PulumiConfig {
	var pulumiCfg PulumiConfig
	cfg := config.New(ctx, "")
	cfg.RequireSecretObject("main", &pulumiCfg)

	return &pulumiCfg
}

func Merge(nodeConfig *Node, defaults *Defaults) (*Node, error) {
	if err := mergo.Merge(nodeConfig, defaults.Global, mergo.WithAppendSlice); err != nil {
		return nodeConfig, err
	}

	switch role := nodeConfig.Role; role {
	case "server":
		if err := mergo.Merge(nodeConfig, defaults.Servers, mergo.WithAppendSlice); err != nil {
			return nodeConfig, err
		}
	case "agent":
		if err := mergo.Merge(nodeConfig, defaults.Agents, mergo.WithAppendSlice); err != nil {
			return nodeConfig, err
		}
	}
	return nodeConfig, nil
}
