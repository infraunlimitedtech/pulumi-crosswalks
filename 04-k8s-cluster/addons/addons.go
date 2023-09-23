package addons

import (
	"k8s-cluster/packages/kilo"
	"k8s-cluster/spec"
	nginx "k8s-cluster/addons/nginx-ingress"
	metricServer "k8s-cluster/addons/metric-server"
	kiloAddon "k8s-cluster/addons/kilo"
	"k8s-cluster/addons/cloudflared"
	"k8s-cluster/addons/monitoring"

	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	pulumiConfig "github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

type Addon interface {
	IsEnabled() bool
	Manage(*pulumi.Context, *corev1.Namespace) error
}

type Addons struct {
	Kilo         *kilo.Kilo
	NginxIngress *nginx.NginxIngress
	Monitoring   *monitoring.Monitoring
	Cloudflared  *cloudflared.Cloudflared
}

type Runner struct {
	ctx    *pulumi.Context
	Namespace *corev1.Namespace
	addons []Addon
}

const (
	namespace = "infra-system"
)

func NewRunner(ctx *pulumi.Context, s *spec.ClusterSpec, infraLayerNodeInfo pulumi.AnyOutput) (*Runner, error) {
	// Init vars from stack's config
	var pulumiAddonsCfg Addons
	cfg := pulumiConfig.New(ctx, "")
	cfg.RequireSecretObject("addons", &pulumiAddonsCfg)


	ns, err := corev1.NewNamespace(ctx, namespace, &corev1.NamespaceArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name: pulumi.String(namespace),
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
		NginxIngress: &nginx.NginxIngress{
			Name:    "nginx-ingress-addon",
			Domain:  s.InternalDomainZone,
			Helm:    pulumiAddonsCfg.NginxIngress.Helm,
			KubeAPI: pulumiAddonsCfg.NginxIngress.KubeAPI,
			ClusterIP: pulumiAddonsCfg.NginxIngress.ClusterIP,
			Replicas: pulumiAddonsCfg.NginxIngress.Replicas,
		},
		Monitoring: pulumiAddonsCfg.Monitoring,
		Kilo:       pulumiAddonsCfg.Kilo,
		Cloudflared: pulumiAddonsCfg.Cloudflared,
	}

	return &Runner{
		ctx:    ctx,
		Namespace: ns,
		addons: []Addon{
			nginx.New(a.NginxIngress),
			metricServer.New(),
			kiloAddon.New(a.Kilo, infraLayerNodeInfo),
			cloudflared.New(a.Cloudflared, a.NginxIngress.ClusterIP),
			monitoring.New(a.Monitoring),
		},
	}, nil
}

func (r *Runner) Run() error {
	for _, addon := range r.addons {
		if addon.IsEnabled() {
			if err := addon.Manage(r.ctx, r.Namespace); err != nil {
				return err
			}
		}
	}

	return nil
}
