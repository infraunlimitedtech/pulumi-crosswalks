package k3s

import (
	"fmt"
	"managed-os/config"
	"managed-os/modules/wireguard"
	"managed-os/utils"

	"github.com/pulumi/pulumi-command/sdk/go/command/remote"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	remotefile "github.com/spigell/pulumi-file/sdk/go/file/remote"
	"gopkg.in/yaml.v3"
)

var svcName = "k3s"

func (c *Cluster) configure(wgPeers pulumi.AnyOutput, deps []map[string]pulumi.Resource) (map[string]pulumi.Resource, error) {
	result := make(map[string]pulumi.Resource)

	leaderDeployed, err := remotefile.NewFile(c.Ctx, fmt.Sprintf("%s-K3sCluster", c.Leader.ID), &remotefile.FileArgs{
		Connection: &remotefile.ConnectionArgs{
			Host:       utils.ExtractValueFromPulumiMapMap(c.InfraLayerNodeInfo, c.Leader.ID, "ip"),
			User:       utils.ExtractValueFromPulumiMapMap(c.InfraLayerNodeInfo, c.Leader.ID, "user"),
			PrivateKey: utils.ExtractValueFromPulumiMapMap(c.InfraLayerNodeInfo, c.Leader.ID, "key"),
		},
		UseSudo:  pulumi.Bool(true),
		Path:     pulumi.String(cfgPath),
		Content:  c.renderK3sCfg(c.Leader, wgPeers),
		SftpPath: pulumi.String("/usr/libexec/ssh/sftp-server"),
	}, pulumi.DependsOn(utils.ConvertMapSliceToSliceByKey(deps, c.Leader.ID)), pulumi.RetainOnDelete(true))
	if err != nil {
		return nil, err
	}

	leaderRestared, err := remote.NewCommand(c.Ctx, fmt.Sprintf("%s-RestartK3s", c.Leader.ID), &remote.CommandArgs{
		Connection: &remote.ConnectionArgs{
			Host:       utils.ExtractValueFromPulumiMapMap(c.InfraLayerNodeInfo, c.Leader.ID, "ip"),
			User:       utils.ExtractValueFromPulumiMapMap(c.InfraLayerNodeInfo, c.Leader.ID, "user"),
			PrivateKey: utils.ExtractValueFromPulumiMapMap(c.InfraLayerNodeInfo, c.Leader.ID, "key"),
		},
		Create: pulumi.Sprintf("sudo systemctl enable --now %s", svcName),
		Delete: pulumi.Sprintf("sudo systemctl disable --now %s", svcName),
		Triggers: pulumi.Array{
			leaderDeployed.Md5sum,
			leaderDeployed.Permissions,
			leaderDeployed.Connection,
		},
	}, pulumi.DependsOn([]pulumi.Resource{leaderDeployed}),
		pulumi.DeleteBeforeReplace(true),
	)
	if err != nil {
		return nil, err
	}

	result[c.Leader.ID] = leaderRestared

	for _, node := range c.Followers {
		if node.Role == "agent" {
			svcName = "k3s-agent"
		}

		nodeDeployed, err := remotefile.NewFile(c.Ctx, fmt.Sprintf("%s-ConfigureK3s", node.ID), &remotefile.FileArgs{
			Connection: &remotefile.ConnectionArgs{
				Host:       utils.ExtractValueFromPulumiMapMap(c.InfraLayerNodeInfo, node.ID, "ip"),
				User:       utils.ExtractValueFromPulumiMapMap(c.InfraLayerNodeInfo, node.ID, "user"),
				PrivateKey: utils.ExtractValueFromPulumiMapMap(c.InfraLayerNodeInfo, node.ID, "key"),
			},
			UseSudo:  pulumi.Bool(true),
			Path:     pulumi.String(cfgPath),
			Content:  c.renderK3sCfg(node, wgPeers),
			SftpPath: pulumi.String("/usr/libexec/ssh/sftp-server"),
		}, pulumi.DependsOn(append(utils.ConvertMapSliceToSliceByKey(deps, node.ID), leaderDeployed)),
			pulumi.RetainOnDelete(true))
		if err != nil {
			err = fmt.Errorf("error creating a follower config for node `%s`: %w", node.ID, err)
			return nil, err
		}

		nodeRestarted, err := remote.NewCommand(c.Ctx, fmt.Sprintf("%s-RestartK3s", node.ID), &remote.CommandArgs{
			Connection: &remote.ConnectionArgs{
				Host:       utils.ExtractValueFromPulumiMapMap(c.InfraLayerNodeInfo, node.ID, "ip"),
				User:       utils.ExtractValueFromPulumiMapMap(c.InfraLayerNodeInfo, node.ID, "user"),
				PrivateKey: utils.ExtractValueFromPulumiMapMap(c.InfraLayerNodeInfo, node.ID, "key"),
			},
			Create: pulumi.Sprintf("sudo systemctl enable --now %s", svcName),
			Delete: pulumi.Sprintf("sudo systemctl disable --now %s", svcName),
			Triggers: pulumi.Array{
				nodeDeployed.Md5sum,
				nodeDeployed.Permissions,
				nodeDeployed.Connection,
				nodeDeployed.Path,
			},
		}, pulumi.DependsOn([]pulumi.Resource{nodeDeployed}),
			pulumi.DeleteBeforeReplace(true),
		)
		if err != nil {
			return nil, err
		}

		result[node.ID] = nodeRestarted
	}
	return result, nil
}

func (c *Cluster) renderK3sCfg(node *config.Node, wgClusterInfo pulumi.AnyOutput) pulumi.StringOutput {
	return wgClusterInfo.ApplyT(func(v interface{}) string {
		parsed := v.([]wireguard.Peer)

		peers := wireguard.ToPeers(parsed)

		node.Wireguard.IP = peers.Get(node.ID).PrivateAddr
		leaderIP := peers.Get(c.Leader.ID).PrivateAddr

		s := c.addCoreParams(node, leaderIP)

		k3sRendered, _ := yaml.Marshal(&s)
		return string(k3sRendered)
	}).(pulumi.StringOutput)
}

func (c *Cluster) addCoreParams(cfg *config.Node, leaderIP string) config.K3sConfig {
	k := cfg

	k.K3s.Config.FlannelIface = c.Iface
	k.K3s.Config.NodeIP = cfg.Wireguard.IP

	if k.Role == "server" {
		k.K3s.Config.WriteKubeconfigMode = "0644"
		k.K3s.Config.BindAddress = cfg.Wireguard.IP
	}

	if k.Leader {
		k.K3s.Config.ClusterInit = true
	} else {
		k.K3s.Config.Server = fmt.Sprintf("https://%s:6443", leaderIP)
	}

	return k.K3s.Config
}
