package sep

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/nickng/bibtex"
)

const (
	MainTextTag        = "#main-text"
	MainTitleTag       = "#aueditable h1"
	TableOfContentsTag = "#toc a"
	PreambleTag        = "#preamble"
	Pubinfo            = "#pubinfo"
	CitationInfoURL    = "https://plato.stanford.edu/cgi-bin/encyclopedia/archinfo.cgi?entry=%s"
)

func GetTitle(doc *goquery.Document) string {
	return normalizeWhitespace(doc.Find(MainTitleTag).First().Text())
}

func GetPubInfo(doc *goquery.Document) string {
	htmlText, _ := doc.Find(Pubinfo).First().Html()
	return normalizeWhitespace(htmlText)
}

func GetPreamble(doc *goquery.Document) string {
	htmlText, _ := doc.Find(PreambleTag).First().Html()
	return normalizeWhitespace(htmlText)
}

func Content(doc *goquery.Document) []Chapter {
	var chapters []Chapter
	var currChapter *Chapter

	doc.Find(MainTextTag).ChildrenFiltered("h2, h3, p, blockquote, ul").Each(func(_ int, element *goquery.Selection) {
		tagName, err := parseHTMLTag(goquery.NodeName(element))
		if err != nil {
			return
		}

		switch tagName {
		case Heading2:
			text := normalizeWhitespace(strings.TrimSpace(element.Text()))
			if text == "" {
				return
			}

			chapters = append(chapters, Chapter{
				Title: text,
			})

			currChapter = &chapters[len(chapters)-1]
		case Heading3:
			text := normalizeWhitespace(strings.TrimSpace(element.Text()))
			if text == "" {
				return
			}

			currChapter.Paragraphs = append(currChapter.Paragraphs, Block{Type: tagName, Content: text})
		case Paragraph, BlockQuote:
			content, err := element.Html()
			if err != nil {
				return
			}

			content = normalizeWhitespace(content)

			currChapter.Paragraphs = append(currChapter.Paragraphs, Block{Type: tagName, Content: content})
		case List:
			currBlock := Block{Type: tagName}

			element.ChildrenFiltered("li").Each(func(i int, li *goquery.Selection) {
				content, err := li.Html()
				if err != nil {
					return
				}

				content = normalizeWhitespace(content)
				currBlock.List = append(currBlock.List, content)
			})

			currChapter.Paragraphs = append(currChapter.Paragraphs, currBlock)
		}
	})
	return chapters
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

func GetAuthors(doc *goquery.Document) []string {
	var authors []string

	doc.Find("#article-copyright p").ChildrenFiltered(`a[target="other"]`).Each(func(i int, anchor *goquery.Selection) {
		/*
						 <div id="article-copyright">
			    			<p>
			       				<a href="../../info.html#c">Copyright © 2025</a> by
			           			<br>
			             		<a href="https://elizabethbrake.com" target="other">Elizabeth Brake</a>
			               		<a href="mailto:eebrake%40wisc%2eedu"><em>eebrake<abbr title=" at ">@</abbr>wisc<abbr title=" dot ">.</abbr>edu</em></a>&gt;<br>
			                 	<a href="https://www.st-andrews.ac.uk/philosophy/people/jrm39/" target="other">Joseph Millum</a>
			                  	<a href="mailto:jrm39%40st-andrews%2eac%2euk"><em>jrm39<abbr title=" at ">@</abbr>st-andrews<abbr title=" dot ">.</abbr>ac<abbr title=" dot ">.</abbr>uk</em></a>&gt;
			                </p>
			            </div>
		*/
		author := normalizeWhitespace(anchor.Text())
		if author != "" {
			authors = append(authors, author)
		}
	})

	return authors
}

func GetBibliography(doc *goquery.Document) Bibliography {
	var bib Bibliography

	doc.Find("#bibliography").ChildrenFiltered("ul").Each(func(i int, element *goquery.Selection) {
		element.ChildrenFiltered("li").Each(func(i int, s *goquery.Selection) {
			htmlText, _ := s.Html()
			text := normalizeWhitespace(strings.TrimSpace(htmlText))
			bib.List = append(bib.List, text)
		})
	})

	return bib
}

var subChapter = regexp.MustCompile(`^[0-9]+\.[0-9]+`)

func TableOfContents(doc *goquery.Document, showSubsections bool) []TOCEntry {
	var entries []TOCEntry

	doc.Find(TableOfContentsTag).Each(func(i int, link *goquery.Selection) {
		htmlText, _ := link.Html()
		title := normalizeWhitespace(htmlText)
		if match := subChapter.MatchString(title); match && !showSubsections {
			return
		}
		href, _ := link.Attr("href")
		entries = append(entries, TOCEntry{Title: title, Href: href})
	})
	return entries
}

func GetCitation(ctx context.Context, url string) Citation {
	client := NewFetcher()
	fmt.Println(citationURL(url))
	body, err := client.Fetch(ctx, citationURL(url))
	if err != nil {
		log.Fatal(err)
	}
	defer body.Close()

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		log.Fatal(err)
	}

	t, _ := doc.Find("pre").Html()
	c, err := ParseCitation(t)
	if err != nil {
		log.Fatal(err)
	}
	return c
}

func citationURL(url string) string {
	if url[len(url)-1] == '/' {
		url = url[:len(url)-1]
	}
	s := strings.Split(url, "/")
	extractedName := s[len(s)-1]
	return fmt.Sprintf(CitationInfoURL, extractedName)
}

func ParseCitation(source string) (Citation, error) {
	bib, err := bibtex.Parse(strings.NewReader(source))
	if err != nil {
		return Citation{}, fmt.Errorf("parse BibTeX: %w", err)
	}

	if len(bib.Entries) == 0 {
		return Citation{}, fmt.Errorf("citation contains no entries")
	}

	entry := bib.Entries[0]

	year, err := strconv.Atoi(value(entry, "year"))
	if err != nil {
		return Citation{}, fmt.Errorf("invalid year: %w", err)
	}

	return Citation{
		URL:       cleanURL(value(entry, "howpublished")),
		Year:      year,
		Edition:   cleanEdition(value(entry, "edition")),
		Publisher: value(entry, "publisher"),
	}, nil
}

func value(entry *bibtex.BibEntry, key string) string {
	field, exists := entry.Fields[key]
	if !exists {
		return ""
	}

	return field.String()
}

func cleanURL(url string) string {
	return url[5 : len(url)-1]
}

func cleanEdition(e string) string {
	return strings.ReplaceAll(
		strings.ReplaceAll(e, "{", ""),
		"}", "")
}

func normalizeWhitespace(value string) string {
	return strings.Join(strings.Fields(value), " ")
}
