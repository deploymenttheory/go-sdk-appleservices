package webhooks

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/deploymenttheory/go-sdk-appleservices/axm/client"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"resty.dev/v3"
)

const baseURL = "https://api-business.apple.com"

func setup(t *testing.T) *Webhooks {
	t.Helper()
	tr, err := client.NewTransport("key", "issuer", "unused", client.WithAuth(client.NewAPIKeyAuth("test", "")), client.WithRetryCount(0))
	require.NoError(t, err)
	httpmock.ActivateNonDefault(tr.GetHTTPClient().Client())
	t.Cleanup(func() { httpmock.DeactivateAndReset(); _ = tr.Close() })
	return NewService(tr)
}

func fixture(t *testing.T, name string) string {
	t.Helper()
	b, err := os.ReadFile("testdata/" + name + ".json")
	require.NoError(t, err)
	return string(b)
}

func createRequest() *WebhookCreateRequest {
	return &WebhookCreateRequest{Data: WebhookCreateRequestData{Type: ResourceTypeWebhooks, Attributes: WebhookCreateRequestAttributes{Name: "Test webhook", URL: "https://example.com/apple", Secret: "test-secret", EventTypes: []string{EventTypeOrgDeviceAttributeChanged}}}}
}

func updateRequest() *WebhookUpdateRequest {
	disabled := false
	return &WebhookUpdateRequest{Data: WebhookUpdateRequestData{Type: ResourceTypeWebhooks, ID: "webhook-1", Attributes: &WebhookUpdateRequestAttributes{Enabled: &disabled}}}
}

func deliveryRequest() *WebhookDeliveryCreateRequest {
	return &WebhookDeliveryCreateRequest{Data: WebhookDeliveryCreateRequestData{Type: ResourceTypeWebhookDeliveries, Relationships: WebhookDeliveryCreateRelationships{Template: ResourceRelationship{Data: ResourceIdentifier{Type: ResourceTypeWebhookDeliveries, ID: "delivery-1"}}}}}
}

func pingRequest() *WebhookPingCreateRequest {
	return &WebhookPingCreateRequest{Data: WebhookPingCreateRequestData{Type: ResourceTypeWebhookPings, Relationships: WebhookPingCreateRelationships{Webhook: ResourceRelationship{Data: ResourceIdentifier{Type: ResourceTypeWebhooks, ID: "webhook-1"}}}}}
}

