package main

import (
	"context"
	"fmt"
	"log"
	"sep2epub/internal/epub"
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
	book, err := loadBook(context.Background(), "https://plato.stanford.edu/entries/consciousness/")
	if err != nil {
		log.Fatal(err)
	}
	if err := epub.Generate(book); err != nil {
		log.Fatal(err)
	}
}

// loadBook fetches an entry once and assembles the data needed by the EPUB generator.
func loadBook(ctx context.Context, url string) (sep.Book, error) {
	client := fetcher.New()
	body, err := client.Fetch(ctx, url)
	if err != nil {
		return sep.Book{}, fmt.Errorf("load book: %w", err)
	}
	defer body.Close()

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return sep.Book{}, fmt.Errorf("parse entry: %w", err)
	}

	return sep.Book{
		Title:        sep.GetTitle(doc),
		Authors:      sep.GetAuthors(doc),
		TOC:          sep.TableOfContents(doc, true),
		Chapters:     sep.Content(doc),
		Bibliography: sep.GetBibliography(doc),
	}, nil
}
