module: "github.com/octohelm/cuemod"

require: {
	"github.com/innoai-tech/runtime": "v0.0.0-20230301034018-d0f9cf039113"
	"wagon.octohelm.tech":            "v0.0.0-20200202235959-3d91e2e3161f"
}

replace: {
	"k8s.io/api":          "" @import("go")
	"k8s.io/apimachinery": "" @import("go")
}
