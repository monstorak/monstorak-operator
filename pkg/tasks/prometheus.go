package tasks

import (
	"github.com/monstorak/pkg/client"
	"github.com/monstorak/pkg/manifests"
	"github.com/pkg/errors"
)

type PrometheusTask struct {
	client  *client.Client
	factory *manifests.Factory
	config  *manifests.Config
}

func NewPrometheusTask(client *client.Client, factory *manifests.Factory, config *manifests.Config) *PrometheusTask {
	return &PrometheusTask{
		client:  client,
		factory: factory,
		config:  config,
	}
}

func (t *PrometheusTask) Run() error {
	pm, err := t.factory.PrometheusK8sRules()
	if err != nil {
		return errors.Wrap(err, "initializing Prometheus rules PrometheusRule failed")
	}

	err = t.client.CreateOrUpdatePrometheusRule(pm)
	if err != nil {
		return errors.Wrap(err, "reconciling Prometheus rules PrometheusRule failed")
	}
}
