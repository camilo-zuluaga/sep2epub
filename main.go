package main

import (
	"fmt"
	"os"
	"sep2epub/internal/cmd"
	"sep2epub/internal/sep"
)

func debugPrint(chapters []sep.Chapter, chapterNum int) {
	fmt.Println()
	chapter := chapters[chapterNum]
	fmt.Println(chapter.Title)
	fmt.Println()
	for _, p := range chapter.Paragraphs {
		switch p.Type {

		case sep.Heading3:
			fmt.Println(p.Content)
			fmt.Println()
		case sep.Paragraph:
			fmt.Println(p.Content)
			fmt.Println()
		case sep.BlockQuote:
			fmt.Printf("%q \n", p.Content)
			fmt.Println()
		case sep.List:
			for _, l := range p.List {
				fmt.Printf("-  %s \n", l)
				fmt.Println()
			}
		}
	}
}

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}
