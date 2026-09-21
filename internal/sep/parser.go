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

func Content(doc *goquery.Document) []Chapter {
	var chapters []Chapter
	var currChapter *Chapter

	doc.Find(MainTextTag).ChildrenFiltered("h2, p, blockquote, ul").Each(func(_ int, element *goquery.Selection) {
		tagName := goquery.NodeName(element)
		text := normalizeWhitespace(strings.TrimSpace(element.Text()))

		if text == "" {
			return
		}

		switch tagName {
		case "h2":
			chapters = append(chapters, Chapter{
				Title: text,
			})

			currChapter = &chapters[len(chapters)-1]
		case "h3":
			currChapter.Paragraphs = append(currChapter.Paragraphs, Block{Type: tagName, Content: text})
		case "p":
			currChapter.Paragraphs = append(currChapter.Paragraphs, Block{Type: tagName, Content: text})
		case "blockquote":
			currChapter.Paragraphs = append(currChapter.Paragraphs, Block{Type: tagName, Content: text})
		case "ul":
			currBlock := Block{Type: "list"}

			element.ChildrenFiltered("li").Each(func(i int, s *goquery.Selection) {
				text := normalizeWhitespace(strings.TrimSpace(s.Text()))
				currBlock.List = append(currBlock.List, text)
			})

			currChapter.Paragraphs = append(currChapter.Paragraphs, currBlock)
		}
	})
	return chapters
}

func TableOfContents(doc *goquery.Document, showSubsections bool) {
	main_title := doc.Find(MainTitleTag).Text()
	fmt.Printf("%s\n\n", main_title)

	doc.Find(TableOfContentsTag).Each(func(i int, link *goquery.Selection) {
		title := normalizeWhitespace(link.Text())
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

func normalizeWhitespace(value string) string {
	return strings.Join(strings.Fields(value), " ")
}
