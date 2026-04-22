package update_history

import (
	"bytes"
	"strings"

	"golang.org/x/net/html"
)

// parseUpdateHistory parses the HTML body of the Microsoft Learn Office for Mac
// update history page and extracts the update table rows.
//
// The page contains an HTML table with columns:
// Release date | Version | Install package | Update packages (one per app)
func parseUpdateHistory(body []byte) ([]UpdateHistoryEntry, error) {
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var entries []UpdateHistoryEntry
	var parseNode func(*html.Node)

	parseNode = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "tr" {
			entry, ok := parseTableRow(n)
			if ok {
				entries = append(entries, entry)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			parseNode(c)
		}
	}
	parseNode(doc)

	return entries, nil
}

// parseTableRow extracts an UpdateHistoryEntry from a <tr> element.
// Returns false if the row does not contain update data (e.g. header rows).
func parseTableRow(tr *html.Node) (UpdateHistoryEntry, bool) {
	cells := collectTDCells(tr)
	if len(cells) < 3 {
		return UpdateHistoryEntry{}, false
	}

	dateText := extractText(cells[0])
	versionText := extractText(cells[1])

	// Skip header rows that contain non-date text in the first column.
	if dateText == "" || strings.EqualFold(dateText, "release date") {
		return UpdateHistoryEntry{}, false
	}

	entry := UpdateHistoryEntry{
		ReleaseDate: strings.TrimSpace(dateText),
		Version:     strings.TrimSpace(versionText),
	}

	// Column 2: install package links (suite downloads)
	if len(cells) > 2 {
		links := collectLinks(cells[2])
		if len(links) > 0 {
			entry.BusinessProSuiteDownload = links[0]
		}
		if len(links) > 1 {
			entry.SuiteDownload = links[1]
		}
	}

	// Columns 3+: individual app updater links
	appUpdates := []*string{
		&entry.WordUpdate,
		&entry.ExcelUpdate,
		&entry.PowerPointUpdate,
		&entry.OutlookUpdate,
		&entry.OneNoteUpdate,
	}
	for i, ptr := range appUpdates {
		col := i + 3
		if col < len(cells) {
			links := collectLinks(cells[col])
			if len(links) > 0 {
				*ptr = links[0]
			}
		}
	}

	// Mark as archived if no download links found.
	if entry.BusinessProSuiteDownload == "" && entry.SuiteDownload == "" {
		entry.Archived = true
	}

	return entry, true
}

// collectTDCells returns the direct <td> and <th> child elements of a <tr> node.
func collectTDCells(tr *html.Node) []*html.Node {
	var cells []*html.Node
	for c := tr.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && (c.Data == "td" || c.Data == "th") {
			cells = append(cells, c)
		}
	}
	return cells
}

// extractText recursively extracts the plain text content of an HTML node.
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

// collectLinks returns the href values of all <a> elements within a node.
func collectLinks(n *html.Node) []string {
	var links []string
	var walk func(*html.Node)
	walk = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "a" {
			for _, attr := range node.Attr {
				if attr.Key == "href" && attr.Val != "" {
					links = append(links, attr.Val)
				}
			}
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(n)
	return links
}
