package v1alpha1

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AlertDetails defines the alert related details
type AlertDetails struct {
	Enabled             bool              `json:"enabled,omitempty"`
	LabelSelector       map[string]string `json:"labelSelector,omitempty"`
	PrometheusNamespace string            `json:"prometheusNamespace,omitempty"`
}

func (a AlertDetails) String() string {
	return fmt.Sprintf("Enabled: %t  | LabelSelector: %+v  | PrometheusNamespace: %s",
		a.Enabled, a.LabelSelector, a.PrometheusNamespace)
}

// DashboardDetails defines the alert related details
type DashboardDetails struct {
	Enabled          bool              `json:"enabled,omitempty"`
	LabelSelector    map[string]string `json:"labelSelector,omitempty"`
	GrafanaNamespace string            `json:"grafanaNamespace,omitempty"`
}

func (d DashboardDetails) String() string {
	return fmt.Sprintf("Enabled: %t  | LabelSelector: %+v  | GrafanaNamespace: %s",
		d.Enabled, d.LabelSelector, d.GrafanaNamespace)
}

// StorageAlertsSpec defines the desired state of StorageAlerts
// +k8s:openapi-gen=true
type StorageAlertsSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html

	// StorageName provides the name of the storage type
	StorageName string `json:"storageName,omitempty"`
	// StorageVersion provides the current version of the storage type used
	StorageVersion   string           `json:"storageVersion,omitempty"`
	StorageAlert     AlertDetails     `json:"storageAlert,omitempty"`
	StorageDashboard DashboardDetails `json:"storageDashboard,omitempty"`
}

func (s StorageAlertsSpec) String() string {
	return fmt.Sprintf("StorageName: %s  | StorageVersion: %s  | "+
		"StorageAlert: %s |  StorageDashboard: %s",
		s.StorageName, s.StorageVersion, s.StorageAlert.String(), s.StorageDashboard.String())
}

// StorageAlertsStatus defines the observed state of StorageAlerts
// +k8s:openapi-gen=true
type StorageAlertsStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// StorageAlerts is the Schema for the storagealerts API
// +k8s:openapi-gen=true
type StorageAlerts struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StorageAlertsSpec   `json:"spec,omitempty"`
	Status StorageAlertsStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// StorageAlertsList contains a list of StorageAlerts
type StorageAlertsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []StorageAlerts `json:"items"`
}

func init() {
	SchemeBuilder.Register(&StorageAlerts{}, &StorageAlertsList{})
}
