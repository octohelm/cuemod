package manifest

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
	releasev1alpha "github.com/octohelm/cuemod/pkg/api/release/v1alpha"
)

func TestExtract(t *testing.T) {
	data, _ := os.ReadFile("../testdata/web.json")
	var r releasev1alpha.Release
	_ = json.Unmarshal(data, &r)

	t.Run("Extract", func(t *testing.T) {
		spew.Dump(Process(r, nil))
	})
}
