package auditevents

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

// AuditEventsResponse is the response for a list of audit events.
type AuditEventsResponse struct {
	Data  []AuditEvent `json:"data"`
	Links *Links       `json:"links,omitempty"`
	Meta  *Meta        `json:"meta,omitempty"`
}

// AuditEvent represents a single audit event resource.
type AuditEvent struct {
	ID         string              `json:"id"`
	Type       string              `json:"type"`
	Attributes *AuditEventAttributes `json:"attributes,omitempty"`
}

// AuditEventAttributes contains all attributes of an audit event.
type AuditEventAttributes struct {
	EventDateTime        *time.Time `json:"eventDateTime,omitempty"`
	Type                 string     `json:"type,omitempty"`
	Category             string     `json:"category,omitempty"`
	ActorType            string     `json:"actorType,omitempty"`
	ActorId              string     `json:"actorId,omitempty"`
	ActorName            string     `json:"actorName,omitempty"`
	SubjectType          string     `json:"subjectType,omitempty"`
	SubjectId            string     `json:"subjectId,omitempty"`
	SubjectName          string     `json:"subjectName,omitempty"`
	Outcome              string     `json:"outcome,omitempty"`
	GroupId              string     `json:"groupId,omitempty"`
	EventDataPropertyKey string     `json:"eventDataPropertyKey,omitempty"`

	EventDataDeviceAddedToOrg               *EventDataDeviceAddedToOrg               `json:"eventDataDeviceAddedToOrg,omitempty"`
	EventDataDeviceRemovedFromOrg           *EventDataDeviceRemovedFromOrg           `json:"eventDataDeviceRemovedFromOrg,omitempty"`
	EventDataDeviceAssignedToServer         *EventDataDeviceAssignedToServer         `json:"eventDataDeviceAssignedToServer,omitempty"`
	EventDataDeviceUnassignedFromServer     *EventDataDeviceUnassignedFromServer     `json:"eventDataDeviceUnassignedFromServer,omitempty"`
	EventDataDeviceIsErased                 *EventDataDeviceIsErased                 `json:"eventDataDeviceIsErased,omitempty"`
	EventDataConfigSettingsCreated          *EventDataConfigSettings                 `json:"eventDataConfigSettingsCreated,omitempty"`
	EventDataConfigSettingsUpdated          *EventDataConfigSettings                 `json:"eventDataConfigSettingsUpdated,omitempty"`
	EventDataConfigSettingsDeleted          *EventDataConfigSettings                 `json:"eventDataConfigSettingsDeleted,omitempty"`
	EventDataCollectionCreated              *EventDataCollection                     `json:"eventDataCollectionCreated,omitempty"`
	EventDataCollectionUpdated              *EventDataCollection                     `json:"eventDataCollectionUpdated,omitempty"`
	EventDataCollectionDeleted              *EventDataCollection                     `json:"eventDataCollectionDeleted,omitempty"`
	EventDataSubscriptionCreated            *EventDataSubscription                   `json:"eventDataSubscriptionCreated,omitempty"`
	EventDataSubscriptionUpdated            *EventDataSubscription                   `json:"eventDataSubscriptionUpdated,omitempty"`
	EventDataSubscriptionDeleted            *EventDataSubscription                   `json:"eventDataSubscriptionDeleted,omitempty"`
	EventDataAccountRoleLocationChanged     *EventDataAccountRoleLocationChanged     `json:"eventDataAccountRoleLocationChanged,omitempty"`
	EventDataAccountAdded                   *EventDataAccountAdded                   `json:"eventDataAccountAdded,omitempty"`
	EventDataAccountDeleted                 *EventDataAccountDeleted                 `json:"eventDataAccountDeleted,omitempty"`
	EventDataExternalAccountAssociated      *EventDataExternalAccount                `json:"eventDataExternalAccountAssociated,omitempty"`
	EventDataExternalAccountDisassociated   *EventDataExternalAccount                `json:"eventDataExternalAccountDisassociated,omitempty"`
	EventDataDomainAdded                    *EventDataDomain                         `json:"eventDataDomainAdded,omitempty"`
	EventDataDomainRemoved                  *EventDataDomain                         `json:"eventDataDomainRemoved,omitempty"`
	EventDataDomainVerified                 *EventDataDomain                         `json:"eventDataDomainVerified,omitempty"`
	EventDataApiAccountCreatedWithKey       *EventDataApiAccount                     `json:"eventDataApiAccountCreatedWithKey,omitempty"`
	EventDataApiAccountCreatedWithoutKey    *EventDataApiAccount                     `json:"eventDataApiAccountCreatedWithoutKey,omitempty"`
	EventDataApiAccountDeleted              *EventDataApiAccount                     `json:"eventDataApiAccountDeleted,omitempty"`
	EventDataApiAccountKeyGenerated         *EventDataApiAccount                     `json:"eventDataApiAccountKeyGenerated,omitempty"`
	EventDataApiAccountKeyRevoked           *EventDataApiAccount                     `json:"eventDataApiAccountKeyRevoked,omitempty"`
	EventDataApiAccountNameChanged          *EventDataApiAccount                     `json:"eventDataApiAccountNameChanged,omitempty"`
	EventDataApiAccountRoleLocationChanged  *EventDataApiAccountRoleLocationChanged  `json:"eventDataApiAccountRoleLocationChanged,omitempty"`
	EventDataSubjectHasICloudStoragePurchaseAdded    *EventDataPurchase `json:"eventDataSubjectHasICloudStoragePurchaseAdded,omitempty"`
	EventDataSubjectHasICloudStoragePurchaseRemoved  *EventDataPurchase `json:"eventDataSubjectHasICloudStoragePurchaseRemoved,omitempty"`
	EventDataSubjectHasAppleCarePurchaseAdded        *EventDataPurchase `json:"eventDataSubjectHasAppleCarePurchaseAdded,omitempty"`
	EventDataSubjectHasAppleCarePurchaseRemoved      *EventDataPurchase `json:"eventDataSubjectHasAppleCarePurchaseRemoved,omitempty"`
}

