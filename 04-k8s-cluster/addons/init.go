package addons

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
	"k8s-cluster/spec"
)

type Addons struct {
	ctx          *pulumi.Context
	Namespace    string
	Kilo         *Kilo
	MetalLb      *MetalLb
	NginxIngress *NginxIngress
}

type MetalLb struct {
	Helm               *HelmParams
	Pools              *MetalLbPools
	DefaultNetworkPool string
	ExternalIP         string
	KubeapiIP          string
}

type MetalLbPools struct {
	Default MetalLbPool
	Kilo    MetalLbPool
	Kubeapi MetalLbPool
}

type MetalLbPool struct {
	Network string
}

type NginxIngress struct {
	Name   string
	Domain string
	Helm   *HelmParams
}

type Kilo struct {
	PrivateKey string
	Crds       *CRDS
	Version    string
	Peers      []KiloPeer
}

type KiloPeer struct {
	Name       string
	PublicKey  string
	AllowedIPs []string
}

type HelmParams struct {
	Version string
}

type CRDS struct {
	Path string
}

func Init(ctx *pulumi.Context, s *spec.ClusterSpec) (*Addons, error) {
	namespace := "kube-system"

	// Init vars from stack's config
	var pulumiAddonsCfg Addons
	cfg := config.New(ctx, "")
	cfg.RequireSecretObject("addons", &pulumiAddonsCfg)

	a := &Addons{
		Namespace: namespace,
		ctx:       ctx,
		NginxIngress: &NginxIngress{
			Name:   "nginx-ingress-addon",
			Domain: s.InternalDomainZone,
			Helm:   pulumiAddonsCfg.NginxIngress.Helm,
		},
		MetalLb: pulumiAddonsCfg.MetalLb,
		Kilo:    pulumiAddonsCfg.Kilo,
	}
	return a, nil
}
