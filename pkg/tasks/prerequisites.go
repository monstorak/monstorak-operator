package tasks

import (
	"github.com/monstorak/monstorak/pkg/common"
	"github.com/monstorak/monstorak/pkg/prometheus"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var prereqLog = logf.Log.WithName("tasks_prerequisites")

const (
	CheckPrerequisites  string = "Checking prerequisites"
	CheckNamespace      string = "Check if namespace exists"
	CheckNamespaceLabel string = "Check if namespace has labels"
	CheckServiceMonitor string = "Check if service monitor exists"
)

// Prerequisites needed for Monitoring to work with Storage
func Prerequisites(storageNamespace, serviceMonitor string) error {
	prereqLog.Info(CheckPrerequisites)
	// Check if Namespace exists
	prereqLog.Info(CheckNamespace, "Namespace", storageNamespace)
	_, err := common.GetNamespace(storageNamespace)
	if err != nil {
		return err
	}
	// Check if Namespace has required labels
	prereqLog.Info(CheckNamespaceLabel, "Namespace", storageNamespace)
	label := make(map[string]string)
	label["openshift.io/cluster-monitoring"] = "true"
	err = common.NamespaceHasLabels(storageNamespace, label)
	if err != nil {
		return err
	}
	// Check if ServiceMonitor exists
	prereqLog.Info(CheckServiceMonitor, "Service Monitor", serviceMonitor, "Namespace", storageNamespace)
	_, err = prometheus.ServiceMonitorExists(serviceMonitor, storageNamespace)
	if err != nil {
		return err
	}
	// TODO:Check if Role and RoleBinding exists
	return nil
}
