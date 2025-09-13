Web Service Endpoint
Get Organization Devices
Get a list of devices in an organization that enroll using Automated Device Enrollment.
Apple School Manager API 1.1+
URL
GET https://api-school.apple.com/v1/orgDevices
Query Parameters
fields[orgDevices]
string
The fields to return for included related types.
Possible Values: serialNumber, addedToOrgDateTime, updatedDateTime, deviceModel, productFamily, productType, deviceCapacity, partNumber, orderNumber, color, status, orderDateTime, imei, meid, eid, wifiMacAddress, bluetoothMacAddress, purchaseSourceId, purchaseSourceType, assignedServer
limit
integer
The number of included related resources to return.
Maximum: 1000

response

    {
    "data": [
      {
        "type": "orgDevices",
        "id": "XABC123X0ABC123X0",
        "attributes": {
          "serialNumber": "XABC123X0ABC123X0",
          "addedToOrgDateTime": "2025-04-30T22:05:14.192Z",
          "updatedDateTime": "2025-05-01T15:33:54.164Z",
          "deviceModel": "iMac 21.5\"",
          "productFamily": "Mac",
          "productType": "iMac16,2",
          "deviceCapacity": "750GB",
          "partNumber": "FD311LL/A",
          "orderNumber": "1234567890",
          "color": "SILVER",
          "status": "UNASSIGNED",
          "orderDateTime": "2011-08-15T07:00:00Z",
          "imei": [
            "123456789012345",
            "123456789012346"
          ],
          "meid": [
            "123456789012347"
          ],
          "eid": "89049037640158663184237812557346",
          "purchaseSourceUid": "-2085650007946880",
          "purchaseSourceType": "APPLE"
        },
        "relationships": {
          "assignedServer": {
            "links": {
              "self": "https://api-school.apple.com/v1/orgDevices/XABC123X0ABC123X0/relationships/assignedServer",
              "related": "https://api-school.apple.com/v1/orgDevices/XABC123X0ABC123X0/assignedServer"
            }
          }
        },
        "links": {
          "self": "https://api-school.apple.com/v1/orgDevices/XABC123X0ABC123X0"
        }
      }
    ],
    "links": {
      "self": "https://api-school.apple.com/v1/orgDevices",
      "next": "https://api-school.apple.com/v1/orgDevices?cursor=MDowOjE3NDYxMTM4OTI1OTA6MTc0NjExMzg5MjU5MDp0cnVlOmZhbHNlOjE3NDYxMTM4OTI1OTA"
    },
    "meta": {
      "paging": {
        "nextCursor": "MDowOjE3NDYxMTM4OTI1OTA6MTc0NjExMzg5MjU5MDp0cnVlOmZhbHNlOjE3NDYxMTM4OTI1OTA",
        "limit": 100
      }
    }
  }

Object
OrgDevicesResponse
A response that contains a list of organization device resources.
Apple School Manager API 1.1+
object OrgDevicesResponse
Properties
data
[OrgDevice]
(Required) The resource data.
links
PagedDocumentLinks
(Required) Navigational links that include the self-link.
meta
PagingInformation
Paging information.
See Also
Objects and data types
object OrgDevice
The data structure that represents an organization device resource.
object OrgDeviceResponse
A response that contains a single organization device resource.
object OrgDeviceAssignedServerLinkageResponse
The data and links that describe the relationship between the resources.
object MdmServer
The data structure that represents device management services in organizations.
object MdmServerResponse
A response that contains a single device management service resource.
object MdmServersResponse
A response that contains a list of device management service resources.
object MdmServerDevicesLinkagesResponse
The data and links that describe the relationship between the resources.
object OrgDeviceActivity
The data structure that represents an organization device activity resource.
object OrgDeviceActivityCreateRequest
The request body you use to update the device management service for a device.
object OrgDeviceActivityResponse
A response that contains a single organization device activity resource.
object PagedDocumentLinks
Links related to the response document, including paging links.
object PagingInformation
Paging information for data responses.
object RelationshipLinks
Links related to the response document, including self-links.
object ResourceLinks
Self-links to requested resources.
object ErrorResponse
The error details that an API returns in the response body whenever the API request isn’t successful.

Web Service Endpoint
Get Device Information
Get information about a device in an organization.
Apple School Manager API 1.1+
URL
GET https://api-school.apple.com/v1/orgDevices/{id}
Path Parameters
id
string
(Required) The unique identifier for the resource.
Query Parameters
fields[orgDevices]
string
The fields to return for included related types.
Possible Values: serialNumber, addedToOrgDateTime, updatedDateTime, deviceModel, productFamily, productType, deviceCapacity, partNumber, orderNumber, color, status, orderDateTime, imei, meid, eid, wifiMacAddress, bluetoothMacAddress, purchaseSourceId, purchaseSourceType, assignedServer
Response Codes
200
OrgDeviceResponse
OK
Content-Type: application/json
400
ErrorResponse
Bad Request
An error occurred with your request.

