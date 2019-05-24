package manifests

import (
	"testing"
)

func TestManifests(t *testing.T) {
	f := NewFactory("storage-monitoring")
	_, err := f.PrometheusK8sRules("ceph", "v14.2.1")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetManifestsFromConfig(t *testing.T) {
	_, err := GetManifestsFromConfig()
	if err != nil {
		t.Fatal(err)
	}
}
