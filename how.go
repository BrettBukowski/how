package main

import (
	"code.google.com/p/go.net/html"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
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

func extractLinks(numberToExtract *int, htmlDoc string) []string {
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
					links = append(links, a.Val)
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
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

func main() {
	https := flag.Bool("https", false, "Use https")
	numAnswers := flag.Int("answers", 1, "Number of answers to retrieve")
	// onlyLinks := flag.Bool("links", false, "Only display answer links, not the result text")

	flag.Parse()

	if len(flag.Args()) == 0 {
		flag.PrintDefaults()
		fmt.Println("Supply a question")
		os.Exit(2)
	}

	var url string = searchUrl(https, flag.Args())

	fmt.Println("url: ", url)
	fmt.Println("answers: ", *numAnswers)

	body, err := getResult(url)

	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	// fmt.Println(body)

	links := extractLinks(numAnswers, body)

	for i := 0; i < len(links); i++ {
		fmt.Println(links[i])
	}
}
