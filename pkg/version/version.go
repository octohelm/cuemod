package version

import (
	"fmt"
	"time"

	"golang.org/x/mod/module"
	"golang.org/x/mod/semver"
)

func Convert(version string, t time.Time, rev string, dirty bool) string {
	exact := true
	base, err := module.PseudoVersionBase(version)
	if err == nil {
		version = base
		exact = false
	}
	if version == "" {
		version = "v0.0.0"
		exact = true
	}
	if dirty {
		version += "-dirty"
		exact = false
	}
	return pseudoVersion(version, t, rev, exact)
}

func pseudoVersion(version string, t time.Time, rev string, exact bool) string {
	major := semver.Major(version)
	if major == "" {
		major = "v0"
	}

	if exact {
		build := semver.Build(version)
		segment := fmt.Sprintf("%s-%s", t.UTC().Format(module.PseudoVersionTimestampFormat), rev)
		version = semver.Canonical(version)
		if version == "" {
			version = major + ".0.0"
		}
		return version + "-" + segment + build
	}

	return module.PseudoVersion(
		major,
		version,
		t,
		rev,
	)
}
