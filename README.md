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
```go
// Create JWT auth
auth := client.NewJWTAuth(client.JWTAuthConfig{
    KeyID:      "YOUR_KEY_ID",
    IssuerID:   "YOUR_ISSUER_ID",
    PrivateKey: privateKey,
})

// Create client
appleClient, err := apple.NewClient(client.Config{
    Auth:   auth,
    Logger: logger,
    Debug:  true,
})
defer appleClient.Close()

// Get devices with pagination
devices, err := appleClient.Devices().GetOrganizationDevices(ctx, &devices.GetOrganizationDevicesOptions{
    Model: devices.ModeliPhone,
    PaginationOptions: client.PaginationOptions{Limit: 100},
})

// Iterate through all pages
for result := range devices.Iterator(ctx) {
    fmt.Printf("Device: %s\n", result.Item.SerialNumber)
}
```

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
