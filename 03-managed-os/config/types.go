package config

import "pulumi-crosswalks/utils/hetzner"

type PulumiConfig struct {
	InfraStack string `json:"infra_stack"`
	IDStack    string `json:"id_stack"`
	Defaults   *Defaults
	Nodes      *Nodes
}

type Defaults struct {
	Global  *Node
	Servers *Node
	Agents  *Node
}

type Nodes struct {
	Servers []Node
	Agents  []Node
}

type Node struct {
	ID        string
	Leader    bool
	Wireguard Wireguard
	Firewall  Firewall
	K3s       K3s
	Role      string
}

type K3s struct {
	Version            string
	CleanDataOnUpgrade bool
	Config             K3sConfig
}

type K3sConfig struct {
	Token                     string
	Server                    string   `yaml:",omitempty"`
	FlannelIface              string   `json:"-" yaml:"flannel-iface,omitempty"`
	ClusterCidr               string   `json:"cluster-cidr" yaml:"cluster-cidr,omitempty"`
	ServiceCidr               string   `json:"service-cidr" yaml:"service-cidr,omitempty"`
	ClusterDomain             string   `json:"cluster-domain" yaml:"cluster-domain,omitempty"`
	ClusterDNS                string   `json:"cluster-dns" yaml:"cluster-dns,omitempty"`
	WriteKubeconfigMode       string   `json:"-" yaml:"write-kubeconfig-mode,omitempty"`
	NodeIP                    string   `json:"-" yaml:"node-ip,omitempty"`
	BindAddress               string   `json:"-" yaml:"bind-address,omitempty"`
	ClusterInit               bool     `json:"-" yaml:"cluster-init,omitempty"`
	NodeLabels                []string `json:"node-label" yaml:"node-label,omitempty"`
	FlannelBackend            string   `json:"flannel-backend" yaml:"flannel-backend,omitempty"`
	DisableNetworkPolicy      bool     `json:"disable-network-policy" yaml:"disable-network-policy,omitempty"`
	NodeTaints                []string `json:"node-taint" yaml:"node-taint,omitempty"`
	KubeleteArgs              []string `json:"kubelet-arg" yaml:"kubelet-arg,omitempty"`
	KubeControllerManagerArgs []string `json:"kube-controller-manager-arg" yaml:"kube-controller-manager-arg,omitempty"`
	KubeAPIServerArgs         []string `json:"kube-apiserver-arg" yaml:"kube-apiserver-arg,omitempty"`
	DisableCloudController    bool     `json:"disable-cloud-controller" yaml:"disable-cloud-controller,omitempty"`
	Disable                   []string
}

type Wireguard struct {
	IP              string
	Firewall        *WGFirewall
	MgmtNode        bool             `json:"-"`
	CIDR            string           `json:"cidr"`
	AdditionalPeers []AdditionalPeer `json:"additional_peers" yaml:"additional_peers"`
}

type WGFirewall struct {
	Firewalld *WGFirewalld
}

type WGFirewalld struct {
	Allowed bool
}

type AdditionalPeer struct {
	AllowedIps []string `json:"allowed_ips" yaml:"allowed_ips"`
	Endpoint   string
	PublicKey  string
}

type Firewall struct {
	Hetzner   []hetzner.Firewall
	Firewalld *FWFirewalld
}

type FWFirewalld struct {
	Enabled      bool
	InternalZone *InternalZone
	PublicZone   *PublicZone
}

type InternalZone struct {
	RestrictToSources []*RestrictToSource
}

type RestrictToSource struct {
	CIDR string
	Name string
	Main bool
}

type PublicZone struct {
	RemoveSSHService bool
}
