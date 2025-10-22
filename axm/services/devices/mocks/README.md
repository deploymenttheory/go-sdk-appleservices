# Device Service Mock Data

This directory contains mock JSON responses for testing the Apple Business Manager Device Service API endpoints.

## Files

- `validate_get_organization_devices.json` - Mock response for GET /orgDevices endpoint
- `validate_get_device_information.json` - Mock response for GET /orgDevices/{id} endpoint

## Usage

These JSON files are used by the test suite to validate API response parsing and ensure the client correctly handles the expected response structure from Apple's Business Manager API.

The mock data is based on the official Apple Business Manager API documentation and real API response examples.
