package main

import (
	"code.google.com/p/go.net/html"
	"flag"
	"fmt"
	"io/ioutil"
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

func getResult(url string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", userAgents[rand.Intn(len(userAgents))])

	resp, err := client.Do(req)
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		return string(body), err
	}
	return "An error occurred", err
}

func fetchPage(url string) string {
	page, err := getResult(url)

	if err != nil {
		fmt.Println("Error retrieving ", url)
		os.Exit(2)
	}
	return page
}

func normalizeLink(link string) string {
	r, _ := regexp.Compile("&[A-Za-z0-9]+=[A-Za-z0-9-_]+")
	link = r.ReplaceAllString(link, "")
	r, _ = regexp.Compile("/url\\?q=")
	link = r.ReplaceAllString(link, "")
	return fmt.Sprint(link, "?answertab=votes")
}

func extractLinks(numberToExtract int, htmlDoc string) []string {
	doc, err := html.Parse(strings.NewReader(htmlDoc))
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	var links []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					matched, _ := regexp.MatchString("^(/url\\?q=)?http://(meta.)?stackoverflow.com", a.Val)
					if matched {
						links = append(links, normalizeLink(a.Val))
						break
					}
				}
			}
		}
		for c := n.FirstChild; c != nil && len(links) < numberToExtract; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return links
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
	var f func(n *html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "div" {
			for _, a := range n.Attr {
				if a.Key == "class" && strings.Contains(a.Val, "answer") {
					fmt.Println(n.Data)
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	var page string
	var doc *html.Node
	var err error
	for _, url := range links {
		fmt.Println(url)
		page = fetchPage(url)
		doc, err = html.Parse(strings.NewReader(page))

		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		f(doc)
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
	links := extractLinks(*numAnswers, page)

	if *onlyLinks {
		for i := 0; i < len(links); i++ {
			fmt.Println(links[i])
		}
		os.Exit(0)
	}

	getInstructions(links)
}
