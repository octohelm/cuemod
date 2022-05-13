module: "github.com/octohelm/cuemod"

require: {
	"dagger.io":           "v0.2.8-0.20220512005159-64cb4f755695" @vcs("release-main")
	"k8s.io/api":          "v0.24.0"
	"k8s.io/apimachinery": "v0.24.0"
	"universe.dagger.io":  "v0.2.8-0.20220512005159-64cb4f755695" @vcs("release-main")
}

replace: {
	"dagger.io":          "github.com/morlay/dagger/pkg/dagger.io@release-main"
	"universe.dagger.io": "github.com/morlay/dagger/pkg/universe.dagger.io@release-main"
}

replace: {
	"k8s.io/api":          "" @import("go")
	"k8s.io/apimachinery": "" @import("go")
}
