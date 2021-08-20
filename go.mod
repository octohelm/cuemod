module github.com/octohelm/cuemod

go 1.16

require (
	cuelang.org/go v0.4.0
	github.com/BurntSushi/toml v0.4.1 // indirect
	github.com/cockroachdb/apd/v2 v2.0.2 // indirect
	github.com/cyphar/filepath-securejoin v0.2.3 // indirect
	github.com/davecgh/go-spew v1.1.1
	github.com/emicklei/proto v1.9.1 // indirect
	github.com/fatih/color v1.12.0
	github.com/go-courier/logr v0.0.2
	github.com/go-courier/ptr v1.0.1
	github.com/go-logr/logr v0.4.0
	github.com/golang/glog v0.0.0-20210429001901-424d2337a529 // indirect
	github.com/google/go-cmp v0.5.6 // indirect
	github.com/google/go-jsonnet v0.17.0
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/jsonnet-bundler/jsonnet-bundler v0.4.0
	github.com/mattn/go-isatty v0.0.13 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/onsi/gomega v1.15.0
	github.com/pelletier/go-toml v1.9.3
	github.com/pkg/errors v0.9.1
	github.com/pmezard/go-difflib v1.0.0
	github.com/prometheus/common v0.30.0 // indirect
	github.com/prometheus/procfs v0.7.2 // indirect
	github.com/protocolbuffers/txtpbfmt v0.0.0-20210726093723-1671b78f4579 // indirect
	github.com/spf13/cast v1.4.1 // indirect
	github.com/spf13/cobra v1.2.1
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/objx v0.3.0
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	go.uber.org/zap v1.19.0 // indirect
	golang.org/x/crypto v0.0.0-20210813211128-0a44fdfbc16e // indirect
	golang.org/x/mod v0.5.0
	golang.org/x/net v0.0.0-20210813160813-60bc85c4be6d // indirect
	golang.org/x/oauth2 v0.0.0-20210810183815-faf39c7919d5 // indirect
	golang.org/x/sys v0.0.0-20210816032535-30e4713e60e3 // indirect
	golang.org/x/term v0.0.0-20210615171337-6886f2dfbf5b
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/tools v0.1.5
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/square/go-jose.v2 v2.6.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	helm.sh/helm/v3 v3.6.3
	k8s.io/api v0.22.0
	k8s.io/apiextensions-apiserver v0.22.0
	k8s.io/apimachinery v0.22.1
	k8s.io/client-go v0.22.0
	k8s.io/klog/v2 v2.10.0 // indirect
	k8s.io/kube-openapi v0.0.0-20210527164424-3c818078ee3d // indirect
	k8s.io/utils v0.0.0-20210802155522-efc7438f0176 // indirect
	sigs.k8s.io/controller-runtime v0.9.6
	sigs.k8s.io/yaml v1.2.0
)

replace (
	github.com/go-logr/logr => github.com/go-logr/logr v0.4.0
	github.com/go-logr/zapr => github.com/go-logr/zapr v0.4.0
)
