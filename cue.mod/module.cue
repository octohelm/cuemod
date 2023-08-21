module: "github.com/octohelm/cuemod"

require: {
	"github.com/innoai-tech/runtime": "v0.0.0-20230807071635-a566ade1c374"
	"wagon.octohelm.tech":            "v0.0.0-20200202235959-8fa253acacb2"
}

replace: {
	"k8s.io/api":          "" @import("go")
	"k8s.io/apimachinery": "" @import("go")
}
