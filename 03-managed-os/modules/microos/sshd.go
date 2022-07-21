package microos

import (
	"bytes"
	"fmt"
	"managed-os/utils"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi-command/sdk/go/command/remote"
	"github.com/spigell/pulumi-file/sdk/go/file"
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

		deployed, err := file.NewRemote(o.Ctx, fmt.Sprintf("%s-ConfigureSSHD-%s", node.ID, name), &file.RemoteArgs{
			Connection: &file.ConnectionArgs{
				Address:    pulumi.Sprintf("%s:22", utils.ExtractValueFromPulumiMapMap(o.InfraLayerNodeInfo, node.ID, "ip")),
				User:       utils.ExtractValueFromPulumiMapMap(o.InfraLayerNodeInfo, node.ID, "user"),
				PrivateKey: utils.ExtractValueFromPulumiMapMap(o.InfraLayerNodeInfo, node.ID, "key"),
			},
			Hooks: &file.HooksArgs{
				CommandAfterCreate: pulumi.String("sudo systemctl reload sshd"),
				CommandAfterUpdate: pulumi.String("sudo systemctl reload sshd"),
			},
			UseSudo: pulumi.Bool(true),
			Path:    pulumi.Sprintf("/etc/ssh/sshd_config.d/%s.conf", name),
			Content: pulumi.String(b.String()),
		}, pulumi.RetainOnDelete(true), pulumi.DependsOn([]pulumi.Resource{cleaned}))
		if err != nil {
			return nil, err
		}
		m[node.ID] = deployed
	}

	return m, nil
}
