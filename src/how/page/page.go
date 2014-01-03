package page

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"how/request"
	"math/rand"
	"net/http"
	"regexp"
)

var userAgents = []string{
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.7; rv:11.0) Gecko/20100101 Firefox/11.0",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:22.0) Gecko/20100 101 Firefox/22.0",
	"Mozilla/5.0 (Windows NT 6.1; rv:11.0) Gecko/20100101 Firefox/11.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_4) AppleWebKit/536.5 (KHTML, like Gecko) Chrome/19.0.1084.46 Safari/536.5",
	"Mozilla/5.0 (Windows; Windows NT 6.1) AppleWebKit/536.5 (KHTML, like Gecko) Chrome/19.0.1084.46 Safari/536.5",
}

// URL and Title of each result link
type Page struct {
	Url   string
	Title string
}

// Get the document for the link.
// Requests the page with a randomized user-agent
// so Google doesn't get suspicious. ಠ_ಠ
func (p *Page) Fetch() (*goquery.Document, error) {
	var response *http.Response
	var doc *goquery.Document
	var e error

	if response, e = request.Get(p.Url, map[string]string{
		"User-Agent": userAgents[rand.Intn(len(userAgents))],
	}); e != nil {
		return nil, e
	}

	if doc, e = goquery.NewDocumentFromResponse(response); e != nil {
		panic(e.Error())
	}
	return doc, nil
}

// Google wraps its redirect tracker around non-https
// result links. Remove that. And make sure Stackoverflow's
// answer tab is focused.
func (p *Page) NormalizeResultUrl() {
	r, _ := regexp.Compile("&[A-Za-z0-9]+=[A-Za-z0-9-_]+")
	p.Url = r.ReplaceAllString(p.Url, "")
	r, _ = regexp.Compile("/url\\?q=")
	p.Url = r.ReplaceAllString(p.Url, "")
	p.Url = fmt.Sprint(p.Url, "?answertab=votes")
}
