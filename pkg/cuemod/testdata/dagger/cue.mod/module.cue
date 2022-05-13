module: "dagger-example"

require: {
	"dagger.io":          "v0.2.10"
	"universe.dagger.io": "v0.2.10"
}

replace: {
	"dagger.io":          "github.com/dagger/dagger/pkg/dagger.io"
	"universe.dagger.io": "github.com/dagger/dagger/pkg/universe.dagger.io"
}
