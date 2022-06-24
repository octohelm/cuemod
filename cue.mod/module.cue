module: "github.com/octohelm/cuemod"

require: {
	"dagger.io":                      "v0.3.0"
	"github.com/innoai-tech/runtime": "v0.0.0-20220624011524-77ab1f5dc464"
	"k8s.io/api":                     "v0.24.1"
	"universe.dagger.io":             "v0.3.0"
}

require: {
	"k8s.io/apimachinery": "v0.24.1" @indirect()
}

replace: {
	"k8s.io/api":          "" @import("go")
	"k8s.io/apimachinery": "" @import("go")
}
