module: "github.com/x/b"

require: {
	"github.com/grafana/jsonnet-libs":                        "v0.0.0" @vcs("1e36fec")
	"github.com/grafana/jsonnet-libs/grafana":                "v0.0.0" @vcs("1e36fec")
	"github.com/istio/istio":                                 "v0.0.0-20210315064903-f88f93ff2b81"
	"github.com/istio/istio/manifests/charts/istio-operator": "v0.0.0-20210315064903-f88f93ff2b81"
	"github.com/jsonnet-libs/k8s-alpha":                      "v0.0.0-20210118111845-5e0d0738721f"
	"github.com/rancher/local-path-provisioner":              "v0.0.19"
	"github.com/rancher/local-path-provisioner/deploy/chart": "v0.0.19"
	"github.com/x/a":                                         "v0.0.0"
}

require: {
	"github.com/grafana/jsonnet-libs/ksonnet-util": "v0.0.0"                             @vcs("1e36fec") @indirect()
	"github.com/jsonnet-libs/k8s-alpha/1.19":       "v0.0.0-20210118111845-5e0d0738721f" @indirect()
	"k8s.io/api":                                   "v0.24.0"                            @indirect()
	"k8s.io/apimachinery":                          "v0.24.0"                            @indirect()
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
