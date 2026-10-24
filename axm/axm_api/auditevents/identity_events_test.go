package auditevents

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Every 2.6 union variant must survive decoding, including explicitly empty objects.
// These synthetic payloads follow Apple's AuditEvent schemas; they contain no live data.
func TestGetAuditEvents_IdentityPayloads(t *testing.T) {
	fixture, err := os.ReadFile("mocks/validate_get_identity_audit_events.json")
	require.NoError(t, err)
	var raw struct {
		Data []json.RawMessage `json:"data"`
	}
	require.NoError(t, json.Unmarshal(fixture, &raw))
	cases := []struct {
		eventType, field string
		want             *AuditEventAttributes
	}{
		{AuditEventTypeExternalAccountInvited, FieldEventDataExternalAccountInvited, &AuditEventAttributes{EventDataExternalAccountInvited: &EventDataExternalAccountInvited{}}},
		{AuditEventTypeAccountPasswordChanged, FieldEventDataAccountPasswordChanged, &AuditEventAttributes{EventDataAccountPasswordChanged: &EventDataAccountPasswordChanged{ActivityID: "password-change-1"}}},
		{AuditEventTypeBatchPasswordResetInitiated, FieldEventDataBatchPasswordResetInitiated, &AuditEventAttributes{EventDataBatchPasswordResetInitiated: &EventDataBatchPasswordResetInitiated{ActivityID: "password-reset-1"}}},
		{AuditEventTypeAccountImageSet, FieldEventDataAccountImageSet, &AuditEventAttributes{EventDataAccountImageSet: &EventDataAccountImageSet{}}},
		{AuditEventTypeAccountImageReplaced, FieldEventDataAccountImageReplaced, &AuditEventAttributes{EventDataAccountImageReplaced: &EventDataAccountImageReplaced{}}},
		{AuditEventTypeAccountImageRemoved, FieldEventDataAccountImageRemoved, &AuditEventAttributes{EventDataAccountImageRemoved: &EventDataAccountImageRemoved{}}},
		{AuditEventTypeAppleServicesSettingsUpdated, FieldEventDataAppleServicesSettingsUpdated, &AuditEventAttributes{EventDataAppleServicesSettingsUpdated: &EventDataAppleServicesSettingsUpdated{SettingsEnabled: []string{"ICLOUD"}, SettingsDisabled: []string{"FACETIME"}}}},
		{AuditEventTypeRoleCreated, FieldEventDataRoleCreated, &AuditEventAttributes{EventDataRoleCreated: &EventDataRoleCreated{Privileges: []string{"VIEW_DEVICES"}}}},
		{AuditEventTypeRoleUpdated, FieldEventDataRoleUpdated, &AuditEventAttributes{EventDataRoleUpdated: &EventDataRoleUpdated{PrivilegesAdded: []string{"VIEW_AUDIT_EVENTS"}, PrivilegesRemoved: []string{"MANAGE_DEVICES"}}}},
		{AuditEventTypeTermsAccepted, FieldEventDataTermsAccepted, &AuditEventAttributes{EventDataTermsAccepted: &EventDataTermsAccepted{TermsName: "Example Terms"}}},
		{AuditEventTypeIDPCreated, FieldEventDataIdpCreated, &AuditEventAttributes{EventDataIdpCreated: &EventDataIdpCreated{VendorName: IDPVendorMicrosoftOIDC, Protocol: IDPProtocolOIDC}}},
		{AuditEventTypeIDPUpdated, FieldEventDataIdpUpdated, &AuditEventAttributes{EventDataIdpUpdated: &EventDataIdpUpdated{VendorName: IDPVendorGoogle}}},
		{AuditEventTypeIDPDeleted, FieldEventDataIdpDeleted, &AuditEventAttributes{EventDataIdpDeleted: &EventDataIdpDeleted{VendorName: IDPVendorGenericIDP}}},
		{AuditEventTypeIDPAccountSyncEnabled, FieldEventDataIdpAccountSyncEnabled, &AuditEventAttributes{EventDataIdpAccountSyncEnabled: &EventDataIdpAccountSyncEnabled{VendorName: IDPVendorMicrosoft, ActivityID: "sync-enable-1"}}},
		{AuditEventTypeIDPAccountSyncDisabled, FieldEventDataIdpAccountSyncDisabled, &AuditEventAttributes{EventDataIdpAccountSyncDisabled: &EventDataIdpAccountSyncDisabled{VendorName: IDPVendorMicrosoft, ActivityID: "sync-disable-1"}}},
		{AuditEventTypeDomainFederationEnabled, FieldEventDataDomainFederationEnabled, &AuditEventAttributes{EventDataDomainFederationEnabled: &EventDataDomainFederationEnabled{ActivityID: "federation-enable-1", SCIMActivityID: "scim-1"}}},
		{AuditEventTypeDomainFederationDisabled, FieldEventDataDomainFederationDisabled, &AuditEventAttributes{EventDataDomainFederationDisabled: &EventDataDomainFederationDisabled{ActivityID: "federation-disable-1"}}},
	}
	require.Len(t, raw.Data, len(cases))
	for i, tc := range cases {
		t.Run(tc.eventType, func(t *testing.T) {
			svc := setupMockClient(t)
			httpmock.RegisterResponder(http.MethodGet, "https://api-business.apple.com/v1/auditEvents", func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, tc.eventType, req.URL.Query().Get("filter[type]"))
				assert.Equal(t, "type,eventDataPropertyKey,"+tc.field, req.URL.Query().Get("fields[auditEvents]"))
				assert.Equal(t, "2026-10-24T00:00:00Z", req.URL.Query().Get("filter[startTimestamp]"))
				assert.Equal(t, "2026-10-24T01:00:00Z", req.URL.Query().Get("filter[endTimestamp]"))
				return httpmock.NewStringResponse(200, `{"data":[`+string(raw.Data[i])+`]}`), nil
			})
			result, resp, err := svc.GetV1(context.Background(), &RequestQueryOptions{FilterStartTimestamp: "2026-10-24T00:00:00Z", FilterEndTimestamp: "2026-10-24T01:00:00Z", FilterType: tc.eventType, Fields: []string{FieldType, FieldEventDataPropertyKey, tc.field}})
			require.NoError(t, err)
			assert.Equal(t, 200, resp.StatusCode())
			require.Len(t, result.Data, 1)
			tc.want.Type = tc.eventType
			tc.want.EventDataPropertyKey = tc.field
			assert.Equal(t, tc.want, result.Data[0].Attributes)
			encoded, err := json.Marshal(result.Data[0])
			require.NoError(t, err)
			assert.JSONEq(t, string(raw.Data[i]), string(encoded))
		})
	}
}

