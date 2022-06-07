package modfile

import (
	"fmt"
	"strings"
)

func ParsePathMayWithVersion(v string) (*PathMayWithVersion, error) {
	if len(v) == 0 {
		return nil, fmt.Errorf("invalid %s", v)
	}

	parts := strings.Split(v, "@")

	i := parts[0]

	if i != "" && i[0] == '.' {
		return &PathMayWithVersion{Path: i}, nil
	}

	if len(parts) > 1 {
		vv := strings.Split(parts[1], "#")
		if len(vv) > 1 {
			return &PathMayWithVersion{Path: i, Version: vv[0], VcsVersion: vv[1]}, nil
		}
		return &PathMayWithVersion{Path: i, Version: vv[0]}, nil
	}
	return &PathMayWithVersion{Path: i}, nil

}

type PathMayWithVersion struct {
	Path       string
	Version    string
	VcsVersion string
}

func (r *PathMayWithVersion) UnmarshalText(text []byte) error {
	rp, err := ParsePathMayWithVersion(string(text))
	if err != nil {
		return err
	}
	*r = *rp
	return nil
}

func (r PathMayWithVersion) MarshalText() (text []byte, err error) {
	return []byte(r.String()), nil
}

func (r PathMayWithVersion) IsLocalReplace() bool {
	return len(r.Path) > 0 && r.Path[0] == '.'
}

func (r PathMayWithVersion) String() string {
	if r.IsLocalReplace() {
		return r.Path
	}

	b := strings.Builder{}
	b.WriteString(r.Path)

	if r.Version != "" || r.VcsVersion != "" {
		b.WriteString("@")
		b.WriteString(r.Version)

		if r.VcsVersion != "" {
			b.WriteString("#")
			b.WriteString(r.VcsVersion)
		}
	}
	return b.String()
}
