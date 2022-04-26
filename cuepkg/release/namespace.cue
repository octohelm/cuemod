package release

import (
	core_v1 "k8s.io/api/core/v1"
)

_namespace: {
	#namespace: string

	spec: namespace: core_v1.#Namespace & {
		metadata: name: "\(#namespace)"
	}
}
