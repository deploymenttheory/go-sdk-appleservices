package usergroups

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

// UserGroupsResponse is the response for a list of user groups.
type UserGroupsResponse struct {
	Data  []UserGroup `json:"data"`
	Links *Links      `json:"links,omitempty"`
	Meta  *Meta       `json:"meta,omitempty"`
}

// UserGroupResponse is the response for a single user group.
type UserGroupResponse struct {
	Data  UserGroup `json:"data"`
	Links *Links    `json:"links,omitempty"`
}

// UserGroup represents a user group resource.
type UserGroup struct {
	ID            string                  `json:"id"`
	Type          string                  `json:"type"`
	Attributes    *UserGroupAttributes    `json:"attributes,omitempty"`
	Relationships *UserGroupRelationships `json:"relationships,omitempty"`
	Links         *ResourceLinks          `json:"links,omitempty"`
}

// UserGroupAttributes contains the attributes of a user group.
type UserGroupAttributes struct {
	OuId             string     `json:"ouId,omitempty"`
	Name             string     `json:"name,omitempty"`
	Type             string     `json:"type,omitempty"`
	TotalMemberCount int        `json:"totalMemberCount,omitempty"`
	CreatedDateTime  *time.Time `json:"createdDateTime,omitempty"`
	UpdatedDateTime  *time.Time `json:"updatedDateTime,omitempty"`
	Status           string     `json:"status,omitempty"`
}

// UserGroupRelationships contains relationship links for a user group.
type UserGroupRelationships struct {
	Users *RelationshipData `json:"users,omitempty"`
}

// RelationshipData holds the links for a relationship.
type RelationshipData struct {
	Links *ResourceLinks `json:"links,omitempty"`
}

// UserGroupUsersLinkagesResponse is the response for user ID linkages of a user group.
type UserGroupUsersLinkagesResponse struct {
	Data  []UserGroupUserLinkage `json:"data"`
	Links *Links                 `json:"links,omitempty"`
	Meta  *Meta                  `json:"meta,omitempty"`
}

// UserGroupUserLinkage represents a user linkage (type + ID only).
type UserGroupUserLinkage struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// RequestQueryOptions represents query parameters for user group endpoints.
type RequestQueryOptions struct {
	// Fields specifies which fields to return. Use Field* constants.
	Fields []string
	// Limit is the number of resources to return (max 1000).
	Limit int
}
