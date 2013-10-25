package main

import (
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var userAgents = []string{
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.7; rv:11.0) Gecko/20100101 Firefox/11.0",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:22.0) Gecko/20100 101 Firefox/22.0",
	"Mozilla/5.0 (Windows NT 6.1; rv:11.0) Gecko/20100101 Firefox/11.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_4) AppleWebKit/536.5 (KHTML, like Gecko) Chrome/19.0.1084.46 Safari/536.5",
	"Mozilla/5.0 (Windows; Windows NT 6.1) AppleWebKit/536.5 (KHTML, like Gecko) Chrome/19.0.1084.46 Safari/536.5",
}

// URL and Title of each
// result link
type link struct {
	url   string
	title string
}

// Get the document for the link.
// Requests the page with a randomized user-agent
// so Google doesn't get suspicious. ಠ_ಠ
func (l *link) FetchPage() *goquery.Document {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", l.url, nil)
	req.Header.Add("User-Agent", userAgents[rand.Intn(len(userAgents))])
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	var doc *goquery.Document
	var e error

	if doc, e = goquery.NewDocumentFromReader(resp.Body); e != nil {
		panic(e.Error())
	}
	return doc
}

// Google wraps its redirect tracker around non-https
// result links. Remove that. And make sure Stackoverflow's
// answer tab is focused.
func (l *link) NormalizeResultUrl() {
	r, _ := regexp.Compile("&[A-Za-z0-9]+=[A-Za-z0-9-_]+")
	l.url = r.ReplaceAllString(l.url, "")
	r, _ = regexp.Compile("/url\\?q=")
	l.url = r.ReplaceAllString(l.url, "")
	l.url = fmt.Sprint(l.url, "?answertab=votes")
}

// Gleam result links out of Google's search result page.
func extractLinks(doc *goquery.Document, numberToExtract int) []link {
	var links []link

	doc.Find("#res a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		matched, _ := regexp.MatchString("^(/url\\?q=)?http://(meta.)?stackoverflow.com", href)
		if matched {
			resultLink := link{href, s.Text()}
			resultLink.NormalizeResultUrl()
			links = append(links, resultLink)
		}
	})

	return links[0:numberToExtract]
}

// Build the Google search URL.
func searchLink(https bool, query []string) link {
	body := "google.com/search?q=site:stackoverflow.com"

	var pre string
	if https {
		pre = "https://encrypted."
	} else {
		pre = "http://"
	}

	return link{
		fmt.Sprint(pre, body, "%20", strings.Join(query, "%20")),
		strings.Join(query, " "),
	}
}

// Extract the top answer out of the
// Stackoverflow page.
func printInstructions(links []link) {
	numberOfResults := len(links)

	for i, link := range links {
		border := strings.Repeat("*", len(link.url)+5)
		if numberOfResults > 1 {
			fmt.Printf("%s\n%d. %s\n%s\n\n", border, i+1, link.url, border)
		} else {
			fmt.Printf("%s\n%s\n%s\n\n", border, link.url, border)
		}

		text := convertLinksToMarkdown(link.FetchPage().Find(".answer").First().Find(".post-text"))

		fmt.Println(text)
	}
}

// For every <a> in the selection, replace with a MD-style link.
func convertLinksToMarkdown(selection *goquery.Selection) string {
	text := selection.Text()

	selection.Find("a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		linkText := s.Text()

		var mdLink = []string{
			"[",
			linkText,
			"](",
			href,
			")",
		}

		if strings.EqualFold(linkText, href) {
			mdLink = []string{
				"<", linkText, ">",
			}
		}

		text = strings.Replace(text, linkText, strings.Join(mdLink, ""), 1)
	})

	return text
}

func main() {
	https := flag.Bool("https", false, "Use https")
	numAnswers := flag.Int("answers", 1, "Number of answers to retrieve")
	onlyLinks := flag.Bool("links", false, "Only display answer links, not the result text")

	flag.Parse()

	if len(flag.Args()) == 0 {
		flag.PrintDefaults()
		fmt.Println("Supply a question")
		os.Exit(2)
	}

	link := searchLink(*https, flag.Args())
	page := link.FetchPage()
	links := extractLinks(page, *numAnswers)

	if *onlyLinks {
		for i, link := range links {
			fmt.Printf("%d. [%s](%s)\n", i+1, link.title, link.url)
		}
		os.Exit(0)
	}

	printInstructions(links)
}
