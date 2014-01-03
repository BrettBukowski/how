package page

import (
	"github.com/franela/goblin"
	"testing"
)

func Test(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("link#NormalizeResultUrl", func() {
		g.It("Should extract random URL parameters", func() {
			bananas := Page{"blah&music=song&a=b&SE3e=3-we_Rr", "title"}
			bananas.NormalizeResultUrl()
			g.Assert(bananas.Url).Equal("blah?answertab=votes")
		})

		g.It("Should extract url wrapper", func() {
			bananas := Page{"/url?q=blah", "title"}
			bananas.NormalizeResultUrl()
			g.Assert(bananas.Url).Equal("blah?answertab=votes")
		})

		g.It("Adds answertab param", func() {
			bananas := Page{"blah", "title"}
			bananas.NormalizeResultUrl()
			g.Assert(bananas.Url).Equal("blah?answertab=votes")
		})
	})
}