Content-Type: application/json
401
ErrorResponse
Unauthorized
Content-Type: application/json
403
ErrorResponse
Forbidden
Request not authorized.

Content-Type: application/json
404
ErrorResponse
Not Found
Content-Type: application/json
429
ErrorResponse
Too Many Requests

Content-Type: application/json
Example
Request
Response
    {
      "data": {
        "type": "orgDevices",
        "id": "XABC123X0ABC123X0",
        "attributes": {
          "addedToOrgDateTime": "2025-04-30T22:05:14.192Z",
          "updatedDateTime": "2025-05-01T15:33:54.164Z",
          "deviceModel": "iMac 21.5\"",
          "productFamily": "Mac",
          "productType": "iMac16,2",
          "deviceCapacity": "750GB",
          "partNumber": "FD311LL/A",
          "orderNumber": "1234567890",
          "color": "",
          "status": "UNASSIGNED",
          "orderDateTime": "2011-08-15T07:00:00Z",
          "imei": [
            "123456789012345",
            "123456789012346"
          ],
          "meid": [
            "123456789012347"
          ],
          "eid": "89049037640158663184237812557346",
          "purchaseSourceUid": "-2085650007946880",
          "purchaseSourceType": "APPLE"
        },
        "relationships": {
          "assignedServer": {
            "links": {
              "self": "https://api-school.apple.com/v1/orgDevices/XABC123X0ABC123X0/relationships/assignedServer",
              "related": "https://api-school.apple.com/v1/orgDevices/XABC123X0ABC123X0/assignedServer"
            }
          }
        },
        "links": {
          "self": "https://api-school.apple.com/v1/orgDevices/XABC123X0ABC123X0"
        }
      },
      "links": {
        "self": "https://api-school.apple.com/v1/orgDevices/XABC123X0ABC123X0"
      }
    }

Web Service Endpoint
Get Device Management Services
Get a list of device management services in an organization.
Apple School Manager API 1.1+
URL
GET https://api-school.apple.com/v1/mdmServers
Query Parameters
fields[mdmServers]
string
The fields to return for included related types.
Possible Values: serverName, serverType, createdDateTime, updatedDateTime, devices
limit
integer
The number of included related resources to return.
Maximum: 1000
Response Codes
200
MdmServersResponse
OK
Content-Type: application/json
400
ErrorResponse
Bad Request
An error occurred with your request.

Content-Type: application/json
401
ErrorResponse
Unauthorized
Content-Type: application/json
403
ErrorResponse
Forbidden
Request not authorized.

Content-Type: application/json
429
ErrorResponse
Too Many Requests

Content-Type: application/json
Example
Request
Response
  {
    "data": [
      {
        "type": "mdmServers",
        "id": "1F97349736CF4614A94F624E705841AD",
        "attributes": {
          "serverName": "Test Device Management Service",
          "serverType": "MDM",
          "createdDateTime": "2025-05-01T03:21:44.685Z",
          "updatedDateTime": "2025-05-01T03:21:46.284Z"
        },
        "relationships": {
          "devices": {
            "links": {
              "self": "https://api-school.apple.com/v1/mdmServers/1F97349736CF4614A94F624E705841AD/relationships/devices"
            }
          }
        }
      }
    ],
    "links": {
      "self": "https://api-school.apple.com/v1/mdmServers"
    },
    "meta": {
      "paging": {
        "limit": 100
      }
    }
  }

  Web Service Endpoint
Get the Device Serial Numbers for a Device Management Service
Get a list of device serial numbers assigned to a device management service.
Apple School Manager API 1.1+
URL
GET https://api-school.apple.com/v1/mdmServers/{id}/relationships/devices
Path Parameters
id
string
(Required) The unique identifier for the resource.
Query Parameters
limit
integer
The number of included related resources to return.
Maximum: 1000
Response Codes
200
MdmServerDevicesLinkagesResponse
OK
Content-Type: application/json
400
ErrorResponse
Bad Request
An error occurred with your request.

Content-Type: application/json
401
ErrorResponse
Unauthorized
Content-Type: application/json
403
ErrorResponse
Forbidden
Request not authorized.

Content-Type: application/json
404
ErrorResponse
Not Found
Content-Type: application/json
429
ErrorResponse
Too Many Requests

