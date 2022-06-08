module: "github.com/octohelm/cuemod"

require: {
	"dagger.io":          "v0.2.18-0.20220608061142-088cc896def9"
	"k8s.io/api":         "v0.24.1"
	"universe.dagger.io": "v0.2.18-0.20220608061142-088cc896def9"
}

require: {
	"github.com/innoai-tech/runtime": "v0.0.0-20220608062226-8ea0e7f327e6" @indirect()
	"k8s.io/apimachinery":            "v0.24.1"                            @indirect()
}

replace: {
	"dagger.io":          "github.com/morlay/dagger/pkg/dagger.io@v0.2.18-0.20220608061142-088cc896def9#release-main"
	"universe.dagger.io": "github.com/morlay/dagger/pkg/universe.dagger.io@v0.2.18-0.20220608061142-088cc896def9#release-main"
}

replace: {
	"k8s.io/api":          "" @import("go")
	"k8s.io/apimachinery": "" @import("go")
}
