# Microsoft Mac Apps Version Tracker Examples

This directory contains examples demonstrating how to use the Microsoft Mac Apps Version Tracker API client.

## Available Examples

### GetLatestApps

Get all latest Microsoft Mac application versions available.

```bash
cd GetLatestApps
go run main.go
```

Features demonstrated:

- Retrieving all available Microsoft Mac applications
- Displaying application metadata (name, version, size, etc.)
- Parsing timestamps
- Calculating total storage requirements
- Working with application components

### GetAppByBundleID

Find a specific Microsoft application by its bundle ID.

```bash
cd GetAppByBundleID
go run main.go
```

Features demonstrated:

- Finding applications by bundle ID
- Using bundle ID constants
- Checking for multiple applications
- Finding apps with specific components (e.g., AutoUpdate)
- Error handling for non-existent apps
- JSON output formatting

### GetAppByName

Find a specific Microsoft application by its name.

```bash
cd GetAppByName
go run main.go
```

Features demonstrated:

- Finding applications by name
- Using application name constants
- Analyzing update timestamps
- Generating deployment scripts
- Working with Office 365 suite components
- Comparing versions across applications

## API Documentation

For more information about the Microsoft Mac Apps API, visit:

- API Endpoint: <https://appledevicepolicy.tools/api/latest>
- Documentation: <https://appledevicepolicy.tools/microsoft-apps>

## Prerequisites

```bash
go get github.com/deploymenttheory/go-api-sdk-apple/msapps
```

## Common Use Cases

### Check for Updates

Use `GetAppByBundleID` or `GetAppByName` to check if newer versions are available for your deployed applications.

### Generate Download Scripts

The examples show how to generate shell scripts to download specific applications for deployment.

### Monitor Application Sizes

Use `GetLatestApps` to monitor the storage requirements for Microsoft applications.

### Track Components

Find which applications include specific components like Microsoft AutoUpdate or licensing packages.

## Notes

- The API does not require authentication
- Data is updated regularly (check the `generated` field in responses)
- All applications include SHA256 hashes for integrity verification
- Download URLs are official Microsoft CDN links

