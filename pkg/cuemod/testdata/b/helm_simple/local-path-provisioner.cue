package helm_simple

import (
	"github.com/rancher/local-path-provisioner/deploy/chart"
)

"local-path-provisioner": {
	chart
	release: name:      "local-path-provisioner"
	release: namespace: "local-path-provisioner"
} @translate("helm")
