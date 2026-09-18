package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ExampleScrape() {
	// Request the HTML page.
	res, err := http.Get("https://plato.stanford.edu/entries/consciousness/")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	TOC(doc, true)
	// article, err := doc.Find("#preamble, #toc ul li").Html()
	// if err != nil {
	// 	log.Fatal(err)
	//
}

func TOC(doc *goquery.Document, showSubsections bool) {
	doc.Find("#toc a").Each(func(i int, link *goquery.Selection) {
		title := normalizeWhitespaceAndQuotes(link.Text())
		// href, _ := link.Attr("href")

		if match, _ := regexp.MatchString("[0-9]+.[0-9]+", title); match {
			if !showSubsections {
				return
			}
			fmt.Printf("	%s \n", title)
		} else {
			fmt.Printf("%s \n", title)
		}
		// fmt.Printf("[%02d] title=%q \n", i, title)
	})
}

func normalizeWhitespaceAndQuotes(value string) string {
	return strings.Join(strings.Fields(value), " ")
}

func main() {
	ExampleScrape()
}
