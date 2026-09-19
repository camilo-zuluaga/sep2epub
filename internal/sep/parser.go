package sep

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	MainTextTag        = "#main-text"
	MainTitleTag       = "#aueditable h1"
	TableOfContentsTag = "#toc a"
)

type Chapter struct {
	Title      string
	Paragraphs []string
}

func Content(doc *goquery.Document) {
	var chapters []Chapter
	var currChapter *Chapter

	doc.Find(MainTextTag).ChildrenFiltered("h2, p").Each(func(_ int, element *goquery.Selection) {
		tagName := goquery.NodeName(element)
		text := strings.TrimSpace(element.Text())

		if text == "" {
			return
		}

		switch tagName {
		case "h2":
			chapters = append(chapters, Chapter{
				Title: text,
			})

			currChapter = &chapters[len(chapters)-1]
		case "p":
			currChapter.Paragraphs = append(currChapter.Paragraphs, text)
		}
	})
	fmt.Println()
	chapter := chapters[0]
	fmt.Println(chapter.Title)
	fmt.Println()
	for _, p := range chapter.Paragraphs {
		fmt.Println(p)
		fmt.Println()
	}
}

func TableOfContents(doc *goquery.Document, showSubsections bool) {
	main_title := doc.Find(MainTitleTag).Text()
	fmt.Printf("%s\n\n", main_title)

	doc.Find(TableOfContentsTag).Each(func(i int, link *goquery.Selection) {
		title := normalizeWhitespaceAndQuotes(link.Text())
		if match, _ := regexp.MatchString("[0-9]+.[0-9]+", title); match {
			if !showSubsections {
				return
			}
			fmt.Printf("	%s \n", title)
		} else {
			fmt.Printf("%s \n", title)
		}
	})
}

func normalizeWhitespaceAndQuotes(value string) string {
	return strings.Join(strings.Fields(value), " ")
}
