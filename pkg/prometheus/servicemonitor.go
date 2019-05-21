package prometheus

import (
	monitoringv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	monitoringclient "github.com/coreos/prometheus-operator/pkg/client/versioned"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var smLog = logf.Log.WithName("common_serviceMonitor")

func newMonitoringClient() (*monitoringclient.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", "")
	if err != nil {
		return nil, err
	}
	monitoringClient, err := monitoringclient.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return monitoringClient, err
}

func CreateOrUpdateServiceMonitors(serviceMonitorName, port, namespace string, labels map[string]string) error {
	serviceMonitor := createServiceMonitor(serviceMonitorName, port, namespace, labels)
	serviceMonitorClient, err := newMonitoringClient()
	oldServiceMonitor, err := serviceMonitorClient.Monitoring().ServiceMonitors(namespace).Get(serviceMonitorName, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		_, err = serviceMonitorClient.Monitoring().ServiceMonitors(namespace).Create(serviceMonitor)
		if err != nil {
			smLog.Error(err, "Could not create service monitor", "ServiceMonitor :", serviceMonitor)
			return err
		}
		smLog.Info("Created service monitor")
		return err
	}
	if err != nil {
		smLog.Error(err, "Could not retrieve service monitor", "ServiceMonitor :", serviceMonitor)
		return err
	}
	serviceMonitor.ResourceVersion = oldServiceMonitor.ResourceVersion
	_, err = serviceMonitorClient.Monitoring().ServiceMonitors(namespace).Update(serviceMonitor)
	if err != nil {
		smLog.Error(err, "Could not update service monitor", "ServiceMonitor :", serviceMonitor)
		return err
	}
	return err
}

func createServiceMonitor(serviceMonitorName, port, namespace string, labels map[string]string) *monitoringv1.ServiceMonitor {
	svcMonitor := serviceMonitor(serviceMonitorName, namespace)
	labelSelector := metav1.LabelSelector{
		MatchLabels: labels,
	}
	endpoint := monitoringv1.Endpoint{
		Port:     port,
		Path:     "/metrics",
		Interval: "5s",
	}
	svcMonitor.Spec = monitoringv1.ServiceMonitorSpec{
		NamespaceSelector: monitoringv1.NamespaceSelector{
			MatchNames: []string{namespace},
		},
		Selector:  labelSelector,
		Endpoints: []monitoringv1.Endpoint{endpoint},
	}
	return svcMonitor
}

func serviceMonitor(serviceMonitorName string, namespace string) *monitoringv1.ServiceMonitor {
	label := make(map[string]string)
	label["app"] = serviceMonitorName
	return &monitoringv1.ServiceMonitor{
		TypeMeta: metav1.TypeMeta{
			Kind:       monitoringv1.DefaultCrdKinds.ServiceMonitor.Kind,
			APIVersion: monitoringv1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceMonitorName,
			Namespace: namespace,
			Labels:    label,
		},
	}
}
