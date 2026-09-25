package cmd

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/spf13/cobra"
)

var outputPath string

var rootCmd = &cobra.Command{
	Use:   "sep2epub <url>",
	Short: "Convert Standford Enclyclopedia entries to EPUB",
	Args:  cobra.ExactArgs(1),
	RunE:  runConversion,
}

func init() {
	rootCmd.Flags().StringVarP(&outputPath, "output", "o", "", "output EPUB")
}

func runConversion(cmd *cobra.Command, args []string) error {
	url, err := validateSEPURL(args[0])
	if err != nil {
		return err
	}

	cmd.Printf("URL: %s\n", url)
	cmd.Printf("Output: %s\n", outputPath)

	return nil
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

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
