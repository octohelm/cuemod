package modutil

import (
	"context"
	"testing"

	. "github.com/onsi/gomega"
)

func TestDownload(t *testing.T) {
	ctx := context.Background()

	pkgs := map[string]string{
		"k8s.io/api": "v0.24.0",
		"github.com/rancher/local-path-provisioner": "v0.0.19",
		"github.com/grafana/jsonnet-libs":           "master",
	}

	for p, v := range pkgs {
		t.Run("download "+p+"@"+v, func(t *testing.T) {
			e := Get(ctx, p, v, true)
			t.Log(e.Path, e.Version, e.Dir, e.Sum)
			NewWithT(t).Expect(e.Error).To(BeEmpty())
		})
	}
}
