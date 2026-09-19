package main

import (
	"log"
	"net/http"
	"sep2epub/internal/sep"

	"github.com/PuerkitoBio/goquery"
)

func ExampleScrape() {
	res, err := http.Get("https://plato.stanford.edu/entries/consciousness/")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	sep.TableOfContents(doc, true)
	sep.Content(doc)
}

func main() {
	ExampleScrape()
}
