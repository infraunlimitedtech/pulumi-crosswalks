package vagrant

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"managed-infrastructure/utils"
)

type Infra struct {
	nodes map[string]map[string]interface{}
}

// Filled by hands.
var vagrantNodes = []map[string]string{
	{
		"id": "k3s-server01",
		"ip": "192.168.99.135",
	},
	{
		"id": "k3s-server02",
		"ip": "192.168.99.136",
	},
	{
		"id": "k3s-server03",
		"ip": "192.168.99.137",
	},
	{
		"id": "k3s-agent01",
		"ip": "192.168.99.140",
	},
	{
		"id": "k3s-agent02",
		"ip": "192.168.99.141",
	},
}

func Init(sshCreds pulumi.Output) *Infra {
	nodes := make(map[string]map[string]interface{})

	for _, node := range vagrantNodes {
		nodes[node["id"]] = map[string]interface{}{
			"key":  utils.ExtractFromExportedMap(sshCreds, "privatekey"),
			"user": utils.ExtractFromExportedMap(sshCreds, "user"),
			"ip":   node["ip"],
		}
	}
	return &Infra{
		nodes: nodes,
	}
}

func (v *Infra) GetNodes() map[string]map[string]interface{} {
	return v.nodes
}
