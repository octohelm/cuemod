package main

import (
	"dagger.io/dagger"

	"universe.dagger.io/bash"
	"universe.dagger.io/alpine"
)

#outputs: {
	output: "/hello.txt"
}

dagger.#Plan & {
	client: filesystem: {
		"build/output.txt": write: contents: actions.test.export.files[#outputs.output]
	}

	actions: {
		_image: alpine.#Build & {
			packages: bash: {}
		}

		test: bash.#Run & {
			input: _image.output
			script: contents: "echo Hello World! > /hello.txt"
			export: files: "/hello.txt": string
		}
	}
}
