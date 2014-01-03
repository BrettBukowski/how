package version

import (
	"github.com/franela/goblin"
	"testing"
)

func Test(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("Version#Check", func() {
		g.It("Doesn't panic", func() {
			result := NewerVersionAvailable()
			// Always on the latest version.
			g.Assert(result).Equal(false)
		})
	})

	g.Describe("Version#NewestVersion", func() {
		g.It("Is the same as the Version", func() {
			result := NewestVersion()
			// Always on the latest version.
			g.Assert(result).Equal(Version)
		})
	})

	g.Describe("Version#Update", func() {
		g.It("Returns false when there's no update available", func() {
			result := Update()
			g.Assert(result).Equal(false)
		})
	})
}
