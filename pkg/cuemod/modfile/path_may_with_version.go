package modfile

import (
	"fmt"
	"strings"
)

func ParsePathMayWithVersion(v string) (*VersionedPathIdentity, error) {
	if len(v) == 0 {
		return nil, fmt.Errorf("invalid %s", v)
	}

	parts := strings.Split(v, "@")

	i := parts[0]

	if i != "" && i[0] == '.' {
		return &VersionedPathIdentity{Path: i}, nil
	}

	if len(parts) > 1 {
		vv := strings.Split(parts[1], "#")
		if len(vv) > 1 {
			return &VersionedPathIdentity{Path: i, ModVersion: ModVersion{
				Version: vv[0], VcsRef: vv[1],
			}}, nil
		}
		return &VersionedPathIdentity{Path: i, ModVersion: ModVersion{
			Version: vv[0],
		}}, nil
	}
	return &VersionedPathIdentity{Path: i}, nil

}

type VersionedPathIdentity struct {
	Path string
	ModVersion
}

func (r *VersionedPathIdentity) UnmarshalText(text []byte) error {
	rp, err := ParsePathMayWithVersion(string(text))
	if err != nil {
		return err
	}
	*r = *rp
	return nil
}

func (r VersionedPathIdentity) MarshalText() (text []byte, err error) {
	return []byte(r.String()), nil
}

func (r VersionedPathIdentity) IsLocalReplace() bool {
	return len(r.Path) > 0 && r.Path[0] == '.'
}

func (r VersionedPathIdentity) String() string {
	if r.IsLocalReplace() {
		return r.Path
	}

	b := strings.Builder{}
	b.WriteString(r.Path)

	if r.Version != "" || r.VcsRef != "" {
		b.WriteString("@")
		b.WriteString(r.Version)

		if r.VcsRef != "" {
			b.WriteString("#")
			b.WriteString(r.VcsRef)
		}
	}
	return b.String()
}
