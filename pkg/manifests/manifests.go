package manifests

import (
	"bytes"
	"io"

	monv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

var (
	CephRules = "jsonnet/manifests/ceph-prometheus-rules.yaml"
)

func MustAssetReader(asset string) io.Reader {
	return bytes.NewReader(MustAsset(asset))
}

type Factory struct {
	namespace string
}

func NewFactory(namespace string) *Factory {
	return &Factory{
		namespace: namespace,
	}
}

func (f *Factory) PrometheusK8sRules() (*monv1.PrometheusRule, error) {
	r, err := f.NewPrometheusRule(MustAssetReader(CephRules))
	if err != nil {
		return nil, err
	}

	r.Namespace = f.namespace

	return r, nil
}

func (f *Factory) NewPrometheusRule(manifest io.Reader) (*monv1.PrometheusRule, error) {
	p, err := NewPrometheusRule(manifest)
	if err != nil {
		return nil, err
	}

	if p.GetNamespace() == "" {
		p.SetNamespace(f.namespace)
	}

	return p, nil
}

func NewPrometheusRule(manifest io.Reader) (*monv1.PrometheusRule, error) {
	p := monv1.PrometheusRule{}
	err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}
