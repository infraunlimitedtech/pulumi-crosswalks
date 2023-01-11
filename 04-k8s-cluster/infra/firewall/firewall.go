package firewall

import (
	"errors"
	"fmt"
	"strconv"
	"pulumi-crosswalks/utils"
	"pulumi-crosswalks/utils/firewalld"
	"pulumi-crosswalks/utils/hetzner"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var (
	errUnknownFirewallRuleType = errors.New("unknown firewall rule type")
)

func Manage(ctx *pulumi.Context, nodes pulumi.AnyOutput, rules []utils.FirewallRule) error {
	hrules := make([]hetzner.Rule, 0)
	for _, rule := range rules {
		switch v := rule.(type) {
		case hetzner.Rule:
			hrules = append(hrules, v)

		case firewalld.Rule:
			continue
			
		default:
			return errUnknownFirewallRuleType
		}
	}

	firewalls := []hetzner.Firewall{{ 
		Rules: hrules, Name: fmt.Sprintf("%s-%s", ctx.Project(), ctx.Stack()),
	}}

	f, err := hetzner.NewFirewalls(ctx, firewalls)
	if err != nil {
		return err
	}

	m := make(map[string]pulumi.IntArrayOutput)
	for _, firewall := range f.Items {
		ids := nodes.ApplyT(func(nodes interface{}) []int {
			ids := make([]int, 0)
			for _, value := range nodes.(map[string]interface{}) {
				node, _ := value.(map[string]interface{})
				s, _ := strconv.Atoi(node["id"].(string))
				ids = append(ids, s)
			}
			return ids
		}).(pulumi.IntArrayOutput)


		m[firewall.Name] = ids
	}
	return f.Attach(m)
}
