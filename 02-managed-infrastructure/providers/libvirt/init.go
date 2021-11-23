package libvirt

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Config []HypervisorConfig

type HypervisorConfig struct {
	Name        string
	URI         string
	NetworkCIDR string `json:"networkcidr"`
	Machines    []Machine
}

type Machine struct {
	ID string
}

type Infra struct {
	nodes map[string]map[string]interface{}
}

func Init(ctx *pulumi.Context, user string, key pulumi.Output, cfg *Config) (*Infra, error) {
	nodes := make(map[string]map[string]interface{})
	for _, hypevisor := range *cfg {
		computedInfo, err := manageLibvirtHost(ctx, hypevisor)

		if err != nil {
			return nil, err
		}

		for k, v := range computedInfo {
			v["key"] = key
			v["user"] = user
			nodes[k] = v
		}
	}
	return &Infra{
		nodes: nodes,
	}, nil
}

func (v *Infra) GetNodes() map[string]map[string]interface{} {
	return v.nodes
}
