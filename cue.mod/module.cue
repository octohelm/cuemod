module: "github.com/octohelm/cuemod"

require: {
	"dagger.io":                      "v0.2.18"
	"github.com/innoai-tech/runtime": "v0.0.0-20220610150018-093e8031f7ea"
	"k8s.io/api":                     "v0.24.1"
	"universe.dagger.io":             "v0.2.18"
}

require: {
	"k8s.io/apimachinery": "v0.24.1" @indirect()
}

replace: {
	"k8s.io/api":          "" @import("go")
	"k8s.io/apimachinery": "" @import("go")
}
