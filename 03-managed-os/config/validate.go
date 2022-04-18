package config

import (
	"errors"
	"fmt"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var (
	errNoLeader    = errors.New("there is no a leader. Please set it in config")
	errAgentLeader = errors.New("agent can't be a leader")
	errManyLeaders = errors.New("there is more than one leader")
)

func Validate(nodes []*Node, infraOutputs pulumi.AnyOutput) error {
	_ = infraOutputs.ApplyT(func(v interface{}) error {
		for _, node := range nodes {
			outputs := v.(map[string]interface{})
			_, ok := outputs[node.ID].(map[string]interface{})
			if !ok {
				panic(fmt.Sprintf("No node `%s` in Infra Layer!", node.ID))
			}
		}
		return nil
	})
	leaderFounded := false
	for _, node := range nodes {
		if node.Leader {
			if node.Role == "agent" {
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
