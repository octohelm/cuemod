package modutil

import (
	"context"
	"testing"

	. "github.com/onsi/gomega"
)

func TestDownload(t *testing.T) {
	ctx := context.Background()

	pkgs := map[string]string{
		//"github.com/rancher/local-path-provisioner": "v0.0.19",
		//"github.com/jsonnet-libs/k8s-alpha":         "v0.0.0-20210118111845-5e0d0738721f",
		"github.com/grafana/jsonnet-libs": "master",
	}

	for p, v := range pkgs {
		t.Run("download "+p+"@"+v, func(t *testing.T) {
			e := Get(ctx, p, v, true)
			t.Log(e.Path, e.Version, e.Dir, e.Sum)
			NewWithT(t).Expect(e.Error).To(BeEmpty())
		})
	}
}
