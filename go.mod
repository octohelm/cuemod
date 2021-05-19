module github.com/octohelm/cuemod

go 1.16

require (
	cuelang.org/go v0.4.0-beta.2
	github.com/Masterminds/sprig/v3 v3.2.2 // indirect
	github.com/cockroachdb/apd/v2 v2.0.2 // indirect
	github.com/davecgh/go-spew v1.1.1
	github.com/emicklei/proto v1.9.0 // indirect
	github.com/fatih/color v1.11.0
	github.com/go-courier/logr v0.0.2
	github.com/go-courier/ptr v1.0.1
	github.com/go-logr/logr v0.4.0
	github.com/golang/glog v0.0.0-20210429001901-424d2337a529 // indirect
	github.com/google/go-jsonnet v0.17.0
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/googleapis/gnostic v0.5.5 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/json-iterator/go v1.1.11 // indirect
	github.com/jsonnet-bundler/jsonnet-bundler v0.4.0
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/onsi/gomega v1.12.0
	github.com/pelletier/go-toml v1.9.1
	github.com/pkg/errors v0.9.1
	github.com/pmezard/go-difflib v1.0.0
	github.com/protocolbuffers/txtpbfmt v0.0.0-20210430143850-408574485efa // indirect
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/objx v0.3.0
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a // indirect
	golang.org/x/mod v0.4.2
	golang.org/x/net v0.0.0-20210510120150-4163338589ed // indirect
	golang.org/x/oauth2 v0.0.0-20210514164344-f6687ab2804c // indirect
	golang.org/x/sys v0.0.0-20210514084401-e8d321eab015 // indirect
	golang.org/x/term v0.0.0-20210503060354-a79de5458b56
	golang.org/x/tools v0.1.1
	google.golang.org/appengine v1.6.7 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	helm.sh/helm/v3 v3.5.1
	k8s.io/api v0.20.5
	k8s.io/apiextensions-apiserver v0.20.5
	k8s.io/apimachinery v0.21.1
	k8s.io/client-go v0.20.5
	k8s.io/kube-openapi v0.0.0-20210421082810-95288971da7e // indirect
	k8s.io/utils v0.0.0-20210305010621-2afb4311ab10 // indirect
	sigs.k8s.io/controller-runtime v0.8.3
	sigs.k8s.io/structured-merge-diff/v4 v4.1.1 // indirect
	sigs.k8s.io/yaml v1.2.0
)

//replace helm.sh/helm/v3 => helm.sh/helm/v3 v3.5.1
