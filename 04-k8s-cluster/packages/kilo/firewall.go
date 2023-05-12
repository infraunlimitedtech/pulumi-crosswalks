package kilo

import (
	"pulumi-crosswalks/utils"
	"pulumi-crosswalks/utils/hetzner"
	"strconv"
)

func (k *StartedKilo) GetRequiredFirewallRules() []utils.FirewallRule {
	rules := make([]utils.FirewallRule, 0)
	if k.Firewalls.Hetzner != nil && k.Firewalls.Hetzner.Managed {
		rules = append(rules, hetzner.Rule{
			Direction: "in",
			Protocol:  "udp",
			SourceIps: []string{"0.0.0.0/0"},
			Port:      strconv.Itoa(k.Port),
		},
		)
	}

	return rules
}
