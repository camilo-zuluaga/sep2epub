package main

import (
	"context"
	"fmt"
	"log"
	"sep2epub/internal/fetcher"
	"sep2epub/internal/sep"

	"github.com/PuerkitoBio/goquery"
)

func debugPrint(chapters []sep.Chapter, chapterNum int) {
	fmt.Println()
	chapter := chapters[chapterNum+1]
	fmt.Println(chapter.Title)
	fmt.Println()
	for _, p := range chapter.Paragraphs {
		fmt.Println(p)
		fmt.Println()
	}
}

func main() {
	client := fetcher.New()
	body, err := client.Fetch(context.Background(), "https://plato.stanford.edu/entries/consciousness/")
	if err != nil {
		log.Fatal(err)
	}
	defer body.Close()

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		log.Fatal(err)
	}

	sep.TableOfContents(doc, true)
	sep.Content(doc)
}
