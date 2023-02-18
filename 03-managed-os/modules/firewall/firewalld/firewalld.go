package firewalld

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"fmt"
	"pulumi-crosswalks/utils"
	"github.com/pulumi/pulumi-command/sdk/go/command/remote"
	"github.com/pulumi/pulumi-command/sdk/go/command/local"
)

type Firewalld struct {
	ctx                *pulumi.Context
	NodeInfo *Node
	InternalIface string
	DependsOn []pulumi.Resource
}

type Node struct {
	ID string
	Host pulumi.StringOutput
	User pulumi.StringOutput
	PrivateKey pulumi.StringOutput
}

func New(ctx *pulumi.Context,
	nodeInfo pulumi.AnyOutput,
	nodeID string,
	deps []map[string]pulumi.Resource,
) (*Firewalld, error) {
	fwd := &Firewalld{
		ctx: ctx,
		NodeInfo: &Node{
			ID: nodeID,
			Host:       utils.ExtractValueFromPulumiMapMap(nodeInfo, nodeID, "ip"),
			User:       utils.ExtractValueFromPulumiMapMap(nodeInfo, nodeID, "user"),
			PrivateKey: utils.ExtractValueFromPulumiMapMap(nodeInfo, nodeID, "key"),
		},
		DependsOn: utils.ConvertMapSliceToSliceByKey(deps, nodeID),
	}

	enabled, err := remote.NewCommand(fwd.ctx, fmt.Sprintf("EnableFirewalld-%s", fwd.NodeInfo.ID), &remote.CommandArgs{
		Connection: &remote.ConnectionArgs{
			Host:       fwd.NodeInfo.Host,
			User:       fwd.NodeInfo.User,
			PrivateKey: fwd.NodeInfo.PrivateKey,
		},

		Create: pulumi.String("sudo systemctl enable --now firewalld"),
		Delete: pulumi.String("sudo systemctl disable --now firewalld"),
	}, pulumi.DependsOn(fwd.DependsOn))

	if err != nil {
		return nil, err
	}

	fwd.DependsOn = append(fwd.DependsOn, enabled)

	return fwd, err
}

func (f *Firewalld) Reload() error {
	triggers, err := depsToCastedArray(f.DependsOn) 
	if err != nil {
		return err
	}

	_, err = remote.NewCommand(f.ctx, fmt.Sprintf("ReloadFirewalld-%s", f.NodeInfo.ID), &remote.CommandArgs{
		Connection: &remote.ConnectionArgs{
			Host:       f.NodeInfo.Host,
			User:       f.NodeInfo.User,
			PrivateKey: f.NodeInfo.PrivateKey,
		},
		Create: pulumi.String("sudo firewall-cmd --reload"),
		Triggers: triggers,
	})

	if err != nil {
		return err
	}

	return nil
}

func GetRequiredPkgs() []string {
	return []string{"firewalld"}
}

func GetWhitelistedIfaces() []string {
//	return []string{"cni0", "flannel1.0", f.InternalIface}
	return []string{"cni0", "flannel.1", "kubewg0"}
}

func depsToCastedArray(deps []pulumi.Resource) (pulumi.Array, error) {
	v := make(pulumi.Array, len(deps))

	for _, d := range deps {
		switch r := d.(type) {
		case *remote.Command:
			v = append(v, r)
		case *local.Command:
			v = append(v, r)
		default:
			return nil, fmt.Errorf("unknown rule type: %T with content %+v", r, r)
		}
	}
	return v, nil
}

