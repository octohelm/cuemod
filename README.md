# CUE Mod

[![GoDoc Widget](https://godoc.org/github.com/octohelm/cuemod?status.svg)](https://pkg.go.dev/github.com/octohelm/cuemod)
[![Go Report Card](https://goreportcard.com/badge/github.com/octohelm/cuemod)](https://goreportcard.com/report/github.com/octohelm/cuemod)

**ALPHA VERSION**

**May deprecated when [cue modules official supported](https://github.com/cue-lang/cue/issues/851)**

dependency management for [CUE](https://cuelang.org/) without committing `cue.mod`

## Requirements

* `git` or other vcs tool supported by go for vcs downloading.

## Install

```shell
go install github.com/octohelm/cuemod/cmd/cuem@latest
```

## Usage

### Quick Start

```bash 
mkdir -p ./demo && cd ./demo

cat << EOT > kube.cue
package kube

import (
   apps_v1 "k8s.io/api/apps/v1"
)

deployment: [string]: apps_v1.#Deployment

_labels: { "app": "nginx" }

deployment: nginx: spec: selector: matchLabels: _labels
deployment: nginx: spec: template: metadata: labels: _labels
deployment: nginx: spec: template: spec: {
	containers: [{
		name: "nginx"
		image: "nginx:1.11.10-alpine"
	}]
}
EOT
 
cuem eval -o kube.yaml ./kube.cue
# build, will automately install deps if not exists or generator if needed.

cuem eval -o ./kube.single-file.cue ./kube.cue
# will bundle to one single cue file
```

### Dependency management

```bash
# auto added deps
cuem get ./...

# upgrade deps
cuem get -u ./...

# install dep with special version
cuem get github.com/grafana/jsonnet-libs@latest
```

## Features

* Dependency management based on go modules
    * all dependency codes will download under `$(go env GOMODCACHE)`
    * `GOPROXY` supported to speed up downloading
* Extract to cue pkg from other language or schema spec.
    * `golang` supported
    * k8s `crd` json supported

## Spec `cue.mod/module.cue`

```cue
// module name
// for sub mod import, <module>/path/to/sub
// NOTICE: the module name should be a valid repo name
module: "github.com/octohelm/cuemod"

require: {
	 // @vsc("release-main"), when upgrade, should use vcs version for upgrade.
	"dagger.io":          "v0.2.8-0.20220512005159-64cb4f755695" @vcs("release-main")
	"k8s.io/api":         "v0.24.0"
	"universe.dagger.io": "v0.2.8-0.20220512005159-64cb4f755695" @vcs("release-main")
}

require: {
	"k8s.io/apimachinery": "v0.24.0" @indirect()
}

replace: {
	// replace module with spec version
	"dagger.io":          "github.com/morlay/dagger/pkg/dagger.io@release-main"
	"universe.dagger.io": "github.com/morlay/dagger/pkg/universe.dagger.io@release-main"
	
	 // **notice** only works for current mod
    "github.com/x/a": "../a"
}

replace: {
	// declare import method for special import path
	"k8s.io/api":          "" @import("go")
	"k8s.io/apimachinery": "" @import("go")
}
```

### Known issues

#### pkg name may not same as path

Some path like `github.com/istio/istio/manifests/charts/istio-operator`, the `istio-operator` is not a valid identifier
in `cue-lang`. Should import with `github.com/istio/istio/manifests/charts/istio-operator:istio_operator`
