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

Support for the [Apple Business Manager API](https://developer.apple.com/documentation/applebusinessmanagerapi), enabling:

- Device enrollment management
- Content purchase and distribution
- User and location management

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
