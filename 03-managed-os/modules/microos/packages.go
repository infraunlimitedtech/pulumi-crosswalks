package microos

import (
	"fmt"
	"managed-os/utils"
	"strings"

	"github.com/pulumi/pulumi-command/sdk/go/command/remote"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (o *Cluster) InstallRequiredPkgs() (map[string]pulumi.Resource, error) {
	pkgs := make(map[string]pulumi.Resource)

	for _, node := range o.Nodes {
		cmd := fmt.Sprintf("sudo transactional-update -n pkg install %s", strings.Join(o.RequiredPkgs, " "))

		installed, err := remote.NewCommand(o.Ctx, fmt.Sprintf("%s-installPackages", node.ID), &remote.CommandArgs{
			Connection: &remote.ConnectionArgs{
				Host:       utils.ExtractValueFromPulumiMapMap(o.InfraLayerNodeInfo, node.ID, "ip"),
				User:       utils.ExtractValueFromPulumiMapMap(o.InfraLayerNodeInfo, node.ID, "user"),
				PrivateKey: utils.ExtractValueFromPulumiMapMap(o.InfraLayerNodeInfo, node.ID, "key"),
			},
			Create: pulumi.String(cmd),
		})
		if err != nil {
			err = fmt.Errorf("error install needed packages: %w", err)
			return nil, err
		}

		pkgs[node.ID] = installed
	}

	return pkgs, nil
}
