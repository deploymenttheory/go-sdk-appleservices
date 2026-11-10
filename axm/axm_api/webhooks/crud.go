package webhooks

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/deploymenttheory/go-sdk-appleservices/axm/client"
	"github.com/deploymenttheory/go-sdk-appleservices/axm/constants"
	"resty.dev/v3"
)

// Webhooks handles webhook configuration, delivery inspection, redelivery and pings.
type Webhooks struct{ client client.Client }

// NewService creates a service backed by the shared AXM transport.
func NewService(c client.Client) *Webhooks { return &Webhooks{client: c} }

func (s *Webhooks) request(ctx context.Context) *client.RequestBuilder {
	return s.client.NewRequest(ctx).SetHeader("Accept", constants.ApplicationJSON).SetHeader("Content-Type", constants.ApplicationJSON)
}

func webhookPath(id string) string { return constants.EndpointWebhooks + "/" + url.PathEscape(id) }

// GetV1 fetches all webhook pages, starting at opts.Cursor when provided.
// https://developer.apple.com/documentation/applebusinessapi/get-webhooks
func (s *Webhooks) GetV1(ctx context.Context, opts *RequestQueryOptions) (*WebhooksResponse, *resty.Response, error) {
	if opts == nil {
		opts = &RequestQueryOptions{}
	}
	params := s.client.QueryBuilder().AddStringSlice("fields[webhooks]", opts.Fields).AddInt("limit", min(opts.Limit, 1000)).AddString("cursor", opts.Cursor)
	result := WebhooksResponse{Data: []Webhook{}}
	resp, err := s.request(ctx).SetQueryParams(params.Build()).GetPaginated(constants.EndpointWebhooks, func(data []byte) error {
		var page WebhooksResponse
		if err := json.Unmarshal(data, &page); err != nil {
			return fmt.Errorf("failed to unmarshal webhook page: %w", err)
		}
		result.Data = append(result.Data, page.Data...)
		result.Links, result.Meta = page.Links, page.Meta
		return nil
	})
	if err != nil {
		return nil, resp, err
	}
	return &result, resp, nil
}

// GetByWebhookIDV1 reads a webhook. Only opts.Fields applies to this operation.
// https://developer.apple.com/documentation/applebusinessapi/get-webhook-information
func (s *Webhooks) GetByWebhookIDV1(ctx context.Context, webhookID string, opts *RequestQueryOptions) (*WebhookResponse, *resty.Response, error) {
	if webhookID == "" {
		return nil, nil, fmt.Errorf("webhook ID is required")
	}
	if opts == nil {
		opts = &RequestQueryOptions{}
	}
	params := s.client.QueryBuilder().AddStringSlice("fields[webhooks]", opts.Fields)
	var result WebhookResponse
	resp, err := s.request(ctx).SetQueryParams(params.Build()).SetResult(&result).Get(webhookPath(webhookID))
	if err != nil {
		return nil, resp, err
	}
	return &result, resp, nil
}

// CreateV1 creates a webhook. The secret is write-only.
// https://developer.apple.com/documentation/applebusinessapi/create-a-webhook
func (s *Webhooks) CreateV1(ctx context.Context, req *WebhookCreateRequest) (*WebhookResponse, *resty.Response, error) {
	if req == nil {
		return nil, nil, fmt.Errorf("request is required")
	}
	a := req.Data.Attributes
	if req.Data.Type == "" || a.Name == "" || a.URL == "" || a.Secret == "" || len(a.EventTypes) == 0 {
		return nil, nil, fmt.Errorf("type, name, url, secret and eventTypes are required")
	}
	var result WebhookResponse
	resp, err := s.request(ctx).SetBody(req).SetResult(&result).Post(constants.EndpointWebhooks)
	if err != nil {
		return nil, resp, err
	}
	return &result, resp, nil
}

// UpdateByWebhookIDV1 updates only provided attributes, including explicit false values.
// https://developer.apple.com/documentation/applebusinessapi/update-a-webhook
func (s *Webhooks) UpdateByWebhookIDV1(ctx context.Context, webhookID string, req *WebhookUpdateRequest) (*WebhookResponse, *resty.Response, error) {
	if webhookID == "" {
		return nil, nil, fmt.Errorf("webhook ID is required")
	}
	if req == nil {
		return nil, nil, fmt.Errorf("request is required")
	}
	if req.Data.Type == "" || req.Data.ID == "" {
		return nil, nil, fmt.Errorf("type and request ID are required")
	}
	if req.Data.ID != webhookID {
		return nil, nil, fmt.Errorf("request ID must match webhook ID")
	}
	var result WebhookResponse
	resp, err := s.request(ctx).SetBody(req).SetResult(&result).Patch(webhookPath(webhookID))
	if err != nil {
		return nil, resp, err
	}
	return &result, resp, nil
}

