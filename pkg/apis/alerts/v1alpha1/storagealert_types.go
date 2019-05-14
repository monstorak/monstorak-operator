package v1alpha1

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AlertDetails defines the alert related details
type AlertDetails struct {
	LabelSelector       map[string]string `json:"labelSelector,omitempty"`
	PrometheusNamespace string            `json:"prometheusNamespace,omitempty"`
}

func (a AlertDetails) String() string {
	return fmt.Sprintf("LabelSelector: %+v  | PrometheusNamespace: %s",
		a.LabelSelector, a.PrometheusNamespace)
}

// StorageAlertSpec defines the desired state of StorageAlert
// +k8s:openapi-gen=true
type StorageAlertSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html

	// StorageName provides the name of the storage type
	StorageName string `json:"storageName,omitempty"`
	// StorageVersion provides the current version of the storage type used
	StorageVersion string       `json:"storageVersion,omitempty"`
	StorageAlert   AlertDetails `json:"storageAlert,omitempty"`
}

// StorageAlertStatus defines the observed state of StorageAlert
// +k8s:openapi-gen=true
type StorageAlertStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
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
