package version

import _ "embed"

//go:embed version
var version string

var Version = "v" + version
