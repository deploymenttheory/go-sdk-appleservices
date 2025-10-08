GET https://api-business.apple.com/v1/mdmServers

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
              "self": "https://api-business.apple.com/v1/mdmServers/1F97349736CF4614A94F624E705841AD/relationships/devices"
            }
          }
        }
      }
    ],
    "links": {
      "self": "https://api-business.apple.com/v1/mdmServers"
    },
    "meta": {
      "paging": {
        "limit": 100
      }
    }
  }

GET https://api-business.apple.com/v1/mdmServers/{id}/relationships/devices

 {
    "data": [
      {
        "type": "orgDevices",
        "id": "XABC123X0ABC123X0"
      }
    ],
    "links": {
      "self": "https://api-business.apple/v1/mdmServers/1F97349736CF4614A94F624E705841AD/relationships/devices",
      "next": "https://api-business.apple/v1/mdmServers/1F97349736CF4614A94F624E705841AD/relationships/devices?cursor=MDowOjE3NDYxMTg0NjkzOTc6MTc0NjExODQ2OTM5Nzp0cnVlOmZhbHNlOjE3NDYxMTg0NjkzOTc",
      "related": "https://api-business.apple/v1/mdmServers/1F97349736CF4614A94F624E705841AD/devices"
    },
    "meta": {
      "paging": {
        "nextCursor": "MDowOjE3NDYxMTg0NjkzOTc6MTc0NjExODQ2OTM5Nzp0cnVlOmZhbHNlOjE3NDYxMTg0NjkzOTc",
        "limit": 100
      }
    }
  }

GET https://api-business.apple.com/v1/orgDevices/{id}/relationships/assignedServer

{
      "data": {
        "type": "mdmServers",
        "id": "1F97349736CF4614A94F624E705841AD"
      },
      "links": {
        "self": "https://api-business.apple.com/v1/orgDevices/DVVS36G1YD3JKQNI/relationships/assignedServer",
        "related": "https://api-business.apple.com/v1/orgDevices/DVVS36G1YD3JKQNI/assignedServer"
      }
    }

GET https://api-business.apple.com/v1/orgDevices/{id}/assignedServer

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
            "self": "https://api-business.apple.com/v1/mdmServers/1F97349736CF4614A94F624E705841AD/relationships/devices"
          }
        }
      }
    },
    "links": {
      "self": "https://api-business.apple.com/v1/orgDevices/DVVS36G1YD3JKQNI/assignedServer"
    }
  }

  POST https://api-business.apple.com/v1/orgDeviceActivities

  curl -X POST https://api-business.apple.com/v1/orgDeviceActivities \
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
      "self": "https://api-business.apple.com/v1/orgDeviceActivities/b1481656-b267-480d-b284-a809eed8b041"
    }
  },
  "links": {
    "self": "https://api-business.apple.com/v1/orgDeviceActivities"
  }
}

GET https://api-business.apple.com/v1/orgDeviceActivities/{id}

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
          "self": "https://api-business.apple.com/v1/orgDeviceActivities/84d7f133-b4a4-41be-ad0a-c2e4e53ea624"
        }
      },
      "links": {
        "self": "https://api-business.apple.com/v1/orgDeviceActivities/84d7f133-b4a4-41be-ad0a-c2e4e53ea624"
      }
    }