module: "dagger-example"

require: {
	"dagger.io":          "v0.2.16"
	"universe.dagger.io": "v0.2.16"
}

replace: {
	"dagger.io":          "github.com/dagger/dagger/pkg/dagger.io@v0.2.16"
	"universe.dagger.io": "github.com/dagger/dagger/pkg/universe.dagger.io@v0.2.16"
}
