package common

import (
	"fmt"

	apiV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	clientcmd "k8s.io/client-go/tools/clientcmd"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var nsLog = logf.Log.WithName("common_namespace")

const (
	NamespaceNotFound      string = "Namespace not found"
	NamespaceUpdateFailed  string = "Namespace could not be updated"
	NamespaceMissingLabels string = "Namespace is missing labels"
	NamespaceLabelFailed   string = "Namespace could not be labelled"
)

func newCoreV1Client() (*v1.CoreV1Client, error) {
	client, err := clientcmd.BuildConfigFromFlags("", "")
	if err != nil {
		return nil, err
	}
	return v1.NewForConfig(client)
}

func GetNamespace(namespace string) (*apiV1.Namespace, error) {
	nsLog.WithValues("Namespace", namespace)
	coreClient, err := newCoreV1Client()
	if err != nil {
		return nil, err
	}

	namespaceClient := coreClient.Namespaces()
	getOptions := metav1.GetOptions{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Namespace",
			APIVersion: apiV1.SchemeGroupVersion.String(),
		},
		ResourceVersion: "1",
	}

	ns, err := namespaceClient.Get(namespace, getOptions)
	if err != nil {
		nsLog.Error(err, NamespaceNotFound)
	}
	return ns, err
}

func UpdateNamespace(namespace *apiV1.Namespace) (*apiV1.Namespace, error) {
	nsLog.WithValues("Namespace", namespace.GetName())
	coreClient, err := newCoreV1Client()
	if err != nil {
		return nil, err
	}

	namespaceClient := coreClient.Namespaces()
	ns, err := namespaceClient.Update(namespace)
	if err != nil {
		nsLog.Error(err, NamespaceUpdateFailed)
	}
	return ns, err
}

func NamespaceHasLabels(namespace string, labels map[string]string) error {
	nsLog.WithValues("Namespace", namespace)
	ns, err := GetNamespace(namespace)
	if err != nil {
		return err
	}
	nsLabels := ns.GetLabels()
	for k, v := range labels {
		if nsLabels[k] != v {
			err = fmt.Errorf("No matching label found for %s=%s", k, v)
			nsLog.Error(err, NamespaceMissingLabels)
			return err
		}
	}
	return nil
}

func AddLabelToNamespace(namespace string, label map[string]string) error {
	nsLog.WithValues("Namespace", namespace, "Label", label)
	ns, err := GetNamespace(namespace)
	if err != nil {
		nsLog.Error(err, NamespaceLabelFailed)
		return err
	}
	ns.ObjectMeta.SetLabels(label)
	_, err = UpdateNamespace(ns)
	if err != nil {
		nsLog.Error(err, NamespaceLabelFailed)
		return err
	}
	return err
}
