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
func (l *link) FetchPage() (*goquery.Document, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", l.url, nil)
	req.Header.Add("User-Agent", userAgents[rand.Intn(len(userAgents))])

	var resp *http.Response
	var e error

	if resp, e = client.Do(req); e != nil {
		return nil, e
	}

	defer resp.Body.Close()

	var doc *goquery.Document

	if doc, e = goquery.NewDocumentFromReader(resp.Body); e != nil {
		panic(e.Error())
	}
	return doc, nil
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
func searchLink(query []string, https bool) link {
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
	results := make(chan string, numberOfResults)

	for _, link := range links {
		go fetchFormattedResult(link, results)
	}

	for i := 0; i < numberOfResults; i++ {
		select {
		case msg := <-results:
			fmt.Println(msg)
		}
	}
}

// For every <a> in the selection, replace with a MD-style link.
// Stuffs the results into the supplied channel.
func fetchFormattedResult(l link, results chan<- string) {
	doc, _ := l.FetchPage()
	selection := doc.Find(".answer").First().Find(".post-text")
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
	results <- linkHeading(l) + text
}

// Header for every link
func linkHeading(l link) string {
	border := strings.Repeat("*", len(l.url)+5)
	return fmt.Sprintf("%s\n%s\n%s\n\n", border, l.url, border)
}

func main() {
	https := flag.Bool("https", false, "Use https")
	numAnswers := flag.Int("answers", 1, "Number of answers to retrieve")
	onlyLinks := flag.Bool("links", false, "Only display answer links, not the result text")

	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		flag.PrintDefaults()
		fmt.Println("Supply a question")
		os.Exit(2)
	}

	link := searchLink(args, *https)
	page, e := link.FetchPage()
	if e != nil {
		fmt.Println("Unable to fetch results due to connection problem.")
		os.Exit(-1)
	}
	links := extractLinks(page, *numAnswers)

	if *onlyLinks {
		for i, link := range links {
			fmt.Printf("%d. [%s](%s)\n", i+1, link.title, link.url)
		}
		os.Exit(0)
	}

	printInstructions(links)
}
