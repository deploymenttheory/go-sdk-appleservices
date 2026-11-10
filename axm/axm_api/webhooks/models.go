// Package webhooks manages Apple Business API webhook subscriptions and deliveries.
package webhooks

import (
	"time"

	"github.com/deploymenttheory/go-sdk-appleservices/axm/client"
)

type Meta = client.Meta
type Links = client.Links
type Paging = client.Paging

// ResourceLinks contains resource and relationship navigation links.
type ResourceLinks struct {
	Self    string `json:"self,omitempty"`
	Related string `json:"related,omitempty"`
}

// ResourceIdentifier identifies a related Apple API resource.
type ResourceIdentifier struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// ResourceRelationship contains a resource linkage and optional navigation links.
type ResourceRelationship struct {
	Data  ResourceIdentifier `json:"data"`
	Links *ResourceLinks     `json:"links,omitempty"`
}

type Webhook struct {
	ID            string                `json:"id"`
	Type          string                `json:"type"`
	Attributes    *WebhookAttributes    `json:"attributes,omitempty"`
	Relationships *WebhookRelationships `json:"relationships,omitempty"`
	Links         *ResourceLinks        `json:"links,omitempty"`
}

type WebhookAttributes struct {
	Name            string     `json:"name,omitempty"`
	Enabled         bool       `json:"enabled,omitempty"`
	EventTypes      []string   `json:"eventTypes,omitempty"`
	URL             string     `json:"url,omitempty"`
	CreatedDateTime *time.Time `json:"createdDateTime,omitempty"`
	UpdatedDateTime *time.Time `json:"updatedDateTime,omitempty"`
}

type WebhookRelationships struct {
	Deliveries *WebhookDeliveriesRelationship `json:"deliveries,omitempty"`
}

type WebhookDeliveriesRelationship struct {
	Links *ResourceLinks `json:"links,omitempty"`
}

type WebhooksResponse struct {
	Data  []Webhook `json:"data"`
	Links *Links    `json:"links,omitempty"`
	Meta  *Meta     `json:"meta,omitempty"`
}

type WebhookResponse struct {
	Data  Webhook `json:"data"`
	Links *Links  `json:"links,omitempty"`
}

type WebhookCreateRequest struct {
	Data WebhookCreateRequestData `json:"data"`
}

type WebhookCreateRequestData struct {
	Type       string                         `json:"type"`
	Attributes WebhookCreateRequestAttributes `json:"attributes"`
}

// WebhookCreateRequestAttributes includes enabled even when false.
// The signing secret is write-only and never returned in WebhookAttributes.
type WebhookCreateRequestAttributes struct {
	Name       string   `json:"name"`
	Enabled    bool     `json:"enabled"`
	EventTypes []string `json:"eventTypes"`
	Secret     string   `json:"secret"`
	URL        string   `json:"url"`
}

type WebhookUpdateRequest struct {
	Data WebhookUpdateRequestData `json:"data"`
}

type WebhookUpdateRequestData struct {
	Type       string                          `json:"type"`
	ID         string                          `json:"id"`
	Attributes *WebhookUpdateRequestAttributes `json:"attributes,omitempty"`
}

// WebhookUpdateRequestAttributes uses pointers so omission and explicit updates differ.
type WebhookUpdateRequestAttributes struct {
	Name       *string   `json:"name,omitempty"`
	Enabled    *bool     `json:"enabled,omitempty"`
	EventTypes *[]string `json:"eventTypes,omitempty"`
	Secret     *string   `json:"secret,omitempty"`
	URL        *string   `json:"url,omitempty"`
}

type WebhookDelivery struct {
	ID            string                        `json:"id"`
	Type          string                        `json:"type"`
	Attributes    *WebhookDeliveryAttributes    `json:"attributes,omitempty"`
	Relationships *WebhookDeliveryRelationships `json:"relationships,omitempty"`
}

type WebhookDeliveryAttributes struct {
	CreatedDateTime *time.Time                   `json:"createdDateTime,omitempty"`
	DeliveryState   string                       `json:"deliveryState,omitempty"`
	ErrorMessage    string                       `json:"errorMessage,omitempty"`
	Redelivery      bool                         `json:"redelivery,omitempty"`
	SentDateTime    *time.Time                   `json:"sentDateTime,omitempty"`
	Request         *WebhookDeliveryRequest      `json:"request,omitempty"`
	Response        *WebhookDeliveryHTTPResponse `json:"response,omitempty"`
}

type WebhookDeliveryRequest struct {
	URL string `json:"url,omitempty"`
}

type WebhookDeliveryHTTPResponse struct {
	HTTPStatusCode int    `json:"httpStatusCode,omitempty"`
	Body           string `json:"body,omitempty"`
}

type WebhookDeliveryRelationships struct {
	Event *ResourceRelationship `json:"event,omitempty"`
}

type WebhookEvent struct {
	ID         string                  `json:"id"`
	Type       string                  `json:"type"`
	Attributes *WebhookEventAttributes `json:"attributes,omitempty"`
}

type WebhookEventAttributes struct {
	EventType string `json:"eventType,omitempty"`
	// Payload contains JSON encoded as a string, as returned by Apple.
	Payload         string     `json:"payload,omitempty"`
	IsPingEvent     bool       `json:"isPingEvent,omitempty"`
	CreatedDateTime *time.Time `json:"createdDateTime,omitempty"`
}

type WebhookDeliveriesResponse struct {
	Data     []WebhookDelivery `json:"data"`
	Included []WebhookEvent    `json:"included,omitempty"`
	Links    *Links            `json:"links,omitempty"`
	Meta     *Meta             `json:"meta,omitempty"`
}

type WebhookDeliveryResponse struct {
	Data     WebhookDelivery `json:"data"`
	Included []WebhookEvent  `json:"included,omitempty"`
	Links    *Links          `json:"links,omitempty"`
}

type WebhookDeliveryCreateRequest struct {
	Data WebhookDeliveryCreateRequestData `json:"data"`
}

type WebhookDeliveryCreateRequestData struct {
	Type          string                             `json:"type"`
	Relationships WebhookDeliveryCreateRelationships `json:"relationships"`
}

// WebhookDeliveryCreateRelationships references the delivery to replay.
type WebhookDeliveryCreateRelationships struct {
	Template ResourceRelationship `json:"template"`
}

type WebhookPing struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type WebhookPingResponse struct {
	Data  WebhookPing `json:"data"`
	Links *Links      `json:"links,omitempty"`
}

type WebhookPingCreateRequest struct {
	Data WebhookPingCreateRequestData `json:"data"`
}

type WebhookPingCreateRequestData struct {
	Type          string                         `json:"type"`
	Relationships WebhookPingCreateRelationships `json:"relationships"`
}

type WebhookPingCreateRelationships struct {
	Webhook ResourceRelationship `json:"webhook"`
}

// RequestQueryOptions selects webhook fields and the starting page.
// Limit controls page size, not the total number returned. Maximum: 1000.
type RequestQueryOptions struct {
	Fields []string
	Limit  int
	Cursor string
}

// DeliveryQueryOptions supports the documented delivery filters and event inclusion.
// Date filters are serialized unchanged; use RFC3339 timestamp strings.
// Apple's live API currently requires both created-date filters, despite marking
// them optional in the endpoint documentation.
type DeliveryQueryOptions struct {
	Fields                         []string
	EventFields                    []string
	Include                        []string
	FilterDeliveryState            []string
	FilterCreatedDateTimeAtOrAfter []string
	FilterCreatedDateTimeAtBefore  []string
	Limit                          int
	Cursor                         string
}
