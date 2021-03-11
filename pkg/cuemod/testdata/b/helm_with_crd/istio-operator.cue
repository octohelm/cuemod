package helm_with_crd

import (
	"github.com/istio/istio/manifests/charts/istio-operator:chart"
)

"istio-operator": {
	chart
	values: hub:        "docker.io/querycapistio"
	release: name:      "istio-operator"
	release: namespace: release.name
} @translate("helm")
