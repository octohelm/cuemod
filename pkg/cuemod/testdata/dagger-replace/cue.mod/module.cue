module: "dagger-example"

require: {
	"dagger.io":          "v0.2.17-0.20220607061721-387f0ef2334f" @vcs("release-main")
	"universe.dagger.io": "v0.2.17-0.20220607061721-387f0ef2334f" @vcs("release-main")
}

replace: {
	"dagger.io":          "github.com/morlay/dagger/pkg/dagger.io@v0.2.17-0.20220607061721-387f0ef2334f"
	"universe.dagger.io": "github.com/morlay/dagger/pkg/universe.dagger.io@v0.2.17-0.20220607061721-387f0ef2334f"
}
