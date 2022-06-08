module: "github.com/x/b"

require: {
	"github.com/x/a": "v0.0.0"
}

require: {
	"k8s.io/api":          "v0.24.1" @indirect()
	"k8s.io/apimachinery": "v0.24.1" @indirect()
}

replace: {
	// local replace
	"github.com/x/a": "../a"
}

replace: {
	"k8s.io/api":          "" @import("go")
	"k8s.io/apimachinery": "" @import("go")
}
