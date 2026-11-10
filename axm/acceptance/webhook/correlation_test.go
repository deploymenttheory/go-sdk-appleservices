package webhook

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/deploymenttheory/go-sdk-appleservices/axm/axm_api/webhooks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSucceededPingDelivery(t *testing.T) {
	var received notification
	require.NoError(t, json.Unmarshal([]byte(`{"data":{"id":"ping-1","type":"webhookPingCreated","attributes":{"createdDateTime":"2026-11-10T12:00:00Z"}}}`), &received))
	for _, tc := range []struct {
		name   string
		change func(*webhooks.WebhookDeliveriesResponse, *notification)
		want   bool
	}{
		{name: "stored template omits ping ID", want: true},
		{name: "stored payload has matching ID", change: func(r *webhooks.WebhookDeliveriesResponse, _ *notification) {
			r.Included[0].Attributes.Payload = `{"data":{"id":"ping-1","type":"webhookPingCreated"}}`
		}, want: true},
		{name: "different explicit ID", change: func(r *webhooks.WebhookDeliveriesResponse, _ *notification) {
			r.Included[0].Attributes.Payload = `{"data":{"id":"other","type":"webhookPingCreated","attributes":{"createdDateTime":"2026-11-10T12:00:00Z"}}}`
		}},
		{name: "different timestamp", change: func(_ *webhooks.WebhookDeliveriesResponse, n *notification) {
			n.Data.Attributes.CreatedDateTime = n.Data.Attributes.CreatedDateTime.Add(time.Second)
		}},
		{name: "missing timestamps", change: func(r *webhooks.WebhookDeliveriesResponse, n *notification) {
			r.Included[0].Attributes.Payload = `{"data":{"id":"","type":"webhookPingCreated"}}`
			n.Data.Attributes.CreatedDateTime = time.Time{}
		}},
		{name: "wrong notification type", change: func(_ *webhooks.WebhookDeliveriesResponse, n *notification) { n.Data.Type = "other" }},
		{name: "invalid JSON", change: func(r *webhooks.WebhookDeliveriesResponse, _ *notification) {
			r.Included[0].Attributes.Payload = "invalid"
		}},
		{name: "non ping event", change: func(r *webhooks.WebhookDeliveriesResponse, _ *notification) {
			r.Included[0].Attributes.IsPingEvent = false
		}},
		{name: "unrelated event", change: func(r *webhooks.WebhookDeliveriesResponse, _ *notification) {
			r.Data[0].Relationships.Event.Data.ID = "other"
		}},
		{name: "pending delivery", change: func(r *webhooks.WebhookDeliveriesResponse, _ *notification) {
			r.Data[0].Attributes.DeliveryState = "PENDING"
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			result := &webhooks.WebhookDeliveriesResponse{
				Data: []webhooks.WebhookDelivery{{
					ID: "delivery-1", Attributes: &webhooks.WebhookDeliveryAttributes{DeliveryState: webhooks.DeliveryStateSucceeded},
					Relationships: &webhooks.WebhookDeliveryRelationships{Event: &webhooks.ResourceRelationship{Data: webhooks.ResourceIdentifier{ID: "event-1"}}},
				}},
				Included: []webhooks.WebhookEvent{{
					ID: "event-1", Attributes: &webhooks.WebhookEventAttributes{IsPingEvent: true,
						Payload: `{"data":{"id":"","type":"webhookPingCreated","attributes":{"createdDateTime":"2026-11-10T12:00:00Z"}}}`},
				}},
			}
			note := received
			if tc.change != nil {
				tc.change(result, &note)
			}
			delivery, ok := succeededPingDelivery(result, note)
			assert.Equal(t, tc.want, ok)
			if tc.want {
				assert.Equal(t, "delivery-1", delivery.ID)
			}
		})
	}
}
