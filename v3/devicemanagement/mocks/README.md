# Device Management Service Mock Data

This directory contains mock JSON responses for testing the Apple Business Manager Device Management Service API endpoints.

## Files

- `validate_get_device_management_services.json` - Mock response for GET /mdmServers endpoint
- `validate_get_mdm_server_device_linkages.json` - Mock response for GET /mdmServers/{id}/relationships/devices endpoint
- `validate_get_assigned_server_linkage.json` - Mock response for GET /orgDevices/{id}/relationships/assignedServer endpoint
- `validate_get_assigned_server_info.json` - Mock response for GET /orgDevices/{id}/assignedServer endpoint
- `validate_assign_devices_response.json` - Mock response for POST /orgDeviceActivities (assign devices) endpoint
- `validate_unassign_devices_response.json` - Mock response for POST /orgDeviceActivities (unassign devices) endpoint

## Usage

These JSON files are used by the test suite to validate API response parsing and ensure the client correctly handles the expected response structure from Apple's Business Manager API.

The mock data is based on the official Apple Business Manager API documentation and real API response examples.
