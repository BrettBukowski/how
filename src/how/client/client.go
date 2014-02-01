package client

import (
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"how/page"
	"how/version"
	"os"
	"regexp"
	"strings"
)

// Gleam result links out of Google's search result page.
func extractPages(doc *goquery.Document, numberToExtract int) []page.Page {
	var links []page.Page

	doc.Find("#res a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")

		if acceptResultUrl(href) {
			resultPage := page.Page{href, s.Text()}
			resultPage.NormalizeResultUrl()
			links = append(links, resultPage)
		}
	})

	numberOfLinks := len(links)

	if numberToExtract > numberOfLinks {
		numberToExtract = numberOfLinks - 1
	}

	if numberToExtract <= 0 {
		return links
	}

	return links[0:numberToExtract]
}

// Determines whether the given url is a
// legit Stackoverflow answer link.
func acceptResultUrl(url string) bool {
	matched, _ := regexp.MatchString("^(/url\\?q=)?http://(meta.)?stackoverflow.com", url)
	return matched && !strings.Contains(url, "stackoverflow.com/questions/tagged/") && !strings.Contains(url, "stackoverflow.com/tags")
}

// Build the Google search URL.
func searchPage(query []string, https bool) page.Page {
	body := "google.com/search?q=site:stackoverflow.com"

	var pre string
	if https {
		pre = "https://encrypted."
	} else {
		pre = "http://"
	}

	return page.Page{
		fmt.Sprint(pre, body, "%20", strings.Join(query, "%20")),
		strings.Join(query, " "),
	}
}

// Extract the top answer out of the
// Stackoverflow page.
func printInstructions(links []page.Page) {
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
func fetchFormattedResult(l page.Page, results chan<- string) {
	doc, _ := l.Fetch()
	answer := doc.Find(".answer")

	if answer.Nodes == nil {
		results <- fmt.Sprintf("No answers given for <%s>", l.Url)
		return
	}

	selection := answer.First().Find(".post-text")
	text := selection.Text()

	selection.Find("a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		linkText := s.Text()

		var mdPage = []string{
			"[",
			linkText,
			"](",
			href,
			")",
		}

		if strings.EqualFold(linkText, href) {
			mdPage = []string{
				"<", linkText, ">",
			}
		}

		text = strings.Replace(text, linkText, strings.Join(mdPage, ""), 1)
	})
	results <- linkHeading(l) + text
}

// Header for every link
func linkHeading(l page.Page) string {
	border := strings.Repeat("*", len(l.Url)+5)
	return fmt.Sprintf("%s\n%s\n%s\n\n", border, l.Url, border)
}

func Main() {
	https := flag.Bool("https", false, "Use https")
	numAnswers := flag.Int("answers", 1, "Number of answers to retrieve")
	onlyPages := flag.Bool("links", false, "Only display answer links, not the result text")
	showVersion := flag.Bool("version", false, "Display the version number and exit")
	update := flag.Bool("update", false, "Update to the latest version")

	flag.Parse()
	args := flag.Args()

	if *showVersion {
		fmt.Printf("%1.1f\n", version.Version)
		os.Exit(0)
	}

	if *update {
		if version.NewerVersionAvailable() {
			version.Update()
		} else {
			fmt.Println("You're on the latest version.")
		}
		os.Exit(0)
	}

	if len(args) == 0 {
		flag.PrintDefaults()
		fmt.Println("Supply a question")
		os.Exit(2)
	}

	link := searchPage(args, *https)
	page, e := link.Fetch()
	if e != nil {
		fmt.Println("Unable to fetch results due to connection problem.")
		os.Exit(-1)
	}
	links := extractPages(page, *numAnswers)

	if *onlyPages {
		for i, link := range links {
			fmt.Printf("%d. [%s](%s)\n", i+1, link.Title, link.Url)
		}
		os.Exit(0)
	}

	printInstructions(links)
}
