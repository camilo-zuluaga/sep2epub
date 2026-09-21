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
	chapter := chapters[chapterNum]
	fmt.Println(chapter.Title)
	fmt.Println()
	for _, p := range chapter.Paragraphs {
		switch p.Type {

		case sep.Heading3:
			fmt.Println(p.Content)
			fmt.Println()
		case sep.Paragraph:
			fmt.Println(p.Content)
			fmt.Println()
		case sep.BlockQuote:
			fmt.Printf("%q \n", p.Content)
			fmt.Println()
		case sep.List:
			for _, l := range p.List {
				fmt.Printf("-  %s \n", l)
				fmt.Println()
			}
		}
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
	l := sep.Content(doc)
	sep.GetBibliography(doc)
	debugPrint(l, 1)
}
