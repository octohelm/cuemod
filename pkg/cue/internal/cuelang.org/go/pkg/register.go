// Copyright 2020 CUE Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pkg

import (
	_ "github.com/octohelm/cuemod/pkg/cue/internal/cuelang.org/go/pkg/crypto/md5"
	_ "github.com/octohelm/cuemod/pkg/cue/internal/cuelang.org/go/pkg/crypto/sha1"
	_ "github.com/octohelm/cuemod/pkg/cue/internal/cuelang.org/go/pkg/crypto/sha256"
	_ "github.com/octohelm/cuemod/pkg/cue/internal/cuelang.org/go/pkg/crypto/sha512"
	_ "github.com/octohelm/cuemod/pkg/cue/internal/cuelang.org/go/pkg/encoding/base64"
	_ "github.com/octohelm/cuemod/pkg/cue/internal/cuelang.org/go/pkg/encoding/csv"
	_ "github.com/octohelm/cuemod/pkg/cue/internal/cuelang.org/go/pkg/encoding/hex"
	_ "github.com/octohelm/cuemod/pkg/cue/internal/cuelang.org/go/pkg/encoding/json"
	_ "github.com/octohelm/cuemod/pkg/cue/internal/cuelang.org/go/pkg/encoding/yaml"
	_ "github.com/octohelm/cuemod/pkg/cue/internal/cuelang.org/go/pkg/html"

	_ "cuelang.org/go/pkg/tool"
	_ "github.com/octohelm/cuemod/pkg/cue/internal/cuelang.org/go/pkg/list"
	_ "github.com/octohelm/cuemod/pkg/cue/internal/cuelang.org/go/pkg/math"
	_ "github.com/octohelm/cuemod/pkg/cue/internal/cuelang.org/go/pkg/math/bits"
	_ "github.com/octohelm/cuemod/pkg/cue/internal/cuelang.org/go/pkg/net"
	_ "github.com/octohelm/cuemod/pkg/cue/internal/cuelang.org/go/pkg/path"
	_ "github.com/octohelm/cuemod/pkg/cue/internal/cuelang.org/go/pkg/regexp"
	_ "github.com/octohelm/cuemod/pkg/cue/internal/cuelang.org/go/pkg/strconv"
	_ "github.com/octohelm/cuemod/pkg/cue/internal/cuelang.org/go/pkg/strings"
	_ "github.com/octohelm/cuemod/pkg/cue/internal/cuelang.org/go/pkg/struct"
	_ "github.com/octohelm/cuemod/pkg/cue/internal/cuelang.org/go/pkg/text/tabwriter"
	_ "github.com/octohelm/cuemod/pkg/cue/internal/cuelang.org/go/pkg/text/template"
	_ "github.com/octohelm/cuemod/pkg/cue/internal/cuelang.org/go/pkg/time"
	_ "github.com/octohelm/cuemod/pkg/cue/internal/cuelang.org/go/pkg/tool/cli"
	_ "github.com/octohelm/cuemod/pkg/cue/internal/cuelang.org/go/pkg/tool/exec"
	_ "github.com/octohelm/cuemod/pkg/cue/internal/cuelang.org/go/pkg/tool/file"
	_ "github.com/octohelm/cuemod/pkg/cue/internal/cuelang.org/go/pkg/tool/http"
	_ "github.com/octohelm/cuemod/pkg/cue/internal/cuelang.org/go/pkg/tool/os"
)
