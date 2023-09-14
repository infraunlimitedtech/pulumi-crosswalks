package firewall

import (
	"managed-os/config"
	"managed-os/modules/firewall/firewalld"
	"managed-os/modules/wireguard"
	"pulumi-crosswalks/utils/hetzner"
	"strconv"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type FirewallCfg struct {
	Ctx                *pulumi.Context
	Nodes              []*config.Node
	InternalIface      string
	InfraLayerNodeInfo pulumi.AnyOutput
}

func (f *FirewallCfg) Manage(deps []map[string]pulumi.Resource) error {
	firewalls := make([]hetzner.Firewall, 0)

	for _, node := range f.Nodes {
		firewalldRules := make([]firewalld.Rule, 0)
		if len(node.Firewall.Hetzner) > 0 {
			for _, firewall := range node.Firewall.Hetzner {
				if !hetzner.Contains(firewalls, firewall.Name) {
					firewalls = append(firewalls, firewall)
				}
			}
		}
		if node.Firewall.Firewalld.Enabled {
			fwd, err := firewalld.New(f.Ctx, f.InfraLayerNodeInfo, node.ID, deps)
			firewalldRules = append(firewalldRules, firewalld.GetFirewallDSSHRule())

			for _, iface := range firewalld.GetWhitelistedIfaces() {
				firewalldRules = append(firewalldRules, &firewalld.InterfaceRule{
					Name:      "auto-iface-to-trusted-" + iface,
					Zone:      "trusted",
					Interface: iface,
				})
			}

			if len(node.Firewall.Firewalld.InternalZone.RestrictToSources) > 0 {
				for _, source := range node.Firewall.Firewalld.InternalZone.RestrictToSources {
					firewalldRules = append(firewalldRules, &firewalld.SourceRule{
						Name:   "auto-source-for-internal:" + source.Name,
						Zone:   "internal",
						Source: source.CIDR,
						Main:   source.Main,
					})
				}
			}

			if node.Wireguard.Firewall.Firewalld.Allowed {
				firewalldRules = append(firewalldRules, wireguard.GetRequiredFirewalldRule())
			}
			if len(firewalldRules) > 0 {
				if err != nil {
					return err
				}
				err = fwd.AddRules(firewalldRules)
				if err != nil {
					return err
				}
			}

			if node.Firewall.Firewalld.PublicZone.RemoveSSHService {
				err = fwd.RemoveService(&firewalld.ServiceRule{
					Name:    "auto-remove-ssh-from-public",
					Zone:    "public",
					Service: "ssh",
				})

				if err != nil {
					return err
				}

			}

			err = fwd.Reload()
			if err != nil {
				return err
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
