module: "github.com/octohelm/cuemod"

require: {
	"dagger.io":          "v0.2.18-0.20220608023333-dfa7f38ab73d"
	"k8s.io/api":         "v0.24.1"
	"universe.dagger.io": "v0.2.18-0.20220608023333-dfa7f38ab73d"
}

require: {
	"k8s.io/apimachinery": "v0.24.1" @indirect()
}

replace: {
	"dagger.io":          "github.com/morlay/dagger/pkg/dagger.io@v0.2.18-0.20220608023333-dfa7f38ab73d#release-main"
	"universe.dagger.io": "github.com/morlay/dagger/pkg/universe.dagger.io@v0.2.18-0.20220608023333-dfa7f38ab73d#release-main"
}

replace: {
	"k8s.io/api":          "" @import("go")
	"k8s.io/apimachinery": "" @import("go")
}
