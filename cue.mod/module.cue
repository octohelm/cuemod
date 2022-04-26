module: "github.com/octohelm/cuemod"

require: {
	"github.com/morlay/dagger": "v0.2.8-0.20220506022722-cdd29e3cfad8" @vcs("poc-multi-arch")
	"k8s.io/api":               "v0.23.5"
	"k8s.io/apimachinery":      "v0.23.5"
}

replace: {
	"dagger.io":          "github.com/morlay/dagger/pkg/dagger.io@poc-multi-arch"
	"universe.dagger.io": "github.com/morlay/dagger/pkg/universe.dagger.io@poc-multi-arch"
}

replace: {
	"k8s.io/api":          "" @import("go")
	"k8s.io/apimachinery": "" @import("go")
}
