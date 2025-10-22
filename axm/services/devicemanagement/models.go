package devicemanagement

import "time"

// Shared types for pagination and links
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

// MDMServer represents an MDM server in the Apple Business Manager system
type MDMServer struct {
	ID            string                  `json:"id"`
	Type          string                  `json:"type"`
	Attributes    *MDMServerAttributes    `json:"attributes,omitempty"`
	Relationships *MDMServerRelationships `json:"relationships,omitempty"`
}

// MDMServerAttributes contains the MDM server attributes
type MDMServerAttributes struct {
	ServerName      string     `json:"serverName,omitempty"`
	ServerType      string     `json:"serverType,omitempty"`
	CreatedDateTime *time.Time `json:"createdDateTime,omitempty"`
	UpdatedDateTime *time.Time `json:"updatedDateTime,omitempty"`
	Devices         []string   `json:"devices,omitempty"`
}

// MDMServerRelationships contains the MDM server relationships
type MDMServerRelationships struct {
	Devices *MDMServerDevicesRelationship `json:"devices,omitempty"`
}

// MDMServerDevicesRelationship contains the devices relationship links
type MDMServerDevicesRelationship struct {
	Links *MDMServerDevicesLinks `json:"links,omitempty"`
}

// MDMServerDevicesLinks contains the navigation links for devices
type MDMServerDevicesLinks struct {
	Self string `json:"self,omitempty"`
}

// MDMServersResponse represents the response for getting MDM servers
type MDMServersResponse struct {
	Data  []MDMServer `json:"data"`
	Meta  *Meta       `json:"meta,omitempty"`
	Links *Links      `json:"links,omitempty"`
}

// RequestQueryOptions represents the query parameters for getting MDM servers
type RequestQueryOptions struct {
	// Field selection - fields to return for mdmServers
	// Possible values: serverName, serverType, createdDateTime, updatedDateTime, devices
	Fields []string `json:"fields,omitempty"`

	// Limit the number of included related resources to return (max 1000)
	Limit int `json:"limit,omitempty"`
}

// MDMServerDeviceLinkage represents a device linkage in the MDM server relationships
type MDMServerDeviceLinkage struct {
	Type string `json:"type"` // Should be "orgDevices"
	ID   string `json:"id"`   // Device ID
}

// MDMServerDevicesLinkagesResponse represents the response for getting device linkages for an MDM server
type MDMServerDevicesLinkagesResponse struct {
	Data  []MDMServerDeviceLinkage `json:"data"`
	Links *Links                   `json:"links,omitempty"`
	Meta  *Meta                    `json:"meta,omitempty"`
}

// GetMDMServerDeviceLinkagesOptions represents the query parameters for getting MDM server device linkages
type GetMDMServerDeviceLinkagesOptions struct {
	// Limit the number of included related resources to return (max 1000)
	Limit int `json:"limit,omitempty"`
}

// OrgDeviceAssignedServerLinkage represents the linkage between a device and its assigned server
type OrgDeviceAssignedServerLinkage struct {
	Type string `json:"type"` // Should be "mdmServers"
	ID   string `json:"id"`   // MDM Server ID
}

// OrgDeviceAssignedServerLinkageResponse represents the response for getting assigned server linkage
type OrgDeviceAssignedServerLinkageResponse struct {
	Data  OrgDeviceAssignedServerLinkage `json:"data"`
	Links *AssignedServerLinks           `json:"links,omitempty"`
}

// AssignedServerLinks contains linkage navigation links
type AssignedServerLinks struct {
	Self    string `json:"self,omitempty"`
	Related string `json:"related,omitempty"`
}

// MDMServerResponse represents the response for getting a single MDM server
type MDMServerResponse struct {
	Data  MDMServer `json:"data"`
	Links *Links    `json:"links,omitempty"`
}

// GetAssignedServerInfoOptions represents the query parameters for getting assigned server info
type GetAssignedServerInfoOptions struct {
	// Field selection - fields to return for mdmServers
	// Possible values: serverName, serverType, createdDateTime, updatedDateTime, devices
	Fields []string `json:"fields,omitempty"`
}

// OrgDeviceActivity represents a device activity (assign/unassign operations)
type OrgDeviceActivity struct {
	ID         string                       `json:"id"`
	Type       string                       `json:"type"`
	Attributes *OrgDeviceActivityAttributes `json:"attributes,omitempty"`
	Links      *OrgDeviceActivityLinks      `json:"links,omitempty"`
}

// OrgDeviceActivityAttributes contains the activity attributes
type OrgDeviceActivityAttributes struct {
	Status          string     `json:"status,omitempty"`
	SubStatus       string     `json:"subStatus,omitempty"`
	CreatedDateTime *time.Time `json:"createdDateTime,omitempty"`
	ActivityType    string     `json:"activityType,omitempty"`
}

// OrgDeviceActivityLinks contains activity navigation links
type OrgDeviceActivityLinks struct {
	Self string `json:"self,omitempty"`
}

// OrgDeviceActivityResponse represents the response for creating an org device activity
type OrgDeviceActivityResponse struct {
	Data  OrgDeviceActivity `json:"data"`
	Links *Links            `json:"links,omitempty"`
}

// OrgDeviceActivityCreateRequest represents the request for creating a device activity
type OrgDeviceActivityCreateRequest struct {
	Data OrgDeviceActivityData `json:"data"`
}

// OrgDeviceActivityData contains the activity data for the request
type OrgDeviceActivityData struct {
	Type          string                               `json:"type"`
	Attributes    OrgDeviceActivityCreateAttributes    `json:"attributes"`
	Relationships OrgDeviceActivityCreateRelationships `json:"relationships"`
}

// OrgDeviceActivityCreateAttributes contains the activity creation attributes
type OrgDeviceActivityCreateAttributes struct {
	ActivityType string `json:"activityType"`
}

// OrgDeviceActivityCreateRelationships contains the relationships for activity creation
type OrgDeviceActivityCreateRelationships struct {
	MDMServer *OrgDeviceActivityMDMServerRelationship `json:"mdmServer,omitempty"`
	Devices   *OrgDeviceActivityDevicesRelationship   `json:"devices,omitempty"`
}

// OrgDeviceActivityMDMServerRelationship represents the MDM server relationship
type OrgDeviceActivityMDMServerRelationship struct {
	Data OrgDeviceActivityMDMServerLinkage `json:"data"`
}

// OrgDeviceActivityMDMServerLinkage represents the MDM server linkage
type OrgDeviceActivityMDMServerLinkage struct {
	Type string `json:"type"` // Should be "mdmServers"
	ID   string `json:"id"`   // MDM Server ID
}

// OrgDeviceActivityDevicesRelationship represents the devices relationship
type OrgDeviceActivityDevicesRelationship struct {
	Data []OrgDeviceActivityDeviceLinkage `json:"data"`
}

// OrgDeviceActivityDeviceLinkage represents a device linkage
type OrgDeviceActivityDeviceLinkage struct {
	Type string `json:"type"` // Should be "orgDevices"
	ID   string `json:"id"`   // Device ID
}