// EventDataDeviceAddedToOrg contains data for a device added to org event.
type EventDataDeviceAddedToOrg struct {
	SerialNumber       string `json:"serialNumber,omitempty"`
	PurchaseSourceType string `json:"purchaseSourceType,omitempty"`
	PurchaseSourceId   string `json:"purchaseSourceId,omitempty"`
}

// EventDataDeviceRemovedFromOrg contains data for a device removed from org event.
type EventDataDeviceRemovedFromOrg struct {
	SerialNumber string `json:"serialNumber,omitempty"`
}

// EventDataDeviceAssignedToServer contains data for a device assigned to server event.
type EventDataDeviceAssignedToServer struct {
	SerialNumber string `json:"serialNumber,omitempty"`
	ServerName   string `json:"serverName,omitempty"`
}

// EventDataDeviceUnassignedFromServer contains data for a device unassigned from server event.
type EventDataDeviceUnassignedFromServer struct {
	SerialNumber string `json:"serialNumber,omitempty"`
	ServerName   string `json:"serverName,omitempty"`
}

// EventDataDeviceIsErased contains data for a device erased event.
type EventDataDeviceIsErased struct {
	SerialNumber string `json:"serialNumber,omitempty"`
}

// EventDataConfigSettings contains data for configuration settings events.
type EventDataConfigSettings struct {
	ConfigName string `json:"configName,omitempty"`
	ConfigType string `json:"configType,omitempty"`
}

// EventDataCollection contains data for collection events.
type EventDataCollection struct {
	CollectionName string `json:"collectionName,omitempty"`
}

// EventDataSubscription contains data for subscription events.
type EventDataSubscription struct {
	SubscriptionName string `json:"subscriptionName,omitempty"`
}

// EventDataAccountRoleLocationChanged contains data for account role/location change events.
type EventDataAccountRoleLocationChanged struct {
	AccountName string `json:"accountName,omitempty"`
	OldRole     string `json:"oldRole,omitempty"`
	NewRole     string `json:"newRole,omitempty"`
}

// EventDataAccountAdded contains data for account added events.
type EventDataAccountAdded struct {
	AccountName string `json:"accountName,omitempty"`
	AccountType string `json:"accountType,omitempty"`
}

// EventDataAccountDeleted contains data for account deleted events.
type EventDataAccountDeleted struct {
	AccountName string `json:"accountName,omitempty"`
}

// EventDataExternalAccount contains data for external account association events.
type EventDataExternalAccount struct {
	ExternalAccountName string `json:"externalAccountName,omitempty"`
	Provider            string `json:"provider,omitempty"`
}

// EventDataDomain contains data for domain events.
type EventDataDomain struct {
	DomainName string `json:"domainName,omitempty"`
}

// EventDataApiAccount contains data for API account events.
type EventDataApiAccount struct {
	AccountName string `json:"accountName,omitempty"`
}

// EventDataApiAccountRoleLocationChanged contains data for API account role/location change events.
type EventDataApiAccountRoleLocationChanged struct {
	AccountName string `json:"accountName,omitempty"`
	OldRole     string `json:"oldRole,omitempty"`
	NewRole     string `json:"newRole,omitempty"`
}

// EventDataPurchase contains data for purchase-related events.
type EventDataPurchase struct {
	ProductName string `json:"productName,omitempty"`
}

// RequestQueryOptions represents query parameters for the audit events endpoint.
type RequestQueryOptions struct {
	// FilterStartTimestamp is the ISO8601 formatted start timestamp (Required).
	FilterStartTimestamp string
	// FilterEndTimestamp is the ISO8601 formatted end timestamp (Required).
	FilterEndTimestamp string
	// FilterActorID filters by actor ID. Only one actor ID is supported.
	FilterActorID string
	// FilterSubjectID filters by subject ID. Only one subject ID is supported.
	FilterSubjectID string
	// FilterType filters by event type. Only one type is supported. Use AuditEventType* constants.
	FilterType string
	// Fields specifies which fields to return. Use Field* constants.
	Fields []string
	// Limit is the number of resources to return (max 1000).
	Limit int
}
