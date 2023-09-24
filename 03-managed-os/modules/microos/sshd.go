package microos

import (
	"bytes"
	"fmt"
	"managed-os/utils"

	"github.com/pulumi/pulumi-command/sdk/go/command/remote"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	remotefile "github.com/spigell/pulumi-file/sdk/go/file/remote"
)

func (o *Cluster) ConfigureSSHD(name string, cfg map[string]string) (map[string]pulumi.Resource, error) {
	m := make(map[string]pulumi.Resource)

	for _, node := range o.Nodes {
		b := new(bytes.Buffer)
		for k, v := range cfg {
			fmt.Fprintf(b, "%s %s\n", k, v)
		}

		cleaned, err := remote.NewCommand(o.Ctx, fmt.Sprintf("%s-RemoveSSHDDefaultConfig", node.ID), &remote.CommandArgs{
			Connection: &remote.ConnectionArgs{
				Host:       utils.ExtractValueFromPulumiMapMap(o.InfraLayerNodeInfo, node.ID, "ip"),
				User:       utils.ExtractValueFromPulumiMapMap(o.InfraLayerNodeInfo, node.ID, "user"),
				PrivateKey: utils.ExtractValueFromPulumiMapMap(o.InfraLayerNodeInfo, node.ID, "key"),
			},
			Create: pulumi.String("sudo rm -rfv /etc/ssh/sshd_config"),
		})
		if err != nil {
			return nil, err
		}

		deployed, err := remotefile.NewFile(o.Ctx, fmt.Sprintf("%s-ConfigureSSHD-%s", node.ID, name), &remotefile.FileArgs{
			Connection: &remotefile.ConnectionArgs{
				Host:       utils.ExtractValueFromPulumiMapMap(o.InfraLayerNodeInfo, node.ID, "ip"),
				User:       utils.ExtractValueFromPulumiMapMap(o.InfraLayerNodeInfo, node.ID, "user"),
				PrivateKey: utils.ExtractValueFromPulumiMapMap(o.InfraLayerNodeInfo, node.ID, "key"),
			},
			UseSudo:  pulumi.Bool(true),
			Path:     pulumi.Sprintf("/etc/ssh/sshd_config.d/%s.conf", name),
			Content:  pulumi.String(b.String()),
			SftpPath: pulumi.String("/usr/libexec/ssh/sftp-server"),
		}, pulumi.RetainOnDelete(true), pulumi.DependsOn([]pulumi.Resource{cleaned}))
		if err != nil {
			return nil, err
		}

		restarted, err := remote.NewCommand(o.Ctx, fmt.Sprintf("%s-RestartSSHD", node.ID), &remote.CommandArgs{
			Connection: &remote.ConnectionArgs{
				Host:       utils.ExtractValueFromPulumiMapMap(o.InfraLayerNodeInfo, node.ID, "ip"),
				User:       utils.ExtractValueFromPulumiMapMap(o.InfraLayerNodeInfo, node.ID, "user"),
				PrivateKey: utils.ExtractValueFromPulumiMapMap(o.InfraLayerNodeInfo, node.ID, "key"),
			},
			Create: pulumi.String("sudo systemctl restart sshd"),
			Triggers: pulumi.Array{
				deployed.Md5sum,
				deployed.Permissions,
				deployed.Connection,
				deployed.Path,
			},
		}, pulumi.DependsOn([]pulumi.Resource{deployed}),
			pulumi.DeleteBeforeReplace(true),
		)
		if err != nil {
			return nil, err
		}

		m[node.ID] = restarted
	}

	return m, nil
}
