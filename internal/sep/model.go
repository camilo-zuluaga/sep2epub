package sep

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
