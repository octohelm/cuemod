package release

import (
	core_v1 "k8s.io/api/core/v1"
	networking_v1 "k8s.io/api/networking/v1"
)

_network: {
	#namespace: string

	spec: {
		services: [Name = _]: core_v1.#Service & {
			metadata: name:      Name
			metadata: namespace: #namespace
		}

		ingresses: [Name = _]: networking_v1.#Ingress & {
			metadata: name:      Name
			metadata: namespace: #namespace
		}
	}
}
