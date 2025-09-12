# GetOrgDevices Examples

This directory contains examples demonstrating how to use the `GetOrgDevices` function from the Apple School and Business Manager API with the modern header-based API.

## Prerequisites

Before running these examples, you need to:

1. Set up your Apple Business Manager or Apple School Manager account
2. Generate a private key and get your Organization ID and Key ID from Apple
3. Choose one of these configuration methods:

### Method 1: Environment Variables (Quick Start)
```bash
export APPLE_ORG_ID="your-organization-id"
export APPLE_KEY_ID="your-key-id"
export APPLE_PRIVATE_KEY="-----BEGIN PRIVATE KEY-----
your-private-key-content-here
-----END PRIVATE KEY-----"

# Or reference a private key file:
export APPLE_PRIVATE_KEY_PATH="/path/to/your/private-key.pem"
```

### Method 2: JSON Configuration File (Recommended for Production)
Create a `config.json` file based on `config.example.json`:

```json
{
  "baseUrl": "https://axm-adm-enroll.apple.com",
  "orgId": "your-organization-id-here",
  "keyId": "your-key-id-here",
  "privateKeyPath": "/secure/path/to/your/private-key.pem",
  "timeoutSeconds": 30,
  "retryCount": 3,
  "retryDelayMs": 1000,
  "userAgent": "go-api-sdk-apple-example/1.0.0",
  "debug": true
}
```

### Method 3: Hybrid Approach
Use a config file with environment variable overrides for sensitive data:

```bash
# Base config in config.json, but override sensitive values
export APPLE_ORG_ID="production-org-id"
export APPLE_PRIVATE_KEY_PATH="/production/keys/private-key.pem"
```

## Modern API Example

### Get Organization Devices with Headers (`get_org_devices_with_headers.go`) - **Recommended**

This is the comprehensive example that demonstrates the modern API with explicit header control.

```bash
go run get_org_devices_with_headers.go
```

**Features:**
- **Custom Headers**: Demonstrates how headers are defined and used for each request
- **Multiple Query Patterns**: Shows various ways to query devices
  - Get all devices with automatic pagination
  - Field filtering to reduce response size
  - Status filtering (e.g., only available devices)
  - Pagination control with custom limits
- **Private Key Loading**: Shows both environment variable and file-based key loading
- **Error Handling**: Comprehensive error handling and logging
- **JSON Output**: Formatted JSON output of device data

**Key Concepts Demonstrated:**

1. **Header Definition**: Each request explicitly defines headers:
   ```go
   headers := map[string]string{
       "Content-Type": "application/json",
       "Accept":       "application/json",
   }
   ```

2. **QueryBuilder Usage**: Fluent interface for building queries:
   ```go
   queryBuilder := axmService.NewQueryBuilder().
       Fields("orgDevices", []string{"serialNumber", "deviceModel"}).
       Limit(10).
       Filter("status", "AVAILABLE")
   ```

3. **Private Key Loading**: Helper function for loading keys from files:
   ```go
   privateKey, err := client.LoadPrivateKeyFromFileWithValidation(privateKeyPath)
   ```

### Get Organization Devices with Config File (`get_org_devices_with_config_file.go`) - **Production Ready**

Shows how to use JSON configuration files for managing settings, including secure private key handling.

```bash
# First, create your config.json based on config.example.json
cp config.example.json config.json
# Edit config.json with your actual values

go run get_org_devices_with_config_file.go
```

**Features:**
- **JSON Configuration**: Load all settings from `config.json` file
- **Environment Overrides**: Allow env vars to override config file values  
- **Secure Key Management**: Reference private keys by file path, not inline
- **Config File Generation**: Programmatically create and save config files
- **Production Ready**: Proper file permissions and secure practices
- **Multiple Config Methods**: Pure file, env overrides, or hybrid approaches

**Key Configuration Methods:**

1. **Pure Config File**: All settings in JSON
   ```go
   config, err := client.LoadConfigFromFile("config.json")
   ```

2. **Config with Environment Overrides**: File + env vars
   ```go
   config, err := client.LoadConfigFromFileWithEnvOverrides("config.json")
   ```

3. **Save Config Programmatically**: Generate config files
   ```go
   client.SaveConfigToFile(config, "config.json", "/path/to/key.pem")
   ```

## Example Output

When you run the example, you'll see output like:

```
=== Example 1: Get All Organization Devices ===
Retrieved 150 devices
First device: {Type:orgDevices ID:abc123 SerialNumber:F9K2V3... ...}

=== Example 2: Get Devices with Field Filtering ===
Retrieved 10 devices with filtered fields
Device 1:
  Serial: F9K2V3H8K9L0
  Model: iPhone15,2
  Family: iPhone
  Status: AVAILABLE
  Added: 2024-01-15T10:30:00Z
...
```

## API Fields Available

The following fields are available for filtering with `fields[orgDevices]`:

- `serialNumber` - Device serial number
- `addedToOrgDateTime` - When device was added to organization
- `updatedDateTime` - Last update timestamp
- `deviceModel` - Device model identifier
- `productFamily` - Product family (iPhone, iPad, Mac, etc.)
- `productType` - Specific product type
- `deviceCapacity` - Storage capacity
- `partNumber` - Apple part number
- `orderNumber` - Purchase order number
- `color` - Device color
- `status` - Device enrollment status
- `orderDateTime` - Purchase date
- `imei` - IMEI (for cellular devices)
- `meid` - MEID (for cellular devices)
- `eid` - eID (for eSIM devices)
- `wifiMacAddress` - Wi-Fi MAC address
- `bluetoothMacAddress` - Bluetooth MAC address
- `purchaseSourceId` - Purchase source identifier
- `purchaseSourceType` - Type of purchase source
- `assignedServer` - Assigned MDM server

## Query Parameters

### Field Filtering
```go
.Fields("orgDevices", []string{"serialNumber", "deviceModel", "status"})
```

### Pagination
```go
.Limit(50) // Number of devices per request (max 200)
```

### Filtering
```go
.Filter("status", "AVAILABLE")     // Only available devices
.Filter("productFamily", "iPhone") // Only iPhones
```

### Sorting
```go
.Sort("addedToOrgDateTime")        // Sort by date added
.Sort("-deviceModel")              // Sort by model (descending)
```

## Error Handling

The example includes comprehensive error handling for:
- Missing environment variables
- Authentication failures
- API request errors
- JSON parsing errors
- File loading errors (when using private key files)

## Rate Limiting

The Apple Business Manager API has rate limits. The SDK includes built-in retry logic, but in production:

1. Monitor rate limit headers in responses
2. Implement exponential backoff for retries
3. Consider caching frequently accessed data
4. Use field filtering to reduce response sizes

## Private Key Management

The example shows two ways to provide the private key:

1. **Environment Variable** (for development):
   ```bash
   export APPLE_PRIVATE_KEY="-----BEGIN PRIVATE KEY-----..."
   ```

2. **File Path** (recommended for production):
   ```bash
   export APPLE_PRIVATE_KEY_PATH="/secure/path/to/key.pem"
   ```
   
   The `LoadPrivateKeyFromFileWithValidation()` helper provides:
   - File reading with error handling
   - Basic PEM format validation
   - Support for multiple key formats (RSA, PKCS#8, EC)

## Next Steps

After running this example:

1. **Explore Filtering**: Try different field combinations and filters
2. **Handle Large Datasets**: Test with organizations that have many devices
3. **Integrate with Your System**: Adapt the patterns for your specific use case
4. **Monitor Performance**: Use field filtering to optimize for your needs