package main

import (
	"github.com/franela/goblin"
	"strings"
	"testing"
)

func Test(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("#searchUrl", func() {
		g.It("Should build a http link", func() {
			l := searchLink([]string{"how", "do", "i"}, false)
			g.Assert(l.url).Equal("http://google.com/search?q=site:stackoverflow.com%20how%20do%20i")
			g.Assert(l.title).Equal("how do i")
		})

		g.It("Should build a https link", func() {
			l := searchLink([]string{"how", "do", "i"}, true)
			g.Assert(l.url).Equal("https://encrypted.google.com/search?q=site:stackoverflow.com%20how%20do%20i")
			g.Assert(l.title).Equal("how do i")
		})
	})

	g.Describe("#linkHeading", func() {
		g.It("Should add a border around the link url", func() {
			l := link{"url", "title"}
			heading := linkHeading(l)
			g.Assert(heading).Equal("********\nurl\n********\n\n")
		})
	})

	g.Describe("link#NormalizeResultUrl", func() {
		g.It("Should extract random URL parameters", func() {
			bananas := link{"blah&music=song&a=b&SE3e=3-we_Rr", "title"}
			bananas.NormalizeResultUrl()
			g.Assert(bananas.url).Equal("blah?answertab=votes")
		})

		g.It("Should extract url wrapper", func() {
			bananas := link{"/url?q=blah", "title"}
			bananas.NormalizeResultUrl()
			g.Assert(bananas.url).Equal("blah?answertab=votes")
		})

		g.It("Adds answertab param", func() {
			bananas := link{"blah", "title"}
			bananas.NormalizeResultUrl()
			g.Assert(bananas.url).Equal("blah?answertab=votes")
		})
	})

	g.Describe("link#FetchPage", func() {
		g.It("Should make a successful http request to get the page", func() {
			l := searchLink([]string{"how", "do", "i"}, false)
			g.Assert(l.FetchPage().Find("body").Length()).Equal(1)
		})
	})

	g.Describe("#fetchFormattedResult", func() {
		g.It("Should grab the results from the page", func() {
			l := link{
				"http://stackoverflow.com/questions/11810218/how-to-set-and-get-fields-in-golang-structs?answertab=votes",
				"how to set and get fields in golang structs?",
			}
			result := make(chan string, 1)
			go fetchFormattedResult(l, result)
			resultText := <-result
			g.Assert(strings.Contains(resultText, linkHeading(l))).Equal(true)
			g.Assert(strings.Count(resultText, "\n") > 6).Equal(true)
		})
	})
}
