# Go API SDK for Apple Services

[![Go Reference](https://pkg.go.dev/badge/github.com/deploymenttheory/go-api-sdk-apple.svg)](https://pkg.go.dev/github.com/deploymenttheory/go-api-sdk-apple)
[![Go Report Card](https://goreportcard.com/badge/github.com/deploymenttheory/go-api-sdk-apple)](https://goreportcard.com/report/github.com/deploymenttheory/go-api-sdk-apple)
[![License](https://img.shields.io/github/license/deploymenttheory/go-api-sdk-apple)](https://github.com/deploymenttheory/go-api-sdk-apple/blob/main/LICENSE)

This repo offers a collection of Go based SDKs and tools for interacting with various Apple API services and device management services, including:

- iTunes Search API
- Apple Business Manager / Apple School Manager API
- Apple Device Management API (MDM)

## Features

- Clean, idiomatic Go API
- Fluent builder patterns for constructing API requests
- Comprehensive error handling
- Configurable logging with zap
- Automatic retries with configurable parameters
- Extensive test coverage
- Complete examples for all supported operations

## Supported Services

### iTunes Search API

The SDK provides a complete implementation of the [iTunes Search API](https://performance-partners.apple.com/search-api), allowing you to:

- Search for content across iTunes, App Store, iBooks Store, and Mac App Store
- Look up content by ID, UPC, EAN, ISRC, or ISBN
- Filter results by media type, entity, country, and more

### Apple Business Manager / Apple School Manager API

Complete implementation of the [Apple Business Manager API](https://developer.apple.com/documentation/applebusinessmanagerapi) with modern Go practices:

**Devices API:**
- Get organization devices with filtering and pagination
- Get detailed device information by serial number
- Support for all device types (iPhone, iPad, Mac, Apple TV, Apple Watch)

**Device Management Services API:**
- List device management services in an organization
- Get device serial numbers assigned to services
- Assign/unassign devices to/from management services
- Get device management service assignments and information
- Track device activity operations

**Key Features:**
- **Centralized Architecture**: Unified error handling, pagination, and query building
- **Resty v3 Integration**: Built on latest Resty v3 with best practices
- **JWT Authentication**: Built-in Apple JWT token generation and management
- **Structured Logging**: Comprehensive request/response logging with zap
- **Type Safety**: Full generics support with structured response models
- **Pagination**: Automatic pagination with iterators and collectors
- **Context Support**: Context-aware operations for timeouts and cancellation

**Quick Start:**

Get started quickly with the GitLab-style client pattern:

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/deploymenttheory/go-api-sdk-apple/axm"
    "github.com/deploymenttheory/go-api-sdk-apple/axm/client"
    "github.com/deploymenttheory/go-api-sdk-apple/axm/services/devices"
    "github.com/deploymenttheory/go-api-sdk-apple/axm/services/devicemanagement"
)

func main() {
    // Method 1: Direct client creation (GitLab-style)
    // Parse private key from PEM format
    privateKey, err := client.ParsePrivateKey([]byte(privateKeyPEM))
    if err != nil {
        log.Fatalf("Failed to parse private key: %v", err)
    }
    
    client, err := axm.NewClient("your-key-id", "your-issuer-id", privateKey)
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }

    // Method 2: From environment variables
    // client, err := axm.NewClientFromEnv()
    // Expects: APPLE_KEY_ID, APPLE_ISSUER_ID, APPLE_PRIVATE_KEY_PATH

    // Method 3: From file
    // client, err := axm.NewClientFromFile("key-id", "issuer-id", "/path/to/key.p8")

    ctx := context.Background()

    // Get organization devices - GitLab-style access
    response, err := client.Devices.GetOrganizationDevices(ctx, &devices.RequestQueryOptions{
        Fields: []string{
            devices.FieldSerialNumber,
            devices.FieldDeviceModel,
            devices.FieldStatus,
        },
        Limit: 100,
    })
    if err != nil {
        log.Fatalf("Error getting devices: %v", err)
    }

    fmt.Printf("Found %d devices\n", len(response.Data))

    // Get device management services - GitLab-style access
    services, err := client.DeviceManagement.GetDeviceManagementServices(ctx, &devicemanagement.RequestQueryOptions{
        Limit: 10,
    })
    if err != nil {
        log.Fatalf("Error getting services: %v", err)
    }

    fmt.Printf("Found %d MDM servers\n", len(services.Data))
}
```

📖 **[Complete Quick Start Guide →](./examples/axm/quick_start.md)**

The guide covers 6 different client setup methods, from simple environment variables to advanced builder patterns with full customization options.

### Apple Device Management API

Integration with [Apple Device Management](https://developer.apple.com/documentation/devicemanagement) for:

- Mobile Device Management (MDM) operations
- Configuration profile management
- App and book distribution
- Declarative device management

## Examples

Explore the [examples directory](./examples) for comprehensive examples of using the SDK with different Apple services.

## Documentation

For detailed documentation, see:

- [Go Reference Documentation](https://pkg.go.dev/github.com/deploymenttheory/go-api-sdk-apple)
- [iTunes Search API Documentation](https://performance-partners.apple.com/search-api)
- [Apple Business Manager API Documentation](https://developer.apple.com/documentation/applebusinessmanagerapi)
- [Apple Device Management Documentation](https://developer.apple.com/documentation/devicemanagement)

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](./CONTRIBUTING.md) for guidelines.

## License

This project is licensed under the [MIT License](./LICENSE).
