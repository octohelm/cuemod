package main

import (
	"strings"

	"wagon.octohelm.tech/core"
	"github.com/innoai-tech/runtime/cuepkg/golang"
)

client: core.#Client & {
	env: {
		GH_USERNAME: string | *""
		GH_PASSWORD: core.#Secret
	}
}

pkg: version: core.#Version & {}

settings: core.#Setting & {
	registry: "ghcr.io": auth: {
		username: client.env.GH_USERNAME
		secret:   client.env.GH_PASSWORD
	}
}

actions: go: golang.#Project & {
	source: {
		path: "."
		include: [
			"cmd/",
			"pkg/",
			"internal/",
			"go.mod",
			"go.sum",
		]
	}

	version: "\(pkg.version.output)"

	goos: ["linux", "darwin"]
	goarch: ["amd64", "arm64"]
	main: "./cmd/cuem"
	ldflags: [
		"-s -w",
		"-X \(go.module)/pkg/version.version=\(go.version)",
	]

	build: pre: [
		"go mod download",
	]

	ship: {
		name: "\(strings.Replace(go.module, "github.com/", "ghcr.io/", -1))/\(go.binary)"
		from: "gcr.io/distroless/static-debian11:debug"
	}
}
