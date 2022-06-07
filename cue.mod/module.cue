module: "github.com/octohelm/cuemod"

require: {
	"dagger.io":          "v0.2.17-0.20220607061721-387f0ef2334f"
	"k8s.io/api":         "v0.24.1"
	"universe.dagger.io": "v0.2.17-0.20220607061721-387f0ef2334f"
}

require: {
	"k8s.io/apimachinery": "v0.24.1" @indirect()
}

replace: {
	"dagger.io":          "github.com/morlay/dagger/pkg/dagger.io@v0.2.17-0.20220607061721-387f0ef2334f#release-main"
	"universe.dagger.io": "github.com/morlay/dagger/pkg/universe.dagger.io@v0.2.17-0.20220607061721-387f0ef2334f#release-main"
}

replace: {
	"k8s.io/api":          "" @import("go")
	"k8s.io/apimachinery": "" @import("go")
}
