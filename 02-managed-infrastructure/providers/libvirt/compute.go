package libvirt

import (
	"fmt"
	"managed-infrastructure/utils"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type ComputeConfig struct {
	Images     Images
	Network    Network
	Storage    Storage
	Hypevisors []*HypervisorConfig `json:"hypervisors"`
}

type HypervisorConfig struct {
	Name        string
	URI         string
	NetworkCIDR string `json:"networkcidr"`
	Network     HypervisorNetwork
	Machines    []*Machine
}

type Images struct {
	Base       string
	Combustion string
}
type Network struct {
	Name string
}

type Storage struct {
	Name string
}

type HypervisorNetwork struct {
	CIDR string
}

type Machine struct {
	ID string
	IP string
}

type ComputedInfra struct {
	nodes map[string]map[string]interface{}
}

func ManageCompute(ctx *pulumi.Context, sshCreds pulumi.Output, cfg *ComputeConfig) (*ComputedInfra, error) {
	nodes := make(map[string]map[string]interface{})
	fmt.Println(len(cfg.Hypevisors))
	for _, hypevisor := range cfg.Hypevisors {
		computedInfo, err := cfg.manage(ctx, hypevisor)
		if err != nil {
			return nil, err
		}

		for k, v := range computedInfo {
			v["key"] = utils.ExtractFromExportedMap(sshCreds, "privatekey")
			v["user"] = utils.ExtractFromExportedMap(sshCreds, "user")
			nodes[k] = v
		}
	}

	return &ComputedInfra{
		nodes: nodes,
	}, nil
}

func (v *ComputedInfra) GetNodes() map[string]map[string]interface{} {
	return v.nodes
}
