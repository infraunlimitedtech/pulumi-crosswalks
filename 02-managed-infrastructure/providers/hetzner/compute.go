package hetzner

import (
	"github.com/pulumi/pulumi-hcloud/sdk/go/hcloud"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type ComputeConfig struct{}

type ComputedInfra struct {
	nodes map[string]map[string]interface{}
}

var c = `#cloud-config
ssh_authorized_keys:
- ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBJgyB/EEX/fsSICjyHha9Pnt1IM7brDsFelakF1hTNdKjA+qdvojKWSNraGN81ewf4nxexV6E6e5fEeyr2IIcAQ=`

func ManageCompute(ctx *pulumi.Context) (*ComputedInfra, error) {
	nodes := make(map[string]map[string]interface{})
	err := manage(ctx)

	if err != nil {
		return nil, err
	}

	return &ComputedInfra{
		nodes: nodes,
	}, nil
}

func manage(ctx *pulumi.Context) error {
	_, err := hcloud.NewServer(ctx, "myServer", &hcloud.ServerArgs{
		ServerType: pulumi.String("cpx11"),
		Location:   pulumi.String("hel1"),
		// testing image with microos
		Image:    pulumi.String("65179453"),
		UserData: pulumi.String(c),
	})
	if err != nil {
		return err
	}
	return nil
}

func (v *ComputedInfra) GetNodes() map[string]map[string]interface{} {
	return v.nodes
}
