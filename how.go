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

func fetchPage(url string) *goquery.Document {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
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

func normalizeLink(link string) string {
	r, _ := regexp.Compile("&[A-Za-z0-9]+=[A-Za-z0-9-_]+")
	link = r.ReplaceAllString(link, "")
	r, _ = regexp.Compile("/url\\?q=")
	link = r.ReplaceAllString(link, "")
	return fmt.Sprint(link, "?answertab=votes")
}

func extractLinks(doc *goquery.Document, numberToExtract int) []string {
	var links []string

	doc.Find("#res a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		matched, _ := regexp.MatchString("^(/url\\?q=)?http://(meta.)?stackoverflow.com", href)
		if matched {
			links = append(links, normalizeLink(href))
		}
	})

	return links[0:numberToExtract]
}

func searchUrl(https *bool, query []string) string {
	body := "google.com/search?q=site:stackoverflow.com"

	var pre string
	if *https {
		pre = "https://encrypted."
	} else {
		pre = "http://"
	}

	return fmt.Sprint(pre, body, "%20", strings.Join(query, "%20"))
}

func getInstructions(links []string) {
	var doc *goquery.Document
	for _, url := range links {
		fmt.Println(url)
		doc = fetchPage(url)
		answer := doc.Find(".answer").First().Find(".post-text")
		text := answer.Text()
		answer.Find("a").Each(func(i int, s *goquery.Selection) {
			href, _ := s.Attr("href")
			linkText := s.Text()

			var mdLink = []string{
				"[",
				linkText,
				"](",
				href,
				")",
			}
			text = strings.Replace(text, linkText, strings.Join(mdLink, ""), 1)
		})

		fmt.Println("")
		fmt.Println(text)
	}
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

	url := searchUrl(https, flag.Args())
	page := fetchPage(url)
	links := extractLinks(page, *numAnswers)

	if *onlyLinks {
		for i := 0; i < len(links); i++ {
			fmt.Println(links[i])
		}
		os.Exit(0)
	}

	getInstructions(links)
}
