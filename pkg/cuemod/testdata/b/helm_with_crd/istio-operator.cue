package helm_with_crd

import (
	"extension/helm"

	"github.com/istio/istio/manifests/charts/istio-operator:chart"
)

"istio-operator": helm.Template(chart, {
	hub: "docker.io/querycapistio"
}, {
	name:      "istio-operator"
	namespace: name
})
