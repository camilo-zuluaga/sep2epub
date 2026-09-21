package sep

type Chapter struct {
	Title      string
	Paragraphs []Block
}

type Block struct {
	Type    string
	Content string
	List    []string
}
