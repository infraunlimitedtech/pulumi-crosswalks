package microos

import (
	"managed-os/config"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Cluster struct {
	Ctx                *pulumi.Context
	Nodes              []*config.Node
	InfraLayerNodeInfo pulumi.AnyOutput
	RequiredPkgs       []string
}
