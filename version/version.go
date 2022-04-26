package version

const (
	DevelopmentVersion = "devel"
)

var (
	Version  = DevelopmentVersion
	Revision = "-"
)

func FullVersion() string {
	if len([]byte(Revision)) > 7 {
		return Version + ".sha+" + Revision[0:7]
	}
	return Version + ".sha+" + Revision
}
