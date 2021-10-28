package services

import (
	globalconfigurationsrv1 "cluster-resources/crds/generated/nginxinc/kubernetes-ingress/k8s/v1alpha1"

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
	Name           string
	UDPDNSListener *globalconfigurationsrv1.GlobalConfigurationSpecListenersArgs
	Helm           *HelmParams
}

type ConsulInfra struct {
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

	i := &Infra{
		Namespace: namespace,
		ctx:       ctx,
		DNS: &DNSInfra{
			InternalDomain: "intra.infraunlimited.tech",
		},
		LB: &LoadBalancersInfra{
			NginxIngress: &NginxIngressInfra{
				Name: "nginx-ingress",
				UDPDNSListener: &globalconfigurationsrv1.GlobalConfigurationSpecListenersArgs{
					Name:     pulumi.String("dns-udp"),
					Port:     pulumi.Int(53),
					Protocol: pulumi.String("UDP"),
				},
				Helm: pulumiCfg.LB.NginxIngress.Helm,
			},
		},
		Consul: pulumiCfg.Consul,
	}
	return i, nil
}