func TestWebhookOperations(t *testing.T) {
	s := setup(t)
	ctx := context.Background()
	for _, tc := range []struct {
		method, path, body string
		status             int
	}{
		{"GET", "/v1/webhooks", `{"data":[]}`, 200},
		{"GET", "/v1/webhooks/webhook-1", fixture(t, "webhook"), 200},
		{"POST", "/v1/webhooks", fixture(t, "webhook"), 201},
		{"PATCH", "/v1/webhooks/webhook-1", fixture(t, "webhook"), 200},
		{"DELETE", "/v1/webhooks/webhook-1", "", 204},
		{"GET", "/v1/webhooks/webhook-1/deliveries", `{"data":[]}`, 200},
		{"POST", "/v1/webhookDeliveries", fixture(t, "delivery"), 201},
		{"POST", "/v1/webhookPings", `{"data":{"type":"webhookPings","id":"ping-1"}}`, 201},
	} {
		httpmock.RegisterResponder(tc.method, baseURL+tc.path, func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, "application/json", req.Header.Get("Accept"))
			assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
			if tc.method == "POST" || tc.method == "PATCH" {
				b, err := io.ReadAll(req.Body)
				require.NoError(t, err)
				switch tc.path {
				case "/v1/webhooks":
					assert.JSONEq(t, `{"data":{"type":"webhooks","attributes":{"name":"Test webhook","enabled":false,"eventTypes":["ORG_DEVICE_ATTRIBUTE_CHANGED"],"secret":"test-secret","url":"https://example.com/apple"}}}`, string(b))
				case "/v1/webhooks/webhook-1":
					assert.JSONEq(t, `{"data":{"type":"webhooks","id":"webhook-1","attributes":{"enabled":false}}}`, string(b))
				case "/v1/webhookDeliveries":
					assert.JSONEq(t, `{"data":{"type":"webhookDeliveries","relationships":{"template":{"data":{"type":"webhookDeliveries","id":"delivery-1"}}}}}`, string(b))
				case "/v1/webhookPings":
					assert.JSONEq(t, `{"data":{"type":"webhookPings","relationships":{"webhook":{"data":{"type":"webhooks","id":"webhook-1"}}}}}`, string(b))
				}
			}
			resp := httpmock.NewStringResponse(tc.status, tc.body)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		})
	}
	list, _, err := s.GetV1(ctx, nil)
	require.NoError(t, err)
	assert.Empty(t, list.Data)
	item, _, err := s.GetByWebhookIDV1(ctx, "webhook-1", nil)
	require.NoError(t, err)
	assert.Equal(t, "Test webhook", item.Data.Attributes.Name)
	assert.False(t, item.Data.Attributes.Enabled)
	require.NotNil(t, item.Data.Attributes.CreatedDateTime)
	assert.NotEmpty(t, item.Data.Relationships.Deliveries.Links.Related)
	_, resp, err := s.CreateV1(ctx, createRequest())
	require.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode())
	_, _, err = s.UpdateByWebhookIDV1(ctx, "webhook-1", updateRequest())
	require.NoError(t, err)
	resp, err = s.DeleteByWebhookIDV1(ctx, "webhook-1")
	require.NoError(t, err)
	assert.Equal(t, 204, resp.StatusCode())
	deliveries, _, err := s.GetDeliveriesByWebhookIDV1(ctx, "webhook-1", nil)
	require.NoError(t, err)
	assert.Empty(t, deliveries.Data)
	delivery, _, err := s.CreateDeliveryV1(ctx, deliveryRequest())
	require.NoError(t, err)
	assert.Equal(t, 200, delivery.Data.Attributes.Response.HTTPStatusCode)
	assert.True(t, delivery.Data.Attributes.Redelivery)
	assert.Equal(t, "event-1", delivery.Data.Relationships.Event.Data.ID)
	assert.True(t, delivery.Included[0].Attributes.IsPingEvent)
	assert.JSONEq(t, `{"data":{"type":"webhookPingCreated","id":"ping-1"}}`, delivery.Included[0].Attributes.Payload)
	ping, _, err := s.CreatePingV1(ctx, pingRequest())
	require.NoError(t, err)
	assert.Equal(t, "ping-1", ping.Data.ID)
	assert.Equal(t, 8, httpmock.GetTotalCallCount())
}

