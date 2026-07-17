package organizationalunits

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
	Self    string `json:"self,omitempty"`
	Related string `json:"related,omitempty"`
}

// OrganizationalUnitsResponse is the response for a list of organizational units.
type OrganizationalUnitsResponse struct {
	Data  []OrganizationalUnit `json:"data"`
	Links *Links               `json:"links,omitempty"`
	Meta  *Meta                `json:"meta,omitempty"`
}

// OrganizationalUnitResponse is the response for a single organizational unit.
type OrganizationalUnitResponse struct {
	Data  OrganizationalUnit `json:"data"`
	Links *Links             `json:"links,omitempty"`
}

// OrganizationalUnit represents an organizational unit resource.
type OrganizationalUnit struct {
	ID            string                           `json:"id"`
	Type          string                           `json:"type"`
	Attributes    *OrganizationalUnitAttributes    `json:"attributes,omitempty"`
	Relationships *OrganizationalUnitRelationships `json:"relationships,omitempty"`
	Links         *ResourceLinks                   `json:"links,omitempty"`
}

// OrganizationalUnitAttributes contains the attributes of an organizational unit.
type OrganizationalUnitAttributes struct {
	Name            string     `json:"name,omitempty"`
	Description     string     `json:"description,omitempty"`
	CreatedDateTime *time.Time `json:"createdDateTime,omitempty"`
	UpdatedDateTime *time.Time `json:"updatedDateTime,omitempty"`
}

// OrganizationalUnitRelationships contains relationship links for an organizational unit.
type OrganizationalUnitRelationships struct {
	Users *RelationshipData `json:"users,omitempty"`
}

// RelationshipData holds the links for a relationship.
type RelationshipData struct {
	Links *ResourceLinks `json:"links,omitempty"`
}

// OrganizationalUnitUsersLinkagesResponse is the response for user ID linkages of an organizational unit.
type OrganizationalUnitUsersLinkagesResponse struct {
	Data  []OrganizationalUnitUserLinkage `json:"data"`
	Links *Links                          `json:"links,omitempty"`
	Meta  *Meta                           `json:"meta,omitempty"`
}

// OrganizationalUnitUserLinkage represents a user linkage (type + ID only).
type OrganizationalUnitUserLinkage struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// RequestQueryOptions represents query parameters for organizational unit endpoints.
type RequestQueryOptions struct {
	// Fields specifies which fields to return. Use Field* constants.
	Fields []string
	// Limit is the number of resources to return (max 1000).
	Limit int
}
