package cve_history

import (
	"bytes"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// cvePattern matches CVE identifiers of the form CVE-YYYY-NNNNN.
var cvePattern = regexp.MustCompile(`CVE-\d{4}-\d{4,}`)

// parseCVEHistory parses the HTML body of the Microsoft Learn Office for Mac
// release notes page and extracts CVE information per release.
//
// The page uses <h2> headings for release dates and nested <h3>/<li> elements
// for security update sections. CVE IDs appear in anchor text and href values.
func parseCVEHistory(body []byte) ([]CVEEntry, error) {
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var entries []CVEEntry
	var current *CVEEntry

	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "h2":
				// Each <h2> marks a new release date section.
				text := strings.TrimSpace(extractText(n))
				if text != "" && !strings.EqualFold(text, "release history") {
					if current != nil {
						entries = append(entries, *current)
					}
					current = &CVEEntry{ReleaseDate: text}
				}

			case "h3":
				// <h3> may contain the version string.
				if current != nil {
					text := strings.TrimSpace(extractText(n))
					if strings.HasPrefix(text, "Version") || strings.HasPrefix(text, "16.") {
						current.Version = text
					}
				}

			case "a":
				// Collect CVE IDs from anchor text and href values.
				if current != nil {
					for _, attr := range n.Attr {
						if attr.Key == "href" {
							for _, cve := range cvePattern.FindAllString(attr.Val, -1) {
								current.CVEs = appendUnique(current.CVEs, cve)
							}
						}
					}
					text := strings.TrimSpace(extractText(n))
					for _, cve := range cvePattern.FindAllString(text, -1) {
						current.CVEs = appendUnique(current.CVEs, cve)
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(doc)

	// Flush the last entry.
	if current != nil {
		entries = append(entries, *current)
	}

	return entries, nil
}

// extractText recursively extracts plain text from an HTML node.
func extractText(n *html.Node) string {
	if n == nil {
		return ""
	}
	var sb strings.Builder
	var walk func(*html.Node)
	walk = func(node *html.Node) {
		if node.Type == html.TextNode {
			sb.WriteString(node.Data)
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(n)
	return sb.String()
}

// appendUnique adds s to slice if it is not already present.
func appendUnique(slice []string, s string) []string {
	for _, existing := range slice {
		if existing == s {
			return slice
		}
	}
	return append(slice, s)
}
