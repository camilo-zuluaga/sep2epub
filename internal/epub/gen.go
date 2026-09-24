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

	imagePath, err := e.AddImage("internal/epub/logo/sep-man-blackshadow.svg", "sep-logo.svg")
	if err != nil {
		return fmt.Errorf("add cover image: %w", err)
	}

	for _, s := range sections(book, imagePath) {
		if _, err := e.AddSection(s.body, s.title, s.filename, ""); err != nil {
			return fmt.Errorf("add section %q: %w", s.title, err)
		}
	}
	if err := e.Write(fmt.Sprintf("%s.epub", book.Title)); err != nil {
		return fmt.Errorf("err writing EPUB: %w", err)
	}
	return nil
}

func sections(book sep.Book, imagePath string) []section {
	s := []section{
		{title: "", body: coverHTML(book, imagePath), filename: "cover.xhtml"},
		{title: "Preamble", body: preambleHTML(book), filename: "preamble.xhtml"},
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
		builder.WriteString(entry.Title)
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

	toc := tocHTML(b.TOC)
	builder.WriteString(toc)

	return builder.String()
}

func coverHTML(book sep.Book, imagePath string) string {
	var builder strings.Builder

	title := html.EscapeString(book.Title)
	sourceURL := html.EscapeString(book.MetaInfo.URL)
	authors := html.EscapeString(strings.Join(book.Authors, ", "))

	builder.WriteString(`
	<div style="
		max-width: 38em;
		margin: 0 auto;
		padding: 2em 1.5em;
		text-align: center;
		font-family: serif;
		line-height: 1.4;
	">

		<p style="font-size: 0.85em; margin-bottom: 0.3em;">
			EPUB version of the entry
		</p>
	`)

	builder.WriteString(`
		<p style="
			font-size: 1.1em;
			margin-top: 0;
			margin-bottom: 0.5em;
		">`)
	builder.WriteString(title)
	builder.WriteString("</p>\n")

	if book.MetaInfo.URL != "" {
		builder.WriteString(`
			<p style="
				font-size: 0.75em;
				margin: 0.4em 0;
			">
				<a href="`)
		builder.WriteString(sourceURL)
		builder.WriteString(`">`)
		builder.WriteString(sourceURL)
		builder.WriteString(`</a>
			</p>
		`)
	}

	if book.MetaInfo.Edition != "" {
		builder.WriteString(`
			<p style="
				font-size: 0.85em;
				margin: 0.5em 0 1.8em;
			">
				from the `)
		builder.WriteString(html.EscapeString(book.MetaInfo.Edition))
		builder.WriteString(` edition of the
			</p>
		`)
	}

	builder.WriteString(`
		<h1 style="
			font-size: 2em;
			font-weight: normal;
			font-variant: small-caps;
			letter-spacing: 0.06em;
			line-height: 1.3;
			margin: 0;
		">
			Stanford Encyclopedia<br/>
			of Philosophy
		</h1>
	`)

	builder.WriteString(`
		<div style="
			height: 10em;
			margin: 1.5em auto;
			display: flex;
			align-items: center;
			justify-content: center;
		">
			<img src="`)
	builder.WriteString(html.EscapeString(imagePath))
	builder.WriteString(`"
				alt="Stanford Encyclopedia of Philosophy logo"
				style="
					display: block;
					max-width: 7em;
					max-height: 9em;
					width: auto;
					height: auto;
					margin: 0 auto;
				"
			/>
		</div>
	`)

	builder.WriteString(`
		<div style="
			font-size: 0.72em;
			line-height: 1.45;
			margin: 0 auto 2em;
		">
			<p>
				Co-Principal Editors:
				Edward N. Zalta and Uri Nodelman<br/>

				Associate Editors:
				Colin Allen, Hannah Kim, and Paul Oppenheimer<br/>

				Faculty Sponsors:
				R. Lanier Anderson and Thomas Icard<br/>

				Editorial Board:
				<a href="https://plato.stanford.edu/board.html">
					https://plato.stanford.edu/board.html
				</a><br/>

				Library of Congress ISSN: 1095-5054
			</p>
		</div>
	`)

	builder.WriteString(`
		<div style="
			font-size: 0.72em;
			font-style: italic;
			line-height: 1.4;
			margin: 2em 0;
		">
			<p>
				Stanford Encyclopedia of Philosophy<br/>
				Copyright © 2026 by the publisher<br/>
				The Metaphysics Research Lab<br/>
				Department of Philosophy<br/>
				Stanford University, Stanford, CA 94305
			</p>
		</div>
	`)

	builder.WriteString(`
		<div style="
			font-size: 0.8em;
			line-height: 1.45;
			margin-top: 2em;
		">
			<p>
				<span style="font-size: 1.15em;">`)
	builder.WriteString(title)
	builder.WriteString("</span><br/>\n")

	if book.MetaInfo.Year > 0 {
		builder.WriteString("Copyright © ")
		builder.WriteString(fmt.Sprintf("%d", book.MetaInfo.Year))

		if len(book.Authors) == 1 {
			builder.WriteString(" by the author<br/>\n")
		} else {
			builder.WriteString(" by the authors<br/>\n")
		}
	}

	if authors != "" {
		builder.WriteString(authors)
		builder.WriteString("<br/>\n")
	}

	builder.WriteString(`
				All rights reserved.
			</p>

			<p>
				Copyright policy:
				<a href="https://plato.stanford.edu/info/copyright/">
					https://plato.stanford.edu/info/copyright/
				</a>
			</p>
		</div>
	</div>
	`)

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
