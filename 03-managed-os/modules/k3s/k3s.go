package k3s

import (
	"fmt"
	"managed-os/config"
	"managed-os/utils"

	"github.com/pulumi/pulumi-command/sdk/go/command/remote"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	remotefile "github.com/spigell/pulumi-file/sdk/go/file/remote"
)

const (
	cfgPath = "/etc/rancher/k3s/config.yaml"
)

type Cluster struct {
	Iface              string
	ServerURL          string
	InfraLayerNodeInfo pulumi.AnyOutput
	Followers          []*config.Node
	Leader             *config.Node
	Ctx                *pulumi.Context
}

type CreatedCluster struct {
	Kubeconfig *pulumi.StringOutput
}

func GetRequiredPkgs() []string {
	return []string{"k3s-selinux"}
}

func GetRequirdSSHDConfig() map[string]string {
	return map[string]string{"AcceptEnv": "INSTALL_K3S_*"}
}

func (c *Cluster) Manage(WgPeers pulumi.AnyOutput, deps []map[string]pulumi.Resource) (*CreatedCluster, error) {
	installed, err := c.install(deps)
	if err != nil {
		return nil, err
	}

	configured, err := c.configure(WgPeers, []map[string]pulumi.Resource{installed})
	if err != nil {
		return nil, err
	}

	// Need to improve the restart
	// err = c.restart(deps)

	kubeConfig, err := c.grabKubeConfig([]map[string]pulumi.Resource{configured})
	if err != nil {
		return nil, err
	}

	return &CreatedCluster{
		Kubeconfig: kubeConfig,
	}, nil
}

func (c *Cluster) restart(deps []map[string]pulumi.Resource) error {
	nodes := c.Followers
	nodes = append(nodes, c.Leader)

	for _, node := range nodes {

		triggers, err := depsToCastedArray(utils.ConvertMapSliceToSliceByKey(deps, node.ID))
		if err != nil {
			return err
		}

		_, err = remote.NewCommand(c.Ctx, fmt.Sprintf("RestartK3s-%s", node.ID), &remote.CommandArgs{
			Connection: &remote.ConnectionArgs{
				Host:       utils.ExtractValueFromPulumiMapMap(c.InfraLayerNodeInfo, node.ID, "ip"),
				User:       utils.ExtractValueFromPulumiMapMap(c.InfraLayerNodeInfo, node.ID, "user"),
				PrivateKey: utils.ExtractValueFromPulumiMapMap(c.InfraLayerNodeInfo, node.ID, "key"),
			},
			Create:   pulumi.String("sudo systemctl restart k3s*"),
			Triggers: triggers,
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func depsToCastedArray(deps []pulumi.Resource) (pulumi.Array, error) {
	v := make(pulumi.Array, len(deps))

	for _, d := range deps {
		switch r := d.(type) {
		case *remote.Command:
			v = append(v, r)
		case *remotefile.File:
			v = append(v, r)
		default:
			return nil, fmt.Errorf("unknown rule type: %T with content %+v", r, r)
		}
	}
	return v, nil
}
