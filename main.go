package main

import (
	"fmt"
	"log"
	"net/http"
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

	// article, err := doc.Find("#preamble, #toc ul li").Html()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//

	doc.Find("#toc a").Each(func(i int, link *goquery.Selection) {
		title := normalizeWhitespace(link.Text())
		// href, _ := link.Attr("href")

		fmt.Printf("[%02d] title=%q \n", i, title)
	})
}

func normalizeWhitespace(value string) string {
	return strings.Join(strings.Fields(value), " ")
}

func main() {
	ExampleScrape()
}
