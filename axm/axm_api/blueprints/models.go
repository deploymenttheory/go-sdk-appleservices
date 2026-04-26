package blueprints

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

// RelationshipLinks holds the self, include, and related links for a relationship object.
type RelationshipLinks struct {
	Self    string `json:"self,omitempty"`
	Include string `json:"include,omitempty"`
	Related string `json:"related,omitempty"`
}

// BlueprintRelationshipLink wraps RelationshipLinks inside the "links" key.
type BlueprintRelationshipLink struct {
	Links *RelationshipLinks `json:"links,omitempty"`
}

// BlueprintResponse is the response for a single Blueprint resource.
type BlueprintResponse struct {
	Data  Blueprint `json:"data"`
	Links *Links    `json:"links,omitempty"`
}

// Blueprint represents a Blueprint resource.
type Blueprint struct {
	ID            string                  `json:"id"`
	Type          string                  `json:"type"`
	Attributes    *BlueprintAttributes    `json:"attributes,omitempty"`
	Relationships *BlueprintRelationships `json:"relationships,omitempty"`
	Links         *ResourceLinks          `json:"links,omitempty"`
}

// BlueprintAttributes contains the attributes of a Blueprint.
type BlueprintAttributes struct {
	Name                string     `json:"name,omitempty"`
	Description         string     `json:"description,omitempty"`
	Status              string     `json:"status,omitempty"`
	AppLicenseDeficient bool       `json:"appLicenseDeficient,omitempty"`
	CreatedDateTime     *time.Time `json:"createdDateTime,omitempty"`
	UpdatedDateTime     *time.Time `json:"updatedDateTime,omitempty"`
}

// BlueprintRelationships contains the relationship links returned in a Blueprint resource.
type BlueprintRelationships struct {
	Apps           *BlueprintRelationshipLink `json:"apps,omitempty"`
	Configurations *BlueprintRelationshipLink `json:"configurations,omitempty"`
	Packages       *BlueprintRelationshipLink `json:"packages,omitempty"`
	OrgDevices     *BlueprintRelationshipLink `json:"orgDevices,omitempty"`
	Users          *BlueprintRelationshipLink `json:"users,omitempty"`
	UserGroups     *BlueprintRelationshipLink `json:"userGroups,omitempty"`
}

// BlueprintLinkage is a single resource linkage (type + id) used in relationship payloads.
type BlueprintLinkage struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// BlueprintLinkageData wraps a list of linkages under the "data" key.
type BlueprintLinkageData struct {
	Data []BlueprintLinkage `json:"data"`
}

// --- Create request ---

// BlueprintCreateRequest is the request body for POST /v1/blueprints.
type BlueprintCreateRequest struct {
	Data BlueprintCreateRequestData `json:"data"`
}

// BlueprintCreateRequestData is the top-level data object for a create request.
type BlueprintCreateRequestData struct {
	Type          string                          `json:"type"` // must be "blueprints"
	Attributes    BlueprintCreateRequestAttributes `json:"attributes"`
	Relationships *BlueprintRequestRelationships  `json:"relationships,omitempty"`
}

// BlueprintCreateRequestAttributes contains attributes for creating a Blueprint.
// name is required.
type BlueprintCreateRequestAttributes struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// --- Update request ---

// BlueprintUpdateRequest is the request body for PATCH /v1/blueprints/{id}.
type BlueprintUpdateRequest struct {
	Data BlueprintUpdateRequestData `json:"data"`
}

// BlueprintUpdateRequestData is the top-level data object for an update request.
type BlueprintUpdateRequestData struct {
	Type          string                          `json:"type"` // must be "blueprints"
	ID            string                          `json:"id"`
	Attributes    BlueprintUpdateRequestAttributes `json:"attributes,omitempty"`
	Relationships *BlueprintRequestRelationships  `json:"relationships,omitempty"`
}

// BlueprintUpdateRequestAttributes contains attributes for updating a Blueprint.
// Only provided fields are updated.
type BlueprintUpdateRequestAttributes struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// BlueprintRequestRelationships is the relationships block used in create and update requests.
type BlueprintRequestRelationships struct {
	Apps           *BlueprintLinkageData `json:"apps,omitempty"`
	Configurations *BlueprintLinkageData `json:"configurations,omitempty"`
	Packages       *BlueprintLinkageData `json:"packages,omitempty"`
	OrgDevices     *BlueprintLinkageData `json:"orgDevices,omitempty"`
	Users          *BlueprintLinkageData `json:"users,omitempty"`
	UserGroups     *BlueprintLinkageData `json:"userGroups,omitempty"`
}

