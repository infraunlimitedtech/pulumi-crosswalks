package hetzner

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"pulumi-crosswalks/utils/hetzner"
	"strconv"
	"strings"

	"github.com/imdario/mergo"

	"github.com/pulumi/pulumi-hcloud/sdk/go/hcloud"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	// Base snapshot for microOs.
	provider          = "hetzner"
	defaultImage      = "108334614"
	defaultServerType = "cpx11"
	defaultLocation   = "hel1"
)

type ComputeConfig struct {
	Configuration *ConfigurationConfig
	sshCreds      map[string]string
}

type ConfigurationConfig struct {
	Firewall    []hetzner.Firewall
	ClusterName string
	Servers     *ServersConfig
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

func ManageCompute(ctx *pulumi.Context, sshCreds map[string]string, cfg *ComputeConfig) (*ComputedInfra, error) {
	nodes := make(map[string]map[string]interface{})
	cfg.sshCreds = sshCreds
	cfg.Configuration.ClusterName = ctx.Stack()

	computedInfo, err := cfg.manage(ctx)
	if err != nil {
		return nil, err
	}

	for k, v := range computedInfo {
		v["key"] = pulumi.ToSecret(sshCreds["privatekey"])
		v["user"] = sshCreds["user"]
		nodes[k] = v
	}

	return &ComputedInfra{
		nodes: nodes,
	}, nil
}

func (i *ComputeConfig) manage(ctx *pulumi.Context) (map[string]map[string]interface{}, error) {
	nodes := make(map[string]map[string]interface{})

	var fwRules pulumi.IntArray

	if len(i.Configuration.Firewall) > 0 {
		firewalls, err := hetzner.NewFirewalls(ctx, i.Configuration.Firewall)
		if err != nil {
			return nodes, err
		}

		for _, fw := range firewalls.Items {
			conv := fw.GetID().ToStringOutput().ApplyT(func(id string) (int, error) {
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
			GrowPart: &GrowPartConfig{
				Devices: []string{
					"/var",
				},
			},
			Users: []*UserCloudConfig{
				{
					Name: i.sshCreds["user"],
					Sudo: "ALL=(ALL) NOPASSWD:ALL",
					SSHAuthorizedKeys: []string{
						i.sshCreds["publickey"],
					},
					// Need to be encrypted
					// k3s4ever
					Passwd: "$6$y184E3Ic0PiyOHcN$BguLDhU9rb6uqv6P.g/22ViZnzapXM5ukg/zYASxA3zB43gx30XG73OGZhEt07GSKc4RIsefMcaSNsoXmrqxI1",
				},
			},
		}

		ud, err := userdata.render()
		if err != nil {
			return nodes, err
		}

		serverName := fmt.Sprintf("%s.%s.%s", srv.ID, i.Configuration.ClusterName, "infraunlimited.tech")

		args := &hcloud.ServerArgs{
			ServerType: pulumi.String(srv.ServerType),
			Location:   pulumi.String(srv.Location),
			Name:       pulumi.String(serverName),
			UserData:   pulumi.String(ud),
		}

		switch image := srv.Image; image {
		case "automation-api":
			var automationAPIResponce AutomationAPIResponce
			url := url.URL{
				Scheme:   "http",
				Host:     os.Getenv("AUTOMATION_API_HTTP_ADDR"),
				Path:     "hetzner/snapshots",
				RawQuery: fmt.Sprintf("server=%s", serverName),
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

		created, err := hcloud.NewServer(ctx, serverName, args)
		if err != nil {
			return nodes, err
		}
		nodes[srv.ID] = make(map[string]interface{})
		nodes[srv.ID]["provider"] = provider
		nodes[srv.ID]["ip"] = created.Ipv4Address
		nodes[srv.ID]["id"] = created.ID()
	}

	return nodes, nil
}

func (v *ComputedInfra) GetNodes() map[string]map[string]interface{} {
	return v.nodes
}
