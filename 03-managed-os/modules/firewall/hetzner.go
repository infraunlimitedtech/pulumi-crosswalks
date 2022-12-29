package firewall

import (
	"managed-os/config"
	"pulumi-crosswalks/utils/hetzner"
	"strconv"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type FirewallCfg struct {
	Ctx                *pulumi.Context
	Nodes              []*config.Node
	InfraLayerNodeInfo pulumi.AnyOutput
}

func (f *FirewallCfg) Manage() error {
	firewalls := make([]hetzner.Firewall, 0)

	for _, node := range f.Nodes {
		if len(node.Firewall.Hetzner) > 0 {
			for _, firewall := range node.Firewall.Hetzner {
				if !hetzner.Contains(firewalls, firewall.Name) {
					firewalls = append(firewalls, firewall)
				}
			}
		}
	}

	n, err := hetzner.NewFirewalls(f.Ctx, firewalls)

	nodesByFirewallRules := make(map[string]pulumi.IntArrayOutput)
	for _, node := range f.Nodes {
		if err != nil {
			return err
		}

		if len(node.Firewall.Hetzner) > 0 {
			for _, firewall := range node.Firewall.Hetzner {
				ids := f.InfraLayerNodeInfo.ApplyT(func(nodes interface{}) []int {
					ids := make([]int, 0)
					for _, value := range nodes.(map[string]interface{}) {
						node, _ := value.(map[string]interface{})
						parsed, _ := strconv.Atoi(node["id"].(string))
						ids = append(ids, parsed)
					}
					return ids
				}).(pulumi.IntArrayOutput)

				nodesByFirewallRules[firewall.Name] = ids
			}
		}
	}
	err = n.Attach(nodesByFirewallRules)
	if err != nil {
		return err
	}
	return nil
}
