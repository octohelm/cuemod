module: "github.com/x/examples"

require: {
	"github.com/octohelm/cuem": "v0.0.0-20210520091405-7e9ddaa903c7"
}

require: {
	"k8s.io/api":          "v0.23.5" @indirect()
	"k8s.io/apimachinery": "v0.23.5" @indirect()
}
