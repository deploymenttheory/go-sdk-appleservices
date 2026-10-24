package devices

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetActivationLockStatusByDeviceIDV1(t *testing.T) {
	for _, tc := range []struct {
		name, attributes, lockType string
		locked                     *bool
	}{
		{"MDM", `{"isLocked":true,"lockType":"MDM"}`, ActivationLockTypeMDM, boolPointer(true)},
		{"User", `{"isLocked":true,"lockType":"USER"}`, ActivationLockTypeUser, boolPointer(true)},
		{"Unlocked", `{"isLocked":false,"lockType":"NONE"}`, ActivationLockTypeNone, boolPointer(false)},
		{"OmittedLockFlag", `{"lockType":"NONE"}`, ActivationLockTypeNone, nil},
		{"FutureLockType", `{"isLocked":true,"lockType":"FUTURE"}`, "FUTURE", boolPointer(true)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			svc := setupMockClient(t)
			endpoint := "https://api-business.apple.com/v1/orgDevices/device-123/activationLockStatus"
			httpmock.RegisterResponder(http.MethodGet, endpoint, func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, "application/json", req.Header.Get("Accept"))
				assert.Empty(t, req.URL.RawQuery)
				return httpmock.NewJsonResponse(200, json.RawMessage(`{"data":{"id":"device-123","type":"activationLockStatus","attributes":`+tc.attributes+`},"links":{"self":"`+endpoint+`"}}`))
			})
			result, resp, err := svc.GetActivationLockStatusByDeviceIDV1(context.Background(), "device-123")
			require.NoError(t, err)
			require.NotNil(t, resp)
			assert.Equal(t, 200, resp.StatusCode())
			require.NotNil(t, result)
			assert.Equal(t, "device-123", result.Data.ID)
			assert.Equal(t, "activationLockStatus", result.Data.Type)
			require.NotNil(t, result.Data.Attributes)
			assert.Equal(t, tc.locked, result.Data.Attributes.IsLocked)
			assert.Equal(t, tc.lockType, result.Data.Attributes.LockType)
			require.NotNil(t, result.Links)
			assert.Equal(t, endpoint, result.Links.Self)
			encoded, err := json.Marshal(result.Data.Attributes)
			require.NoError(t, err)
			assert.JSONEq(t, tc.attributes, string(encoded))
			assert.Equal(t, 1, httpmock.GetTotalCallCount())
		})
	}
}

func TestGetActivationLockStatusByDeviceIDV1_EmptyID(t *testing.T) {
	svc := setupMockClient(t)
	result, resp, err := svc.GetActivationLockStatusByDeviceIDV1(context.Background(), "")
	require.EqualError(t, err, "device ID is required")
	assert.Nil(t, result)
	assert.Nil(t, resp)
	assert.Zero(t, httpmock.GetTotalCallCount())
}

func TestGetActivationLockStatusByDeviceIDV1_HTTPError(t *testing.T) {
	for _, status := range []int{403, 404, 500} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			svc := setupMockClient(t)
			httpmock.RegisterResponder(http.MethodGet, "https://api-business.apple.com/v1/orgDevices/device-123/activationLockStatus", activationLockJSONResponder(status, `{"errors":[{"code":"TEST_ERROR","title":"Test failure"}]}`))
			result, resp, err := svc.GetActivationLockStatusByDeviceIDV1(context.Background(), "device-123")
			require.Error(t, err)
			assert.Nil(t, result)
			require.NotNil(t, resp)
			assert.Equal(t, status, resp.StatusCode())
		})
	}
}

func TestGetActivationLockStatusByDeviceIDV1_OptionalAttributes(t *testing.T) {
	svc := setupMockClient(t)
	httpmock.RegisterResponder(http.MethodGet, "https://api-business.apple.com/v1/orgDevices/device-123/activationLockStatus", activationLockJSONResponder(200, `{"data":{"id":"device-123","type":"activationLockStatus"}}`))
	result, _, err := svc.GetActivationLockStatusByDeviceIDV1(context.Background(), "device-123")
	require.NoError(t, err)
	assert.Nil(t, result.Data.Attributes)
	assert.Nil(t, result.Links)
}

func TestGetActivationLockStatusByDeviceIDV1_InvalidResponse(t *testing.T) {
	svc := setupMockClient(t)
	httpmock.RegisterResponder(http.MethodGet, "https://api-business.apple.com/v1/orgDevices/device-123/activationLockStatus", activationLockJSONResponder(200, `{"data":{"attributes":{"isLocked":"true"}}}`))
	result, _, err := svc.GetActivationLockStatusByDeviceIDV1(context.Background(), "device-123")
	require.Error(t, err)
	assert.Nil(t, result)
}

func boolPointer(value bool) *bool { return &value }

func activationLockJSONResponder(status int, body string) httpmock.Responder {
	return func(*http.Request) (*http.Response, error) {
		return httpmock.NewJsonResponse(status, json.RawMessage(body))
	}
}
