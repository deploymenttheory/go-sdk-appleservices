package mdm_servers

// MdmServer represents device management services in organizations
type MdmServer struct {
	Type          string              `json:"type"`
	ID            string              `json:"id"`
	Attributes    MdmServerAttributes `json:"attributes"`
	Relationships any                 `json:"relationships,omitempty"`
	Links         ResourceLinks       `json:"links,omitempty"`
}

// MdmServerAttributes contains the MDM server attributes
// Based on Apple API documentation fields: serverName, serverType, createdDateTime, updatedDateTime, devices
type MdmServerAttributes struct {
	ServerName      string `json:"serverName,omitempty"`
	ServerType      string `json:"serverType,omitempty"`
	CreatedDateTime string `json:"createdDateTime,omitempty"`
	UpdatedDateTime string `json:"updatedDateTime,omitempty"`
}

// MdmServerResponse represents a response containing a single MDM server resource
type MdmServerResponse struct {
	Data  MdmServer          `json:"data"`
	Links PagedDocumentLinks `json:"links"`
	Meta  *PagingInformation `json:"meta,omitempty"`
}

// MdmServersResponse represents a response containing multiple MDM server resources
type MdmServersResponse struct {
	Data  []MdmServer        `json:"data"`
	Links PagedDocumentLinks `json:"links"`
	Meta  *PagingInformation `json:"meta,omitempty"`
}

// GetData returns the data slice for pagination
func (r *MdmServersResponse) GetData() interface{} {
	return r.Data
}

// GetNextURL returns the next URL for pagination
func (r *MdmServersResponse) GetNextURL() string {
	return r.Links.Next
}

// AppendData appends data from another response for pagination
func (r *MdmServersResponse) AppendData(data interface{}) {
	if existingData, ok := data.([]MdmServer); ok {
		r.Data = append(r.Data, existingData...)
	}
}

// MdmServerDevicesLinkagesResponse represents the relationship between MDM servers and devices
type MdmServerDevicesLinkagesResponse struct {
	Data  []ResourceLinkage  `json:"data"`
	Links PagedDocumentLinks `json:"links"`
	Meta  *PagingInformation `json:"meta,omitempty"`
}

// ResourceLinkage represents a linkage between resources
type ResourceLinkage struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// ResourceLinks contains self-links to requested resources
type ResourceLinks struct {
	Self string `json:"self,omitempty"`
}

// PagedDocumentLinks contains navigational links for paged responses
type PagedDocumentLinks struct {
	Self  string `json:"self"`
	First string `json:"first,omitempty"`
	Next  string `json:"next,omitempty"`
	Prev  string `json:"prev,omitempty"`
	Last  string `json:"last,omitempty"`
}

// PagingInformation contains metadata about pagination
type PagingInformation struct {
	Paging *PagingDetails `json:"paging,omitempty"`
}

// PagingDetails contains specific paging information
type PagingDetails struct {
	NextCursor string `json:"nextCursor,omitempty"`
	Limit      int    `json:"limit,omitempty"`
}

// APIError represents an API error response
type APIError struct {
	Errors []ErrorDetail `json:"errors,omitempty"`
}

// ErrorDetail contains detailed error information
type ErrorDetail struct {
	ID     string `json:"id,omitempty"`
	Status string `json:"status,omitempty"`
	Code   string `json:"code,omitempty"`
	Title  string `json:"title,omitempty"`
	Detail string `json:"detail,omitempty"`
	Source *struct {
		Pointer   string `json:"pointer,omitempty"`
		Parameter string `json:"parameter,omitempty"`
	} `json:"source,omitempty"`
}
