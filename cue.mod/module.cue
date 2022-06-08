module: "github.com/octohelm/cuemod"

require: {
	"dagger.io":          "v0.2.18-0.20220608064710-85125429d9f3"
	"k8s.io/api":         "v0.24.1"
	"universe.dagger.io": "v0.2.18-0.20220608064710-85125429d9f3"
}

require: {
	"github.com/innoai-tech/runtime": "v0.0.0-20220608072531-1168e6749e1a" @indirect()
	"k8s.io/apimachinery":            "v0.24.1"                            @indirect()
}

replace: {
	"dagger.io":          "github.com/morlay/dagger/pkg/dagger.io@v0.2.18-0.20220608064710-85125429d9f3#release-main"
	"universe.dagger.io": "github.com/morlay/dagger/pkg/universe.dagger.io@v0.2.18-0.20220608064710-85125429d9f3#release-main"
}

replace: {
	"k8s.io/api":          "" @import("go")
	"k8s.io/apimachinery": "" @import("go")
}
