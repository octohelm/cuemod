module: "github.com/x/examples"

require: {
	"github.com/octohelm/cuem": "v0.1.0"
	"k8s.io/api":               "v0.23.5"
	"k8s.io/apimachinery":      "v0.23.5"
}

replace: {
	"k8s.io/api":          "" @import("go")
	"k8s.io/apimachinery": "" @import("go")
}