func TestPaginationAndQueries(t *testing.T) {
	s := setup(t)
	httpmock.RegisterResponder("GET", baseURL+"/v1/webhooks", func(req *http.Request) (*http.Response, error) {
		q := req.URL.Query()
		assert.Equal(t, "1000", q.Get("limit"))
		assert.Equal(t, "name,enabled", q.Get("fields[webhooks]"))
		if q.Get("cursor") == "initial" {
			return httpmock.NewJsonResponse(200, map[string]any{"data": []any{map[string]any{"type": "webhooks", "id": "1"}}, "links": map[string]string{"next": baseURL + "/v1/webhooks?cursor=next"}})
		}
		assert.Equal(t, "next", q.Get("cursor"))
		return httpmock.NewJsonResponse(200, map[string]any{"data": []any{map[string]any{"type": "webhooks", "id": "2"}}, "meta": map[string]any{"paging": map[string]int{"limit": 1000}}})
	})
	opts := &RequestQueryOptions{Fields: []string{FieldName, FieldEnabled}, Limit: 2000, Cursor: "initial"}
	res, _, err := s.GetV1(context.Background(), opts)
	require.NoError(t, err)
	require.Len(t, res.Data, 2)
	assert.Equal(t, 2000, opts.Limit)
	assert.Equal(t, 1000, res.Meta.Paging.Limit)
	httpmock.RegisterResponder("GET", baseURL+"/v1/webhooks/webhook-1", func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "name,enabled", req.URL.Query().Get("fields[webhooks]"))
		assert.Empty(t, req.URL.Query().Get("cursor"))
		assert.Empty(t, req.URL.Query().Get("limit"))
		return httpmock.NewJsonResponse(200, map[string]any{"data": map[string]string{"id": "webhook-1", "type": "webhooks"}})
	})
	item, _, err := s.GetByWebhookIDV1(context.Background(), "webhook-1", opts)
	require.NoError(t, err)
	assert.Nil(t, item.Data.Attributes)
	httpmock.RegisterResponder("GET", baseURL+"/v1/webhooks/webhook-1/deliveries", func(req *http.Request) (*http.Response, error) {
		q := req.URL.Query()
		for k, v := range map[string]string{"fields[webhookDeliveries]": "deliveryState,event", "fields[webhookEvents]": "payload,isPingEvent", "include": "event", "filter[deliveryState]": "SUCCEEDED,PENDING", "filter[createdDateTimeAtOrAfter]": "2026-10-01T00:00:00Z", "filter[createdDateTimeAtBefore]": "2026-11-01T00:00:00Z", "limit": "1000"} {
			assert.Equal(t, v, q.Get(k), k)
		}
		var f WebhookDeliveryResponse
		require.NoError(t, json.Unmarshal([]byte(fixture(t, "delivery")), &f))
		page := WebhookDeliveriesResponse{Data: []WebhookDelivery{f.Data}, Included: f.Included}
		if q.Get("cursor") == "initial" {
			page.Links = &Links{Next: baseURL + "/v1/webhooks/webhook-1/deliveries?cursor=next"}
		} else {
			assert.Equal(t, "next", q.Get("cursor"))
			page.Data[0].ID = "delivery-2"
			page.Included = append(page.Included, WebhookEvent{Type: ResourceTypeWebhookEvents, ID: "event-2"})
		}
		return httpmock.NewJsonResponse(200, page)
	})
	d, _, err := s.GetDeliveriesByWebhookIDV1(context.Background(), "webhook-1", &DeliveryQueryOptions{Fields: []string{FieldDeliveryState, FieldDeliveryEvent}, EventFields: []string{FieldEventPayload, FieldEventIsPingEvent}, Include: []string{IncludeEvent}, FilterDeliveryState: []string{DeliveryStateSucceeded, DeliveryStatePending}, FilterCreatedDateTimeAtOrAfter: []string{"2026-10-01T00:00:00Z"}, FilterCreatedDateTimeAtBefore: []string{"2026-11-01T00:00:00Z"}, Limit: 2000, Cursor: "initial"})
	require.NoError(t, err)
	assert.Len(t, d.Data, 2)
	assert.Len(t, d.Included, 2)
}

func operations(s *Webhooks) []struct {
	name, method, path string
	run                func(context.Context) (*resty.Response, error)
} {
	return []struct {
		name, method, path string
		run                func(context.Context) (*resty.Response, error)
	}{
		{"list", "GET", "/v1/webhooks", func(ctx context.Context) (*resty.Response, error) { _, r, e := s.GetV1(ctx, nil); return r, e }},
		{"read", "GET", "/v1/webhooks/webhook-1", func(ctx context.Context) (*resty.Response, error) {
			_, r, e := s.GetByWebhookIDV1(ctx, "webhook-1", nil)
			return r, e
		}},
		{"create", "POST", "/v1/webhooks", func(ctx context.Context) (*resty.Response, error) {
			_, r, e := s.CreateV1(ctx, createRequest())
			return r, e
		}},
		{"update", "PATCH", "/v1/webhooks/webhook-1", func(ctx context.Context) (*resty.Response, error) {
			_, r, e := s.UpdateByWebhookIDV1(ctx, "webhook-1", updateRequest())
			return r, e
		}},
		{"delete", "DELETE", "/v1/webhooks/webhook-1", func(ctx context.Context) (*resty.Response, error) { return s.DeleteByWebhookIDV1(ctx, "webhook-1") }},
		{"deliveries", "GET", "/v1/webhooks/webhook-1/deliveries", func(ctx context.Context) (*resty.Response, error) {
			_, r, e := s.GetDeliveriesByWebhookIDV1(ctx, "webhook-1", nil)
			return r, e
		}},
		{"redeliver", "POST", "/v1/webhookDeliveries", func(ctx context.Context) (*resty.Response, error) {
			_, r, e := s.CreateDeliveryV1(ctx, deliveryRequest())
			return r, e
		}},
		{"ping", "POST", "/v1/webhookPings", func(ctx context.Context) (*resty.Response, error) {
			_, r, e := s.CreatePingV1(ctx, pingRequest())
			return r, e
		}},
	}
}

