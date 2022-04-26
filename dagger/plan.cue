package main

import (
	"strings"
	"dagger.io/dagger"
	"universe.dagger.io/docker"
	"github.com/octohelm/cuemod/cuepkg/tool"
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
		}

		filesystem: "./": read: {
			contents: dagger.#FS
			exclude: [
				"cue.mod/gen/",
				"cue.mod/pkg/",
				"build/",
			]
		}

		for _os in actions.build.os for _arch in actions.build.arch {
			filesystem: "build/output/\(actions.build.name)_\(_os)_\(_arch)": write: contents: actions.build["\(_os)"]["\(_arch)"].output
		}
	}

	actions: {
		_source: client.filesystem."./".read.contents
		_env: {
			for k, v in client.env if k != "$dagger" {
				"\(k)": v
			}
		}

		_imageName: "ghcr.io/octohelm/cuem"

		_version: [
				if strings.HasPrefix(_env.GIT_REF, "refs/tags/v") {
				strings.TrimPrefix(_env.GIT_REF, "refs/tags/v")
			},
			if strings.HasPrefix(_env.GIT_REF, "refs/heads/") {
				strings.TrimPrefix(_env.GIT_REF, "refs/heads/")
			},
			_env.VERSION,
		][0]

		_tag: _version

		info: tool.#GoModInfo & {
			source: _source
		}

		_archs: ["amd64", "arm64"]

		build: tool.#GoBuild & {
			source: _source
			arch:   _archs
			os: ["linux", "darwin", "windows"]
			env: _env & {
				CGO_ENABLED: "0"
			}
			ldflags: [
				"-s -w",
				"-X \(info.module)/version.Version=\(_version)",
				"-X \(info.module)/version.Revision=\(_env.GIT_SHA)",
			]
			package: "./cmd/cuem"
		}

		image: {
			for _arch in _archs {
				"\(_arch)": docker.#Dockerfile & {
					source: build.linux["\(_arch)"].output
					dockerfile: contents: """
						FROM alpine:3
						COPY ./cuem /bin/cuem
						ENTRYPOINT ["/bin/cuem"]
						"""
					platform: "linux/\(_arch)"
					label: {
						"org.opencontainers.image.source":   "https://\(info.module)"
						"org.opencontainers.image.revision": "\(_env.GIT_SHA)"
					}
				}
			}
		}

		push: docker.#Push & {
			dest: "\(_imageName):\(_tag)"
			images: {
				for _arch in _archs {
					"linux/\(_arch)": image["\(_arch)"].output
				}
			}
			auth: {
				username: _env.GH_USERNAME
				secret:   _env.GH_PASSWORD
			}
		}
	}
}
