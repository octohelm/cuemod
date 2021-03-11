package modfile

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestLoadModFile(t *testing.T) {
	mod := ModFile{}

	_, err := LoadModFile("../testdata/b", &mod)

	NewWithT(t).Expect(err).To(BeNil())
	NewWithT(t).Expect(mod.Module).To(Equal("github.com/x/b"))

	t.Log(mod.String())
}
