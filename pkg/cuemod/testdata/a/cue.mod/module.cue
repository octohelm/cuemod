module: "github.com/x/a"

require: {
	"k8s.io/api": "v0.20.5"
}

require: {
	"k8s.io/apimachinery": "v0.20.5" @indirect()
}
