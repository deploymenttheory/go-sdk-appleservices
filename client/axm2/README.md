# AXM2 Client - Direct Pattern Apple Business/School Manager API Client

A simplified, direct-pattern client for Apple Business Manager (ABM) and Apple School Manager (ASM) APIs built with **Resty v3**.

## Key Features

✅ **Direct Client Pattern** - No complex interfaces or service layers  
✅ **Context-First** - All operations take `context.Context`  
✅ **Built-in Pagination** - Automatic pagination handling with `SetResult()`  
✅ **Smart Authentication** - JWT with automatic refresh using Resty v3 middleware  
✅ **Intelligent Retry** - Exponential backoff with Apple AXM-specific conditions  
✅ **Structured Errors** - Proper error wrapping with `SetError()`  
✅ **Custom Content Handling** - Apple AXM-specific JSON encoders/decoders  
✅ **Zero Interfaces** - Simple, concrete client implementation  
✅ **Resty v3 Patterns** - Modern middleware and automatic unmarshaling  

## Quick Start

```go
package main

import (
    "context"
    "log"
    
    "github.com/deploymenttheory/go-api-sdk-apple/client/axm2"
)

func main() {
    // Create client with direct configuration (Resty v3)
    client, err := axm2.NewClient(axm2.Config{
        APIType:    axm2.APITypeABM, // or axm2.APITypeASM
        ClientID:   "BUSINESSAPI.your-client-id",
        KeyID:      "your-key-id",
        PrivateKey: "-----BEGIN PRIVATE KEY-----\n...\n-----END PRIVATE KEY-----",
        Debug:      true,
    })
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close() // Resty v3 requires explicit close

    ctx := context.Background()

    // Get devices with field selection (context-first pattern)
    devices, err := client.GetOrgDevices(ctx, 
        client.NewQueryBuilder().
            Limit(50).
            Fields("orgDevices", []string{"serialNumber", "deviceModel", "status"}),
    )
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Found %d devices", len(devices))
}
```

## Configuration from File

```go
// Load from JSON config file
config, err := axm2.LoadConfigFromFile("config.json")
if err != nil {
    log.Fatal(err)
}

client, err := axm2.NewClient(config)
```

Example `config.json`:
```json
{
  "api_type": "abm",
  "client_id": "BUSINESSAPI.your-client-id",
  "key_id": "your-key-id", 
  "private_key": "-----BEGIN PRIVATE KEY-----\n...\n-----END PRIVATE KEY-----",
  "debug": true
}
```

## Available Methods

### Device Management
- `GetOrgDevices(ctx, queryBuilder)` - Get all devices with pagination
- `GetOrgDevice(ctx, deviceID, queryBuilder)` - Get specific device

### Query Building
```go
queryBuilder := client.NewQueryBuilder().
    Limit(100).                             // Limit results (max 1000)
    Fields("orgDevices", []string{...}).    // Select specific fields (note: plural "orgDevices")
    Include([]string{"assignedServer"})     // Include relationships

// Note: Apple AXM API doesn't support filter or sort parameters for devices
// Available fields: serialNumber, addedToOrgDateTime, updatedDateTime, deviceModel,
// productFamily, productType, deviceCapacity, partNumber, orderNumber, color, 
// status, orderDateTime, imei, meid, eid, wifiMacAddress, bluetoothMacAddress,
// purchaseSourceId, purchaseSourceType, assignedServer
```

### Authentication
- `IsAuthenticated()` - Check if token is valid
- `ForceReauthenticate()` - Force token refresh

## Architecture Benefits

This client follows the **Direct Client Pattern** used by modern SDKs like AWS SDK v2 and Stripe:

1. **No Interface Complexity** - Direct methods on concrete client
2. **Context Everywhere** - Proper cancellation and timeout support  
3. **Built-in Pagination** - No manual pagination logic needed
4. **Smart Retry** - Automatic token refresh on 401 errors
5. **Type Safety** - Strongly typed responses with proper JSON mapping

## Migration from AXM v1

**Old Pattern (v1):**
```go
// Complex service layer
axmClient, _ := client.NewAXMClient(config)
axmService := axm.NewClient(axmClient) 
devices, _ := axmService.GetOrgDevices(queryBuilder)
```

**New Pattern (v2 with Resty v3):**
```go  
// Direct client with Resty v3 patterns
client, _ := axm2.NewClient(config)
defer client.Close() // v3 requires explicit close
devices, _ := client.GetOrgDevices(ctx, queryBuilder)
```

## Resty v3 Improvements

This client leverages [Resty v3](https://resty.dev/) features:

- **Middleware System**: Uses `AddRequestMiddleware()` and `AddResponseMiddleware()` for token injection
- **Automatic Unmarshaling**: `SetResult()` and `SetError()` for type-safe responses  
- **Smart Retry Logic**: `AddRetryConditions()` and `AddRetryHooks()` for intelligent retries
- **Content-Type Handling**: `AddContentTypeEncoder()` and `AddContentTypeDecoder()` for custom formats
- **Resource Management**: Explicit `client.Close()` for proper cleanup
- **Enhanced Performance**: v3 offers improved memory efficiency over v2

## Retry Mechanism

The client includes intelligent retry logic based on [Resty v3's retry mechanism](https://resty.dev/docs/retry-mechanism/):

### Default Retry Conditions
- **401 Unauthorized**: Automatically refreshes JWT token and retries
- **429 Too Many Requests**: Respects rate limiting with exponential backoff  
- **5xx Server Errors**: Retries on server errors (except 501 Not Implemented)
- **Exponential Backoff**: Default 1s min wait, 10s max wait, 3 retries

### Configuration
```go
client, _ := axm2.NewClient(axm2.Config{
    RetryCount:     5,                    // Number of retries
    RetryMinWait:   2 * time.Second,      // Minimum wait between retries
    RetryMaxWait:   30 * time.Second,     // Maximum wait between retries  
    EnableRetryLog: true,                 // Log retry attempts
})

// Or override per use case
client.SetRetryConfig(5, 2*time.Second, 30*time.Second)
```

### Retry Hooks
The client automatically:
- Refreshes JWT tokens on 401 errors before retry
- Logs retry attempts with detailed context
- Respects `Retry-After` headers from Apple's API

## Content-Type Handling

The client includes custom content-type encoders and decoders based on [Resty v3's extensible content handling](https://resty.dev/docs/content-type-encoder-and-decoder/):

### Built-in Handlers
- **JSON Encoding**: Optimized for Apple APIs (no HTML escaping, compact format)
- **JSON Decoding**: Strict parsing with unknown field detection
- **Text/Plain**: Handles Apple's plain text error responses  
- **Problem+JSON**: RFC 7807 Problem Details format support

### Custom Content Types
```go
// Add custom encoder for a specific content-type
client.AddContentTypeEncoder("application/vnd.apple.axm+json", func(w io.Writer, v any) error {
    // Custom encoding logic for Apple's proprietary format
    return json.NewEncoder(w).Encode(v)
})

// Add custom decoder for Apple's custom response formats
client.AddContentTypeDecoder("application/vnd.apple.error+json", func(r io.Reader, v any) error {
    // Custom decoding logic for Apple's error format
    return json.NewDecoder(r).Decode(v)
})
```

### Error Response Handling
The client automatically handles various Apple error response formats:
- Standard JSON API errors
- Plain text error messages  
- RFC 7807 Problem Details
- Mixed content-type responses

## Examples

See the `examples/axm2/` directory for complete working examples:
- `basic_usage/` - Simple device retrieval
- `config_file_usage/` - Loading configuration from JSON file