func TestErrorsAndCancellation(t *testing.T) {
	s := setup(t)
	for _, op := range operations(s) {
		t.Run(op.name, func(t *testing.T) {
			httpmock.RegisterResponder(op.method, baseURL+op.path, httpmock.NewJsonResponderOrPanic(403, map[string]any{"errors": []any{map[string]string{"status": "403", "code": "FORBIDDEN", "detail": "denied"}}}))
			r, err := op.run(context.Background())
			require.Error(t, err)
			require.NotNil(t, r)
			var apiErr *client.APIError
			require.ErrorAs(t, err, &apiErr)
			assert.Equal(t, "FORBIDDEN", apiErr.Code)
			before := httpmock.GetTotalCallCount()
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			_, err = op.run(ctx)
			assert.ErrorIs(t, err, context.Canceled)
			assert.Equal(t, before, httpmock.GetTotalCallCount())
		})
	}
}

func TestDeliveryQueryErrorSource(t *testing.T) {
	s := setup(t)
	httpmock.RegisterResponder("GET", baseURL+"/v1/webhooks/webhook-1/deliveries", httpmock.NewJsonResponderOrPanic(400, map[string]any{
		"errors": []any{map[string]any{"status": "400", "code": "PARAMETER_ERROR.REQUIRED", "detail": "Required parameter missing.", "source": map[string]string{"parameter": "filter[createdDateTimeAtOrAfter]"}}},
	}))
	_, resp, err := s.GetDeliveriesByWebhookIDV1(context.Background(), "webhook-1", nil)
	require.NotNil(t, resp)
	assert.Equal(t, 400, resp.StatusCode())
	var apiErr *client.APIError
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, "PARAMETER_ERROR.REQUIRED", apiErr.Code)
	require.NotNil(t, apiErr.Source)
	require.NotNil(t, apiErr.Source.Parameter)
	assert.Equal(t, "filter[createdDateTimeAtOrAfter]", apiErr.Source.Parameter.Parameter)
}

func TestValidation(t *testing.T) {
	s := setup(t)
	ctx := context.Background()
	cases := []func() error{
		func() error { _, _, e := s.GetByWebhookIDV1(ctx, "", nil); return e }, func() error { _, _, e := s.CreateV1(ctx, nil); return e }, func() error { _, _, e := s.CreateV1(ctx, &WebhookCreateRequest{}); return e },
		func() error { _, _, e := s.UpdateByWebhookIDV1(ctx, "", nil); return e }, func() error { _, _, e := s.UpdateByWebhookIDV1(ctx, "webhook-1", nil); return e }, func() error { _, _, e := s.UpdateByWebhookIDV1(ctx, "webhook-1", &WebhookUpdateRequest{}); return e }, func() error { _, _, e := s.UpdateByWebhookIDV1(ctx, "other", updateRequest()); return e },
		func() error { _, e := s.DeleteByWebhookIDV1(ctx, ""); return e }, func() error { _, _, e := s.GetDeliveriesByWebhookIDV1(ctx, "", nil); return e },
		func() error { _, _, e := s.CreateDeliveryV1(ctx, nil); return e }, func() error { _, _, e := s.CreateDeliveryV1(ctx, &WebhookDeliveryCreateRequest{}); return e },
		func() error { _, _, e := s.CreatePingV1(ctx, nil); return e }, func() error { _, _, e := s.CreatePingV1(ctx, &WebhookPingCreateRequest{}); return e },
	}
	for _, run := range cases {
		require.Error(t, run())
	}
	assert.Zero(t, httpmock.GetTotalCallCount())
	for _, field := range []string{"name", "url", "secret", "events", "type"} {
		req := createRequest()
		switch field {
		case "name":
			req.Data.Attributes.Name = ""
		case "url":
			req.Data.Attributes.URL = ""
		case "secret":
			req.Data.Attributes.Secret = ""
		case "events":
			req.Data.Attributes.EventTypes = nil
		case "type":
			req.Data.Type = ""
		}
		_, _, err := s.CreateV1(ctx, req)
		require.Error(t, err)
	}
}

