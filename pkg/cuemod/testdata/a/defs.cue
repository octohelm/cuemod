package a

import (
	"k8s.io/api/core/v1"

	v0 "github.com/octohelm/cuemod-versioned-example/cuepkg"
	v2 "github.com/octohelm/cuemod-versioned-example/v2/cuepkg"

	root "github.com/octohelm/cuemod-versioned-example:example"
)

services: [Name=_]: v1.#Service & {
	metadata: name: Name
	metadata: labels: app:     Name
	metadata: labels: version: "\(v0.#Version) \(v2.#Version) \(root.#Version)"
}
