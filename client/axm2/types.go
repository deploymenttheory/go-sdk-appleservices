package axm2

// Type aliases that map client interfaces to service package types
// This allows the client interfaces to reference these types while
// keeping the actual type definitions in the service packages

import (
	"github.com/deploymenttheory/go-api-sdk-apple/services/axm/mdm_servers"
	"github.com/deploymenttheory/go-api-sdk-apple/services/axm/org_device_activities"
	"github.com/deploymenttheory/go-api-sdk-apple/services/axm/org_devices"
)

// Organization Devices types
type OrgDevice = org_devices.OrgDevice
type OrgDeviceAttributes = org_devices.OrgDeviceAttributes
type OrgDevicesResponse = org_devices.OrgDevicesResponse
type OrgDeviceResponse = org_devices.OrgDeviceResponse

// MDM Servers types
type MdmServer = mdm_servers.MdmServer
type MdmServerAttributes = mdm_servers.MdmServerAttributes
type MdmServersResponse = mdm_servers.MdmServersResponse
type MdmServerResponse = mdm_servers.MdmServerResponse

// Organization Device Activities types
type OrgDeviceActivity = org_device_activities.OrgDeviceActivity
type OrgDeviceActivityAttributes = org_device_activities.OrgDeviceActivityAttributes
type OrgDeviceActivityResponse = org_device_activities.OrgDeviceActivityResponse

// Common types (using org_devices as the canonical source)
type ResourceLinkage = org_devices.ResourceLinkage
type PagedDocumentLinks = org_devices.PagedDocumentLinks
type PagingInformation = org_devices.PagingInformation
type PagingDetails = org_devices.PagingDetails
