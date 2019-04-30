package manifests

import (
	"testing"
)

func TestManifests(t *testing.T) {
	f := NewFactory("storage-monitoring")
	_, err := f.PrometheusK8sRules()
	if err != nil {
		t.Fatal(err)
	}
}
