package main

import (
	"fmt"
	"managed-os/config"
	"managed-os/modules/firewall/firewalld"
	"managed-os/modules/k3s"
	"managed-os/modules/wireguard"
	"managed-os/utils"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		pulumiCfg := config.ParseConfig(ctx)

		infraStack, err := pulumi.NewStackReference(ctx, pulumiCfg.InfraStack, nil)
		if err != nil {
			return err
		}

		idStack, err := pulumi.NewStackReference(ctx, pulumiCfg.IDStack, nil)
		if err != nil {
			return err
		}

		infraLayerNodeInfo := infraStack.GetOutput(pulumi.String("infra:nodes:info"))

		wgInfo := idStack.GetOutput(pulumi.String("identity:organization")).ApplyT(func(v interface{}) pulumi.AnyOutput {
			selfStack, err := pulumi.NewStackReference(ctx, fmt.Sprintf("%s/%s/%s", v.(string), ctx.Project(), ctx.Stack()), nil)
			if err != nil {
				panic("error get wireguard keys from stack")
			}
			return selfStack.GetOutput(pulumi.String("os:wireguard:info"))
		}).(pulumi.AnyOutput)

		cluster, err := defineCluster(&cluster{
			ctx:                ctx,
			PulumiConfig:       pulumiCfg,
			InfraLayerNodeInfo: infraLayerNodeInfo,
			RequiredPkgs: append(
				append(
					k3s.GetRequiredPkgs(),
					wireguard.GetRequiredPkgs()...,
				),
				firewalld.GetRequiredPkgs()...,
			),
			InternalIface: "kubewg0",
			WgInfo:        wgInfo,
		})
		if err != nil {
			err = fmt.Errorf("error init cluster: %w", err)
			ctx.Log.Error(err.Error(), nil)
			return err
		}

		pkgs, err := cluster.Os.InstallRequiredPkgs()
		if err != nil {
			err = fmt.Errorf("error while installing required packages: %w", err)
			ctx.Log.Error(err.Error(), nil)
			return err
		}

		cfg, err := cluster.Os.ConfigureSSHD("k3s", k3s.GetRequirdSSHDConfig())
		if err != nil {
			err = fmt.Errorf("error configure sshd service for k3s cluster: %w", err)
			ctx.Log.Error(err.Error(), nil)
			return err
		}

		reboot, _ := cluster.Os.Reboot([]map[string]pulumi.Resource{pkgs, cfg})

		wgCluster, err := cluster.Wireguard.Manage([]map[string]pulumi.Resource{reboot})
		if err != nil {
			err = fmt.Errorf("error creating a wireguard cluster: %w", err)
			ctx.Log.Error(err.Error(), nil)
			return err
		}

		k3sCluster, err := cluster.K3s.Manage(wgCluster.Peers, []map[string]pulumi.Resource{wgCluster.Resources})
		if err != nil {
			ctx.Log.Error(err.Error(), nil)
			return err
		}

		err = cluster.Firewalls.Manage([]map[string]pulumi.Resource{reboot})
		if err != nil {
			ctx.Log.Error(err.Error(), nil)
			return err
		}

		ctx.Export("os:wireguard:info", wgCluster.ConvertPeersToMapMap())
		ctx.Export("os:wireguard:config", pulumi.ToSecret(wgCluster.MasterConfig))

		ctx.Export("os:vpn:address", pulumi.Unsecret(
			pulumi.Sprintf("%s:%d", utils.ExtractValueFromPulumiMapMap(infraLayerNodeInfo, cluster.K3s.Leader.ID, "ip"), wgCluster.ListenPort)),
		)

		ctx.Export("os:k3s:kubeconfig", pulumi.ToSecret(k3sCluster.Kubeconfig))

		return nil
	})
}
