package tool

import (
	"regexp"
	"path"
	"strings"

	"dagger.io/dagger"
	"dagger.io/dagger/core"
	"universe.dagger.io/docker"
)

#GoModInfo: {
	source: dagger.#FS

	_readGoMod: core.#ReadFile & {
		input: source
		path:  "go.mod"
	}

	go:     regexp.FindSubmatch(#"go (.+)\n"#, _readGoMod.contents)[1]
	module: regexp.FindSubmatch(#"module (.+)\n"#, _readGoMod.contents)[1]
}

#GoBuild: {
	source:  dagger.#FS
	package: string
	env: [Key=string]: string | dagger.#Secret
	os: [...string]
	arch: [...string]
	ldflags: *["-x -w"] | [...string]
	name:    path.Base(package)

	gomod: #GoModInfo & {
		"source": source
	}

	_image: docker.#Pull & {
		source: "golang:\(gomod.go)-bullseye"
	}

	_sourcePath: "/go/src"

	_cacheMounts: {
		_paths: {
			mod_cache:   "/go/pkg/mod"
			build_cache: "/root/.cache/go-build"
		}

		for n, p in _paths {
			"\(p)": core.#Mount & {
				dest:     p
				contents: core.#CacheDir & {
					id: "go_\(n)"
				}
			}
		}
	}

	_#go: docker.#Run & {
		input:   _image.output
		workdir: _sourcePath
		"env":   env
		mounts:  _cacheMounts & {
			"source": core.#Mount & {
				dest:     _sourcePath
				contents: source
			}
		}
		command: name: "go"
	}

	_dep: _#go & {
		command: args: [
			"mod",
			"download",
			"-x",
		]
	}

	for _os in os {
		"\(_os)": {
			for _arch in arch {
				"\(_arch)": {
					_build: _#go & {
						env: {
							GOOS:   _os
							GOARCH: _arch
						}
						command: {
							args: [
								"\(package)",
							]
							flags: {
								build:      true
								"-ldflags": strings.Join(ldflags, " ")
								"-o":       "/output/\(name)"
							}
						}
						export: directories: "/output": _
					}

					output: _build.export.directories."/output"
				}
			}
		}
	}
}
