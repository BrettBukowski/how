package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	// "net/http"
)

func getResult() {

}

func searchUrl(https *bool, query string) string {
	var pre string
	if *https {
		pre = "https://encrypted."
	} else {
		pre = "http://"
	}

	body := "google.com/search?q=site:stackoverflow.com%20"
	result := fmt.Sprint(pre, body, query)

	return result
}

func main() {
	https := flag.Bool("https", false, "Use https")
	numAnswers := flag.Int("answers", 1, "Number of answers to retrieve")
	// onlyLinks := flag.Bool("links", false, "Only display answer links, not the result text")

	flag.Parse()

	if len(flag.Args()) == 0 {
		flag.PrintDefaults()
		os.Exit(2)
	}

	var url string = searchUrl(https, strings.Join(flag.Args(), "%20"))

	fmt.Println("url: ", url)

	fmt.Println("answers: ", *numAnswers)
}
