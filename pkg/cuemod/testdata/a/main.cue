package a

import (
	"encoding/json"
	"extension/toml"
	"k8s.io/api/core/v1"
)

services: [Name=_]: v1.#Service & {
	metadata: name: Name
	metadata: labels: app: Name
}

services: test: {
	metadata: annotations: configAsJson: json.Marshal({a: int:               1})
	metadata: annotations: configAsToml: toml.FromJSON(json.Marshal({a: int: 1}))
}
