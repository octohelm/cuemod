module: "github.com/x/a"

require: {
	"k8s.io/api":          "v0.20.5"
	"k8s.io/apimachinery": "v0.20.5"
}

replace: {
	"k8s.io/api":          "" @import("go")
	"k8s.io/apimachinery": "" @import("go")
}
