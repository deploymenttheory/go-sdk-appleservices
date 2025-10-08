GET https://api-business.apple.com/v1/orgDevices

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
            "12345678901237"
          ],
          "eid": "89049037640158663184237812557346",
          "purchaseSourceUid": "-2085650007946880",
          "purchaseSourceType": "APPLE"
        },
        "relationships": {
          "assignedServer": {
            "links": {
              "self": "https://api-business.apple.com/v1/orgDevices/XABC123X0ABC123X0/relationships/assignedServer",
              "related": "https://api-business.apple.com/v1/orgDevices/XABC123X0ABC123X0/assignedServer"
            }
          }
        },
        "links": {
          "self": "https://api-business.apple.com/v1/orgDevices/XABC123X0ABC123X0"
        }
      }
    ],
    "links": {
      "self": "https://api-business.apple.com/v1/orgDevices",
      "next": "https://api-business.apple.com/v1/orgDevices?cursor=MDowOjE3NDYxMTM4OTI1OTA6MTc0NjExMzg5MjU5MDp0cnVlOmZhbHNlOjE3NDYxMTM4OTI1OTA"
    },
    "meta": {
      "paging": {
        "nextCursor": "MDowOjE3NDYxMTM4OTI1OTA6MTc0NjExMzg5MjU5MDp0cnVlOmZhbHNlOjE3NDYxMTM4OTI1OTA",
        "limit": 100
      }
    }
  }

GET https://api-business.apple.com/v1/orgDevices/{id}

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
              "self": "https://api-business.apple.com/v1/orgDevices/XABC123X0ABC123X0/relationships/assignedServer",
              "related": "https://api-business.apple.com/v1/orgDevices/XABC123X0ABC123X0/assignedServer"
            }
          }
        },
        "links": {
          "self": "https://api-business.apple.com/v1/orgDevices/XABC123X0ABC123X0"
        }
      },
      "links": {
        "self": "https://api-business.apple.com/v1/orgDevices/XABC123X0ABC123X0"
      }
    }