package hetzner

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"managed-infrastructure/utils"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/imdario/mergo"

	"github.com/pulumi/pulumi-hcloud/sdk/go/hcloud"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	// Base snapshot for microOs.
	defaultImage      = "75477320"
	defaultServerType = "cpx11"
	defaultLocation   = "hel1"
)

type ComputeConfig struct {
	Configuration *ConfigurationConfig

	sshCreds pulumi.Output
}

type ConfigurationConfig struct {
	Firewall *FirewallConfig
	Servers  *ServersConfig
}

type FirewallConfig struct {
	Rules []*FirewallRule
}

type FirewallRule struct {
	Enabled  bool
	Name     string
	Port     string
	Protocol string
	Allowed  []string
}

type ServersConfig struct {
	Defaults *Machine
	Machines []*Machine
}

type Machine struct {
	ServerType string `json:"server_type"`
	ID         string
	Location   string
	Image      string
}

type ComputedInfra struct {
	nodes map[string]map[string]interface{}
}

type AutomationAPIResponce struct {
	Body   Body `json:"body"`
	Error  string
	Status string
	ID     string
}

type Body struct {
	ID int
}

func ManageCompute(ctx *pulumi.Context, sshCreds pulumi.Output, cfg *ComputeConfig) (*ComputedInfra, error) {
	nodes := make(map[string]map[string]interface{})
	cfg.sshCreds = sshCreds

	computedInfo, err := cfg.manage(ctx)
	if err != nil {
		return nil, err
	}

	for k, v := range computedInfo {
		v["key"] = utils.ExtractFromExportedMap(sshCreds, "privatekey")
		v["user"] = utils.ExtractFromExportedMap(sshCreds, "user")
		nodes[k] = v
	}

	return &ComputedInfra{
		nodes: nodes,
	}, nil
}

func (i *ComputeConfig) manage(ctx *pulumi.Context) (map[string]map[string]interface{}, error) {
	nodes := make(map[string]map[string]interface{})

	var fwRules pulumi.IntArray

	if len(i.Configuration.Firewall.Rules) > 0 {
		for _, rule := range i.Configuration.Firewall.Rules {
			args := &hcloud.FirewallRuleArgs{
				Direction: pulumi.String("in"),
				Protocol:  pulumi.String(rule.Protocol),
				SourceIps: pulumi.ToStringArray(rule.Allowed),
			}
			if rule.Port != "" {
				args.Port = pulumi.String(rule.Port)
			}

			fwRule, err := hcloud.NewFirewall(ctx, rule.Name, &hcloud.FirewallArgs{
				Name: pulumi.String(rule.Name),
				Rules: hcloud.FirewallRuleArray{
					args,
				},
			})
			if err != nil {
				return nodes, err
			}

			conv := fwRule.ID().ToStringOutput().ApplyT(func(id string) (int, error) {
				return strconv.Atoi(strings.Split(id, "-")[0])
			}).(pulumi.IntOutput)

			fwRules = append(fwRules, conv)
		}
	}

	for _, srv := range i.Configuration.Servers.Machines {
		if err := mergo.Merge(srv, i.Configuration.Servers.Defaults, mergo.WithAppendSlice); err != nil {
			return nodes, err
		}

		if srv.ServerType == "" {
			srv.ServerType = defaultServerType
		}

		if srv.Location == "" {
			srv.Location = defaultLocation
		}

		userdata := &CloudConfig{
			FQDN: fmt.Sprintf("%s.%s", srv.ID, "infraunlimited.tech"),
			GrowPart: &GrowPartConfig{
				Devices: []string{
					"/var",
				},
			},
		}

		args := &hcloud.ServerArgs{
			ServerType: pulumi.String(srv.ServerType),
			Location:   pulumi.String(srv.Location),
			Name:       pulumi.String(srv.ID),
			UserData: i.sshCreds.ApplyT(func(v interface{}) string {
				m := v.(map[string]interface{})
				userdata.Users = []*UserCloudConfig{
					{
						Name: m["user"].(string),
						Sudo: "ALL=(ALL) NOPASSWD:ALL",
						SSHAuthorizedKeys: []string{
							m["publickey"].(string),
						},
					},
				}
				rendered, _ := userdata.render()
				return rendered
			}).(pulumi.StringOutput),
		}

		switch image := srv.Image; image {
		case "automation-api":
			var automationAPIResponce AutomationAPIResponce
			url := url.URL{
				Scheme:   "http",
				Host:     os.Getenv("AUTOMATION_API_HTTP_ADDR"),
				Path:     "hetzner/snapshots",
				RawQuery: fmt.Sprintf("server=%s", srv.ID),
			}

			req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url.String(), nil)
			if err != nil {
				return nodes, err
			}

			cli := &http.Client{}

			resp, err := cli.Do(req)
			if err != nil {
				return nodes, err
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nodes, err
			}

			err = json.Unmarshal(body, &automationAPIResponce)

			if err != nil {
				return nodes, err
			}

			switch resp.StatusCode {
			case http.StatusOK:
				args.Image = pulumi.String(strconv.Itoa(automationAPIResponce.Body.ID))

			case http.StatusNotFound:
				args.Image = pulumi.String(defaultImage)

			default:
				return nodes, fmt.Errorf("bad status code, error: %s", automationAPIResponce.Error)
			}

			resp.Body.Close()

		case "":
			args.Image = pulumi.String(defaultImage)
			ctx.Log.Warn(fmt.Sprintf("Will use default image for %s", srv.ID), nil)

		default:
			args.Image = pulumi.String(image)
		}

		if len(fwRules) > 0 {
			args.FirewallIds = fwRules
		}

		created, err := hcloud.NewServer(ctx, srv.ID, args)
		if err != nil {
			return nodes, err
		}
		nodes[srv.ID] = make(map[string]interface{})
		nodes[srv.ID]["ip"] = created.Ipv4Address
		nodes[srv.ID]["id"] = created.ID()
	}

	return nodes, nil
}

func (v *ComputedInfra) GetNodes() map[string]map[string]interface{} {
	return v.nodes
}
