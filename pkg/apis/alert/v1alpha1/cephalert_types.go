package v1alpha1

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AlertDetails defines the alert related details
type AlertDetails struct {
	Enabled             string            `json:"enabled,omitempty"`
	LabelSelector       map[string]string `json:"labelSelector,omitempty"`
	PrometheusNamespace string            `json:"prometheusNamespace,omitempty"`
}

func (a AlertDetails) String() string {
	return "Enabled: " + a.Enabled +
		"  LabelSelector: " + fmt.Sprintf("%+v", a.LabelSelector) +
		"  PrometheusNamespace: " + a.PrometheusNamespace
}

// DashboardDetails defines the alert related details
type DashboardDetails struct {
	Enabled          string            `json:"enabled,omitempty"`
	LabelSelector    map[string]string `json:"labelSelector,omitempty"`
	GrafanaNamespace string            `json:"grafanaNamespace,omitempty"`
}

func (d DashboardDetails) String() string {
	return "Enabled: " + d.Enabled +
		"  LabelSelector: " + fmt.Sprintf("%+v", d.LabelSelector) +
		"  GrafanaNamespace: " + d.GrafanaNamespace
}

// CephAlertSpec defines the desired state of CephAlert
// +k8s:openapi-gen=true
type CephAlertSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
	StorageName      string           `json:"storageName,omitempty"`
	StorageVersion   string           `json:"storageVersion,omitempty"`
	StorageAlert     AlertDetails     `json:"storageAlert,omitempty"`
	StorageDashboard DashboardDetails `json:"storageDashboard,omitempty"`
}

func (c CephAlertSpec) String() string {
	return "StorageName: " + c.StorageName +
		"  StorageVersion: " + c.StorageVersion +
		"  StorageAlert: " + c.StorageAlert.String() +
		"  StorageDashboard: " + c.StorageDashboard.String()
}

// CephAlertStatus defines the observed state of CephAlert
// +k8s:openapi-gen=true
type CephAlertStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CephAlert is the Schema for the cephalerts API
// +k8s:openapi-gen=true
type CephAlert struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CephAlertSpec   `json:"spec,omitempty"`
	Status CephAlertStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CephAlertList contains a list of CephAlert
type CephAlertList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CephAlert `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CephAlert{}, &CephAlertList{})
}
