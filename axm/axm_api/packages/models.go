package packages

import "time"

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

// PackagesResponse is the response for a list of packages.
type PackagesResponse struct {
	Data  []Package `json:"data"`
	Links *Links    `json:"links,omitempty"`
	Meta  *Meta     `json:"meta,omitempty"`
}

// PackageResponse is the response for a single package.
type PackageResponse struct {
	Data  Package `json:"data"`
	Links *Links  `json:"links,omitempty"`
}

// Package represents a package resource.
type Package struct {
	ID         string             `json:"id"`
	Type       string             `json:"type"`
	Attributes *PackageAttributes `json:"attributes,omitempty"`
	Links      *ResourceLinks     `json:"links,omitempty"`
}

// PackageAttributes contains the attributes of a package.
type PackageAttributes struct {
	Name            string     `json:"name,omitempty"`
	URL             string     `json:"url,omitempty"`
	Hash            string     `json:"hash,omitempty"`
	BundleIds       []string   `json:"bundleIds,omitempty"`
	Description     string     `json:"description,omitempty"`
	Version         string     `json:"version,omitempty"`
	CreatedDateTime *time.Time `json:"createdDateTime,omitempty"`
	UpdatedDateTime *time.Time `json:"updatedDateTime,omitempty"`
}

// RequestQueryOptions represents query parameters for package endpoints.
type RequestQueryOptions struct {
	// Fields specifies which fields to return. Use Field* constants.
	Fields []string
	// Limit is the number of resources to return (max 1000).
	Limit int
}
