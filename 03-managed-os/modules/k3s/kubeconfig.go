package k3s

import (
	"fmt"
	"managed-os/utils"

	"github.com/pulumi/pulumi-command/sdk/go/command/remote"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"k8s.io/client-go/tools/clientcmd"
)

func (c *Cluster) grabKubeConfig(deps []map[string]pulumi.Resource) (*pulumi.StringOutput, error) {
	grabbed, err := remote.NewCommand(c.Ctx, fmt.Sprintf("%s-GrabKubeConfig", c.Leader.ID), &remote.CommandArgs{
		Connection: &remote.ConnectionArgs{
			Host:       utils.ExtractValueFromPulumiMapMap(c.InfraLayerNodeInfo, c.Leader.ID, "ip"),
			User:       utils.ExtractValueFromPulumiMapMap(c.InfraLayerNodeInfo, c.Leader.ID, "user"),
			PrivateKey: utils.ExtractValueFromPulumiMapMap(c.InfraLayerNodeInfo, c.Leader.ID, "key"),
		},
		Create: pulumi.String("sudo cat /etc/rancher/k3s/k3s.yaml"),
	}, pulumi.DependsOn(utils.ConvertMapSliceToSliceByKey(deps, c.Leader.ID)))
	if err != nil {
		err = fmt.Errorf("error grab kubeconfig: %w", err)
		return nil, err
	}

	k := grabbed.Stdout.ApplyT(func (v interface{}) string {
		kubeconfig, err := clientcmd.Load([]byte(v.(string)))
		if err != nil {
			panic("Failed to parse kubeconfig")
		}

		ctxName := fmt.Sprintf("%s-direct", c.Ctx.Stack())

		kubeconfig.Contexts[ctxName] = kubeconfig.Contexts["default"]
		delete(kubeconfig.Contexts, "default")
		kubeconfig.CurrentContext = ctxName

		w, _ := clientcmd.Write(*kubeconfig)

		return string(w)

	}).(pulumi.StringOutput)

	return &k, nil
//	return &grabbed.Stdout, nil
}
