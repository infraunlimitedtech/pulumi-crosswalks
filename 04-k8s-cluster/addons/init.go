package addons

import (
	"k8s-cluster/spec"

	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
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
	Kubeapi MetalLbPool
}

type MetalLbPool struct {
	Network string
}

type NginxIngress struct {
	Name    string
	Domain  string
	KubeAPI NginxKubeAPI
	Replica int
	Helm    *HelmParams
}

type NginxKubeAPI struct {
	ClusterIP string
}

type Kilo struct {
	PrivateKey  string
	Crds        *CRDS
	Version     string
	Peers       []KiloPeer
	Firewalls   *Firewalls
	ExternalIPs []string
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

type Firewalls struct {
	Hetzner   *Firewall
	Firewalld *Firewall
}

type Firewall struct {
	Managed bool
}

func Init(ctx *pulumi.Context, s *spec.ClusterSpec) (*Addons, error) {
	namespace := "infra-system"

	_, err := corev1.NewNamespace(ctx, namespace, &corev1.NamespaceArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name: pulumi.String(namespace),
		},
	})
	if err != nil {
		return nil, err
	}

	// Init vars from stack's config
	var pulumiAddonsCfg Addons
	cfg := config.New(ctx, "")
	cfg.RequireSecretObject("addons", &pulumiAddonsCfg)

	_, err = corev1.NewLimitRange(ctx, namespace, &corev1.LimitRangeArgs{
		ApiVersion: pulumi.String("v1"),
		Kind:       pulumi.String("LimitRange"),
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String("mem-range"),
			Namespace: pulumi.String(namespace),
		},
		Spec: &corev1.LimitRangeSpecArgs{
			Limits: corev1.LimitRangeItemArray{
				&corev1.LimitRangeItemArgs{
					Default: pulumi.StringMap{
						"memory": pulumi.String("128Mi"),
						"cpu":    pulumi.String("200m"),
					},
					DefaultRequest: pulumi.StringMap{
						"memory": pulumi.String("64Mi"),
						"cpu":    pulumi.String("100m"),
					},
					Type: pulumi.String("Container"),
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	a := &Addons{
		Namespace: namespace,
		ctx:       ctx,
		NginxIngress: &NginxIngress{
			Name:    "nginx-ingress-addon",
			Domain:  s.InternalDomainZone,
			Helm:    pulumiAddonsCfg.NginxIngress.Helm,
			KubeAPI: pulumiAddonsCfg.NginxIngress.KubeAPI,
		},
		MetalLb: pulumiAddonsCfg.MetalLb,
		Kilo:    pulumiAddonsCfg.Kilo,
	}
	return a, nil
}
