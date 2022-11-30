module: "github.com/octohelm/cuemod"

require: {
	"dagger.io":                      "v0.3.0"
	"github.com/innoai-tech/runtime": "v0.0.0-20221114082425-7a5e0cdc3035"
	"k8s.io/api":                     "v0.25.4"
	"universe.dagger.io":             "v0.3.0"
}

require: {
	"k8s.io/apimachinery": "v0.25.4" @indirect()
}

replace: {
	"k8s.io/api":          "" @import("go")
	"k8s.io/apimachinery": "" @import("go")
}
