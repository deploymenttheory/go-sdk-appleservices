# Apple Business Manager API Client Specification

## API Characteristics & Design Implications

### Authentication
- **Pattern**: OAuth 2.0 with JWT Client Assertion (RFC 7523)
- **Token Endpoint**: `https://account.apple.com/auth/oauth2/v2/token`
- **Flow**: 
  1. Generate JWT client assertion signed with private key (ES256)
  2. Exchange JWT for access token using `client_credentials` grant
  3. Use access token in `Authorization: Bearer {token}` header
  
- **JWT Claims**:
  - `iss`: Team ID (Issuer ID)
  - `sub`: Team ID (Subject - same as issuer for Apple)
  - `aud`: Token endpoint URL
  - `iat`: Issued at time
  - `exp`: Expiration (max 180 days as per Apple docs)
  - `jti`: Unique identifier (timestamp-based)
  
- **Key Requirements**:
  - Private key: ECDSA (ES256) or RSA (RS256) from Apple Developer portal
  - Key ID (kid) in JWT header
  - PEM format private key (.p8 file)
  
- **Token Lifetime**: Dynamic (returned in `expires_in` field)
- **Scope**: `business.api` or `school.api`
- **Impact on Client**:
  - JWT generation and signing logic
  - Token manager with automatic refresh logic
  - Refresh before expiry (5 minute buffer)
  - Thread-safe token acquisition for concurrent requests
  - 401 responses trigger automatic token refresh
  - Contrast with VirusTotal: VirusTotal uses simple static API key
  - Contrast with Nexthink: Nexthink uses standard OAuth2 client credentials without JWT

### API Versioning
- **Pattern**: Version in URL path
- **URL Structure**: `https://api-business.apple.com/v{version}/{resource}`
- **Current Version**: `v1`
- **Impact on Client**:
  - API version constant in each service package (`APIVersion = "/v1"`)
  - Base URL excludes version: `https://api-business.apple.com`
  - Each service endpoint prepends version (e.g., `APIVersion + "/orgDevices"`)
  - Function names include version suffix for future-proofing (e.g., `GetOrganizationDevicesV1`)
  - Future v2 endpoints can coexist with v1 in same service
  - Contrast with Workbrew: Workbrew uses header-based versioning

### Data Model
- **Standard**: JSON:API Specification (jsonapi.org)
- **Structure**:
  ```json
  {
    "data": {
      "type": "orgDevices",
      "id": "...",
      "attributes": { ... },
      "relationships": { ... }
    },
    "meta": { ... },
    "links": { "self": "...", "next": "..." }
  }
  ```
  
- **Resource Types**: `orgDevices`, `mdmServers`, `orgDeviceActivities`, `appleCareCoverage`
- **Relationships**: Linkages between resources (e.g., device → assigned MDM server)
- **Impact on Client**:
  - Strongly typed models per resource
  - Separate types for attributes, relationships, linkages
  - Metadata and links at response root level
  - Field selection via sparse fieldsets
  - Contrast with Nexthink: Nexthink uses simple JSON, not JSON:API
  - Contrast with VirusTotal: VirusTotal uses custom envelope format

### Pagination
- **Pattern**: Cursor-based pagination per JSON:API spec
- **Parameters**:
  - `limit`: Number of items per page (max 1000)
  - `cursor`: Opaque cursor from `links.next`
  
- **Response Structure**:
  ```json
  {
    "data": [...],
    "meta": {
      "totalCount": N
    },
    "links": {
      "self": "...",
      "next": "...",
      "prev": "..."
    }
  }
  ```
  
- **Impact on Client**:
  - `GetPaginated` helper method in HTTP client
  - Automatic pagination until no more `links.next`
  - Page callback function for streaming large datasets
  - Enforced 1000-item limit per request
  - Contrast with Nexthink: Nexthink has no pagination, uses export for large data

### Resource Identification
- **Pattern**: Type-safe resource linkages
- **ID Types**:
  - Devices: Serial numbers (primary), UUIDs (relationships)
  - MDM Servers: UUIDs
  - Activities: UUIDs
  
- **Relationships**: JSON:API relationship objects
  ```json
  {
    "data": {
      "type": "mdmServers",
      "id": "uuid"
    }
  }
  ```
  
- **Impact on Client**:
  - Separate methods for fetching resources vs relationships
  - `GetDeviceInformationByDeviceID` for full device object
  - `GetAssignedDeviceManagementServiceIDForADevice` for relationship linkage only
  - Clear type safety in request/response structs

### Sparse Fieldsets
- **Pattern**: Client-specified field selection per JSON:API spec
- **Query Parameter**: `fields[{type}]={field1,field2,...}`
- **Example**: `fields[orgDevices]=serialNumber,deviceModel,status`
- **Impact on Client**:
  - `RequestQueryOptions` struct with `Fields` slice
  - Query builder constructs `fields[...]` parameters
  - Reduces payload size for large device lists
  - All fields returned if not specified
  - Field constants defined per resource type

