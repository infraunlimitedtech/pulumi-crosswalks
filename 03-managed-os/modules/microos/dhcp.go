package microos

import (
	"fmt"
	"managed-os/utils"

	"github.com/pulumi/pulumi-command/sdk/go/command/remote"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (o *Cluster) ConfigureDHCPClient() (map[string]pulumi.Resource, error) {
	res := make(map[string]pulumi.Resource)

	for _, node := range o.Nodes {
		cmd := "sudo sed -i.bak 's/SET_HOSTNAME=\"yes\"/SET_HOSTNAME=\"no\"/' /etc/sysconfig/network/dhcp"

		configured, err := remote.NewCommand(o.Ctx, fmt.Sprintf("%s-ConfigureDHCPClient", node.ID), &remote.CommandArgs{
			Connection: &remote.ConnectionArgs{
				Host:       utils.ExtractValueFromPulumiMapMap(o.InfraLayerNodeInfo, node.ID, "ip"),
				User:       utils.ExtractValueFromPulumiMapMap(o.InfraLayerNodeInfo, node.ID, "user"),
				PrivateKey: utils.ExtractValueFromPulumiMapMap(o.InfraLayerNodeInfo, node.ID, "key"),
			},
			Create: pulumi.String(cmd),
		})
		if err != nil {
			err = fmt.Errorf("error command execution: %w", err)
			return nil, err
		}

		res[node.ID] = configured
	}

	return res, nil
}
