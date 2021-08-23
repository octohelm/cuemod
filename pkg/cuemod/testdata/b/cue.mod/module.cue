module: "github.com/x/b"

require: {
	"github.com/grafana/jsonnet-libs":           "v0.0.0" @vcs("1e36fec")
	"github.com/istio/istio":                    "v0.0.0-20210315064903-f88f93ff2b81"
	"github.com/rancher/local-path-provisioner": "v0.0.19"
	"github.com/x/a":                            "v0.0.0"
}

require: {
	"github.com/jsonnet-libs/k8s-alpha": "v0.0.0-20210118111845-5e0d0738721f" @indirect()
	"k8s.io/api":                        "v0.20.5"                            @indirect()
	"k8s.io/apimachinery":               "v0.20.5"                            @indirect()
}

replace: {
	"github.com/grafana/jsonnet-libs": "@1e36fec"
	// helm with crd
	"github.com/istio/istio/manifests/charts/istio-operator": ""
	// lock version
	"github.com/rancher/local-path-provisioner": "@v0.0.19"
	// declare gen method
	"github.com/rancher/local-path-provisioner/deploy/chart": ""
	// local replace
	"github.com/x/a": "../a"
}
