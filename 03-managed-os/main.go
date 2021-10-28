package main

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/spigell/pulumi-k3os/sdk/go/k3os"
)

const (
	wgListenPort = 51820
	serverStr    = "server"
	agentStr     = "agent"
)

var (
	errNoLeader    = errors.New("There is no leader. Please set it in config")
	errAgentLeader = errors.New("Agent can't be a leader")
	errManyLeaders = errors.New("There is more than one leader")
)

type cluster struct {
	serverURL string
	leader    Node
	followers []Node
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		pulumiCfg := parseConfig(ctx)

		infraStack, err := pulumi.NewStackReference(ctx, pulumiCfg.InfraStack, nil)
		if err != nil {
			return err
		}

		nodesInfo := infraStack.GetOutput(pulumi.String("infra:nodes:info"))

		cluster, err := newCluster(pulumiCfg.Nodes, nodesInfo)
		if err != nil {
			err = errors.Wrap(err, "Error init cluster")
			ctx.Log.Error(err.Error(), nil)
			return err
		}

		config, err := buildNodeConfig(cluster, pulumiCfg, cluster.leader)
		if err != nil {
			err = errors.Wrap(err, "Error creating a leader config")
			ctx.Log.Error(err.Error(), nil)
			return err
		}

		leader, err := k3os.NewNode(ctx, cluster.leader.ID, &k3os.NodeArgs{
			Connection: &k3os.ConnectionArgs{
				Addr: pulumi.Sprintf("%s:22", extractConnectionArg(nodesInfo, cluster.leader.ID, "ip")),
				User: extractConnectionArg(nodesInfo, cluster.leader.ID, "user"),
				Key:  extractConnectionArg(nodesInfo, cluster.leader.ID, "key"),
			},
			NodeConfiguration: &k3os.NodeConfigurationArgs{
				Hostname: pulumi.String(cluster.leader.ID),
				RunCmd:   pulumi.ToStringArray(config.Runcmd),
				BootCmd:  pulumi.ToStringArray(config.Bootcmd),
				WriteFiles: &k3os.CloudInitFilesArray{
					&k3os.CloudInitFilesArgs{
						Content:  pulumi.String(config.WriteFiles[0].Content),
						Path:     pulumi.String(config.WriteFiles[0].Path),
						Encoding: pulumi.String(config.WriteFiles[0].Encoding),
					},
				},
				K3OS: &k3os.K3OSArgs{
					Token:   pulumi.String(config.K3os.Token),
					K3sArgs: pulumi.ToStringArray(config.K3os.K3sArgs),
					Labels:  pulumi.ToStringMap(config.K3os.Labels),
				},
			},
		})
		if err != nil {
			return err
		}

		ctx.Export("infra:vpn:address", pulumi.Unsecret(
			pulumi.Sprintf("%s:%d", extractConnectionArg(nodesInfo, cluster.leader.ID, "ip"), wgListenPort)),
		)

		for _, node := range cluster.followers {
			config, err := buildNodeConfig(cluster, pulumiCfg, node)
			if err != nil {
				err = errors.Wrap(err, "Error creating a follower config")
				ctx.Log.Error(err.Error(), nil)
				return err
			}

			_, err = k3os.NewNode(ctx, node.ID, &k3os.NodeArgs{
				Connection: &k3os.ConnectionArgs{
					Addr: pulumi.Sprintf("%s:22", extractConnectionArg(nodesInfo, node.ID, "ip")),
					User: extractConnectionArg(nodesInfo, node.ID, "user"),
					Key:  extractConnectionArg(nodesInfo, node.ID, "key"),
				},
				NodeConfiguration: &k3os.NodeConfigurationArgs{
					Hostname: pulumi.String(node.ID),
					BootCmd:  pulumi.ToStringArray(config.Bootcmd),
					RunCmd:   pulumi.ToStringArray(config.Runcmd),
					WriteFiles: &k3os.CloudInitFilesArray{
						&k3os.CloudInitFilesArgs{
							Content:  pulumi.String(config.WriteFiles[0].Content),
							Path:     pulumi.String(config.WriteFiles[0].Path),
							Encoding: pulumi.String(config.WriteFiles[0].Encoding),
						},
					},
					K3OS: &k3os.K3OSArgs{
						Token:     pulumi.String(config.K3os.Token),
						ServerUrl: pulumi.String(cluster.serverURL),
						K3sArgs:   pulumi.ToStringArray(config.K3os.K3sArgs),
						Labels:    pulumi.ToStringMap(config.K3os.Labels),
					},
				},
			}, pulumi.DependsOn([]pulumi.Resource{leader}))
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func newCluster(n Nodes, infraOutputs pulumi.AnyOutput) (*cluster, error) {
	var leader Node
	var serverURL string
	followers := make([]Node, 0)

	for _, agent := range n.Agents {
		agent.kind = agentStr
		followers = append(followers, agent)
	}
	for _, server := range n.Servers {
		server.kind = serverStr
		if server.Leader {
			leader = server
			serverURL = fmt.Sprintf("https://%s:6443", server.Wireguard.PrivateAddr)
		} else {
			followers = append(followers, server)
		}
	}

	if err := validate(append(followers, leader), infraOutputs); err != nil {
		return nil, err
	}

	return &cluster{
		leader:    leader,
		serverURL: serverURL,
		followers: followers,
	}, nil

}

func validate(nodes []Node, infraOutputs pulumi.AnyOutput) error {

	_ = infraOutputs.ApplyT(func(v interface{}) error {
		for _, node := range nodes {
			outputs := v.(map[string]interface{})
			m, ok := outputs[node.ID].(map[string]interface{})
			if !ok {
				panic(fmt.Sprintf("No node `%s` in Infra Layer!", node.ID))
			}
			if node.PublicIP != m["ip"] {
				panic(fmt.Sprintf("Please set correct PublicIP `%s` for node `%s` into Pulumi config! Current:`%s`",
					m["ip"], node.ID, node.PublicIP),
				)
			}
		}
		return nil
	})
	leaderFounded := false
	for _, node := range nodes {
		if node.Leader {
			if node.kind == agentStr {
				return errAgentLeader
			}

			if !leaderFounded {
				leaderFounded = true
			} else {
				return errManyLeaders
			}
		}
	}
	if !leaderFounded {
		return errNoLeader
	}
	return nil
}

func extractConnectionArg(p pulumi.AnyOutput, nodeID, key string) pulumi.StringOutput {
	return p.ApplyT(func(v interface{}) string {
		nodes := v.(map[string]interface{})
		node, ok := nodes[nodeID].(map[string]interface{})
		if !ok {
			panic(fmt.Sprintf("Can't find values for node `%s`. It managed by Pulumi Infra Layer?", nodeID))
		}
		return node[key].(string)
	}).(pulumi.StringOutput)
}
