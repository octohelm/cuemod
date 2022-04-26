package release

import (
	core_v1 "k8s.io/api/core/v1"
)

_storage: {
	#namespace: string

	spec: persistentVolumeClaims: [Name = _]: core_v1.#PersistentVolumeClaim & {
		metadata: name:      Name
		metadata: namespace: #namespace
	}
}
