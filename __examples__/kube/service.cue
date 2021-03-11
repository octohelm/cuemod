package kube

import (
	"k8s.io/api/core/v1"
)

service: [Name=_]: v1.#Service & {
	metadata: name: Name
}

service: nginx: spec: {
	type:           "LoadBalancer"
	loadBalancerIP: "1.3.4.5"
	ports: [{
		name: "http"
		port: 80
	}, {
		name: "https"
		port: 443
	}]
}
