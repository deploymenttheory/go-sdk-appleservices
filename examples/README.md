# iTunes Search API Examples

This directory contains various examples demonstrating how to use the iTunes Search API Go client.

## Available Examples

Each example is a standalone Go program that can be run independently:

### 1. Basic Search (`basic_search.go`)
Demonstrates basic music search functionality.
```bash
go run basic_search.go
```

### 2. Music Videos Search (`music_videos_search.go`)
Shows how to search specifically for music videos.
```bash
go run music_videos_search.go
```

### 3. Apps Search (`apps_search.go`)
Demonstrates searching for iOS/macOS applications.
```bash
go run apps_search.go
```

### 4. Country-Specific Search (`country_search.go`)
Shows how to search within a specific country/region.
```bash
go run country_search.go
```

### 5. Lookup by ID (`lookup_by_id.go`)
Demonstrates looking up items by their iTunes ID.
```bash
go run lookup_by_id.go
```

### 6. Lookup by UPC (`lookup_by_upc.go`)
Shows how to lookup items using Universal Product Code (UPC).
```bash
go run lookup_by_upc.go
```

### 7. Advanced JSON Output (`advanced_json.go`)
Demonstrates working with raw JSON responses and complete data structures.
```bash
go run advanced_json.go
```

### 8. All Examples (`itunes_search_example.go`)
Runs all examples in sequence for a comprehensive demonstration.
```bash
go run itunes_search_example.go
```

## Features Demonstrated

- **Colored Logging**: All examples use colored Zap logging in debug mode
- **Error Handling**: Proper error handling and logging
- **Parameter Configuration**: Various search parameters and filters
- **Response Processing**: Different ways to handle and display results
- **Client Configuration**: Debug mode, timeouts, and retry logic

## Configuration

Each example uses the following client configuration:
- Debug mode enabled for detailed request/response logging
- Colored console output
- 30-second timeout
- 3 retry attempts with 1-second delay

## Prerequisites

Make sure you have Go 1.21+ installed and run:
```bash
go mod tidy
```

This will download all required dependencies including:
- `github.com/go-resty/resty/v2` for HTTP requests
- `go.uber.org/zap` for structured logging