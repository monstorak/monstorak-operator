package prometheus

import (
	monv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var ruleLog = logf.Log.WithName("prometheus_prometheusRule")

const (
	PrometheusRuleCreateFailed   string = "Prometheus Rules could not be created"
	PrometheusRuleRetrieveFailed string = "Prometheus Rules could not be retrieved"
	PrometheusRuleUpdateFailed   string = "Prometheus Rules could not be updated"
)

// CreateOrUpdatePrometheusRule Creates or Updates prometheusRule object
func CreateOrUpdatePrometheusRule(p *monv1.PrometheusRule) error {
	ruleLog.WithValues("Prometheus Rule", p.GetName(), "Namespace", p.GetNamespace())
	mclient, err := newMonitoringClient()
	pclient := mclient.MonitoringV1().PrometheusRules(p.GetNamespace())
	oldRule, err := pclient.Get(p.GetName(), metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		_, err := pclient.Create(p)
		if err != nil {
			ruleLog.Error(err, PrometheusRuleCreateFailed)
			return err
		}
		return err
	}
	if err != nil {
		ruleLog.Error(err, PrometheusRuleRetrieveFailed)
		return err
	}
	p.ResourceVersion = oldRule.ResourceVersion
	_, err = pclient.Update(p)
	if err != nil {
		ruleLog.Error(err, PrometheusRuleUpdateFailed)
		return err
	}
	return err
}
