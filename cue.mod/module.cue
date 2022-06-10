module: "github.com/octohelm/cuemod"

require: {
	"dagger.io":                      "v0.2.18-0.20220608075319-28308bda2857"
	"github.com/innoai-tech/runtime": "v0.0.0-20220610150018-093e8031f7ea"
	"k8s.io/api":                     "v0.24.1"
	"universe.dagger.io":             "v0.2.18-0.20220608075319-28308bda2857"
}

require: {
	"k8s.io/apimachinery": "v0.24.1" @indirect()
}

replace: {
	"dagger.io":          "github.com/morlay/dagger/pkg/dagger.io@v0.2.18-0.20220608075319-28308bda2857#release-main"
	"universe.dagger.io": "github.com/morlay/dagger/pkg/universe.dagger.io@v0.2.18-0.20220608075319-28308bda2857#release-main"
}

replace: {
	"k8s.io/api":          "" @import("go")
	"k8s.io/apimachinery": "" @import("go")
}
