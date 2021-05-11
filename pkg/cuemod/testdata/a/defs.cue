package a

import (
	"k8s.io/api/core/v1"
)

services: [Name=_]: v1.#Service & {
	metadata: name: Name
	metadata: labels: app: Name
}
