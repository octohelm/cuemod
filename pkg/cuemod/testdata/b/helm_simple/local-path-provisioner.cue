package helm_simple

import (
	"extension/helm"

	"github.com/rancher/local-path-provisioner/deploy/chart"
)

"local-path-provisioner": helm.Template(chart, {

}, {
	name:      "local-path-provisioner"
	namespace: "local-path-provisioner"
})
