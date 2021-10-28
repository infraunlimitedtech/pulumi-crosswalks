package services

import (
	"fmt"

	nginxv1alpha1 "cluster-resources/crds/generated/nginxinc/kubernetes-ingress/k8s/v1alpha1"

	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	helmv3 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v3"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (infra *Infra) RunConsulStack() error {
	serviceName := "consul"

	gossipSecret, err := corev1.NewSecret(infra.ctx, "gossip-secret", &corev1.SecretArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Namespace: pulumi.String(infra.Namespace),
		},
		Type: pulumi.String("Opaque"),
		StringData: pulumi.StringMap{
			"key": pulumi.String(infra.Consul.GossipSecret),
		},
	})
	if err != nil {
		return err
	}

	_, err = helmv3.NewChart(infra.ctx, "consul-stack", helmv3.ChartArgs{
		Chart:     pulumi.String("consul"),
		Version:   pulumi.String("v0.32.1"),
		Namespace: pulumi.String(infra.Namespace),
		FetchArgs: helmv3.FetchArgs{
			Repo: pulumi.String("https://helm.releases.hashicorp.com"),
		},
		Values: pulumi.Map{
			"global": pulumi.Map{
				"name":       pulumi.String(serviceName),
				"domain":     pulumi.String(infra.DNS.InternalDomain),
				"datacenter": pulumi.String("CHANGE_ME"),
				"recursors": pulumi.Array{
					pulumi.String("8.8.8.8"),
					pulumi.String("1.1.1.1"),
				},

				"gossipEncryption": pulumi.Map{
					"secretName": gossipSecret.Metadata.Name(),
					"secretKey":  pulumi.String("key"),
				},
			},
			"server": pulumi.Map{
				"replicas": pulumi.Int(3),
				"disruptionBudget": pulumi.Map{
					"enabled": pulumi.Bool(false),
				},
				"storage": pulumi.String("1Gi"),
				"resources": pulumi.Map{
					"requests": pulumi.Map{
						"memory": pulumi.String("50Mi"),
						"cpu":    pulumi.String("50m"),
					},
					"limits": pulumi.Map{
						"memory": pulumi.String("50Mi"),
						"cpu":    pulumi.String("50m"),
					},
				},
				"extraConfig": pulumi.String("{    \"log_level\": \"DEBUG\"  }"),
			},
			"client": pulumi.Map{
				"resources": pulumi.Map{
					"requests": pulumi.Map{
						"memory": pulumi.String("50Mi"),
						"cpu":    pulumi.String("50m"),
					},
					"limits": pulumi.Map{
						"memory": pulumi.String("50Mi"),
						"cpu":    pulumi.String("50m"),
					},
				},
			},
			"syncCatalog": pulumi.Map{
				"enabled":               pulumi.Bool(true),
				"addK8SNamespaceSuffix": pulumi.Bool(false),
				"syncClusterIPServices": pulumi.Bool(false),
				"toK8S":                 pulumi.Bool(false),
				"resources": pulumi.Map{
					"requests": pulumi.Map{
						"memory": pulumi.String("20Mi"),
						"cpu":    pulumi.String("20m"),
					},
					"limits": pulumi.Map{
						"memory": pulumi.String("30Mi"),
						"cpu":    pulumi.String("30m"),
					},
				},
			},
		},
	})
	if err != nil {
		return err
	}

	nginxIngressDNSUpstreamName := "dns"

	_, err = nginxv1alpha1.NewTransportServer(infra.ctx, "nginx-transport-server-dns", &nginxv1alpha1.TransportServerArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String("consul-dns"),
			Namespace: pulumi.String(infra.Namespace),
		},
		Spec: &nginxv1alpha1.TransportServerSpecArgs{
			Listener: &nginxv1alpha1.TransportServerSpecListenerArgs{
				Name:     infra.LB.NginxIngress.UDPDNSListener.Name,
				Protocol: infra.LB.NginxIngress.UDPDNSListener.Protocol,
			},
			Upstreams: &nginxv1alpha1.TransportServerSpecUpstreamsArray{
				&nginxv1alpha1.TransportServerSpecUpstreamsArgs{
					Name: pulumi.String(nginxIngressDNSUpstreamName),
					// Same as https://github.com/hashicorp/consul-helm/blob/master/templates/dns-service.yaml#L6
					Service: pulumi.String(fmt.Sprintf("%s-dns", serviceName)),
					Port:    pulumi.Int(53),
				},
			},
			Action: &nginxv1alpha1.TransportServerSpecActionArgs{
				Pass: pulumi.String(nginxIngressDNSUpstreamName),
			},
		},
	})

	if err != nil {
		return err
	}

	_, err = corev1.NewService(infra.ctx, "nginx-lb-dns", &corev1.ServiceArgs{
		ApiVersion: pulumi.String("v1"),
		Kind:       pulumi.String("Service"),
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String("dns"),
			Namespace: pulumi.String(infra.Namespace),
		},
		Spec: &corev1.ServiceSpecArgs{
			Selector: pulumi.StringMap{
				"app": pulumi.String(infra.LB.NginxIngress.Name),
			},
			Type:           pulumi.String("LoadBalancer"),
			LoadBalancerIP: pulumi.String("192.168.75.1"),
			Ports: corev1.ServicePortArray{
				&corev1.ServicePortArgs{
					Name:     infra.LB.NginxIngress.UDPDNSListener.Name,
					Protocol: infra.LB.NginxIngress.UDPDNSListener.Protocol,
					Port:     pulumi.Int(53),
				},
			},
		},
	}, pulumi.DeleteBeforeReplace(true))
	if err != nil {
		return err
	}
	return nil
}
