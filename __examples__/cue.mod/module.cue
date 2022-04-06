module: "github.com/x/examples"

require: {
	"github.com/octohelm/cuem": "v0.0.0-20210401081148-fa62bbf6a07a"
	"k8s.io/api":               "v0.20.5"
}

require: {
	"k8s.io/apimachinery": "v0.20.5" @indirect()
}
