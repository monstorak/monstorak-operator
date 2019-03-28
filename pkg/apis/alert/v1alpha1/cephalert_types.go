package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// // Alert defines whether alert is enabled or not
// type Alert struct {
// 	Enabled             bool              `json:"enabled"`
// 	LabelSelector       map[string]string `json:"labelSelector"`
// 	PrometheusNamespace string            `json:"prometheusNamespace"`
// }

// // Dashboard defines whether dashboard is enabled or not
// type Dashboard struct {
// 	Enabled          bool              `json:"enabled"`
// 	LabelSelector    map[string]string `json:"labelSelector"`
// 	GrafanaNamespace string            `json:"grafanaNamespace"`
// }

// // StorageType defines the type of storage which should be monitored for alerts
// type StorageType struct {
// 	TypeName string `json:"name"`
// 	Version  string `json:"version"`
// 	// Alert     Alert     `json:"alert"`
// 	// Dashboard Dashboard `json:"dashboard"`
// }

// CephAlertSpec defines the desired state of CephAlert
// +k8s:openapi-gen=true
type CephAlertSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
	TypeName string `json:"name"`
	Version  string `json:"version"`
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
