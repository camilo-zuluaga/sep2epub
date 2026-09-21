package sep

type Chapter struct {
	Title      string
	Paragraphs []Block
}

type Block struct {
	Type    HTMLTag
	Content string
	List    []string
}
