package update_history

// UpdateHistoryResponse holds the full Office for Mac update history.
type UpdateHistoryResponse struct {
	Entries []UpdateHistoryEntry
}

// UpdateHistoryEntry represents a single row from the Office for Mac update history table
// on the Microsoft Learn documentation page.
type UpdateHistoryEntry struct {
	// ReleaseDate is the release date string as scraped from the table (e.g. "April 15, 2025").
	ReleaseDate string

	// Version is the Office version string (e.g. "16.108").
	Version string

	// BusinessProSuiteDownload is the download URL for the full BusinessPro suite installer.
	BusinessProSuiteDownload string

	// SuiteDownload is the download URL for the standard Office suite installer.
	SuiteDownload string

	// WordUpdate is the download URL for the Word individual updater.
	WordUpdate string

	// ExcelUpdate is the download URL for the Excel individual updater.
	ExcelUpdate string

	// PowerPointUpdate is the download URL for the PowerPoint individual updater.
	PowerPointUpdate string

	// OutlookUpdate is the download URL for the Outlook individual updater.
	OutlookUpdate string

	// OneNoteUpdate is the download URL for the OneNote individual updater.
	OneNoteUpdate string

	// Archived indicates whether the download links are no longer available.
	Archived bool
}
