package k3s

import (
	"managed-os/config"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	cfgPath = "/etc/rancher/k3s/config.yaml"
)

type Cluster struct {
	Iface              string
	ServerURL          string
	InfraLayerNodeInfo pulumi.AnyOutput
	Followers          []*config.Node
	Leader             *config.Node
	Ctx                *pulumi.Context
}

type CreatedCluster struct {
	Kubeconfig *pulumi.StringOutput
}

func GetRequiredPkgs() []string {
	return []string{"k3s-selinux"}
}

func GetRequirdSSHDConfig() map[string]string {
	return map[string]string{"AcceptEnv": "INSTALL_K3S_*"}
}

func (c *Cluster) Manage(WgPeers pulumi.AnyOutput, deps []map[string]pulumi.Resource) (*CreatedCluster, error) {
	installed, err := c.install(deps)
	if err != nil {
		return nil, err
	}

	configured, err := c.configure(WgPeers, []map[string]pulumi.Resource{installed})
	if err != nil {
		return nil, err
	}

	kubeConfig, err := c.grabKubeConfig([]map[string]pulumi.Resource{configured})
	if err != nil {
		return nil, err
	}

	return &CreatedCluster{
		Kubeconfig: kubeConfig,
	}, nil
}