### Device Assignment Operations
- **Pattern**: Asynchronous activity-based mutations
- **Process**:
  1. POST `/orgDeviceActivities` with activity type and device IDs
  2. Returns activity resource with status
  3. Poll activity status (IN_PROGRESS → COMPLETED/FAILED)
  
- **Activity Types**:
  - `ASSIGN_DEVICES`: Assign devices to MDM server
  - `UNASSIGN_DEVICES`: Remove devices from MDM server
  
- **Batch Limits**: No explicit limit documented, practical limits apply
- **Impact on Client**:
  - Two-step process: submit activity + monitor status
  - Activity response includes per-device success/failure details
  - No built-in polling helper (client polls manually)
  - Contrast with Nexthink: Nexthink operations are synchronous or use export pattern

### AppleCare Coverage
- **Pattern**: Read-only sub-resource of devices
- **Endpoint**: `/orgDevices/{id}/appleCareCoverage`
- **Data**: Coverage status, payment type, dates, contract info
- **Pagination**: Supports pagination (a device can have multiple coverage entries)
- **Impact on Client**:
  - Separate service method for AppleCare data
  - Paginated response handling
  - Coverage-specific field selection

### Response Formats
- **Supported**: JSON only
- **Content-Type**: `application/json`
- **Accept**: `application/json` required
- **Impact on Client**:
  - Single response parsing path
  - No CSV/XML support
  - Type-safe JSON structs

### Rate Limiting
- **Pattern**: Standard HTTP 429
- **Headers**: Not explicitly documented by Apple
- **Impact on Client**:
  - Automatic retry with exponential backoff (3 retries default)
  - Configurable retry count and wait time
  - Response object accessible even on error for header inspection
  - Contrast with Nexthink: Nexthink provides explicit rate limit headers

### Validation Requirements
- **Serial Numbers**: Device serial numbers for device operations
- **UUIDs**: MDM server IDs, activity IDs
- **Activity Types**: Must be valid enum (`ASSIGN_DEVICES`, `UNASSIGN_DEVICES`)
- **Batch Operations**: Submit lists of device IDs
- **Impact on Client**:
  - Validation before API calls to fail fast
  - Clear error messages for invalid inputs
  - Type safety through enums and constants

### Error Handling
- **Pattern**: JSON:API error format
- **Structure**:
  ```json
  {
    "errors": [
      {
        "status": "400",
        "code": "INVALID_REQUEST",
        "title": "Invalid Request",
        "detail": "..."
      }
    ]
  }
  ```
  
- **Status Codes**:
  - 400: Invalid request parameters
  - 401: Invalid or expired access token
  - 403: Insufficient permissions
  - 404: Resource not found
  - 429: Rate limit exceeded
  - 500: Server error
  
- **Impact on Client**:
  - Custom error handler parses JSON:API errors
  - Error details extracted and formatted
  - Automatic token refresh on 401
  - Response object returned even on error

### Base URL
- **Pattern**: Fixed base URL
- **URL**: `https://api-business.apple.com`
- **No Dynamic Construction**: Unlike Nexthink's instance/region pattern
- **Impact on Client**:
  - Simple constant base URL
  - Optional override via `WithBaseURL` option for testing
  - Contrast with Nexthink: Nexthink constructs URL from instance + region

### HTTP Client Configuration
- **Transport Options**: Comprehensive configuration via functional options
- **Available Options**:
  - TLS: Custom certificates, mTLS, minimum TLS version, skip verify
  - Proxy: HTTP/HTTPS/SOCKS5 proxy support
  - Timeouts: Request timeout, retry wait times
  - Headers: Global headers, custom user agent
  - Logging: Structured logging with zap
  - Debug: Request/response inspection
  - Authentication: Custom auth providers
  
- **Impact on Client**:
  - Production-ready security defaults (TLS 1.2+)
  - Enterprise proxy support
  - Observability through logging
  - Testability through mock auth providers
  - Contrast with simpler SDKs: More comprehensive enterprise features

### Query Builder
- **Pattern**: Fluent interface for URL parameters
- **Capabilities**:
  - Add strings, integers, string slices
  - Build final query string
  - URL encoding handled automatically
  
- **Impact on Client**:
  - Type-safe query construction
  - Consistent parameter formatting
  - Used for field selection, pagination parameters

### JSON:API Compliance
- **Full Compliance**: Client implements JSON:API specification
- **Features Implemented**:
  - Resource objects with type/id/attributes/relationships
  - Sparse fieldsets (field selection)
  - Pagination with links
  - Relationship objects and linkages
  - Meta information
  - Error objects
  
- **Impact on Client**:
  - Interoperability with JSON:API tooling
  - Consistent API interaction patterns
  - Self-documenting resource structures
  - Contrast with other SDKs: Apple's commitment to JSON:API standard is unique
