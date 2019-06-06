package tasks

import (
	monv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/monstorak/monstorak-operator/pkg/manifests"
	"github.com/monstorak/monstorak-operator/pkg/prometheus"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var promLog = logf.Log.WithName("tasks_prometheus")

const (
	RetrievePrometheusRule     string = "Retrieving Prometheus Rules"
	CreateUpdatePrometheusRule string = "Creating/Updating Prometheus Rules"
)

func GetPrometheusRule(storageNamespace, storageProvider, storageVersion string) (*monv1.PrometheusRule, error) {
	promLog.Info(RetrievePrometheusRule, "Storage Provider", storageProvider, "Storage Version", storageVersion)
	f := manifests.NewFactory(storageNamespace)
	prometheusK8sRules, err := f.PrometheusK8sRules(storageProvider, storageVersion)
	if err != nil {
		return nil, err
	}
	return prometheusK8sRules, nil
}

func DeployPrometheusRule(storageNamespace string, prometheusRule *monv1.PrometheusRule) error {
	promLog.Info(CreateUpdatePrometheusRule, "Namespace", storageNamespace, "Prometheus Rule", prometheusRule.GetName())
	prometheusRule.Namespace = storageNamespace
	err := prometheus.CreateOrUpdatePrometheusRule(prometheusRule)
	return err
}
