package apps

import "time"

// AppsResponse represents the complete API response
type AppsResponse struct {
	Apps      []App  `json:"apps"`
	Generated string `json:"generated"`
}

// App represents a Microsoft application for Mac
type App struct {
	AppPath        string      `json:"app_path"`
	BundleID       string      `json:"bundle_id"`
	ComponentCount int         `json:"component_count,omitempty"`
	Components     []Component `json:"components,omitempty"`
	Detected       string      `json:"detected"`
	DirectURL      string      `json:"direct_url"`
	DownloadURL    string      `json:"download_url"`
	ETag           string      `json:"etag"`
	InstallKB      *int        `json:"install_kb,omitempty"`
	LastModified   string      `json:"last_modified"`
	Name           string      `json:"name"`
	NumFiles       *int        `json:"num_files,omitempty"`
	PackageID      string      `json:"package_identifier"`
	SHA256         string      `json:"sha256"`
	SizeBytes      int64       `json:"size_bytes"`
	SizeMB         float64     `json:"size_mb"`
	Type           string      `json:"type"`
	Version        string      `json:"version"`
}

// Component represents a component bundled with an application
type Component struct {
	AppPath   *string `json:"app_path"`
	BundleID  *string `json:"bundle_id"`
	Name      string  `json:"name"`
	PackageID string  `json:"package_identifier"`
	Version   string  `json:"version"`
}

// ParseDetectedTime parses the detected timestamp
func (a *App) ParseDetectedTime() (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05.999999", a.Detected)
}

// ParseLastModifiedTime parses the last modified timestamp
func (a *App) ParseLastModifiedTime() (time.Time, error) {
	return time.Parse(time.RFC1123, a.LastModified)
}

// ParseGeneratedTime parses the generated timestamp from the response
func (r *AppsResponse) ParseGeneratedTime() (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", r.Generated)
}
