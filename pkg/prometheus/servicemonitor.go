package prometheus

import (
	monitoringv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var smLog = logf.Log.WithName("prometheus_serviceMonitor")

const (
	ServiceMonitorNotFound string = "Service Monitor could not be found"
)

func ServiceMonitorExists(serviceMonitorName, namespace string) (*monitoringv1.ServiceMonitor, error) {
	smLog.WithValues("Service Monitor", serviceMonitorName, "Namespace", namespace)
	serviceMonitorClient, err := newMonitoringClient()
	serviceMonitor, err := serviceMonitorClient.Monitoring().ServiceMonitors(namespace).Get(serviceMonitorName, metav1.GetOptions{})
	if err != nil {
		smLog.Error(err, ServiceMonitorNotFound)
		return nil, err
	}
	return serviceMonitor, nil
}
