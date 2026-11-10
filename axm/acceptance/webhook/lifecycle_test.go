package webhook

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/deploymenttheory/go-sdk-appleservices/axm"
	acc "github.com/deploymenttheory/go-sdk-appleservices/axm/acceptance"
	"github.com/deploymenttheory/go-sdk-appleservices/axm/axm_api/webhooks"
	"github.com/deploymenttheory/go-sdk-appleservices/axm/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var liveContext context.Context

func TestMain(m *testing.M) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	liveContext = ctx
	code := m.Run()
	stop()
	os.Exit(code)
}

func TestAcceptance_WebhookLifecycle(t *testing.T) {
	if os.Getenv("AXM_WEBHOOK_E2E") != "true" {
		t.Skip("set AXM_WEBHOOK_E2E=true to create a temporary lab webhook")
	}
	acc.SkipIfNotConfigured(t)
	ctx, cancel := context.WithTimeout(liveContext, 8*time.Minute)
	t.Cleanup(cancel)
	options := []axm.ClientOption{axm.WithTimeout(acc.Config.RequestTimeout), axm.WithRetryCount(0)}
	if token := os.Getenv("AXM_WEBHOOK_ACCESS_TOKEN"); token != "" {
		options = append(options, client.WithAuth(client.NewAPIKeyAuth(token, "")))
	}
	c, err := axm.NewClientFromEnv(options...)
	require.NoError(t, err)
	svc := c.AXMAPI.Webhooks
	// Check account permissions before starting infrastructure or creating resources.
	_, _, err = svc.GetV1(ctx, &webhooks.RequestQueryOptions{Limit: 100})
	require.NoError(t, err, "webhook API preflight failed")
	secret := rand.Text()
	receiver := newReceiver(secret)
	origin := startReceiver(t, receiver)
	callback := startTunnel(t, ctx, origin) + "/apple"
	name := "go-sdk-webhook-" + time.Now().UTC().Format("20060102T150405") + "-" + rand.Text()[:8]
	deliveryStart := time.Now().UTC().Add(-time.Minute).Format(time.RFC3339)
	created, resp, err := svc.CreateV1(ctx, &webhooks.WebhookCreateRequest{Data: webhooks.WebhookCreateRequestData{Type: webhooks.ResourceTypeWebhooks, Attributes: webhooks.WebhookCreateRequestAttributes{Name: name, Enabled: true, URL: callback, Secret: secret, EventTypes: []string{webhooks.EventTypeOrgDeviceAttributeChanged}}}})
	require.NoError(t, err)
	require.NotEmpty(t, created.Data.ID)
	id := created.Data.ID
	deleted := false
	t.Cleanup(func() {
		cleanupCtx, stop := context.WithTimeout(context.Background(), 30*time.Second)
		defer stop()
		if !deleted {
			_, cleanupErr := svc.DeleteByWebhookIDV1(cleanupCtx, id)
			if cleanupErr != nil && !axm.IsNotFound(cleanupErr) {
				t.Errorf("test webhook cleanup failed: %v", cleanupErr)
				return
			}
		}
		_, _, readErr := svc.GetByWebhookIDV1(cleanupCtx, id, nil)
		if !axm.IsNotFound(readErr) {
			t.Errorf("test webhook deletion unconfirmed: %v", readErr)
		} else {
			t.Log("Cleanup: test webhook deletion confirmed (404)")
		}
	})
	require.Equal(t, 201, resp.StatusCode())
	t.Log("CreateV1: 201")
	listed, resp, err := svc.GetV1(ctx, &webhooks.RequestQueryOptions{Fields: []string{webhooks.FieldName, webhooks.FieldEnabled}, Limit: 1})
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode())
	found := false
	for _, item := range listed.Data {
		if item.ID == id {
			found = true
		}
	}
	require.True(t, found, "new test webhook must be listed")
	t.Log("GetV1: 200; test webhook found")
	read, resp, err := svc.GetByWebhookIDV1(ctx, id, &webhooks.RequestQueryOptions{Fields: []string{webhooks.FieldName, webhooks.FieldEnabled, webhooks.FieldURL}})
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode())
	require.NotNil(t, read.Data.Attributes)
	assert.Equal(t, callback, read.Data.Attributes.URL)
	t.Log("GetByWebhookIDV1: 200")
	// Verify an explicit false update, then re-enable the test subscription for pings.
	for _, enabled := range []bool{false, true} {
		updated, r, e := svc.UpdateByWebhookIDV1(ctx, id, &webhooks.WebhookUpdateRequest{Data: webhooks.WebhookUpdateRequestData{Type: webhooks.ResourceTypeWebhooks, ID: id, Attributes: &webhooks.WebhookUpdateRequestAttributes{Enabled: &enabled}}})
		require.NoError(t, e)
		require.Equal(t, 200, r.StatusCode())
		require.NotNil(t, updated.Data.Attributes)
		assert.Equal(t, enabled, updated.Data.Attributes.Enabled)
	}
	t.Log("UpdateByWebhookIDV1: 200; disabled and re-enabled")
	ping, resp, err := svc.CreatePingV1(ctx, &webhooks.WebhookPingCreateRequest{Data: webhooks.WebhookPingCreateRequestData{Type: webhooks.ResourceTypeWebhookPings, Relationships: webhooks.WebhookPingCreateRelationships{Webhook: webhooks.ResourceRelationship{Data: webhooks.ResourceIdentifier{Type: webhooks.ResourceTypeWebhooks, ID: id}}}}})
	require.NoError(t, err)
	require.Equal(t, 201, resp.StatusCode())
	require.NotEmpty(t, ping.Data.ID)
	pingCtx, stopPing := context.WithTimeout(ctx, 90*time.Second)
	defer stopPing()
	note, err := receiver.wait(pingCtx, 0, func(n notification) bool { return n.Data.Type == "webhookPingCreated" })
	require.NoError(t, err, "signed ping was not received")
	t.Log("CreatePingV1: 201; signed webhookPingCreated received")
	var original webhooks.WebhookDelivery
	require.NoError(t, poll(pingCtx, func() (bool, error) {
		result, r, e := svc.GetDeliveriesByWebhookIDV1(pingCtx, id, &webhooks.DeliveryQueryOptions{
			Include: []string{webhooks.IncludeEvent}, Limit: 1,
			Fields:                         []string{webhooks.FieldDeliveryState, webhooks.FieldDeliveryEvent},
			EventFields:                    []string{webhooks.FieldEventPayload, webhooks.FieldEventIsPingEvent},
			FilterCreatedDateTimeAtOrAfter: []string{deliveryStart},
			FilterCreatedDateTimeAtBefore:  []string{time.Now().UTC().Format(time.RFC3339)},
		})
		if e != nil {
			return false, deliveryPollingError(r, e)
		}
		if r.StatusCode() != 200 {
			return false, nil
		}
		for _, event := range result.Included {
			if event.Attributes == nil || !event.Attributes.IsPingEvent {
				continue
			}
			var payload notification
			if json.Unmarshal([]byte(event.Attributes.Payload), &payload) != nil || payload.Data.ID != note.Data.ID {
				continue
			}
			for _, d := range result.Data {
				if d.Relationships != nil && d.Relationships.Event != nil && d.Relationships.Event.Data.ID == event.ID && d.Attributes != nil && d.Attributes.DeliveryState == webhooks.DeliveryStateSucceeded {
					original = d
					return true, nil
				}
			}
		}
		return false, nil
	}), "ping delivery did not reach SUCCEEDED")
	t.Log("GetDeliveriesByWebhookIDV1: 200; ping event correlated and SUCCEEDED")
	before := receiver.count()
	replayed, resp, err := svc.CreateDeliveryV1(ctx, &webhooks.WebhookDeliveryCreateRequest{Data: webhooks.WebhookDeliveryCreateRequestData{Type: webhooks.ResourceTypeWebhookDeliveries, Relationships: webhooks.WebhookDeliveryCreateRelationships{Template: webhooks.ResourceRelationship{Data: webhooks.ResourceIdentifier{Type: webhooks.ResourceTypeWebhookDeliveries, ID: original.ID}}}}})
	require.NoError(t, err)
	require.Equal(t, 201, resp.StatusCode())
	require.NotEmpty(t, replayed.Data.ID)
	replayCtx, stopReplay := context.WithTimeout(ctx, 90*time.Second)
	defer stopReplay()
	_, err = receiver.wait(replayCtx, before, func(n notification) bool { return n.Data.ID == note.Data.ID })
	require.NoError(t, err, "signed redelivery was not received")
	require.NoError(t, poll(replayCtx, func() (bool, error) {
		result, r, e := svc.GetDeliveriesByWebhookIDV1(replayCtx, id, &webhooks.DeliveryQueryOptions{
			Fields:                         []string{webhooks.FieldDeliveryState, webhooks.FieldDeliveryRedelivery, webhooks.FieldDeliveryEvent},
			FilterCreatedDateTimeAtOrAfter: []string{deliveryStart},
			FilterCreatedDateTimeAtBefore:  []string{time.Now().UTC().Format(time.RFC3339)},
		})
		if e != nil {
			return false, deliveryPollingError(r, e)
		}
		for _, d := range result.Data {
			if d.ID == replayed.Data.ID && d.Attributes != nil && d.Attributes.Redelivery && d.Attributes.DeliveryState == webhooks.DeliveryStateSucceeded && d.Relationships != nil && d.Relationships.Event != nil && d.Relationships.Event.Data.ID == original.Relationships.Event.Data.ID {
				return true, nil
			}
		}
		return false, nil
	}), "redelivery did not reach SUCCEEDED for the original event")
	t.Log("CreateDeliveryV1: 201; signed redelivery received and SUCCEEDED")
	resp, err = svc.DeleteByWebhookIDV1(ctx, id)
	require.NoError(t, err)
	require.Equal(t, 204, resp.StatusCode())
	deleted = true
	t.Log("DeleteByWebhookIDV1: 204")
}
