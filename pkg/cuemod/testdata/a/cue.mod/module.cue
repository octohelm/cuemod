module: "github.com/x/a"

require: {
	"github.com/octohelm/cuemod-versioned-example":    "v0.0.0-20230822070100-38465a937b3c"
	"github.com/octohelm/cuemod-versioned-example/v2": "v2.0.1"
	"k8s.io/api":                                      "v0.24.1"
	"std.x.io":                                        "v0.3.0"
}

require: {
	"k8s.io/apimachinery": "v0.24.1" @indirect()
}

replace: {
	"k8s.io/api":          "" @import("go")
	"k8s.io/apimachinery": "" @import("go")
}
