package monitoring

import (
	"fmt"
	"k8s-cluster/config"

	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Monitoring struct {
	NodeExporter    *NodeExporter
	VictoriaMetrics *VictoriaMetrics
	VMAlert         *VMAlert
}

type NodeExporter struct {
	Helm *config.HelmParams
}

type VictoriaMetrics struct {
	Helm   *config.HelmParams
	Server *VictoriaMetricsServer
}

type VMAlert struct {
	Enabled      *config.Status
	Alertmanager *VMAlertAlertmanager
	Helm         *config.HelmParams
}

type VMAlertAlertmanager struct {
	Telegram *VMAlertAlertmanagerTelegram
}

type VMAlertAlertmanagerTelegram struct {
	Token  string
	ChatID string
}

type VictoriaMetricsServer struct {
	ClusterIP string
	Port      int
}

type Stack struct {
	ctx             *pulumi.Context
	Namespace       *corev1.Namespace
	NodeExporter    *NodeExporter
	VictoriaMetrics *VictoriaMetrics
	VMAlert         *VMAlert
}

func New(cfg *Monitoring) *Monitoring {
	if cfg == nil {
		cfg = &Monitoring{}
	}

	if cfg.NodeExporter == nil {
		cfg.NodeExporter = &NodeExporter{}
	}

	if cfg.VictoriaMetrics == nil {
		cfg.VictoriaMetrics = &VictoriaMetrics{}
	}

	if cfg.VMAlert == nil {
		cfg.VMAlert = &VMAlert{}
	}

	return cfg
}

func (m *Monitoring) IsEnabled() bool {
	return true
}

func (m *Monitoring) Manage(ctx *pulumi.Context, ns *corev1.Namespace) error {
	namespace := "monitoring"

	// Setup all monitoring services and deployments to mon namespace
	ns, err := corev1.NewNamespace(ctx, namespace, &corev1.NamespaceArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name: pulumi.String(namespace),
		},
	})
	if err != nil {
		return fmt.Errorf("monitoring namespace: %w", err)
	}

	mon := &Stack{
		ctx:             ctx,
		Namespace:       ns,
		VictoriaMetrics: m.VictoriaMetrics,
		NodeExporter:    m.NodeExporter,
		VMAlert:         m.VMAlert,
	}

	err = mon.runNodeExporter()
	if err != nil {
		return fmt.Errorf("node-exporter: %w", err)
	}

	err = mon.runVM()
	if err != nil {
		return fmt.Errorf("victoria-metrics: %w", err)
	}

	err = mon.runVMAlert()
	if err != nil {
		return fmt.Errorf("victoria-metrics-alert: %w", err)
	}

	return nil
}
