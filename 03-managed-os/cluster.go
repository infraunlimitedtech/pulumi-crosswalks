package main

import (
	"managed-os/config"
	"managed-os/modules/firewall"
	"managed-os/modules/k3s"
	"managed-os/modules/microos"
	"managed-os/modules/wireguard"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type cluster struct {
	ctx                *pulumi.Context
	PulumiConfig       *config.PulumiConfig
	InfraLayerNodeInfo pulumi.AnyOutput
	InternalIface      string
	RequiredPkgs       []string
	WgInfo             pulumi.AnyOutput
}

type definedCluster struct {
	Os        *microos.Cluster
	K3s       *k3s.Cluster
	Firewalls *firewall.FirewallCfg
	Wireguard *wireguard.Cluster
}

func defineCluster(c *cluster) (*definedCluster, error) {
	var leader *config.Node
	followers := make([]*config.Node, 0)

	for i := range c.PulumiConfig.Nodes.Agents {
		c.PulumiConfig.Nodes.Agents[i].Role = "agent"
		agent, err := config.Merge(&c.PulumiConfig.Nodes.Agents[i], c.PulumiConfig.Defaults)
		if err != nil {
			return nil, err
		}
		followers = append(followers, agent)
	}
	for i := range c.PulumiConfig.Nodes.Servers {
		c.PulumiConfig.Nodes.Servers[i].Role = "server"
		if c.PulumiConfig.Nodes.Servers[i].Leader {
			var err error
			leader, err = config.Merge(&c.PulumiConfig.Nodes.Servers[i], c.PulumiConfig.Defaults)
			if err != nil {
				return nil, err
			}
		} else {
			server, err := config.Merge(&c.PulumiConfig.Nodes.Servers[i], c.PulumiConfig.Defaults)
			if err != nil {
				return nil, err
			}
			followers = append(followers, server)
		}
	}

	// WHY append ?!
	if err := config.Validate(append(followers, leader), c.InfraLayerNodeInfo); err != nil {
		return nil, err
	}

	allNodes := followers
	allNodes = append(allNodes, leader)

	return &definedCluster{
		Os: &microos.Cluster{
			Ctx:                c.ctx,
			Nodes:              allNodes,
			RequiredPkgs:       c.RequiredPkgs,
			InfraLayerNodeInfo: c.InfraLayerNodeInfo,
		},
		Wireguard: &wireguard.Cluster{
			Ctx:                c.ctx,
			Nodes:              allNodes,
			Iface:              c.InternalIface,
			InfraLayerNodeInfo: c.InfraLayerNodeInfo,
			Info:               c.WgInfo,
		},
		K3s: &k3s.Cluster{
			Ctx:                c.ctx,
			Leader:             leader,
			InfraLayerNodeInfo: c.InfraLayerNodeInfo,
			Iface:              c.InternalIface,
			Followers:          followers,
		},
		Firewalls: &firewall.FirewallCfg{
			Ctx:                c.ctx,
			InternalIface:      c.InternalIface,
			InfraLayerNodeInfo: c.InfraLayerNodeInfo,
			Nodes:              allNodes,
		},
	}, nil
}
