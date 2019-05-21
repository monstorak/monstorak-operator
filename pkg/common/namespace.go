package common

import (
	apiV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	clientcmd "k8s.io/client-go/tools/clientcmd"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var nsLog = logf.Log.WithName("common_namespace")

func newCoreV1Client() (*v1.CoreV1Client, error) {
	client, err := clientcmd.BuildConfigFromFlags("", "")
	if err != nil {
		return nil, err
	}
	return v1.NewForConfig(client)
}

func getNamespace(namespace string) (*apiV1.Namespace, error) {
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
		nsLog.Error(err, "Could not find the namespace", "namespace : ", namespace)
		return nil, err
	}
	return ns, err
}

func updateNamespace(ns *apiV1.Namespace) error {
	coreClient, err := newCoreV1Client()
	if err != nil {
		return err
	}

	namespaceClient := coreClient.Namespaces()
	_, err = namespaceClient.Update(ns)
	if err != nil {
		nsLog.Error(err, "Namespace could not be updated", "Namespace :", ns)
		return err
	}
	return err
}

func AddLabelToNamespace(namespace string, label map[string]string) error {
	ns, err := getNamespace(namespace)
	if err != nil {
		nsLog.Error(err, "Could not add label to namespace. Check if namespace exists.", "namespace : ", namespace, "label :", label)
		return err
	}
	ns.ObjectMeta.SetLabels(label)
	err = updateNamespace(ns)
	if err != nil {
		nsLog.Error(err, "Could not add label to namespace", "namespace : ", namespace, "label :", label)
		return err
	}
	return err
}
