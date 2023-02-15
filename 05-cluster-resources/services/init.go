package services

import (
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

type Infra struct {
	ctx       *pulumi.Context
	Namespace string
	DNS       *DNSInfra
	LB        *LoadBalancersInfra
	Consul    *ConsulInfra
}

type DNSInfra struct {
	InternalDomain string
}

type LoadBalancersInfra struct {
	NginxIngress *NginxIngressInfra
}

type NginxIngressInfra struct {
	Enabled   bool
	Name      string
	ClusterIP string
	Helm      *HelmParams
}

type ConsulInfra struct {
	Enabled      bool
	GossipSecret string
}

type HelmParams struct {
	Version string
}

type CRDS struct {
	Path string
}

func Init(ctx *pulumi.Context) (*Infra, error) {
	namespace := "infra-services"

	// Init vars from stack's config
	var pulumiCfg Infra
	cfg := config.New(ctx, "")
	cfg.RequireSecretObject("infra", &pulumiCfg)

	// Create main NS
	_, err := corev1.NewNamespace(ctx, namespace, &corev1.NamespaceArgs{
		ApiVersion: pulumi.String("v1"),
		Kind:       pulumi.String("Namespace"),
		Metadata: &metav1.ObjectMetaArgs{
			Name: pulumi.String(namespace),
		},
	})
	if err != nil {
		return nil, err
	}

	_, err = corev1.NewResourceQuota(ctx, namespace, &corev1.ResourceQuotaArgs{
		ApiVersion: pulumi.String("v1"),
		Kind:       pulumi.String("ResourceQuota"),
		Metadata: &metav1.ObjectMetaArgs{
			Namespace: pulumi.String(namespace),
			Name:      pulumi.String("cpu-mem-limit"),
		},
		Spec: &corev1.ResourceQuotaSpecArgs{
			Hard: pulumi.StringMap{
				"requests.cpu":    pulumi.String("3"),
				"requests.memory": pulumi.String("3Gi"),
				"limits.cpu":      pulumi.String("4"),
				"limits.memory":   pulumi.String("4Gi"),
			},
		},
	})
	if err != nil {
		return nil, err
	}

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
	ingress := pulumiCfg.LB.NginxIngress
	ingress.Name = "nginx-ingress"

	i := &Infra{
		Namespace: namespace,
		ctx:       ctx,
		DNS: &DNSInfra{
			InternalDomain: "intra.infraunlimited.tech",
		},
		LB: &LoadBalancersInfra{
			NginxIngress: ingress,
		},
		Consul: pulumiCfg.Consul,
	}
	return i, nil
}
