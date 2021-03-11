# CUE Mod

[![GoDoc Widget](https://godoc.org/github.com/octohelm/cuemod?status.svg)](https://pkg.go.dev/github.com/octohelm/cuemod)
[![codecov](https://codecov.io/gh/octohelm/cuemod/branch/master/graph/badge.svg)](https://codecov.io/gh/octohelm/cuemod)
[![Go Report Card](https://goreportcard.com/badge/github.com/octohelm/cuemod)](https://goreportcard.com/report/github.com/octohelm/cuemod)

**WIP**

dependency management for [CUE](https://cuelang.org/) without committing `cue.mod`

## Requirements

* `git` or other vcs tool supported by go for vcs downloading.

## Install

```shell
go install github.com/octohelm/cuemod/cmd/cuem@latest
```

## Usage

### Quick Start

```shell 
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
 
cuem eval ./kube.cue
# build, will automately install deps if not exists or generator if needed.
```

### Dependency management

```
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
* Automate detect witch extractor should be used for generate code to cue.
    * golang supported
    * helm chart supported
    * jsonnet supported

## Spec `mod.cue`

```cue
// module name
// for sub mod import, <module>/path/to/sub
module: "github.com/x/b"

// automately resolve by the jsonnet code `import` or `importstr`
// rules follow go modules
require: {		
    // @vsc("master"), when upgrade, should use vcs version for upgrade.
    "github.com/grafana/jsonnet-libs":           "v0.0.0-20210315182639-887607c77457" @vcs("master")
    "github.com/jsonnet-libs/k8s-alpha":         "v0.0.0-20210118111845-5e0d0738721f" @indirect()
}

replace: {
    // version lock 
    "github.com/rancher/local-path-provisioner": "@v0.0.19"
    // declare gen method for special import path
    "github.com/rancher/local-path-provisioner/deploy/chart": "" @gen("helm")
    // local replace
    "github.com/x/a": "../a"
}
```

### Known issues

#### pkg name may not same as path

Some path like `github.com/istio/istio/manifests/charts/istio-operator`, the `istio-operator` is not a valid identifier
in cuelang. Should import with `github.com/istio/istio/manifests/charts/istio-operator:istio_operator`

#### dep incompatible go mod repo

For some go project like

```
$ go mod download -json github.com/grafana/loki@v2.1.0
{
        "Path": "github.com/grafana/loki",
        "Version": "v2.1.0",
        "Error": "github.com/grafana/loki@v2.1.0: invalid version: module contains a go.mod file, so major version must be compatible: should be v0 or v1, not v2"
}
```

Could config `mod.cue` replace with commit hash of the tag to hack

```cue
replace: {
    "github.com/grafana/loki": "@1b79df3",
}
```

## Plugin Kube [TODO]

like [Tanka](https://tanka.dev), but for cuelang.

return [`tanka.dev/Environment` object](https://tanka.dev/inline-environments#converting-to-an-inline-environment)

```
cd ./__examples__
cuem k show ./kube
cuem k apply ./kube
cuem k delete ./kube
cuem k prune ./kube
```