module: "dagger-example"

require: {
	"dagger.io":          "v0.2.18-0.20220608023333-dfa7f38ab73d"
	"universe.dagger.io": "v0.2.18-0.20220608023333-dfa7f38ab73d"
}

replace: {
	"dagger.io":          "github.com/morlay/dagger/pkg/dagger.io@v0.2.18-0.20220608023333-dfa7f38ab73d#release-main"
	"universe.dagger.io": "github.com/morlay/dagger/pkg/universe.dagger.io@v0.2.18-0.20220608023333-dfa7f38ab73d#release-main"
}
