package hetzner

import (
	"strconv"
	"strings"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"github.com/pulumi/pulumi-hcloud/sdk/go/hcloud"
)

type Firewall struct {
	id    pulumi.IDOutput
	Name  string
	Rules []Rule
}

type Rule struct {
	Protocol    string
	Port        string
	SourceIps   []string
	Direction   string
	Description string
}

type Firewalls struct {
	ctx   *pulumi.Context
	Items []Firewall
}

func NewFirewalls(ctx *pulumi.Context, firewalls []Firewall) (*Firewalls, error) {
	f := &Firewalls{
		ctx: ctx,
	}

	for _, firewall := range firewalls {
		var rules hcloud.FirewallRuleArray

		for _, rule := range firewall.Rules {
			if rule.Protocol == "" {
				rule.Protocol = "tcp"
			}

			rules = append(rules, hcloud.FirewallRuleArgs{
				Direction:   pulumi.String("in"),
				Description: pulumi.String(rule.Description),
				Protocol:    pulumi.String(rule.Protocol),
				Port:        pulumi.String(rule.Port),
				SourceIps:   pulumi.ToStringArray(rule.SourceIps),
			})
		}
		created, err := hcloud.NewFirewall(ctx, firewall.Name, &hcloud.FirewallArgs{
			Name:  pulumi.String(firewall.Name),
			Rules: rules,
		})
		if err != nil {
			return nil, err
		}
		firewall.id = created.ID()
		f.Items = append(f.Items, firewall)
	}

	return f, nil
}

func (l *Firewalls) Attach(ids map[string]pulumi.IntArrayOutput) error {
	for name, v := range ids {
		for _, firewall := range l.Items {
			if firewall.Name == name {
				_, err := hcloud.NewFirewallAttachment(l.ctx, name, &hcloud.FirewallAttachmentArgs{
					FirewallId: firewall.id.ToStringOutput().ApplyT(func(id string) (int, error) {
						return strconv.Atoi(strings.Split(id, "-")[0])
					}).(pulumi.IntOutput),
					ServerIds: v,
				})
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (f *Firewall) GetID() pulumi.IDOutput {
	return f.id
}

func Contains(firewalls []Firewall, name string) bool {
	for _, firewall := range firewalls {
		if firewall.Name == name {
			return true
		}
	}

	return false
}
