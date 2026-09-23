package epub

import (
	"fmt"
	"html"
	"regexp"
	"sep2epub/internal/sep"
	"strings"

	"github.com/go-shiori/go-epub"
)

func Generate(book sep.Book) error {
	e, err := epub.NewEpub(book.Title)
	if err != nil {
		return fmt.Errorf("create EPUB: %w", err)
	}

	e.SetAuthor(strings.Join(book.Authors, ", "))

	_, err = e.AddSection(preambleHTML(book), "Preamble", "", "")
	if err != nil {
		return fmt.Errorf("add chapter %q: %w", "err", err)
	}

	_, err = e.AddSection(tocHTML(book.TOC), "Table of Contents", "", "")
	if err != nil {
		return fmt.Errorf("add chapter %q: %w", "err", err)
	}

	// TODO: Map book.TOC links to the generated chapter and subsection filenames.

	for i, chapter := range book.Chapters {
		body := chapterHTML(chapter)

		_, err := e.AddSection(body, chapter.Title, fmt.Sprintf("chapter-%02d.xhtml", i+1), "")
		if err != nil {
			return fmt.Errorf("add chapter %q: %w", chapter.Title, err)
		}
	}

	_, err = e.AddSection(bibHTML(book.Bibliography), "Bibliography", "", "")
	if err != nil {
		return fmt.Errorf("add chapter %q: %w", "err", err)
	}

	// TODO: Render book.Bibliography as a final section.

	if err := e.Write(fmt.Sprintf("%s.epub", book.Title)); err != nil {
		return fmt.Errorf("err writing EPUB: %w", err)
	}

	return nil
}

func chapterHTML(chapter sep.Chapter) string {
	var builder strings.Builder

	builder.WriteString("<h1>")
	builder.WriteString(html.EscapeString(chapter.Title))
	builder.WriteString("</h1>\n")

	for _, block := range chapter.Paragraphs {
		switch block.Type {
		case sep.Heading2:
			builder.WriteString("<h2>")
			builder.WriteString(html.EscapeString(block.Content))
			builder.WriteString("</h2>\n")
		case sep.Heading3:
			builder.WriteString("<h3>")
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

func tocHTML(toc []sep.TOCEntry) string {
	var builder strings.Builder

	builder.WriteString("<h1>Table of Contents</h1>\n")
	builder.WriteString("<ul>\n")

	inSub := false
	for _, entry := range toc {
		isSub, _ := regexp.MatchString(`^[0-9]+\.[0-9]+`, entry.Title)

		if isSub && !inSub {
			builder.WriteString("<ul>\n")
			inSub = true
		} else if !isSub && inSub {
			builder.WriteString("</ul>\n")
			inSub = false
		}

		builder.WriteString("<li>")
		builder.WriteString(html.EscapeString(entry.Title))
		builder.WriteString("</li>\n")
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
	builder.WriteString(b.Preamble)
	builder.WriteString("</p>\n")

	return builder.String()
}
