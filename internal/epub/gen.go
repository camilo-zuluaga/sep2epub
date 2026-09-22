package epub

import (
	"fmt"
	"html"
	"log"
	"sep2epub/internal/sep"
	"strings"

	"github.com/go-shiori/go-epub"
)

func Generate(book sep.Book) error {
	e, err := epub.NewEpub(book.Title)
	if err != nil {
		log.Println(err)
	}

	e.SetAuthor("Hingle McCringleberry")

	for i, chapter := range book.Chapters {
		body := chapterHTML(chapter)

		_, err := e.AddSection(body, chapter.Title, fmt.Sprintf("chapter-%02d.xhtml", i+1), "")
		if err != nil {
			log.Println(err)
		}

	}

	if err := e.Write("CONSCIOUSNESS.epub"); err != nil {
		return fmt.Errorf("write EPUB: %w", err)
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
			builder.WriteString(html.EscapeString(block.Content))
			builder.WriteString("</p>\n")
		case sep.BlockQuote:
			builder.WriteString("<blockquote>")
			builder.WriteString(html.EscapeString(block.Content))
			builder.WriteString("</blockquote>\n")
		case sep.List:
			builder.WriteString("<ul>")
			for _, item := range block.List {
				builder.WriteString("<li>")
				builder.WriteString(html.EscapeString(item))
				builder.WriteString("</li>\n")
			}
			builder.WriteString("</ul>\n")
		}
	}
	return builder.String()
}