func TestGetAuditEvents_IdentityPayloadsAcrossPages(t *testing.T) {
	svc := setupMockClient(t)
	fixture, err := os.ReadFile("mocks/validate_get_identity_audit_events.json")
	require.NoError(t, err)
	var raw struct {
		Data []json.RawMessage `json:"data"`
	}
	require.NoError(t, json.Unmarshal(fixture, &raw))
	httpmock.RegisterResponder(http.MethodGet, "https://api-business.apple.com/v1/auditEvents", func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "2026-10-24T00:00:00Z", req.URL.Query().Get("filter[startTimestamp]"))
		if req.URL.Query().Get("cursor") == "next-page" {
			return httpmock.NewJsonResponse(200, map[string]any{"data": raw.Data[8:]})
		}
		return httpmock.NewJsonResponse(200, map[string]any{"data": raw.Data[:8], "links": map[string]string{"next": "https://api-business.apple.com/v1/auditEvents?cursor=next-page"}})
	})
	result, _, err := svc.GetV1(context.Background(), &RequestQueryOptions{FilterStartTimestamp: "2026-10-24T00:00:00Z", FilterEndTimestamp: "2026-10-24T01:00:00Z"})
	require.NoError(t, err)
	require.Len(t, result.Data, 17)
	encoded, err := json.Marshal(result.Data)
	require.NoError(t, err)
	expected, err := json.Marshal(raw.Data)
	require.NoError(t, err)
	assert.JSONEq(t, string(expected), string(encoded))
	assert.Equal(t, 2, httpmock.GetTotalCallCount())
}

func TestIdentityProviderProtocolSAML(t *testing.T) {
	var payload EventDataIdpCreated
	require.NoError(t, json.Unmarshal([]byte(`{"vendorName":"GENERIC_IDP","protocol":"SAML"}`), &payload))
	assert.Equal(t, IDPVendorGenericIDP, payload.VendorName)
	assert.Equal(t, IDPProtocolSAML, payload.Protocol)
}
