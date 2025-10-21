# Apple Business Manager API SDK v3

A completely refactored Go SDK for the Apple Business Manager API with direct service access and no wrappers.

## Key V3 Features

✅ **No Wrappers** - Direct service access without intermediate client layers  
✅ **Clean API** - `client.Service.Method()` pattern  
✅ **Type Safety** - Full type safety with proper interfaces  
✅ **Preserved Functionality** - All original CRUD methods and comments maintained  
✅ **Integrated Pagination** - Pagination logic moved into CRUD functions where appropriate  

## Architecture

```
v3/
├── client/           # Core HTTP client and shared utilities
├── devicemanagement/ # Device management service (no wrappers)
├── devices/          # Devices service (no wrappers) 
└── examples/         # Updated examples
```

## Usage

### Before (V2 - Unwanted Pattern)
```go
// V2 pattern with wrapper clients
axmClient, _ := axm.NewClient(config)
dmClient := devicemanagement.NewClient(axmClient)    // Wrapper!
devicesClient := devices.NewClient(axmClient)        // Wrapper!

response, err := dmClient.GetDeviceManagementServices(ctx, opts)
```

### After (V3 - Target Pattern)
```go
// V3 pattern with direct service access
client, err := axm.NewClient(config)

// Direct access - no wrappers!
response, err := client.DeviceManagement.GetDeviceManagementServices(ctx, opts)
devices, err := client.Devices.GetOrganizationDevices(ctx, opts)
```

## Quick Start

```go
package main

import (
    "context"
    
    axm "github.com/deploymenttheory/go-api-sdk-apple/v3"
    "github.com/deploymenttheory/go-api-sdk-apple/v3/devices"
)

func main() {
    // Parse your private key
    privateKey, err := axm.ParsePrivateKey([]byte(privateKeyPEM))
    if err != nil {
        panic(err)
    }

    // Create client with embedded services
    client, err := axm.NewClientBuilder().
        WithJWTAuth(keyID, issuerID, privateKey).
        WithDebug(true).
        Build()
    if err != nil {
        panic(err)
    }

    ctx := context.Background()

    // Direct service access - no wrappers needed!
    devicesResponse, err := client.Devices.GetOrganizationDevices(ctx, &devices.GetOrganizationDevicesOptions{
        Fields: []string{
            devices.FieldSerialNumber,
            devices.FieldDeviceModel,
            devices.FieldStatus,
        },
        Limit: 10,
    })
    
    // Direct device management access
    servers, err := client.DeviceManagement.GetDeviceManagementServices(ctx, nil)
    
    // Chain service calls naturally
    for _, device := range devicesResponse.Data {
        assigned, err := client.DeviceManagement.GetAssignedDeviceManagementServiceIDForADevice(ctx, device.ID)
        // ... handle assigned server
    }
}
```

## Available Services

### Device Management (`client.DeviceManagement`)
- `GetDeviceManagementServices()` - Get MDM servers
- `GetMDMServerDeviceLinkages()` - Get device linkages for MDM server  
- `GetAssignedDeviceManagementServiceIDForADevice()` - Get assigned server ID for device
- `GetAssignedDeviceManagementServiceInformationByDeviceID()` - Get assigned server info
- `AssignDevicesToServer()` - Assign devices to MDM server
- `UnassignDevicesFromServer()` - Unassign devices from MDM server

### Devices (`client.Devices`)
- `GetOrganizationDevices()` - Get organization devices
- `GetDeviceInformationByDeviceID()` - Get specific device information

## Authentication

Same authentication as V2, but with cleaner builder pattern:

```go
client, err := axm.NewClientBuilder().
    WithJWTAuth(keyID, issuerID, privateKey).
    WithBaseURL("https://api-business.apple.com/v1").
    WithTimeout(30 * time.Second).
    WithRetry(3, time.Second).
    WithDebug(true).
    Build()
```

## Migration from V2

1. Update imports:
   ```go
   // Old
   import "github.com/deploymenttheory/go-api-sdk-apple/client/axm"
   import "github.com/deploymenttheory/go-api-sdk-apple/services/axm/devicemanagement"
   
   // New  
   import axm "github.com/deploymenttheory/go-api-sdk-apple/v3"
   ```

2. Replace wrapper clients:
   ```go
   // Old
   axmClient, _ := axm.NewClient(config)
   dmClient := devicemanagement.NewClient(axmClient)
   response, _ := dmClient.GetDeviceManagementServices(ctx, opts)
   
   // New
   client, _ := axm.NewClient(config)  
   response, _ := client.DeviceManagement.GetDeviceManagementServices(ctx, opts)
   ```

3. All method signatures and functionality remain identical - only the access pattern changes!

## Why V3?

- **Eliminates Complexity**: No need to manage multiple client instances
- **Improves DX**: Intuitive `client.service.method()` pattern
- **Type Safety**: Proper interfaces without wrapper overhead  
- **Maintains Compatibility**: All existing CRUD methods preserved
- **Better Performance**: Removes unnecessary abstraction layers
- **Cleaner Code**: Less boilerplate, more focus on business logic

## Examples

See `examples/` directory for complete working examples demonstrating:
- Direct service access patterns
- Mixed service operations  
- Error handling
- Pagination usage
- Authentication setup