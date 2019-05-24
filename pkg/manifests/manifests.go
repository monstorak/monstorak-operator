package manifests

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"

	monv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var (
	manLog = logf.Log.WithName("manifests_manifests")
)

const (
	ManifestConfigFile      string = "manifests_config.json"
	ManifestDoesNotExist    string = "Manifest for requested Storage Provider/Version may not exist"
	PrometheusObjectError   string = "Prometheus object could not be created"
	ManifestConfigReadError string = "Manifests config could not be read"
	ManifestParseError      string = "Manifest could not be parsed"
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

func GetManifestsFromConfig() (map[string]map[string]string, error) {
	var rules map[string]map[string]string
	config, err := ioutil.ReadFile(ManifestConfigFile)
	if err != nil {
		manLog.Error(err, ManifestConfigReadError)
		return nil, err
	}
	err = json.Unmarshal([]byte(config), &rules)
	if err != nil {
		manLog.Error(err, ManifestParseError)
	}
	return rules, err
}

func (f *Factory) PrometheusK8sRules(storageProvider, storageVersion string) (*monv1.PrometheusRule, error) {
	manLog.WithValues("Storage Provider", storageProvider, "Storage Version", storageVersion)
	// Get Rule manifests from config
	rules, err := GetManifestsFromConfig()
	if err != nil {
		return nil, err
	}
	r, err := f.NewPrometheusRule(MustAssetReader(rules[storageProvider][storageVersion]))
	if err != nil {
		manLog.Error(err, ManifestDoesNotExist)
	}
	return r, err
}

func (f *Factory) NewPrometheusRule(manifest io.Reader) (*monv1.PrometheusRule, error) {
	p, err := NewPrometheusRule(manifest)
	if err != nil {
		manLog.Error(err, PrometheusObjectError)
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
		manLog.Error(err, ManifestParseError)
		return nil, err
	}

	return &p, nil
}
