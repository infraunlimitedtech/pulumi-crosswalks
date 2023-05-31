package monitoring

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	helmv3 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v3"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	alertmanagerTelegramReceiver = pulumi.String("ops-telegram")
	alertmanagerDefaultReceiver  = "blackhole"
)

var (
	remotePrometheusAlerts = []string{
		// Very decent alerting rules
		"https://raw.githubusercontent.com/samber/awesome-prometheus-alerts/master/dist/rules/host-and-hardware/node-exporter.yml",
		"https://raw.githubusercontent.com/samber/awesome-prometheus-alerts/master/dist/rules/docker-containers/google-cadvisor.yml",
		"https://raw.githubusercontent.com/samber/awesome-prometheus-alerts/master/dist/rules/prometheus-self-monitoring/embedded-exporter.yml",
	}

	skippedAlerts = []string{
		// This alerts need to be fixed
		"HostContextSwitching",
		"PrometheusAlertmanagerE2eDeadManSwitch",
		"PrometheusAlertmanagerJobMissing",
		"PrometheusJobMissing",
		"PrometheusTimeserieCardinality",
	}
)

func (m *Stack) runVMAlert() error {
	appName := "victoria-metrics-alert"
	alerts := "groups: \n"
	// Need convert to int here. We can't get the pulumi secret as int
	telegramChatID, err := strconv.Atoi(m.VMAlert.Alertmanager.Telegram.ChatID)
	if err != nil {
		return fmt.Errorf("conv: %w", err)
	}

	// TO DO: move the file to monitoring repo
	alertmanagerTelegramTemplate, err := ioutil.ReadFile("addons/monitoring/alertmanager-telegram.tmpl")
	if err != nil {
		return fmt.Errorf("alertmanager template: %w", err)
	}

	// Build regex for skipped alerts
	skippedAlertsRegex := strings.Join(skippedAlerts, "|")

	// Get awesome prometheus rules
	// Make a context with deadline for the request.
	// Max deadline is 5 seconds per rule
	maxDeadline := 5 * len(remotePrometheusAlerts)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(maxDeadline))
	defer cancel()

	cli := &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 3 * time.Second,
			}).Dial,
		},
		Timeout: time.Second * 4,
	}
	for _, rule := range remotePrometheusAlerts {
		req, err := http.NewRequestWithContext(ctx, "GET", rule, nil)
		if err != nil {
			return err
		}
		rules, err := cli.Do(req)
		if err != nil {
			return err
		}
		defer rules.Body.Close()

		body, err := io.ReadAll(rules.Body)
		if err != nil {
			return err
		}
		// Remove 1st `groups:` word. We don't need it. It can be improved.
		alerts += string(bytes.Replace(body, []byte("groups:"), []byte(""), 1))
	}

	// We use a separate config map coz it is easier than parse all rules and create a bunch of pulumi.Map{}
	cm, err := corev1.NewConfigMap(m.ctx, fmt.Sprintf("%s-rules", appName), &corev1.ConfigMapArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Namespace: m.Namespace.Metadata.Name().Elem(),
		},
		Data: pulumi.StringMap{
			// Hardcoded filename in the helm chart.
			"alert-rules.yaml": pulumi.String(alerts),
		},
	})
	if err != nil {
		return fmt.Errorf("rules config map: %w", err)
	}

	_, err = helmv3.NewRelease(m.ctx, appName, &helmv3.ReleaseArgs{
		Name:      pulumi.String(appName),
		Chart:     pulumi.String("victoria-metrics-alert"),
		Namespace: m.Namespace.Metadata.Name().Elem(),
		RepositoryOpts: helmv3.RepositoryOptsArgs{
			Repo: pulumi.String("https://victoriametrics.github.io/helm-charts"),
		},
		Version: pulumi.String(m.VMAlert.Helm.Version),
		Values: pulumi.Map{
			"alertmanager": pulumi.Map{
				"enabled": pulumi.Bool(true),
				"templates": pulumi.Map{
					// The template file will be in same directory as config file
					"telegram.tmpl": pulumi.String(alertmanagerTelegramTemplate),
				},
				"config": pulumi.Map{
					"templates": pulumi.Array{
						// All templates
						pulumi.String("*.tmpl"),
					},
					"global": pulumi.Map{
						"resolve_timeout": pulumi.String("5m"),
					},
					"route": pulumi.Map{
						"receiver": pulumi.String(alertmanagerDefaultReceiver),
						"routes": pulumi.Array{
							pulumi.Map{
								"receiver": alertmanagerTelegramReceiver,
								"group_by": pulumi.Array{
									pulumi.String("alertname"),
								},
								"group_wait":      pulumi.String("30s"),
								"group_interval":  pulumi.String("5m"),
								"repeat_interval": pulumi.String("1h"),
								"matchers": pulumi.Array{
									pulumi.String("severity=~critical|warning"),
								},
								"routes": pulumi.Array{
									pulumi.Map{
										"repeat_interval": pulumi.String("768h"),
										"receiver":        pulumi.String(alertmanagerDefaultReceiver),
										"matchers": pulumi.Array{
											pulumi.Sprintf("alertname=~%s", skippedAlertsRegex),
										},
									},
								},
							},
						},
					},
					"receivers": pulumi.Array{
						pulumi.Map{
							"name": pulumi.String(alertmanagerDefaultReceiver),
						},
						pulumi.Map{
							"name": alertmanagerTelegramReceiver,
							"telegram_configs": pulumi.Array{
								pulumi.Map{
									"chat_id":       pulumi.Int(telegramChatID),
									"bot_token":     pulumi.String(m.VMAlert.Alertmanager.Telegram.Token),
									"send_resolved": pulumi.Bool(true),
									"message":       pulumi.String("{{ template \"telegram.infra.message\" . }}"),
								},
							},
						},
					},
				},
				"tolerations": pulumi.MapArray{
					pulumi.Map{
						"operator": pulumi.String("Exists"),
						"key":      pulumi.String("CriticalAddonsOnly"),
					},
					pulumi.Map{
						"operator": pulumi.String("Exists"),
						"key":      pulumi.String("node-role.kubernetes.io/control-plane"),
					},
				},
			},
			"server": pulumi.Map{
				"datasource": pulumi.Map{
					"url": pulumi.Sprintf("http://victoria-metrics:%d", m.VictoriaMetrics.Server.Port),
				},
				"service": pulumi.Map{
					"annotations": pulumi.Map{
						"prometheus.io/scrape": pulumi.String("true"),
						"prometheus.io/port":   pulumi.String("9100"),
					},
				},
				"configMap": cm.Metadata.Name().Elem(),
				"tolerations": pulumi.MapArray{
					pulumi.Map{
						"operator": pulumi.String("Exists"),
						"key":      pulumi.String("CriticalAddonsOnly"),
					},
					pulumi.Map{
						"operator": pulumi.String("Exists"),
						"key":      pulumi.String("node-role.kubernetes.io/control-plane"),
					},
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("helm: %w", err)
	}
	return nil
}
