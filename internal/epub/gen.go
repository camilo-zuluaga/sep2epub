package epub

import (
	"fmt"
	"html"
	"regexp"
	"sep2epub/internal/sep"
	"strings"
	"unicode"

	"github.com/go-shiori/go-epub"
)

type section struct {
	title, body, filename string
}

func Generate(book sep.Book) error {
	e, err := epub.NewEpub(book.Title)
	if err != nil {
		return fmt.Errorf("create EPUB: %w", err)
	}
	e.SetAuthor(strings.Join(book.Authors, ", "))

	for _, s := range sections(book) {
		if _, err := e.AddSection(s.body, s.title, s.filename, ""); err != nil {
			return fmt.Errorf("add section %q: %w", s.title, err)
		}
	}
	if err := e.Write(fmt.Sprintf("%s.epub", book.Title)); err != nil {
		return fmt.Errorf("err writing EPUB: %w", err)
	}
	return nil
}

func sections(book sep.Book) []section {
	s := []section{
		{title: "Preamble", body: preambleHTML(book), filename: "preamble.xhtml"},
		{title: "Table of Contents", body: tocHTML(book.TOC), filename: "toc.xhtml"},
	}

	for i, chapter := range book.Chapters {
		s = append(s, section{
			title:    chapter.Title,
			body:     chapterHTML(chapter),
			filename: chapterFilename(i),
		})
	}

	return append(s, section{title: "Bibliography", body: bibHTML(book.Bibliography)})
}

func chapterHTML(chapter sep.Chapter) string {
	var builder strings.Builder

	builder.WriteString("<h1>")
	builder.WriteString(html.EscapeString(chapter.Title))
	builder.WriteString("</h1>\n")

	for _, block := range chapter.Paragraphs {
		id := slugify(block.Content)

		switch block.Type {
		case sep.Heading2:
			builder.WriteString("<h2>")
			builder.WriteString(html.EscapeString(block.Content))
			builder.WriteString("</h2>\n")
		case sep.Heading3:
			builder.WriteString(`<h3 id="`)
			builder.WriteString(id)
			builder.WriteString(`">`)
			builder.WriteString(html.EscapeString(block.Content))
			builder.WriteString("</h3>\n")
		case sep.Paragraph:
			builder.WriteString("<p>")
			builder.WriteString(block.Content)
			builder.WriteString("</p>\n")
		case sep.BlockQuote:
			builder.WriteString("<blockquote>")
			builder.WriteString(block.Content)
			builder.WriteString("</blockquote>\n")
		case sep.List:
			builder.WriteString("<ul>")
			for _, item := range block.List {
				builder.WriteString("<li>")
				builder.WriteString(item)
				builder.WriteString("</li>\n")
			}
			builder.WriteString("</ul>\n")
		}
	}
	return builder.String()
}

var subEntry = regexp.MustCompile(`^[0-9]+\.[0-9]+`)

func tocHTML(toc []sep.TOCEntry) string {
	var builder strings.Builder

	builder.WriteString("<h1>Table of Contents</h1>\n")
	builder.WriteString("<ul>\n")

	chapterIndex := -1
	inSub := false
	for _, entry := range toc {
		var href string
		isSub := subEntry.MatchString(entry.Title)

		if isSub && !inSub {
			builder.WriteString("<ul>\n")
			inSub = true
		} else if !isSub && inSub {
			builder.WriteString("</ul>\n")
			inSub = false
		}

		if inSub {
			href = chapterFilename(chapterIndex) + "#" + slugify(entry.Title)
		} else {
			chapterIndex++
			href = chapterFilename(chapterIndex)
		}

		builder.WriteString(`<li><a href="`)
		builder.WriteString(html.EscapeString(href))
		builder.WriteString(`" style="text-decoration:none;">`)
		builder.WriteString(html.EscapeString(entry.Title))
		builder.WriteString("</a></li>\n")
	}

	if inSub {
		builder.WriteString("</ul>\n")
	}

	builder.WriteString("</ul>\n")

	return builder.String()
}

func bibHTML(bib sep.Bibliography) string {
	var builder strings.Builder

	builder.WriteString("<h1>Bibliography</h1>\n")
	builder.WriteString("<ul>")
	for _, ref := range bib.List {
		builder.WriteString("<li>")
		builder.WriteString(ref)
		builder.WriteString("</li>\n")
	}

	builder.WriteString("</ul>\n")

	return builder.String()
}

func preambleHTML(b sep.Book) string {
	var builder strings.Builder

	builder.WriteString("<h1>")
	builder.WriteString(b.Title)
	builder.WriteString("</h1>\n")

	builder.WriteString("<p>")
	builder.WriteString(b.PubInfo)
	builder.WriteString("</p>\n")

	builder.WriteString("<p>")
	builder.WriteString(b.Preamble)
	builder.WriteString("</p>\n")

	return builder.String()
}

func chapterFilename(idx int) string {
	return fmt.Sprintf("chapter-%02d.xhtml", idx+1)
}

func slugify(value string) string {
	value = strings.ToLower(value)
	value = strings.TrimSpace(value)

	var builder strings.Builder
	previousWasSeparator := false

	for _, r := range value {
		switch {
		case unicode.IsLetter(r), unicode.IsDigit(r):
			builder.WriteRune(r)
			previousWasSeparator = false

		case !previousWasSeparator:
			builder.WriteRune('-')
			previousWasSeparator = true
		}
	}

	return strings.Trim(builder.String(), "-")
}
