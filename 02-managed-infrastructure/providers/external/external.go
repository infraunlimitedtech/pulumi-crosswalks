package external

import "github.com/pulumi/pulumi/sdk/v3/go/pulumi"

type Infra struct {
	nodes map[string]map[string]interface{}
}

type ComputeConfig []Machine

type Machine struct {
	ID string
	IP string
}

func Init(sshCreds map[string]string, config *ComputeConfig) *Infra {
	nodes := make(map[string]map[string]interface{})

	for _, node := range *config {
		nodes[node.ID] = map[string]interface{}{
			"id":       node.ID,
			"provider": "static",
			"key":      pulumi.ToSecret(sshCreds["privatekey"]),
			"user":     sshCreds["user"],
			"ip":       node.IP,
		}
	}
	return &Infra{
		nodes: nodes,
	}
}

func (v *Infra) GetNodes() map[string]map[string]interface{} {
	return v.nodes
}
