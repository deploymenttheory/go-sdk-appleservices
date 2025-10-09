# Apple Business Manager API - Quick Start Guide

This guide shows you how to get started with the Apple Business Manager API SDK using different client setup methods.

## 📋 Prerequisites

Before you begin, you'll need:

1. **Apple Business Manager Account** with API access
2. **Private Key** (`.p8` file) downloaded from Apple Business Manager
3. **Key ID** and **Issuer ID** from your Apple Business Manager account

## 🔐 Authentication Setup

The SDK supports both **RSA** and **ECDSA** private keys and implements the full **OAuth 2.0 client credentials flow** required by Apple Business Manager API.

### Get Your Credentials

1. **Log into Apple Business Manager**
2. **Go to Settings → API Keys**
3. **Create a new API key** or use an existing one
4. **Download the private key** (`.p8` file)
5. **Note your Key ID and Issuer ID**

## 🚀 Client Setup Methods

The SDK provides multiple ways to create and configure clients, from simple one-liners to fully customized builders.

### Method 1: Environment Variables (Recommended)

The simplest approach - set environment variables and use the convenience function:

```bash
export APPLE_KEY_ID="your-key-id"
export APPLE_ISSUER_ID="your-issuer-id" 
export APPLE_PRIVATE_KEY_PATH="/path/to/your/private-key.p8"
```

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/deploymenttheory/go-api-sdk-apple/client/axm"
    "github.com/deploymenttheory/go-api-sdk-apple/services/axm/devices"
)

func main() {
    // Create client from environment variables
    axmClient, err := axm.NewClientFromEnv()
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }

    // Create service client
    devicesClient := devices.NewClient(axmClient)

    // Use the client
    ctx := context.Background()
    response, err := devicesClient.GetOrganizationDevices(ctx, nil)
    if err != nil {
        log.Fatalf("Error getting devices: %v", err)
    }

    fmt.Printf("Found %d devices\n", len(response.Data))
}
```

### Method 2: Direct File Paths

Specify credentials directly without environment variables:

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/deploymenttheory/go-api-sdk-apple/client/axm"
    "github.com/deploymenttheory/go-api-sdk-apple/services/axm/devices"
)

func main() {
    // Create client with direct file paths
    axmClient, err := axm.NewClientFromFile(
        "your-key-id",            // Key ID
        "your-issuer-id",         // Issuer ID  
        "/path/to/private-key.p8", // Private key path
    )
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }

    // Create service client
    devicesClient := devices.NewClient(axmClient)

    // Use the client
    ctx := context.Background()
    response, err := devicesClient.GetOrganizationDevices(ctx, nil)
    if err != nil {
        log.Fatalf("Error getting devices: %v", err)
    }

    fmt.Printf("Found %d devices\n", len(response.Data))
}
```

### Method 3: Environment Variables with Custom Options

Use environment variables but customize client behavior:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/deploymenttheory/go-api-sdk-apple/client/axm"
    "github.com/deploymenttheory/go-api-sdk-apple/services/axm/devices"
)

func main() {
    // Create client from environment with custom options
    axmClient, err := axm.NewClientFromEnvWithOptions(
        true,                             // Enable debug logging
        60*time.Second,                   // Request timeout
        "MyApp/1.0.0",                   // Custom user agent
    )
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }

    // Create service client
    devicesClient := devices.NewClient(axmClient)

    // Use the client
    ctx := context.Background()
    response, err := devicesClient.GetOrganizationDevices(ctx, nil)
    if err != nil {
        log.Fatalf("Error getting devices: %v", err)
    }

    fmt.Printf("Found %d devices\n", len(response.Data))
}
```

### Method 4: File Paths with Custom Options

Specify everything directly with custom options:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/deploymenttheory/go-api-sdk-apple/client/axm"
    "github.com/deploymenttheory/go-api-sdk-apple/services/axm/devices"
)

func main() {
    // Create client from files with custom options
    axmClient, err := axm.NewClientFromFileWithOptions(
        "your-key-id",                    // Key ID
        "your-issuer-id",                 // Issuer ID
        "/path/to/private-key.p8",        // Private key path
        true,                             // Enable debug logging
        60*time.Second,                   // Request timeout
        "MyApp/1.0.0",                   // Custom user agent
    )
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }

    // Create service client
    devicesClient := devices.NewClient(axmClient)

    // Use the client
    ctx := context.Background()
    response, err := devicesClient.GetOrganizationDevices(ctx, nil)
    if err != nil {
        log.Fatalf("Error getting devices: %v", err)
    }

    fmt.Printf("Found %d devices\n", len(response.Data))
}
```