func TestMalformedPagination(t *testing.T) {
	s := setup(t)
	for _, path := range []string{"webhooks", "webhooks/webhook-1/deliveries"} {
		httpmock.RegisterResponder("GET", baseURL+"/v1/"+path, func(*http.Request) (*http.Response, error) {
			r := httpmock.NewStringResponse(200, `{"data":`)
			r.Header.Set("Content-Type", "application/json")
			return r, nil
		})
	}
	_, _, err := s.GetV1(context.Background(), nil)
	require.ErrorContains(t, err, "unmarshal")
	_, _, err = s.GetDeliveriesByWebhookIDV1(context.Background(), "webhook-1", nil)
	require.ErrorContains(t, err, "unmarshal")
}

func TestPatchOmissionAndExplicitEmptyValues(t *testing.T) {
	empty := ""
	events := []string{}
	encoded, err := json.Marshal(WebhookUpdateRequest{Data: WebhookUpdateRequestData{
		Type: ResourceTypeWebhooks, ID: "webhook-1",
		Attributes: &WebhookUpdateRequestAttributes{Name: &empty, URL: &empty, Secret: &empty, EventTypes: &events},
	}})
	require.NoError(t, err)
	assert.JSONEq(t, `{"data":{"type":"webhooks","id":"webhook-1","attributes":{"name":"","url":"","secret":"","eventTypes":[]}}}`, string(encoded))
	encoded, err = json.Marshal(WebhookUpdateRequest{Data: WebhookUpdateRequestData{Type: ResourceTypeWebhooks, ID: "webhook-1"}})
	require.NoError(t, err)
	assert.JSONEq(t, `{"data":{"type":"webhooks","id":"webhook-1"}}`, string(encoded))
}

func TestCancellationBetweenServicePages(t *testing.T) {
	for _, deliveries := range []bool{false, true} {
		t.Run(map[bool]string{false: "webhooks", true: "deliveries"}[deliveries], func(t *testing.T) {
			s := setup(t)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			path := "/v1/webhooks"
			if deliveries {
				path += "/webhook-1/deliveries"
			}
			httpmock.RegisterResponder("GET", baseURL+path, func(*http.Request) (*http.Response, error) {
				cancel()
				return httpmock.NewJsonResponse(200, map[string]any{"data": []any{}, "links": map[string]string{"next": baseURL + path + "?cursor=next"}})
			})
			var err error
			if deliveries {
				_, _, err = s.GetDeliveriesByWebhookIDV1(ctx, "webhook-1", nil)
			} else {
				_, _, err = s.GetV1(ctx, nil)
			}
			assert.ErrorIs(t, err, context.Canceled)
			assert.Equal(t, 1, httpmock.GetTotalCallCount())
		})
	}
}

func TestExpiredDeadline(t *testing.T) {
	s := setup(t)
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
	defer cancel()
	for _, op := range operations(s) {
		_, err := op.run(ctx)
		assert.ErrorIs(t, err, context.DeadlineExceeded, op.name)
	}
	assert.Zero(t, httpmock.GetTotalCallCount())
}
