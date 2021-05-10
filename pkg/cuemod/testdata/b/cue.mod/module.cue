module: "github.com/x/b"

require: {
	"github.com/grafana/jsonnet-libs":           "v0.0.0-20210506152003-188bf1e03f0c" @vcs("master")
	"github.com/istio/istio":                    "v0.0.0-20210315064903-f88f93ff2b81"
	"github.com/jsonnet-libs/k8s-alpha":         "v0.0.0-20210118111845-5e0d0738721f" @indirect()
	"github.com/rancher/local-path-provisioner": "v0.0.19"
	"github.com/x/a":                            "v0.0.0"
	"k8s.io/api":                                "v0.20.5" @indirect()
	"k8s.io/apimachinery":                       "v0.20.5" @indirect()
}

replace: {
	// helm with crd
	"github.com/istio/istio/manifests/charts/istio-operator": ""
	// lock version
	"github.com/rancher/local-path-provisioner": "@v0.0.19"
	// declare gen method
	"github.com/rancher/local-path-provisioner/deploy/chart": ""
	// local replace
	"github.com/x/a": "../a"
}
