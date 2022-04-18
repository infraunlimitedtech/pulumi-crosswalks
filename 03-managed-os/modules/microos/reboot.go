package microos

import (
	"fmt"
	"managed-os/utils"

	"github.com/pulumi/pulumi-command/sdk/go/command/local"
	"github.com/pulumi/pulumi-command/sdk/go/command/remote"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (o *Cluster) Reboot(deps []map[string]pulumi.Resource) (map[string]pulumi.Resource, error) {
	m := make(map[string]pulumi.Resource)

	for _, node := range o.Nodes {
		rebooted, err := remote.NewCommand(o.Ctx, fmt.Sprintf("%s-Reboot", node.ID), &remote.CommandArgs{
			Connection: &remote.ConnectionArgs{
				Host:       utils.ExtractValueFromPulumiMapMap(o.InfraLayerNodeInfo, node.ID, "ip"),
				User:       utils.ExtractValueFromPulumiMapMap(o.InfraLayerNodeInfo, node.ID, "user"),
				PrivateKey: utils.ExtractValueFromPulumiMapMap(o.InfraLayerNodeInfo, node.ID, "key"),
			},
			Create: pulumi.String("(sleep 1 && sudo shutdown -r now) &"),
		}, pulumi.DependsOn(utils.ConvertMapSliceToSliceByKey(deps, node.ID)))
		if err != nil {
			err = fmt.Errorf("error reboot node: %w", err)
			return nil, err
		}

		waited, _ := local.NewCommand(o.Ctx, fmt.Sprintf("%s-localWait", node.ID), &local.CommandArgs{
			Create: pulumi.Sprintf("until nc -z %s 22; do sleep 5; done", utils.ExtractValueFromPulumiMapMap(o.InfraLayerNodeInfo, node.ID, "ip")),
			Triggers: pulumi.Array{
				rebooted,
			},
		}, pulumi.DependsOn([]pulumi.Resource{rebooted}),
			pulumi.Timeouts(&pulumi.CustomTimeouts{Create: "10m"}),
		)

		m[node.ID] = waited
	}

	return m, nil
}
