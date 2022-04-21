module: "github.com/x/examples"

require: {
	"github.com/octohelm/cuem": "v0.0.0-20220420103500-0e2c8e58fa59"
}

require: {
	"k8s.io/api":          "v0.23.5" @indirect()
	"k8s.io/apimachinery": "v0.23.5" @indirect()
}

replace: {
	"k8s.io/api":          "" @import("go")
	"k8s.io/apimachinery": "" @import("go")
}
