package sep

type HTMLTag int

const (
	Unknown = iota
	Paragraph
	Heading2
	Heading3
	BlockQuote
	List
)

type Chapter struct {
	Title      string
	Paragraphs []Block
}

type Bibliography struct {
	List []string
}

type TOCEntry struct {
	Title string
	Href  string
}

type Block struct {
	Type    HTMLTag
	Content string
	List    []string
}

// usage for epub generation
type Book struct {
	Title        string
	PubInfo      string
	Preamble     string
	Authors      []string
	TOC          []TOCEntry
	Chapters     []Chapter
	Bibliography Bibliography
}
