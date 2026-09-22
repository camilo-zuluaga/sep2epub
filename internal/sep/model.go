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

type Block struct {
	Type    HTMLTag
	Content string
	List    []string
}

// usage for epub generation
type Book struct {
	Title    string
	Chapters []Chapter
}