### Method 5: Full Builder Pattern (Advanced)

For maximum control, use the full builder pattern:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/deploymenttheory/go-api-sdk-apple/client/axm"
    "github.com/deploymenttheory/go-api-sdk-apple/services/axm/devices"
    "go.uber.org/zap"
)

func main() {
    // Create custom logger
    logger, _ := zap.NewDevelopment()

    // Build client with full control
    axmClient, err := axm.NewClientBuilder().
        WithJWTAuthFromEnv().                              // Load auth from environment
        WithBaseURL("https://api-business.apple.com/v1").  // Custom base URL
        WithTimeout(60*time.Second).                       // Request timeout
        WithRetry(3, 2*time.Second).                      // Retry configuration
        WithDebug(true).                                  // Enable debug logging
        WithUserAgent("MyApp/1.0.0").                    // Custom user agent
        WithLogger(logger).                               // Custom logger
        Build()

    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }

    // Create service client
    devicesClient := devices.NewClient(axmClient)

    // Use the client
    ctx := context.Background()
    response, err := devicesClient.GetOrganizationDevices(ctx, nil)
    if err != nil {
        log.Fatalf("Error getting devices: %v", err)
    }

    fmt.Printf("Found %d devices\n", len(response.Data))
}
```

### Method 6: Hardcoded Credentials (Testing Only)

For testing or when you have the private key as a string:

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/deploymenttheory/go-api-sdk-apple/client/axm"
    "github.com/deploymenttheory/go-api-sdk-apple/services/axm/devices"
)

func main() {
    // Your credentials (for testing only)
    keyID := "your-key-id"
    issuerID := "your-issuer-id"
    privateKeyPEM := `-----BEGIN EC PRIVATE KEY-----
MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQg...
-----END EC PRIVATE KEY-----`

    // Parse the private key
    privateKey, err := axm.ParsePrivateKey([]byte(privateKeyPEM))
    if err != nil {
        log.Fatalf("Failed to parse private key: %v", err)
    }

    // Create client with hardcoded credentials
    axmClient, err := axm.NewClientBuilder().
        WithJWTAuth(keyID, issuerID, privateKey).
        WithDebug(true).
        Build()

    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }

    // Create service client
    devicesClient := devices.NewClient(axmClient)

    // Use the client
    ctx := context.Background()
    response, err := devicesClient.GetOrganizationDevices(ctx, nil)
    if err != nil {
        log.Fatalf("Error getting devices: %v", err)
    }

    fmt.Printf("Found %d devices\n", len(response.Data))
}
```

## 🏗️ Architecture Overview

The SDK follows a two-tier architecture:

### 1. Generic AXM Client (`axm.Client`)
- Handles authentication (OAuth 2.0 with JWT)
- Manages HTTP requests and responses
- Provides pagination, retry logic, and error handling
- Shared across all Apple Business Manager services

### 2. Service-Specific Clients
- **`devices.Client`** - Organization devices operations
- **`devicemanagement.Client`** - MDM server and assignment operations
- Each wraps the generic `axm.Client` for specific API endpoints

```go
// Step 1: Create generic AXM client
axmClient, err := axm.NewClientFromEnv()

// Step 2: Create service-specific clients
devicesClient := devices.NewClient(axmClient)
deviceManagementClient := devicemanagement.NewClient(axmClient)
```

## 🔧 Configuration Options

### Builder Methods

| Method | Description | Default |
|--------|-------------|---------|
| `WithJWTAuth(keyID, issuerID, privateKey)` | Set credentials directly | Required |
| `WithJWTAuthFromFile(keyID, issuerID, path)` | Load credentials from file | Required |
| `WithJWTAuthFromEnv()` | Load credentials from environment | Required |
| `WithBaseURL(url)` | Set API base URL | `https://api-business.apple.com/v1` |
| `WithTimeout(duration)` | Set request timeout | `30s` |
| `WithRetry(count, wait)` | Set retry configuration | `3 retries, 1s wait` |
| `WithUserAgent(agent)` | Set user agent string | `go-api-sdk-apple/1.0.0` |
| `WithDebug(enabled)` | Enable debug logging | `false` |
| `WithLogger(logger)` | Set custom zap logger | Auto-configured |

### Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `APPLE_KEY_ID` | Your Apple API Key ID | ✅ |
| `APPLE_ISSUER_ID` | Your Apple Issuer ID | ✅ |
| `APPLE_PRIVATE_KEY_PATH` | Path to your `.p8` private key file | ✅ |

## 📱 Service Examples

### Devices Service

```go
// Get all devices
devices, err := devicesClient.GetOrganizationDevices(ctx, nil)

// Get devices with filtering
devices, err := devicesClient.GetOrganizationDevices(ctx, &devices.GetOrganizationDevicesOptions{
    Fields: []string{
        devices.FieldSerialNumber,
        devices.FieldDeviceModel,
        devices.FieldStatus,
    },
    Limit: 100,
})

// Get specific device information
device, err := devicesClient.GetDeviceInformationByDeviceID(ctx, "device-id", nil)
```

### Device Management Service

```go
// Get MDM servers
servers, err := dmClient.GetDeviceManagementServices(ctx, nil)

// Get devices assigned to a server
linkages, err := dmClient.GetMDMServerDeviceLinkages(ctx, "server-id", nil)

// Assign devices to server
activity, err := dmClient.AssignDevicesToServer(ctx, "server-id", []string{"device-id-1", "device-id-2"})

// Unassign devices from server
activity, err := dmClient.UnassignDevicesFromServer(ctx, "server-id", []string{"device-id-1"})
```

## 🔍 Error Handling

The SDK provides comprehensive error handling:

```go
devices, err := devicesClient.GetOrganizationDevices(ctx, nil)
if err != nil {
    // Handle different types of errors
    switch {
    case strings.Contains(err.Error(), "authentication"):
        log.Printf("Authentication error: %v", err)
    case strings.Contains(err.Error(), "not found"):
        log.Printf("Resource not found: %v", err)
    default:
        log.Printf("API error: %v", err)
    }
    return
}
```

## 📄 Pagination

Handle paginated responses:

```go
// Get first page
response, err := devicesClient.GetOrganizationDevices(ctx, &devices.GetOrganizationDevicesOptions{
    Limit: 10,
})

// Check for more pages
if axm.HasNextPage(response.Links) {
    fmt.Printf("Next page available: %s\n", response.Links.Next)
}

// Use the client's GetAllPages method for automatic pagination
err = axmClient.GetAllPages(ctx, "/orgDevices", map[string]string{"limit": "10"}, nil, 
    func(pageData []byte) error {
        var page devices.OrgDevicesResponse
        if err := json.Unmarshal(pageData, &page); err != nil {
            return err
        }
        // Process page data
        fmt.Printf("Page has %d devices\n", len(page.Data))
        return nil
    })
```

## 🎯 Best Practices

1. **Use Environment Variables** - Keep credentials secure and out of code
2. **Enable Debug Logging** - Use `WithDebug(true)` during development
3. **Set Appropriate Timeouts** - Configure timeouts based on your use case
4. **Handle Errors Gracefully** - Always check and handle errors appropriately
5. **Use Field Selection** - Request only the fields you need for better performance
6. **Implement Pagination** - Handle large datasets with proper pagination
7. **Context Usage** - Always pass context for timeout and cancellation support

## 🔗 Next Steps

- Explore the [complete examples](./README.md) for detailed usage scenarios
- Check out the [API documentation](https://developer.apple.com/documentation/applebusinessmanagerapi)
- Review the [SDK source code](../../client/axm/) for advanced usage

## 🆘 Troubleshooting

### Common Issues

**Authentication Errors:**
- Verify your Key ID, Issuer ID, and private key file
- Ensure your private key file is in the correct format (`.p8`)
- Check that your Apple Business Manager account has API access

**Network Errors:**
- Verify internet connectivity
- Check if you're behind a corporate firewall
- Ensure the API endpoint is accessible

**Parsing Errors:**
- Verify your private key format (supports both RSA and ECDSA)
- Check that the private key file is not corrupted
- Ensure proper file permissions on the private key file

---

**Happy coding!** 🚀 This guide should get you up and running with the Apple Business Manager API SDK.
