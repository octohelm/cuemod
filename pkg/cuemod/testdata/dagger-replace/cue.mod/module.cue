module: "dagger-example"

require: {
	"dagger.io":          "v0.2.10"
	"universe.dagger.io": "v0.2.10"
}

replace: {
	"dagger.io":          "github.com/morlay/dagger/pkg/dagger.io@release-main"
	"universe.dagger.io": "github.com/morlay/dagger/pkg/universe.dagger.io@release-main"
}
