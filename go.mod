module github.com/octohelm/cuemod

go 1.16

require (
	cuelang.org/go v0.4.0
	github.com/cockroachdb/apd/v2 v2.0.2 // indirect
	github.com/cyphar/filepath-securejoin v0.2.3 // indirect
	github.com/davecgh/go-spew v1.1.1
	github.com/emicklei/proto v1.9.0 // indirect
	github.com/fatih/color v1.12.0
	github.com/go-courier/logr v0.0.2
	github.com/go-courier/ptr v1.0.1
	github.com/go-logr/logr v0.4.0
	github.com/golang/glog v0.0.0-20210429001901-424d2337a529 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/google/go-cmp v0.5.6 // indirect
	github.com/google/go-jsonnet v0.17.0
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/jsonnet-bundler/jsonnet-bundler v0.4.0
	github.com/mattn/go-isatty v0.0.13 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/onsi/gomega v1.13.0
	github.com/pelletier/go-toml v1.9.3
	github.com/pkg/errors v0.9.1
	github.com/pmezard/go-difflib v1.0.0
	github.com/prometheus/common v0.29.0 // indirect
	github.com/prometheus/procfs v0.7.0 // indirect
	github.com/protocolbuffers/txtpbfmt v0.0.0-20210430143850-408574485efa // indirect
	github.com/spf13/cobra v1.2.1
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/objx v0.3.0
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	go.uber.org/atomic v1.8.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	go.uber.org/zap v1.18.1 // indirect
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97 // indirect
	golang.org/x/mod v0.4.2
	golang.org/x/net v0.0.0-20210614182718-04defd469f4e // indirect
	golang.org/x/oauth2 v0.0.0-20210628180205-a41e5a781914 // indirect
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	golang.org/x/term v0.0.0-20210615171337-6886f2dfbf5b
	golang.org/x/tools v0.1.5
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	helm.sh/helm/v3 v3.6.3
	k8s.io/api v0.22.0
	k8s.io/apiextensions-apiserver v0.21.2
	k8s.io/apimachinery v0.22.0
	k8s.io/client-go v0.21.2
	k8s.io/kube-openapi v0.0.0-20210527164424-3c818078ee3d // indirect
	k8s.io/utils v0.0.0-20210709001253-0e1f9d693477 // indirect
	sigs.k8s.io/controller-runtime v0.9.2
	sigs.k8s.io/yaml v1.2.0
)

//replace github.com/go-logr/logr => github.com/go-logr/logr v0.4.0
//
//replace github.com/go-logr/zapr => github.com/go-logr/zapr v0.4.0
