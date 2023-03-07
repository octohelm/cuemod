package version

const (
	DevelopmentVersion = "devel"
)

var (
	version = DevelopmentVersion
)

func FullVersion() string {
	return version
}
