package k3s

import (
	"fmt"
	"managed-os/utils"
	"path"

	"github.com/pulumi/pulumi-command/sdk/go/command/remote"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (c *Cluster) install(deps []map[string]pulumi.Resource) (map[string]pulumi.Resource, error) {
	nodes := c.Followers
	nodes = append(nodes, c.Leader)

	result := make(map[string]pulumi.Resource)

	for _, node := range nodes {
		k3sExec := node.Role

		installed, err := remote.NewCommand(c.Ctx, fmt.Sprintf("%s-installK3s", node.ID), &remote.CommandArgs{
			Connection: &remote.ConnectionArgs{
				Host:       utils.ExtractValueFromPulumiMapMap(c.InfraLayerNodeInfo, node.ID, "ip"),
				User:       utils.ExtractValueFromPulumiMapMap(c.InfraLayerNodeInfo, node.ID, "user"),
				PrivateKey: utils.ExtractValueFromPulumiMapMap(c.InfraLayerNodeInfo, node.ID, "key"),
			},
			Environment: pulumi.StringMap{
				"INSTALL_K3S_SKIP_START":       pulumi.String("true"),
				"INSTALL_K3S_SKIP_SELINUX_RPM": pulumi.String("true"),
				"INSTALL_K3S_VERSION":          pulumi.String(node.K3s.Version),
				"INSTALL_K3S_EXEC":             pulumi.String(k3sExec),
			},
			Create: pulumi.Sprintf("sudo mkdir -p %s && if [[ -e /usr/local/bin/k3s ]]; then restart=true; fi ; curl -sfL https://get.k3s.io | sudo -E sh - 2>&1 >> /tmp/k3s-pulumi.log && if [[ $restart ]]; then sudo systemctl restart k3s*; fi ",
				path.Dir(cfgPath)),
			Delete: pulumi.String("/usr/local/bin/k3s-killall.sh"),
		}, pulumi.ReplaceOnChanges([]string{"create", "environment"}),
			pulumi.DependsOn(utils.ConvertMapSliceToSliceByKey(deps, node.ID)),
			pulumi.RetainOnDelete(!node.K3s.CleanDataOnUpgrade),
		)
		if err != nil {
			err = fmt.Errorf("error install a k3s cluster via script: %w", err)
			return nil, err
		}
		result[node.ID] = installed
	}

	return result, nil
}
