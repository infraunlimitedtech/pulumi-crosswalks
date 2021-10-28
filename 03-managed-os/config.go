package main

import (
	"encoding/base64"
	"fmt"
	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
	k3os "github.com/spigell/pulumi-k3os/provider/pkg/resources"
)

type pulumiConfig struct {
	InfraStack string `json:"infra_stack"`
	Defaults   Defaults
	Nodes      Nodes
}

type Defaults struct {
	Global  k3os.NodeConfig
	Servers k3os.NodeConfig
	Agents  k3os.NodeConfig
}

type Nodes struct {
	Servers []Node
	Agents  []Node
}

type Node struct {
	ID        string
	Leader    bool
	Wireguard Wireguard
	PublicIP  string `json:"public_ip"`
	Config    k3os.NodeConfig
	kind      string
}

var (
	wgFilePath        = "/etc/wireguard/kubewg0.conf"
	additionalModules = []string{"wireguard"}
	additionalK3sArgs = []string{"--flannel-iface=kubewg0"}
	additionalRunCmd  = []string{"sudo wg-quick up kubewg0"}
)

func parseConfig(ctx *pulumi.Context) *pulumiConfig {
	var pulumiCfg pulumiConfig
	cfg := config.New(ctx, "")
	cfg.RequireSecretObject("main", &pulumiCfg)

	return &pulumiCfg

}

func buildNodeConfig(cluster *cluster, pulumiCfg *pulumiConfig, node Node) (*k3os.NodeConfig, error) {
	wgPeers := buildWgPeers(append(cluster.followers, cluster.leader), node)
	wgConfig, err := renderWgConfig(wgPeers, node)
	if err != nil {
		return nil, errors.Wrap(err, "rendering wireguard config")
	}
	wgFile := k3os.CloudInitFile{
		Encoding: "b64",
		Content:  base64.StdEncoding.EncodeToString([]byte(wgConfig)),
		Path:     wgFilePath,
	}

	config, err := mergeVars(node, pulumiCfg.Defaults)
	if err != nil {
		return nil, errors.Wrap(err, "merge variables")
	}

	config.WriteFiles = append(config.WriteFiles, wgFile)
	config.K3os.Modules = append(config.K3os.Modules, additionalModules...)
	config.K3os.K3sArgs = append(config.K3os.K3sArgs, additionalK3sArgs...)
	config.Runcmd = append(config.Runcmd, additionalRunCmd...)
	if node.kind == serverStr {
		config.K3os.K3sArgs = append(config.K3os.K3sArgs, fmt.Sprintf("--bind-address=%s", node.Wireguard.PrivateAddr))
	}
	config.K3os.K3sArgs = append(config.K3os.K3sArgs, fmt.Sprintf("--node-ip=%s", node.Wireguard.PrivateAddr))
	if node.Leader {
		config.K3os.K3sArgs = append(config.K3os.K3sArgs, "--cluster-init")
	}

	return config, nil
}

func mergeVars(node Node, defaults Defaults) (*k3os.NodeConfig, error) {
	nodeConfig := &node.Config
	if err := mergo.Merge(nodeConfig, defaults.Global, mergo.WithAppendSlice); err != nil {

		return nil, err
	}

	switch kind := node.kind; kind {
	case serverStr:
		if err := mergo.Merge(nodeConfig, defaults.Servers, mergo.WithAppendSlice); err != nil {
			return nil, err
		}
	case agentStr:
		if err := mergo.Merge(nodeConfig, defaults.Agents, mergo.WithAppendSlice); err != nil {
			return nil, err
		}
	}
	return nodeConfig, nil

}
