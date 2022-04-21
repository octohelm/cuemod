# CUE Mod

[![GoDoc Widget](https://godoc.org/github.com/octohelm/cuemod?status.svg)](https://pkg.go.dev/github.com/octohelm/cuemod)
[![codecov](https://codecov.io/gh/octohelm/cuemod/branch/main/graph/badge.svg)](https://codecov.io/gh/octohelm/cuemod)
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
* Automate detect witch extractor should be used for generate code to cue.
    * golang supported
    * helm chart supported
    * jsonnet supported
* Post-processing where value with attribute `@translate(<name>)` when final marshalling.

## Spec `cue.mod/module.cue`

```cue
// module name
// for sub mod import, <module>/path/to/sub
// NOTICE: the module name should be a valid repo name
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
    // declare import method for special import path
    "github.com/rancher/local-path-provisioner/deploy/chart": "" @import("helm")
    // local replace
    // **notice** only works for current mod
    "github.com/x/a": "../a"
}
```

## Post-Processing `@translate(<name>)`

Now cuelang not provide api to added custom functions. So we use the attribute to mark which value should be translated
to other formats.

### `toml`

```cue
configmap: xxx: data: "xxx.toml": json.Marshal({ a: 1 }) @translate("toml")
// why json.Marshal here, just let type constraints happy
```

### `helm`

```cue
package localpathprovisioner

import (
	"github.com/rancher/local-path-provisioner/deploy/chart"
)

"local-path-provisioner": {
	chart
	values: {}
	release: name:      "local-path-provisioner"
	release: namespace: "local-path-provisioner"
} @translate("helm")
```

### `jsonnet`

```cue
package grafana

import (
	"encoding/json"

	"github.com/grafana/jsonnet-libs/grafana"
)

"grafana": {
	data: '''
		local grafana = import 'github.com/grafana/jsonnet-libs/grafana/grafana.libsonnet';
		
		{
		    config+:: (import 'config.jsonnet'),
		
		    prometheus_datasource:: grafana.datasource.new('prometheus', $.config.prometheus_url, type='prometheus', default=true),
		
		    grafana: grafana
		         + grafana.withAnonymous()
		         + grafana.addFolder('Example')
		         + grafana.addDatasource('prometheus', $.prometheus_datasource)
		         ,
		}
		'''

	imports: "github.com/grafana/jsonnet-libs/grafana/grafana.libsonnet": grafana["grafana.libsonnet"]
	imports: "config.jsonnet": code: json.Marshal({ prometheus_url: 'http://prometheus' })
} @translate("jsonnet")
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

## Plugin Kube

like [Tanka](https://tanka.dev), but for cuelang.

* only k3s/k8s which support `server-apply`
* make sure the cue file return an object as struct
  below ([full template](https://github.com/octohelm/cuem/blob/main/release/release.cue)):

```cue
apiVersion: "octohelm.tech/v1alpha"
kind:       "Release"

// release name
metadata: name:      "\(#name)"

// release namespace
metadata: namespace: "\(#namespace)"

// must an unique `NAME` of `kubectl config get-contexts`
metadata: labels: context: "\(#context)"

// nested object may contains kube resources
spec: {} 
```

```shell
cd ./__examples__
cuem k show ./clusters/demo/nginx.cue
# patch
cuem k show ./clusters/demo/nginx.cue '{ #values: image: tag: "latest" }'
cuem k apply ./clusters/demo/nginx.cue
cuem k prune ./clusters/demo/nginx.cue
cuem k delete ./clusters/demo/nginx.cue
```