// DeleteByWebhookIDV1 deletes a webhook (204 No Content).
// https://developer.apple.com/documentation/applebusinessapi/delete-a-webhook
func (s *Webhooks) DeleteByWebhookIDV1(ctx context.Context, webhookID string) (*resty.Response, error) {
	if webhookID == "" {
		return nil, fmt.Errorf("webhook ID is required")
	}
	return s.request(ctx).Delete(webhookPath(webhookID))
}

// GetDeliveriesByWebhookIDV1 fetches all delivery pages and their included events.
// https://developer.apple.com/documentation/applebusinessapi/get-all-deliveries-for-a-webhook
func (s *Webhooks) GetDeliveriesByWebhookIDV1(ctx context.Context, webhookID string, opts *DeliveryQueryOptions) (*WebhookDeliveriesResponse, *resty.Response, error) {
	if webhookID == "" {
		return nil, nil, fmt.Errorf("webhook ID is required")
	}
	if opts == nil {
		opts = &DeliveryQueryOptions{}
	}
	params := s.client.QueryBuilder().AddStringSlice("fields[webhookDeliveries]", opts.Fields).
		AddStringSlice("fields[webhookEvents]", opts.EventFields).AddStringSlice("include", opts.Include).
		AddStringSlice("filter[deliveryState]", opts.FilterDeliveryState).
		AddStringSlice("filter[createdDateTimeAtOrAfter]", opts.FilterCreatedDateTimeAtOrAfter).
		AddStringSlice("filter[createdDateTimeAtBefore]", opts.FilterCreatedDateTimeAtBefore).
		AddInt("limit", min(opts.Limit, 1000)).AddString("cursor", opts.Cursor)
	result := WebhookDeliveriesResponse{Data: []WebhookDelivery{}}
	seenEvents := make(map[string]bool)
	resp, err := s.request(ctx).SetQueryParams(params.Build()).GetPaginated(webhookPath(webhookID)+"/deliveries", func(data []byte) error {
		var page WebhookDeliveriesResponse
		if err := json.Unmarshal(data, &page); err != nil {
			return fmt.Errorf("failed to unmarshal delivery page: %w", err)
		}
		result.Data = append(result.Data, page.Data...)
		for _, event := range page.Included {
			key := event.Type + "/" + event.ID
			if !seenEvents[key] {
				result.Included = append(result.Included, event)
				seenEvents[key] = true
			}
		}
		result.Links, result.Meta = page.Links, page.Meta
		return nil
	})
	if err != nil {
		return nil, resp, err
	}
	return &result, resp, nil
}

// CreateDeliveryV1 replays the delivery identified by relationships.template.
// https://developer.apple.com/documentation/applebusinessapi/create-a-webhookdelivery
func (s *Webhooks) CreateDeliveryV1(ctx context.Context, req *WebhookDeliveryCreateRequest) (*WebhookDeliveryResponse, *resty.Response, error) {
	if req == nil {
		return nil, nil, fmt.Errorf("request is required")
	}
	link := req.Data.Relationships.Template.Data
	if req.Data.Type == "" || link.ID == "" || link.Type == "" {
		return nil, nil, fmt.Errorf("type and template type and ID are required")
	}
	var result WebhookDeliveryResponse
	resp, err := s.request(ctx).SetBody(req).SetResult(&result).Post(constants.EndpointWebhookDeliveries)
	if err != nil {
		return nil, resp, err
	}
	return &result, resp, nil
}

// CreatePingV1 sends a test notification to the related webhook.
// https://developer.apple.com/documentation/applebusinessapi/create-a-webhookping
func (s *Webhooks) CreatePingV1(ctx context.Context, req *WebhookPingCreateRequest) (*WebhookPingResponse, *resty.Response, error) {
	if req == nil {
		return nil, nil, fmt.Errorf("request is required")
	}
	link := req.Data.Relationships.Webhook.Data
	if req.Data.Type == "" || link.ID == "" || link.Type == "" {
		return nil, nil, fmt.Errorf("type and webhook type and ID are required")
	}
	var result WebhookPingResponse
	resp, err := s.request(ctx).SetBody(req).SetResult(&result).Post(constants.EndpointWebhookPings)
	if err != nil {
		return nil, resp, err
	}
	return &result, resp, nil
}