Content-Type: application/json
Example
Request
Response
  {
    "data": [
      {
        "type": "orgDevices",
        "id": "XABC123X0ABC123X0"
      }
    ],
    "links": {
      "self": "https://api-school.apple/v1/mdmServers/1F97349736CF4614A94F624E705841AD/relationships/devices",
      "next": "https://api-school.apple/v1/mdmServers/1F97349736CF4614A94F624E705841AD/relationships/devices?cursor=MDowOjE3NDYxMTg0NjkzOTc6MTc0NjExODQ2OTM5Nzp0cnVlOmZhbHNlOjE3NDYxMTg0NjkzOTc",
      "related": "https://api-school.apple/v1/mdmServers/1F97349736CF4614A94F624E705841AD/devices"
    },
    "meta": {
      "paging": {
        "nextCursor": "MDowOjE3NDYxMTg0NjkzOTc6MTc0NjExODQ2OTM5Nzp0cnVlOmZhbHNlOjE3NDYxMTg0NjkzOTc",
        "limit": 100
      }
    }
  }

  Web Service Endpoint
Get the Assigned Device Management Service ID for a Device
Get the assigned device management service ID information for a device.
Apple School Manager API 1.1+
URL
GET https://api-school.apple.com/v1/orgDevices/{id}/relationships/assignedServer
Path Parameters
id
string
(Required) The unique identifier for the resource.
Response Codes
200
OrgDeviceAssignedServerLinkageResponse
OK
Content-Type: application/json
400
ErrorResponse
Bad Request
An error occurred with your request.

Content-Type: application/json
401
ErrorResponse
Unauthorized
Content-Type: application/json
403
ErrorResponse
Forbidden
Request not authorized.

Content-Type: application/json
404
ErrorResponse
Not Found
Content-Type: application/json
429
ErrorResponse
Too Many Requests

Content-Type: application/json
Example
Request
Response
    {
      "data": {
        "type": "mdmServers",
        "id": "1F97349736CF4614A94F624E705841AD"
      },
      "links": {
        "self": "https://api-school.apple.com/v1/orgDevices/XABC123X0ABC123X0/relationships/assignedServer",
        "related": "https://api-school.apple.com/v1/orgDevices/XABC123X0ABC123X0/assignedServer"
      }
    }

    Web Service Endpoint
Get the Assigned Device Management Service Information for a Device
Get the assigned device management service information for a device.
Apple School Manager API 1.1+
URL
GET https://api-school.apple.com/v1/orgDevices/{id}/assignedServer
Path Parameters
id
string
(Required) The unique identifier for the resource.
Query Parameters
fields[mdmServers]
string
The fields to return for included related types.
Possible Values: serverName, serverType, createdDateTime, updatedDateTime, devices
Response Codes
200
MdmServerResponse
OK
Content-Type: application/json
400
ErrorResponse
Bad Request
An error occurred with your request.

Content-Type: application/json
401
ErrorResponse
Unauthorized
Content-Type: application/json
403
ErrorResponse
Forbidden
Request not authorized.

Content-Type: application/json
404
ErrorResponse
Not Found
Content-Type: application/json
429
ErrorResponse
Too Many Requests

Content-Type: application/json
Example
Request
Response
  {
    "data": {
      "type": "mdmServers",
      "id": "1F97349736CF4614A94F624E705841AD",
      "attributes": {
        "serverName": "Test Device Management Service",
        "serverType": "APPLE_CONFIGURATOR",
        "createdDateTime": "2025-05-01T03:21:44.685Z",
        "updatedDateTime": "2025-05-01T03:21:46.284Z"
      },
      "relationships": {
        "devices": {
          "links": {
            "self": "https://api-school.apple.com/v1/mdmServers/1F97349736CF4614A94F624E705841AD/relationships/devices"
          }
        }
      }
    },
    "links": {
      "self": "https://api-school.apple.com/v1/orgDevices/DVVS36G1YD3JKQNI/assignedServer"
    }
  }


  Web Service Endpoint
Assign or Unassign Devices to a Device Management Service
Assign or unassign devices to a device management service.
Apple School Manager API 1.1+
URL
POST https://api-school.apple.com/v1/orgDeviceActivities
HTTP Body
OrgDeviceActivityCreateRequest
Content-Type: application/json
Response Codes
201
OrgDeviceActivityResponse
Created
Content-Type: application/json
400
ErrorResponse
Bad Request
An error occurred with your request.

Content-Type: application/json
401
ErrorResponse
Unauthorized
Content-Type: application/json
403
ErrorResponse
Forbidden
Request not authorized.

Content-Type: application/json
409
ErrorResponse
Conflict
Content-Type: application/json
422
ErrorResponse
Content-Type: application/json
429
ErrorResponse
Too Many Requests

