// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package modfetch

import (
	"github.com/octohelm/cuemod/pkg/modutil/internal/cmd/go/internals/cfg"

	"golang.org/x/mod/module"
)

// allowInsecure reports whether we are allowed to fetch this path in an insecure manner.
func allowInsecure(path string) bool {
	return cfg.Insecure || module.MatchPrefixPatterns(cfg.GOINSECURE, path)
}
