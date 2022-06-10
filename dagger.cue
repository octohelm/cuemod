package main

import (
	"strings"
	"dagger.io/dagger"
	"dagger.io/dagger/core"
	"universe.dagger.io/docker"

	"github.com/innoai-tech/runtime/cuepkg/tool"
	"github.com/innoai-tech/runtime/cuepkg/golang"
	"github.com/innoai-tech/runtime/cuepkg/debian"
)

dagger.#Plan & {
	client: {
		env: {
			VERSION: string | *"dev"
			GIT_SHA: string | *""
			GIT_REF: string | *""

			GOPROXY:   string | *""
			GOPRIVATE: string | *""
			GOSUMDB:   string | *""

			GH_USERNAME: string | *""
			GH_PASSWORD: dagger.#Secret

			LINUX_MIRROR: string | *""
		}

		filesystem: "build/output": write: contents: actions.export.output
	}

	actions: {
		version: (tool.#ResolveVersion & {ref: client.env.GIT_REF, version: "\(client.env.VERSION)"}).output

		src: core.#Source & {
			path: "."
			include: [
				"cmd/",
				"pkg/",
				"go.mod",
				"go.sum",
			]
		}

		info: golang.#Info & {
			"source": src.output
		}

		build: golang.#Build & {
			source: src.output
			go: {
				os: ["linux", "darwin"]
				arch: ["amd64", "arm64"]
				package: "./cmd/cuem"
				ldflags: [
					"-s -w",
					"-X \(info.module)/pkg/version.Version=\(version)",
					"-X \(info.module)/pkg/version.Revision=\(client.env.GIT_SHA)",
				]
			}
			run: env: {
				GOPROXY:   client.env.GOPROXY
				GOPRIVATE: client.env.GOPRIVATE
				GOSUMDB:   client.env.GOSUMDB
			}
			image: mirror: client.env.LINUX_MIRROR
		}

		export: tool.#Export & {
			archive: true
			directories: {
				for _os in build.go.os for _arch in build.go.arch {
					"\(build.go.name)_\(_os)_\(_arch)": build["\(_os)/\(_arch)"].output
				}
			}
		}

		images: {
			for arch in build.go.arch {
				"linux/\(arch)": docker.#Build & {
					steps: [
						debian.#Build & {
							platform: "linux/\(arch)"
							mirror:   client.env.LINUX_MIRROR
							packages: {
								"ca-certificates": _
							}
						},
						docker.#Copy & {
							contents: build["linux/\(arch)"].output
							source:   "./cuem"
							dest:     "/cuem"
						},
						docker.#Set & {
							config: {
								label: {
									"org.opencontainers.image.source":   "https://\(info.module)"
									"org.opencontainers.image.revision": "\(client.env.GIT_SHA)"
								}
								workdir: "/"
								entrypoint: ["/cuem"]
							}
						},
					]
				}
			}
		}

		ship: {
			_push: docker.#Push & {
				dest: "\(strings.Replace(info.module, "github.com/", "ghcr.io/", -1))/cuem:\(version)"
				"images": {
					for p, image in images {
						"\(p)": image.output
					}
				}
				auth: {
					username: client.env.GH_USERNAME
					secret:   client.env.GH_PASSWORD
				}
			}

			result: _push.result
		}
	}
}
