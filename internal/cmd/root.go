package cmd

import (
	"context"
	"fmt"
	"net/url"
	"path/filepath"
	"sep2epub/internal/epub"
	"sep2epub/internal/sep"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cobra"
)

var outputPath string

var rootCmd = &cobra.Command{
	Use:           "sep2epub",
	Short:         "Convert Standford Enclyclopedia entries to EPUB",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE:          runConvert,
}

func runConvert(cmd *cobra.Command, args []string) error {
	url, err := validateSEPURL(args[0])
	if err != nil {
		return err
	}

	output := outputPath
	if output == "" {
		output = defaultOutputPath(url)
	}

	book, err := loadBook(cmd.Context(), url.String())
	if err != nil {
		return err
	}

	if err := epub.Generate(book, output); err != nil {
		return err
	}

	cmd.Printf("Created: %s\n", output)
	return nil
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().StringVarP(&outputPath, "output", "o", "", "output EPUB")
}

func validateSEPURL(rawURL string) (*url.URL, error) {
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	if parsedURL.Scheme != "https" {
		return nil, fmt.Errorf("must be HTTPS: %w", err)
	}

	if !strings.EqualFold(
		parsedURL.Hostname(),
		"plato.stanford.edu",
	) {
		return nil, fmt.Errorf("must belong to plato.stanford.edu")
	}

	return parsedURL, nil
}

func loadBook(ctx context.Context, url string) (sep.Book, error) {
	client := sep.NewFetcher()
	body, err := client.Fetch(ctx, url)
	if err != nil {
		return sep.Book{}, fmt.Errorf("load book: %w", err)
	}
	defer body.Close()

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return sep.Book{}, fmt.Errorf("parse entry: %w", err)
	}

	return sep.Book{
		Title:        sep.GetTitle(doc),
		MetaInfo:     sep.GetCitation(ctx, url),
		PubInfo:      sep.GetPubInfo(doc),
		Preamble:     sep.GetPreamble(doc),
		Authors:      sep.GetAuthors(doc),
		TOC:          sep.TableOfContents(doc, true),
		Chapters:     sep.Content(doc),
		Bibliography: sep.GetBibliography(doc),
	}, nil
}

func defaultOutputPath(sourceURL *url.URL) string {
	path := strings.Trim(sourceURL.Path, "/")
	entryName := filepath.Base(path)

	if entryName == "." || entryName == "/" || entryName == "" {
		return "book.epub"
	}

	return entryName + ".epub"
}
