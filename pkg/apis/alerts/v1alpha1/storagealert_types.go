package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PrometheusSpec defines the prometheus to be used
type PrometheusSpec struct {
	Label     map[string]string `json:"label,omitempty"`
	Namespace string            `json:"namespace,omitempty"`
}

// StorageSpec defines the storages to be monitored
type StorageSpec struct {
	Provider  string   `json:"provider,omitempty"`
	Version   []string `json:"version,omitempty"`
	Namespace string   `json:"namespace,omitempty"`
}

// StorageAlertSpec defines the desired state of StorageAlert
// +k8s:openapi-gen=true
type StorageAlertSpec struct {
	Storage    []StorageSpec  `json:"storage,omitempty"`
	Prometheus PrometheusSpec `json:"prometheus,omitempty"`
}

// StorageAlertStatus defines the observed state of StorageAlert
// +k8s:openapi-gen=true
type StorageAlertStatus struct {
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// StorageAlert is the Schema for the storagealerts API
// +k8s:openapi-gen=true
type StorageAlert struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StorageAlertSpec   `json:"spec,omitempty"`
	Status StorageAlertStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// StorageAlertList contains a list of StorageAlert
type StorageAlertList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []StorageAlert `json:"items"`
}

func init() {
	SchemeBuilder.Register(&StorageAlert{}, &StorageAlertList{})
}
