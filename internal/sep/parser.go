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

type HTMLTag int

const (
	Unknown = iota
	Paragraph
	Heading2
	Heading3
	BlockQuote
	List
)

func Content(doc *goquery.Document) []Chapter {
	var chapters []Chapter
	var currChapter *Chapter

	doc.Find(MainTextTag).ChildrenFiltered("h2, h3, p, blockquote, ul").Each(func(_ int, element *goquery.Selection) {
		tagName, err := parseHTMLTag(goquery.NodeName(element))
		if err != nil {
			return
		}
		text := normalizeWhitespace(strings.TrimSpace(element.Text()))

		if text == "" {
			return
		}

		switch tagName {
		case Heading2:
			chapters = append(chapters, Chapter{
				Title: text,
			})

			currChapter = &chapters[len(chapters)-1]
		case Heading3:
			currChapter.Paragraphs = append(currChapter.Paragraphs, Block{Type: tagName, Content: text})
		case Paragraph:
			currChapter.Paragraphs = append(currChapter.Paragraphs, Block{Type: tagName, Content: text})
		case BlockQuote:
			currChapter.Paragraphs = append(currChapter.Paragraphs, Block{Type: tagName, Content: text})
		case List:
			currBlock := Block{Type: List}

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

func parseHTMLTag(s string) (HTMLTag, error) {
	switch strings.ToLower(s) {
	case "h2":
		return Heading2, nil
	case "h3":
		return Heading3, nil
	case "p":
		return Paragraph, nil
	case "blockquote":
		return BlockQuote, nil
	case "ul":
		return List, nil
	default:
		return Unknown, fmt.Errorf("error parsing HTML tag: %q", s)
	}
}
