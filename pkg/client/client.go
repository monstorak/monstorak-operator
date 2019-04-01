// Copyright 2018 The Cluster Monitoring Operator Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"fmt"
	"time"

	mon "github.com/coreos/prometheus-operator/pkg/apis/monitoring"
	monv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	monitoring "github.com/coreos/prometheus-operator/pkg/client/versioned"
	"github.com/coreos/prometheus-operator/pkg/k8sutil"
	prometheusoperator "github.com/coreos/prometheus-operator/pkg/prometheus"
	"github.com/pkg/errors"
	extensionsobj "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	deploymentCreateTimeout = 5 * time.Minute
)

type Client struct {
	version           string
	namespace         string
	namespaceSelector string
	kclient           kubernetes.Interface
	mclient           monitoring.Interface
	eclient           apiextensionsclient.Interface
}

func New(cfg *rest.Config, version string, namespace string, namespaceSelector string) (*Client, error) {
	mclient, err := monitoring.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	kclient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "creating kubernetes clientset client")
	}

	eclient, err := apiextensionsclient.NewForConfig(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "creating apiextensions client")
	}

	return &Client{
		version:           version,
		namespace:         namespace,
		namespaceSelector: namespaceSelector,
		kclient:           kclient,
		mclient:           mclient,
		eclient:           eclient,
	}, nil
}

func (c *Client) KubernetesInterface() kubernetes.Interface {
	return c.kclient
}

func (c *Client) Namespace() string {
	return c.namespace
}

func (c *Client) NamespacesToMonitor() ([]string, error) {
	namespaces, err := c.kclient.CoreV1().Namespaces().List(metav1.ListOptions{
		LabelSelector: c.namespaceSelector,
	})
	if err != nil {
		return nil, errors.Wrap(err, "listing namespaces failed")
	}

	namespaceNames := make([]string, len(namespaces.Items))
	for i, namespace := range namespaces.Items {
		namespaceNames[i] = namespace.Name
	}

	return namespaceNames, nil
}

func (c *Client) WaitForPrometheusOperatorCRDsReady() error {
	return wait.Poll(time.Second, time.Minute*5, func() (bool, error) {
		err := c.WaitForCRDReady(k8sutil.NewCustomResourceDefinition(monv1.DefaultCrdKinds.Prometheus, mon.GroupName, map[string]string{}, false))
		if err != nil {
			return false, err
		}

		err = c.WaitForCRDReady(k8sutil.NewCustomResourceDefinition(monv1.DefaultCrdKinds.Alertmanager, mon.GroupName, map[string]string{}, false))
		if err != nil {
			return false, err
		}

		err = c.WaitForCRDReady(k8sutil.NewCustomResourceDefinition(monv1.DefaultCrdKinds.ServiceMonitor, mon.GroupName, map[string]string{}, false))
		if err != nil {
			return false, err
		}

		_, err = c.mclient.MonitoringV1().Prometheuses(c.namespace).List(metav1.ListOptions{})
		if err != nil {
			return false, err
		}

		_, err = c.mclient.MonitoringV1().Alertmanagers(c.namespace).List(metav1.ListOptions{})
		if err != nil {
			return false, err
		}

		_, err = c.mclient.MonitoringV1().ServiceMonitors(c.namespace).List(metav1.ListOptions{})
		if err != nil {
			return false, err
		}

		return true, nil
	})
}

func (c *Client) WaitForCRDReady(crd *extensionsobj.CustomResourceDefinition) error {
	return wait.Poll(5*time.Second, 5*time.Minute, func() (bool, error) {
		return c.CRDReady(crd)
	})
}

func (c *Client) CRDReady(crd *extensionsobj.CustomResourceDefinition) (bool, error) {
	crdClient := c.eclient.ApiextensionsV1beta1().CustomResourceDefinitions()

	crdEst, err := crdClient.Get(crd.ObjectMeta.Name, metav1.GetOptions{})
	if err != nil {
		return false, err
	}
	for _, cond := range crdEst.Status.Conditions {
		switch cond.Type {
		case extensionsobj.Established:
			if cond.Status == extensionsobj.ConditionTrue {
				return true, err
			}
		case extensionsobj.NamesAccepted:
			if cond.Status == extensionsobj.ConditionFalse {
				return false, fmt.Errorf("CRD naming conflict (%s): %v", crd.ObjectMeta.Name, cond.Reason)
			}
		}
	}
	return false, err
}

func (c *Client) WaitForPrometheus(p *monv1.Prometheus) error {
	var lastErr error
	if err := wait.Poll(time.Second*10, time.Minute*5, func() (bool, error) {
		p, err := c.mclient.MonitoringV1().Prometheuses(p.GetNamespace()).Get(p.GetName(), metav1.GetOptions{})
		if err != nil {
			return false, errors.Wrap(err, "retrieving Prometheus object failed")
		}
		status, _, err := prometheusoperator.PrometheusStatus(c.kclient.(*kubernetes.Clientset), p)
		if err != nil {
			return false, errors.Wrap(err, "retrieving Prometheus status failed")
		}

		expectedReplicas := *p.Spec.Replicas
		if status.UpdatedReplicas == expectedReplicas && status.AvailableReplicas >= expectedReplicas {
			return true, nil
		}
		lastErr = fmt.Errorf("expected %d replicas, updated %d and available %d", expectedReplicas, status.UpdatedReplicas, status.AvailableReplicas)
		return false, nil
	}); err != nil {
		if err == wait.ErrWaitTimeout && lastErr != nil {
			err = lastErr
		}
		return errors.Wrap(err, "waiting for Prometheus")
	}
	return nil
}

func (c *Client) CreateOrUpdatePrometheusRule(p *monv1.PrometheusRule) error {
	pclient := c.mclient.MonitoringV1().PrometheusRules(p.GetNamespace())
	oldRule, err := pclient.Get(p.GetName(), metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		_, err := pclient.Create(p)
		return errors.Wrap(err, "creating PrometheusRule object failed")
	}
	if err != nil {
		return errors.Wrap(err, "retrieving PrometheusRule object failed")
	}

	p.ResourceVersion = oldRule.ResourceVersion

	_, err = pclient.Update(p)
	return errors.Wrap(err, "updating PrometheusRule object failed")
}
