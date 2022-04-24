module: "dagger-example"

require: {
	"github.com/dagger/dagger": "v0.2.7"
}

replace: {
	"dagger.io":          "github.com/dagger/dagger/pkg/dagger.io"
	"universe.dagger.io": "github.com/dagger/dagger/pkg/universe.dagger.io"
}
