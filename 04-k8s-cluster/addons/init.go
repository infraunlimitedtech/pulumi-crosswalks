package addons

import (
	"k8s-cluster/config"
	"k8s-cluster/packages/kilo"
	"k8s-cluster/spec"

	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	pulumiConfig "github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

type Addons struct {
	ctx          *pulumi.Context
	Namespace    *corev1.Namespace
	Kilo         *kilo.Kilo
	MetalLb      *MetalLb
	NginxIngress *NginxIngress
	Monitoring   *Monitoring
}

type Monitoring struct {
	NodeExporter    *NodeExporter
	VictoriaMetrics *VictoriaMetrics
}

type NodeExporter struct {
	Helm *config.HelmParams
}

type VictoriaMetrics struct {
	Helm   *config.HelmParams
	Server *VictoriaMetricsServer
}

type VictoriaMetricsServer struct {
	ClusterIP string
}

type MetalLb struct {
	Helm               *config.HelmParams
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
	Helm    *config.HelmParams
}

type NginxKubeAPI struct {
	ClusterIP string
}

func Init(ctx *pulumi.Context, s *spec.ClusterSpec) (*Addons, error) {
	namespace := "infra-system"

	ns, err := corev1.NewNamespace(ctx, namespace, &corev1.NamespaceArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name: pulumi.String(namespace),
		},
	})
	if err != nil {
		return nil, err
	}

	// Init vars from stack's config
	var pulumiAddonsCfg Addons
	cfg := pulumiConfig.New(ctx, "")
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
		Namespace: ns,
		ctx:       ctx,
		NginxIngress: &NginxIngress{
			Name:    "nginx-ingress-addon",
			Domain:  s.InternalDomainZone,
			Helm:    pulumiAddonsCfg.NginxIngress.Helm,
			KubeAPI: pulumiAddonsCfg.NginxIngress.KubeAPI,
		},
		Monitoring: pulumiAddonsCfg.Monitoring,
		MetalLb:    pulumiAddonsCfg.MetalLb,
		Kilo:       pulumiAddonsCfg.Kilo,
	}
	return a, nil
}
