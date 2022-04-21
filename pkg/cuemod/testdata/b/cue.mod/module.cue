module: "github.com/x/b"

require: {
	"github.com/grafana/jsonnet-libs":           "v0.0.0" @vcs("1e36fec")
	"github.com/istio/istio":                    "v0.0.0-20210315064903-f88f93ff2b81"
	"github.com/jsonnet-libs/k8s-alpha":         "v0.0.0-20210118111845-5e0d0738721f"
	"github.com/rancher/local-path-provisioner": "v0.0.19"
	"github.com/x/a":                            "v0.0.0"
	"k8s.io/api":                                "v0.20.5"
	"k8s.io/apimachinery":                       "v0.20.5"
}

replace: {
	// lock version
	"github.com/rancher/local-path-provisioner": "@v0.0.19"
	// local replace
	"github.com/x/a": "../a"
}

replace: {
	"github.com/grafana/jsonnet-libs/grafana":                "" @import("jsonnet")
	"github.com/grafana/jsonnet-libs/ksonnet-util":           "" @import("jsonnet")
	"github.com/istio/istio/manifests/charts/istio-operator": "" @import("helm")
	"github.com/jsonnet-libs/k8s-alpha/1.19":                 "" @import("jsonnet")
	"github.com/rancher/local-path-provisioner/deploy/chart": "" @import("helm")
	"k8s.io/api":                                             "" @import("go")
	"k8s.io/apimachinery":                                    "" @import("go")
}
