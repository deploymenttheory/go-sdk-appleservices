package apps

// Shared pagination types

type Meta struct {
	Paging *Paging `json:"paging,omitempty"`
}

type Paging struct {
	Total      int    `json:"total,omitempty"`
	Limit      int    `json:"limit,omitempty"`
	NextCursor string `json:"nextCursor,omitempty"`
}

type Links struct {
	Self  string `json:"self,omitempty"`
	First string `json:"first,omitempty"`
	Next  string `json:"next,omitempty"`
	Prev  string `json:"prev,omitempty"`
	Last  string `json:"last,omitempty"`
}

type ResourceLinks struct {
	Self string `json:"self,omitempty"`
}

// AppsResponse is the response for a list of apps.
type AppsResponse struct {
	Data  []App  `json:"data"`
	Links *Links `json:"links,omitempty"`
	Meta  *Meta  `json:"meta,omitempty"`
}

// AppResponse is the response for a single app.
type AppResponse struct {
	Data  App    `json:"data"`
	Links *Links `json:"links,omitempty"`
}

// App represents an app resource.
type App struct {
	ID         string         `json:"id"`
	Type       string         `json:"type"`
	Attributes *AppAttributes `json:"attributes,omitempty"`
	Links      *ResourceLinks `json:"links,omitempty"`
}

// AppAttributes contains the attributes of an app.
type AppAttributes struct {
	Name        string   `json:"name,omitempty"`
	BundleId    string   `json:"bundleId,omitempty"`
	WebsiteUrl  string   `json:"websiteUrl,omitempty"`
	Version     string   `json:"version,omitempty"`
	SupportedOS []string `json:"supportedOS,omitempty"`
	IsCustomApp bool     `json:"isCustomApp,omitempty"`
	AppStoreUrl string   `json:"appStoreUrl,omitempty"`
}

// RequestQueryOptions represents query parameters for app endpoints.
type RequestQueryOptions struct {
	// Fields specifies which fields to return. Use Field* constants.
	Fields []string
	// Limit is the number of resources to return (max 1000).
	Limit int
}
