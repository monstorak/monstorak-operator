package prometheus

import (
	monitoringclient "github.com/coreos/prometheus-operator/pkg/client/versioned"
	"k8s.io/client-go/tools/clientcmd"
)

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
