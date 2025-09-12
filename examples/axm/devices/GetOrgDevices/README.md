# Basic GetOrgDevices Example

This example demonstrates how to use the Apple School and Business Manager API SDK to retrieve organization devices using JSON configuration.

## Features Demonstrated

- **JSON Configuration Loading**: Load client configuration from a JSON file
- **Authentication Testing**: Automatically test authentication before making API calls
- **Query Building**: Use the fluent QueryBuilder to filter and select device data
- **Structured Logging**: Use zap for structured, production-ready logging
- **Error Handling**: Proper error handling and cleanup
- **Data Processing**: Process and summarize device information

## Prerequisites

1. **Apple Business Manager or Apple School Manager account** with API access
2. **API credentials** including:
   - Client ID (format: `BUSINESSAPI.xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx`)
   - Key ID from your API key
   - Private key file (ECDSA P-256 in PEM format)

## Setup

### 1. Create Configuration File

Copy the example configuration:
```bash
cp config.example.json config.json
```

### 2. Update Configuration

Edit `config.json` with your actual credentials:

```json
{
  "baseUrl": "https://axm-adm-enroll.apple.com",
  "clientId": "BUSINESSAPI.your-actual-client-id-here",
  "keyId": "YOUR_KEY_ID",
  "privateKeyPath": "./your-private-key.pem",
  "scope": "business.api",
  "timeoutSeconds": 30,
  "retryCount": 3,
  "retryDelayMs": 1000,
  "userAgent": "go-api-sdk-apple/1.0.0",
  "debug": true
}
```

### 3. Add Your Private Key

Place your private key file in the same directory or update the `privateKeyPath` in config.json:
```bash
# Example private key file
./your-private-key.pem
```

## Running the Example

```bash
go run main.go
```

## Expected Output

The example will:

1. **Load and validate configuration** from config.json
2. **Test authentication** with Apple's API
3. **Print configuration summary** (excluding sensitive data)
4. **Fetch organization devices** with filtering and field selection:
   - Limited to 50 devices
   - Only "ASSIGNED" status devices
   - Selected fields only (serialNumber, deviceModel, etc.)
   - Sorted by addedToOrgDateTime
5. **Process and summarize** device data:
   - Count by product family (iPhone, iPad, Mac, etc.)
   - Count by status
   - Show details for first 10 devices
6. **Log structured output** using zap logger

## Example Log Output

```
INFO    Basic GetOrgDevices Example
INFO    This example demonstrates loading config from JSON and calling GetOrgDevices
INFO    Successfully authenticated    {"client_id": "BUSINESSAPI.xxx"}
INFO    Configuration Summary {"base_url": "https://axm-adm-enroll.apple.com", ...}
INFO    Fetching organization devices with filters    {"client_id": "BUSINESSAPI.xxx", "base_url": "https://axm-adm-enroll.apple.com"}
INFO    Successfully retrieved organization devices    {"device_count": 25}
INFO    Device summary by product family
INFO    Product family count    {"family": "iPhone", "count": 15}
INFO    Product family count    {"family": "iPad", "count": 10}
INFO    Device summary by status
INFO    Status count    {"status": "ASSIGNED", "count": 25}
INFO    === Device Details ===
INFO    Device details    {"index": 1, "id": "abc123", "serial_number": "C02XY1234567", ...}
INFO    Example completed successfully    {"total_devices_processed": 25}
```

## Key Code Features

### QueryBuilder Usage
```go
queryBuilder := axmService.NewQueryBuilder().
    Limit(50).                                    // Limit results
    Fields("orgDevice", []string{                 // Select specific fields
        "serialNumber",
        "deviceModel", 
        "productFamily",
        "status",
        "addedToOrgDateTime",
    }).
    Filter("status", "ASSIGNED").                 // Filter for assigned devices
    Sort("addedToOrgDateTime")                    // Sort by date added
```

### Configuration Loading
```go
// Load config from JSON file and test authentication
config, axmClient, err := client.LoadAndTestConfig(configPath, logger)
if err != nil {
    logger.Fatal("Failed to load and test config", zap.Error(err))
}
defer axmClient.Close()
```

### Structured Logging
```go
logger.Info("Device details",
    zap.Int("index", i+1),
    zap.String("id", device.ID),
    zap.String("serial_number", device.Attributes.SerialNumber),
    zap.String("device_model", device.Attributes.DeviceModel),
)
```

## Error Handling

The example includes proper error handling for:
- Configuration loading failures
- Authentication failures  
- API request failures
- Empty result sets

## Security Notes

- **Never commit** your actual `config.json` or private key files to version control
- Use **restrictive file permissions** (0600) for configuration files
- The private key should be in **ECDSA P-256 format** as required by Apple
- Consider using environment variables for sensitive data in production

## Troubleshooting

1. **Authentication Errors**: Verify your clientId, keyId, and private key are correct
2. **Empty Results**: Check if your organization actually has devices with "ASSIGNED" status
3. **Network Errors**: Verify connectivity to `https://axm-adm-enroll.apple.com`
4. **Private Key Errors**: Ensure your private key is in ECDSA P-256 PEM format
