package org_device_activities

import "time"

// OrgDeviceActivity represents an organization device activity resource
// This is for getting information about device management actions that were performed
type OrgDeviceActivity struct {
	Type       string                      `json:"type"`
	ID         string                      `json:"id"`
	Attributes OrgDeviceActivityAttributes `json:"attributes"`
	Links      ResourceLinks               `json:"links,omitempty"`
}

// OrgDeviceActivityAttributes contains the device activity attributes
// Based on Apple's API spec with supported fields: status, subStatus, createdDateTime, completedDateTime, downloadUrl
type OrgDeviceActivityAttributes struct {
	Status            string     `json:"status,omitempty"`            // Activity status (e.g., "COMPLETED")
	SubStatus         string     `json:"subStatus,omitempty"`         // Detailed sub-status (e.g., "COMPLETED_WITH_SUCCESS")
	CreatedDateTime   time.Time  `json:"createdDateTime,omitempty"`   // When the activity was created
	CompletedDateTime *time.Time `json:"completedDateTime,omitempty"` // When the activity completed
	DownloadUrl       string     `json:"downloadUrl,omitempty"`       // URL to download activity CSV report
}

// OrgDeviceActivityResponse represents a response containing a single device activity
type OrgDeviceActivityResponse struct {
	Data  OrgDeviceActivity  `json:"data"`
	Links PagedDocumentLinks `json:"links"`
}

// OrgDeviceActivitiesRequest represents the request body for device assignment operations
// Apple API endpoint: POST /v1/orgDeviceActivities
// Follows the <endpoint>Request naming pattern
type OrgDeviceActivitiesRequest struct {
	Data OrgDeviceActivitiesRequestData `json:"data"`
}

// OrgDeviceActivitiesRequestData contains the data for creating device activities
type OrgDeviceActivitiesRequestData struct {
	Type          string                                  `json:"type"`
	Attributes    OrgDeviceActivitiesRequestAttributes    `json:"attributes"`
	Relationships OrgDeviceActivitiesRequestRelationships `json:"relationships"`
}

// OrgDeviceActivitiesRequestAttributes contains attributes for creating device activities
type OrgDeviceActivitiesRequestAttributes struct {
	ActivityType string `json:"activityType"` // "ASSIGN_DEVICES" or "UNASSIGN_DEVICES"
}

// OrgDeviceActivitiesRequestRelationships contains relationships for creating device activities
type OrgDeviceActivitiesRequestRelationships struct {
	MdmServer *MdmServerRelationship `json:"mdmServer,omitempty"` // For assignment operations
	Devices   *DevicesRelationship   `json:"devices"`             // Devices to assign/unassign
}

// MdmServerRelationship represents the mdmServer relationship
type MdmServerRelationship struct {
	Data ResourceLinkage `json:"data"`
}

// DevicesRelationship represents the devices relationship
type DevicesRelationship struct {
	Data []ResourceLinkage `json:"data"`
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

// Activity type constants
const (
	ActivityTypeAssignDevices   = "ASSIGN_DEVICES"
	ActivityTypeUnassignDevices = "UNASSIGN_DEVICES"
)