// --- Relationship linkage request types ---

// BlueprintAppsLinkagesRequest is the request body for add/remove apps relationship endpoints.
type BlueprintAppsLinkagesRequest struct {
	Data []BlueprintLinkage `json:"data"`
}

// BlueprintConfigurationsLinkagesRequest is the request body for add/remove configurations relationship endpoints.
type BlueprintConfigurationsLinkagesRequest struct {
	Data []BlueprintLinkage `json:"data"`
}

// BlueprintPackagesLinkagesRequest is the request body for add/remove packages relationship endpoints.
type BlueprintPackagesLinkagesRequest struct {
	Data []BlueprintLinkage `json:"data"`
}

// BlueprintOrgDevicesLinkagesRequest is the request body for add/remove orgDevices relationship endpoints.
type BlueprintOrgDevicesLinkagesRequest struct {
	Data []BlueprintLinkage `json:"data"`
}

// BlueprintUsersLinkagesRequest is the request body for add/remove users relationship endpoints.
type BlueprintUsersLinkagesRequest struct {
	Data []BlueprintLinkage `json:"data"`
}

// BlueprintUserGroupsLinkagesRequest is the request body for add/remove userGroups relationship endpoints.
type BlueprintUserGroupsLinkagesRequest struct {
	Data []BlueprintLinkage `json:"data"`
}

// --- Relationship linkage response types ---

// BlueprintAppsLinkagesResponse is the response for GET /v1/blueprints/{id}/relationships/apps.
type BlueprintAppsLinkagesResponse struct {
	Data  []BlueprintLinkage `json:"data"`
	Links *Links             `json:"links,omitempty"`
	Meta  *Meta              `json:"meta,omitempty"`
}

// BlueprintConfigurationsLinkagesResponse is the response for GET /v1/blueprints/{id}/relationships/configurations.
type BlueprintConfigurationsLinkagesResponse struct {
	Data  []BlueprintLinkage `json:"data"`
	Links *Links             `json:"links,omitempty"`
	Meta  *Meta              `json:"meta,omitempty"`
}

// BlueprintPackagesLinkagesResponse is the response for GET /v1/blueprints/{id}/relationships/packages.
type BlueprintPackagesLinkagesResponse struct {
	Data  []BlueprintLinkage `json:"data"`
	Links *Links             `json:"links,omitempty"`
	Meta  *Meta              `json:"meta,omitempty"`
}

// BlueprintOrgDevicesLinkagesResponse is the response for GET /v1/blueprints/{id}/relationships/orgDevices.
type BlueprintOrgDevicesLinkagesResponse struct {
	Data  []BlueprintLinkage `json:"data"`
	Links *Links             `json:"links,omitempty"`
	Meta  *Meta              `json:"meta,omitempty"`
}

// BlueprintUsersLinkagesResponse is the response for GET /v1/blueprints/{id}/relationships/users.
type BlueprintUsersLinkagesResponse struct {
	Data  []BlueprintLinkage `json:"data"`
	Links *Links             `json:"links,omitempty"`
	Meta  *Meta              `json:"meta,omitempty"`
}

// BlueprintUserGroupsLinkagesResponse is the response for GET /v1/blueprints/{id}/relationships/userGroups.
type BlueprintUserGroupsLinkagesResponse struct {
	Data  []BlueprintLinkage `json:"data"`
	Links *Links             `json:"links,omitempty"`
	Meta  *Meta              `json:"meta,omitempty"`
}

// --- Query options ---

// RequestQueryOptions represents query parameters for relationship list endpoints.
// limit is the only applicable parameter (max 1000).
type RequestQueryOptions struct {
	Limit int
}

// GetBlueprintQueryOptions represents query parameters for GetByBlueprintIDV1.
type GetBlueprintQueryOptions struct {
	// Fields specifies which fields to return. Use Field* constants.
	Fields []string
	// Include specifies related resources to include. Use Include* constants.
	Include []string
	// Per-relationship limits (max 1000 each).
	LimitApps           int
	LimitConfigurations int
	LimitPackages       int
	LimitOrgDevices     int
	LimitUsers          int
	LimitUserGroups     int
}
