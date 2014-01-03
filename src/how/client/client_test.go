package client

import (
	"github.com/franela/goblin"
	"how/page"
	"strings"
	"testing"
)

func Test(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("#searchUrl", func() {
		g.It("Should build a http link", func() {
			l := searchPage([]string{"how", "do", "i"}, false)
			g.Assert(l.Url).Equal("http://google.com/search?q=site:stackoverflow.com%20how%20do%20i")
			g.Assert(l.Title).Equal("how do i")
		})

		g.It("Should build a https link", func() {
			l := searchPage([]string{"how", "do", "i"}, true)
			g.Assert(l.Url).Equal("https://encrypted.google.com/search?q=site:stackoverflow.com%20how%20do%20i")
			g.Assert(l.Title).Equal("how do i")
		})
	})

	g.Describe("#linkHeading", func() {
		g.It("Should add a border around the link url", func() {
			l := page.Page{"url", "title"}
			heading := linkHeading(l)
			g.Assert(heading).Equal("********\nurl\n********\n\n")
		})
	})

	g.Describe("link#Fetch", func() {
		g.It("Should make a successful http request to get the page", func() {
			l := searchPage([]string{"how", "do", "i"}, false)
			doc, err := l.Fetch()
			g.Assert(err).Equal(nil)
			g.Assert(doc.Find("body").Length()).Equal(1)
		})
	})

	g.Describe("#fetchFormattedResult", func() {
		g.It("Should grab the results from the page", func() {
			l := page.Page{
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