Content-Type: application/json
Example
Request
Response
curl -X POST https://api-school.apple.com/v1/orgDeviceActivities \
 -H "Authorization: Bearer ${ACCESS_TOKEN} \
 -d '{
   "data": {
     "type": "orgDeviceActivities",
     "attributes": {
       "activityType": "ASSIGN_DEVICES"
     },
     "relationships": {
       "mdmServer": {
         "data": {
           "type": "mdmServers",
           "id": "1F97349736CF4614A94F624E705841AD"
         }
       },
       "devices": {
         "data": [
           {
             "type": "orgDevices",
             "id": "XABC123X0ABC123X0"
           }
         ]
       }
     }
   }
 }'


       {
  "data": {
    "type": "orgDeviceActivities",
    "id": "b1481656-b267-480d-b284-a809eed8b041",
    "attributes": {
      "status": "IN_PROGRESS",
      "subStatus": "SUBMITTED",
      "createdDateTime": "2025-05-05T04:15:43.282Z"
    },
    "links": {
      "self": "https://api-school.apple.com/v1/orgDeviceActivities/b1481656-b267-480d-b284-a809eed8b041"
    }
  },
  "links": {
    "self": "https://api-school.apple.com/v1/orgDeviceActivities"
  }
}


Web Service Endpoint
Get Organization Device Activity Information
Get information for an organization device activity that a device management action, such as assigning or unassigning, creates.
Apple School Manager API 1.1+
URL
GET https://api-school.apple.com/v1/orgDeviceActivities/{id}
Path Parameters
id
string
(Required) The unique identifier for the resource.
Query Parameters
fields[orgDeviceActivities]
string
The fields to return for included related types.
Possible Values: status, subStatus, createdDateTime, completedDateTime, downloadUrl
Response Codes
200
OrgDeviceActivityResponse
OK
Content-Type: application/json
400
ErrorResponse
Bad Request
An error occurred with your request.

Content-Type: application/json
401
ErrorResponse
Unauthorized
Content-Type: application/json
403
ErrorResponse
Forbidden
Request not authorized.

Content-Type: application/json
404
ErrorResponse
Not Found
Content-Type: application/json
429
ErrorResponse
Too Many Requests

Content-Type: application/json
Overview
Note

This API supports fetching organization device activity information that the system creates for device management actions, such as assigning or unassigning devices to the device management service in the past 30 days.

Request
Response
    {
      "data": {
        "type": "orgDeviceActivities",
        "id": "84d7f133-b4a4-41be-ad0a-c2e4e53ea624",
        "attributes": {
          "status": "COMPLETED",
          "subStatus": "COMPLETED_WITH_SUCCESS",
          "createdDateTime": "2025-05-01T18:22:27.106Z",
          "completedDateTime": "2025-05-01T18:22:38.894Z",
          "downloadUrl": "https://store.blobstore.apple.com/4a7f73ecd1/1317ef485b3aea95ac80/9b5c028726eacaa1f3e2/ccfa60f6114198d48cc5/88c579b1995efc28cd36?response-content-disposition=attachment%3Bfilename%3D%22ABM-ActivityLog_May-1-2025_14-22-40.csv%22&response-content-type=application%2Foctet-stream&iCloudAccessKeyId=MACOSX_SU_ACCESS_KEY&Expires=1746123820&Signature=wrYW2gn7UZDM9aO3KOMqT8YWCBo%3D"
        },
        "links": {
          "self": "https://api-school.apple.com/v1/orgDeviceActivities/84d7f133-b4a4-41be-ad0a-c2e4e53ea624"
        }
      },
      "links": {
        "self": "https://api-school.apple.com/v1/orgDeviceActivities/84d7f133-b4a4-41be-ad0a-c2e4e53ea624"
      }
    }

    ────────────────────────────────────────────────────────────────────────────────────────────────────╮
│ > now i want you to go through the sdk functions and comapre with this -                             │
│   /Users/dafyddwatkins/GitHub/deploymenttheory/go-api-sdk-apple/client/axm2/resource_docs -          │
│   validate that only the funcs that we name api to the api calls we have avialble from the api set.  │
│   ensure all structs are exactly correct. \                                                          │
│   \                                                                                                  │
│   also we need to use a query builder that aligns with how resty v3 works and functions but is       │
│   compatible with the realities of the api. functions that support queries should have them as a     │
│   param. there should not be seperate funcs for sdk funcs with params. redundant code. pagaination   │
│   also should be applied by default to funcs that support it. and not applied to those that don't.   │
│   pagination should be a part of the client GET command. use the resty provided option, but ensure   │
│   it aligns with how the api docs function and work.    