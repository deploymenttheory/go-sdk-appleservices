package configurations

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

// ConfigurationsResponse is the response for a list of configurations.
// Note: customSettingsValues is always null in list responses per API spec.
type ConfigurationsResponse struct {
	Data  []Configuration `json:"data"`
	Links *Links          `json:"links,omitempty"`
	Meta  *Meta           `json:"meta,omitempty"`
}

// ConfigurationResponse is the response for a single configuration.
type ConfigurationResponse struct {
	Data  Configuration `json:"data"`
	Links *Links        `json:"links,omitempty"`
}

// Configuration represents a configuration resource.
type Configuration struct {
	ID         string                   `json:"id"`
	Type       string                   `json:"type"`
	Attributes *ConfigurationAttributes `json:"attributes,omitempty"`
	Links      *ResourceLinks           `json:"links,omitempty"`
}

// ConfigurationAttributes contains the attributes of a configuration.
type ConfigurationAttributes struct {
	Type                   string                `json:"type,omitempty"`
	Name                   string                `json:"name,omitempty"`
	ConfiguredForPlatforms []string              `json:"configuredForPlatforms,omitempty"`
	CustomSettingsValues   *CustomSettingsValues `json:"customSettingsValues,omitempty"`
	CreatedDateTime        *time.Time            `json:"createdDateTime,omitempty"`
	UpdatedDateTime        *time.Time            `json:"updatedDateTime,omitempty"`
}

// CustomSettingsValues holds the profile content for CUSTOM_SETTING configurations.
type CustomSettingsValues struct {
	ConfigurationProfile string `json:"configurationProfile,omitempty"`
	Filename             string `json:"filename,omitempty"`
}

// ConfigurationCreateRequest is the request body for creating a configuration.
// Only configurations with type CUSTOM_SETTING can be created via the API.
type ConfigurationCreateRequest struct {
	Data ConfigurationCreateRequestData `json:"data"`
}

// ConfigurationCreateRequestData is the data object for a create request.
type ConfigurationCreateRequestData struct {
	Type       string                               `json:"type"` // must be "configurations"
	Attributes ConfigurationCreateRequestAttributes `json:"attributes"`
}

// ConfigurationCreateRequestAttributes contains attributes for creating a configuration.
// configurationProfile is required. filename and configuredForPlatforms are optional.
type ConfigurationCreateRequestAttributes struct {
	Type                   string               `json:"type"` // must be "CUSTOM_SETTING"
	Name                   string               `json:"name"`
	ConfiguredForPlatforms []string             `json:"configuredForPlatforms,omitempty"`
	CustomSettingsValues   CustomSettingsValues `json:"customSettingsValues"`
}

// ConfigurationUpdateRequest is the request body for updating a configuration.
// Only CUSTOM_SETTING configurations can be updated. Only provided fields are changed.
type ConfigurationUpdateRequest struct {
	Data ConfigurationUpdateRequestData `json:"data"`
}

// ConfigurationUpdateRequestData is the data object for an update request.
type ConfigurationUpdateRequestData struct {
	Type       string                               `json:"type"` // must be "configurations"
	ID         string                               `json:"id"`
	Attributes ConfigurationUpdateRequestAttributes `json:"attributes"`
}

// ConfigurationUpdateRequestAttributes contains attributes for updating a configuration.
// At least one of name, configuredForPlatforms, configurationProfile, or filename must be provided.
// If filename is provided it must end in .mobileconfig.
type ConfigurationUpdateRequestAttributes struct {
	Name                   string                `json:"name,omitempty"`
	ConfiguredForPlatforms []string              `json:"configuredForPlatforms,omitempty"`
	CustomSettingsValues   *CustomSettingsValues `json:"customSettingsValues,omitempty"`
}

// RequestQueryOptions represents query parameters for configuration endpoints.
type RequestQueryOptions struct {
	// Fields specifies which fields to return. Use Field* constants.
	Fields []string
	// Limit is the number of resources to return (max 1000).
	Limit int
}
