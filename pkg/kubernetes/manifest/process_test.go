package manifest

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
	releasev1alpha1 "github.com/octohelm/cuemod/pkg/apis/release/v1alpha1"
)

func TestExtract(t *testing.T) {
	data, _ := os.ReadFile("../testdata/web.json")
	var r releasev1alpha1.Release
	_ = json.Unmarshal(data, &r)

	t.Run("Extract", func(t *testing.T) {
		spew.Dump(Process(r, nil))
	})
}
