package version

import (
	"testing"
	"time"

	. "github.com/octohelm/x/testing"
	"golang.org/x/mod/module"
)

func TestConvert(t *testing.T) {
	mustTime := func(s string) time.Time {
		t, _ := time.Parse(module.PseudoVersionTimestampFormat, s)
		return t
	}

	t.Run("when v0.0.0", func(t *testing.T) {
		v := Convert(
			"v0.0.0-20230809023744-57fc047576ed",
			mustTime("20230809023744"),
			"57fc047576ed",
			false,
		)

		Expect(t, v, Equal("v0.0.0-20230809023744-57fc047576ed"))

		t.Run("when dirty", func(t *testing.T) {
			v := Convert(
				"v0.0.0-20230809023744-57fc047576ed",
				mustTime("20230809023744"),
				"57fc047576ed",
				true,
			)

			Expect(t, v, Equal("v0.0.0-dirty.0.20230809023744-57fc047576ed"))
		})
	})

	t.Run("when vX.Y.Z", func(t *testing.T) {
		v := Convert(
			"v1.1.0",
			mustTime("20230809023744"),
			"57fc047576ed",
			false,
		)

		Expect(t, v, Equal("v1.1.0-20230809023744-57fc047576ed"))

		t.Run("when dirty", func(t *testing.T) {
			v := Convert(
				"v1.1.0",
				mustTime("20230809023744"),
				"57fc047576ed",
				true,
			)

			Expect(t, v, Equal("v1.1.0-dirty.0.20230809023744-57fc047576ed"))
		})
	})

	t.Run("when vX.Y.(Z+1)-0", func(t *testing.T) {
		v := Convert(
			"v1.1.1-0.20230809023744-57fc047576ed",
			mustTime("20230809023744"),
			"57fc047576ed",
			false,
		)

		Expect(t, v, Equal("v1.1.1-0.20230809023744-57fc047576ed"))

		t.Run("when dirty", func(t *testing.T) {
			v := Convert(
				"v1.1.1-0.20230809023744-57fc047576ed",
				mustTime("20230809023744"),
				"57fc047576ed",
				true,
			)

			Expect(t, v, Equal("v1.1.0-dirty.0.20230809023744-57fc047576ed"))
		})
	})
}
