# Apple Business Manager API Examples

This directory contains comprehensive examples for using the Apple Business Manager API SDK.

## 📁 Directory Structure

### Device Management Examples
- **`GetDeviceManagementServices/`** - Get MDM servers in your organization
- **`GetMDMServerDeviceLinkages/`** - Get device IDs assigned to an MDM server
- **`GetAssignedDeviceManagementServiceID/`** - Get assigned server ID for a device
- **`GetAssignedDeviceManagementServiceInfo/`** - Get assigned server information for a device
- **`AssignDevicesToServer/`** - Assign devices to an MDM server
- **`UnassignDevicesFromServer/`** - Unassign devices from an MDM server

### Device Examples
- **`GetOrganizationDevices/`** - Get devices in your organization
- **`GetDeviceInformation/`** - Get detailed information for a specific device

## 🚀 Quick Start

### Prerequisites

1. **Apple Business Manager Account** with API access
2. **Private Key** (`.p8` file) from Apple Business Manager
3. **Key ID** and **Issuer ID** from your Apple Business Manager account

### Environment Setup

Set these environment variables:
```bash
export APPLE_KEY_ID="your-key-id"
export APPLE_ISSUER_ID="your-issuer-id" 
export APPLE_PRIVATE_KEY_PATH="/path/to/your/private-key.p8"
```

### Running Examples

Each example can be run independently:

```bash
# Device Management Examples
cd GetDeviceManagementServices && go run main.go
cd GetMDMServerDeviceLinkages && go run main.go
cd GetAssignedDeviceManagementServiceID && go run main.go
cd GetAssignedDeviceManagementServiceInfo && go run main.go
cd AssignDevicesToServer && go run main.go
cd UnassignDevicesFromServer && go run main.go

# Device Examples  
cd GetOrganizationDevices && go run main.go
cd GetDeviceInformation && go run main.go
```

## 🔧 Client Builder Features

The `axm.ClientBuilder` (located in `client/axm/client_builder.go`) provides a fluent interface for configuring **generic AXM clients**:

```go
import (
    "github.com/deploymenttheory/go-api-sdk-apple/client/axm"
    "github.com/deploymenttheory/go-api-sdk-apple/services/axm/devices"
    "github.com/deploymenttheory/go-api-sdk-apple/services/axm/devicemanagement"
)

// Step 1: Build generic AXM client
axmClient, err := axm.NewClientBuilder().
    WithJWTAuthFromEnv().                    // Load auth from environment
    WithBaseURL("https://api-business.apple.com/v1"). // Custom base URL
    WithTimeout(60 * time.Second).           // Request timeout
    WithRetry(3, 2 * time.Second).          // Retry configuration
    WithDebug(true).                        // Enable debug logging
    WithUserAgent("MyApp/1.0.0").          // Custom user agent
    Build()

// Step 2: Create service-specific clients
devicesClient := devices.NewClient(axmClient)
deviceManagementClient := devicemanagement.NewClient(axmClient)
```

### Architecture

- **`client/axm/client_builder.go`** - Generic AXM client builder (creates `axm.Client`)
- **`client/axm/crypto.go`** - RSA/ECDSA private key helpers
- **Service clients** - Wrap the generic `axm.Client` for specific API endpoints

### Public Convenience Functions

The client builder also provides public convenience functions for common patterns:

```go
// Simple client creation from environment variables
client, err := axm.NewClientFromEnv()

// Client creation from files
client, err := axm.NewClientFromFile(keyID, issuerID, privateKeyPath)

// Client with custom options from environment
client, err := axm.NewClientFromEnvWithOptions(
    true,                    // debug
    60 * time.Second,       // timeout
    "MyApp/1.0.0",         // user agent
)

// Client with custom options from files
client, err := axm.NewClientFromFileWithOptions(
    keyID, issuerID, privateKeyPath,  // credentials
    true,                             // debug
    60 * time.Second,                // timeout
    "MyApp/1.0.0",                  // user agent
)
```

## 📋 Example Features

Each example demonstrates:

### ✅ Core Functionality
- **Authentication** - OAuth 2.0 with JWT client assertions
- **Error Handling** - Comprehensive error checking and logging
- **Field Selection** - Using field constants for optimal API requests
- **Pagination** - Handling paginated responses
- **JSON Output** - Pretty-printed API responses

### 🎯 Device Management Examples

#### GetDeviceManagementServices
- Get all MDM servers with default/custom options
- Field selection and pagination
- Error handling for invalid parameters

#### GetMDMServerDeviceLinkages  
- Get device IDs assigned to specific MDM servers
- Pagination with different limits
- URL parameter handling

#### GetAssignedDeviceManagementServiceID
- Check server assignments for individual devices
- Bulk checking across multiple devices
- Assignment status summary

#### GetAssignedDeviceManagementServiceInfo
- Get detailed server information for assigned devices
- Field selection for performance optimization
- Server name resolution and comparison

#### AssignDevicesToServer
- Single and multiple device assignment
- Assignment verification and status tracking
- Error handling for invalid IDs

#### UnassignDevicesFromServer
- Single and multiple device unassignment  
- Unassignment verification and status tracking
- Bulk unassignment operations

### 🔍 Device Examples

#### GetOrganizationDevices
- Get all devices with various field selections
- Pagination demonstration
- Device filtering and limits

#### GetDeviceInformation
- Get detailed information for specific devices
- Field selection optimization
- Multiple device information retrieval

## 🔐 Authentication

All examples support both **RSA** and **ECDSA** private keys and implement the full **OAuth 2.0 client credentials flow** as required by Apple Business Manager API.

### Key Features:
- **Token Caching** - Automatic access token management
- **Token Refresh** - Automatic token renewal before expiration
- **ES256 Signing** - ECDSA signature support (preferred by Apple)
- **RS256 Fallback** - RSA signature support

## 🛠️ Error Handling

Examples include comprehensive error handling for:
- **Invalid Credentials** - Wrong key ID, issuer ID, or private key
- **Invalid Parameters** - Empty/invalid device IDs, server IDs
- **API Errors** - HTTP errors, rate limiting, server errors
- **Network Issues** - Connection timeouts, network failures

## 📊 Pagination

Examples demonstrate Apple's pagination system:
- **Cursor-based Pagination** - Using `nextCursor` for large datasets
- **Limit Parameters** - Controlling page size (max 1000)
- **Link Following** - Using `next` URLs for subsequent pages
- **Metadata Handling** - Processing pagination metadata

## 🎨 Output Examples

Each example provides:
- **Structured Output** - Clear, formatted console output
- **JSON Responses** - Pretty-printed API responses
- **Status Tracking** - Progress indicators and summaries
- **Error Messages** - Detailed error information

## 📝 Notes

- **Rate Limiting**: Apple Business Manager API has rate limits. Examples include appropriate delays.
- **Processing Time**: Device assignments/unassignments may take time to process in Apple's system.
- **Verification**: Always verify operations in the Apple Business Manager portal.
- **Testing**: Examples use real API calls - test with non-production data when possible.

## 🔗 Related Documentation

- [Apple Business Manager API Documentation](https://developer.apple.com/documentation/applebusinessmanagerapi)
- [OAuth 2.0 Client Credentials Flow](https://developer.apple.com/documentation/applebusinessmanagerapi/implementing_oauth_for_the_apple_business_manager_api)
- [SDK Documentation](../../README.md)

---

**Happy coding!** 🚀 These examples provide a solid foundation for integrating with Apple Business Manager API.
