Web Service
Notary API
Submit your macOS software for notarization through a web interface.
Notary API 2.0.0+
Overview
Notarization gives people confidence that your Developer ID-signed macOS software has been checked by Apple for malicious code. In addition to interacting with the notary service through Xcode or the notarytool command-line utility, you can bypass notarytool and interact directly with the service through its REST API. The Notary API is helpful for instances where you need to avoid a macOS dependency when uploading your app to the notary service, and provides endpoints that enable you to:

Prepare the notary service to receive a new version of your software, and get credentials that you use to upload your software to an Amazon S3 endpoint.

Check the status of a submission.

Retrieve a log file that provides details about a submission.

Get a list of your team’s previous submissions.

To learn how notarization works, see Notarizing macOS software before distribution. For details about using the notary service REST API to upload your software, see Submitting software for notarization over the web.

Topics
Essentials
Submitting software for notarization over the web
Eliminate a dependency on macOS in your notarization workflow by interfacing directly with the notary service.
Software submission
Submit Software
Start the process of uploading a new version of your software to the notary service.
object NewSubmissionRequest
Data that you provide when starting a submission to the notary service.
object NewSubmissionResponse
The notary service’s response to a software submission.
Notarization results
Get Submission Status
Fetch the status of a software notarization submission.
object SubmissionResponse
The notary service’s response to a request for the status of a submission.
Get Submission Log
Fetch details about a single completed notarization.
object SubmissionLogURLResponse
The notary service’s response to a request for the log information about a completed submission.
History
Get Previous Submissions
Fetch a list of your team’s previous notarization submissions.
object SubmissionListResponse
The notary service’s response to a request for information about your team’s previous submissions.
Errors
object ErrorResponse
The notary service’s response when an error occurs.

Article
Submitting software for notarization over the web
Eliminate a dependency on macOS in your notarization workflow by interfacing directly with the notary service.
Overview
If you notarize macOS software that you distribute with Developer ID, and if you use a custom notarization workflow, you can use notarytool (included with Xcode) to interact with the notary service. However, if you want to avoid a macOS dependency for this part of the workflow, you can use the notary service’s REST API.

For an overview of the entire custom notarization workflow, see Customizing the notarization workflow. This article describes how to use the notary service API to replace the parts of that article that depend on notarytool, namely uploading your software and checking on the status of your request.

Create a private key
When you access the notary service, you authenticate the access by including a token that you sign with a private key. Use the same key to sign tokens for the notary service that you use for the App Store Connect API. For details on how to create a key, see Creating API Keys for App Store Connect API.

Keep the key secure and private. You only need to create a key once, but if you lose the key, or if the key becomes compromised, revoke it immediately and create a new one. For more information, see Revoking API Keys.

Include a signed token with each API access
The notary service requires a JSON Web Token (JWT) to authorize each API request. JWT is an open standard (RFC 7519) that defines a way to securely transmit information by using a private key to cryptographically sign a message in JSON format. You use the key to create a JWT for each request. For more information, see Generating Tokens for API Requests.

After creating the signed token, include it in the header for all of your notary service API calls. For example, to use the Get Previous Submissions endpoint to get a list of up to 100 of your team’s previous notarization submissions, replace <token> in the following call to the curl command with your signed token:

% curl -v -H "Authorization: Bearer <token>" "https://appstoreconnect.apple.com/notary/v2/submissions"
Create a token that expires a short time after you create it — for example, 20 minutes. You can reuse a token until it expires.

Start a submission
To submit a new version of your software for notarization, start by calling the Submit Software endpoint. This call doesn’t upload your software; it tells the notary service to prepare for a new upload, and returns information that you use to perform the upload. Prepare the body of the request with a name and a hash of the software that you want to notarize, as well as optional information about how to notify you when notarization completes. For example, you can create a body variable in Python for a file named OvernightTextEditor_11.6.8.zip:

#!/usr/bin/env python3


import hashlib


with open("OvernightTextEditor_11.6.8.zip", "rb") as file:
    hash = hashlib.sha256()
    hash.update(file.read())
    sha256 = hash.hexdigest()


body = {
    "submissionName": "OvernightTextEditor_11.6.8.zip",
    "sha256": sha256,
    "notifications": [{"channel": "webhook", "target": "https://example.com" }]
}
The body consists of three keys, one of which is optional:

submissionName
Indicates the name of the file that you plan to submit. Use a unique file name for each submission to make it easier when reviewing log files and other status information about the submission. The example above does this by including a version string as part of the file name.

sha256
A hash that acts as a signature the service can use to match against the software that you upload later. Use the Secure Hashing Algorithm 2 (SHA-2) with a 256-bit digest. The above example uses the hashlib library to compute the hash.

notifications
An optional array that contains a dictionary with a target URL that the notary service accesses when it completes the notarization process. Use notifications to avoid polling the service to find out when notarization completes. If you ask for a notification, be sure to verify its cryptographic signature. For a description of the notification format, see NewSubmissionRequest.Notifications. Omit the notifications key and its associated array if you don’t need a callback.

Use the body, plus the token described in the previous section, to call the endpoint and get a response. The following code continues from the previous Python example:

import requests


token = generate_token() # Defined elsewhere.
resp = requests.post("https://appstoreconnect.apple.com/notary/v2/submissions", json=body, headers={"Authorization": "Bearer " + token})
resp.raise_for_status()
output = resp.json()
The service responds with an identifier that you use to track your submission, and temporary security credentials that you use to upload your software to an Amazon S3 endpoint. An example response looks like this:

{
  "data": {
    "attributes": {
      "awsAccessKeyId": "ASIAIOSFODNN7EXAMPLE",
      "awsSecretAccessKey": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
      "awsSessionToken": "AQoDYXdzEJr...",
      "bucket": "EXAMPLE-BUCKET",
      "object": "EXAMPLE-KEY-NAME"
    },
    "id": "2efe2717-52ef-43a5-96dc-0797e4ca1041",
    "type": "submissionsPostResponse"
  },
  "meta": {
  }
} 
Store the identifier (id) so you can use it to ask the notary service about the status of the submission, and to retrieve information about the outcome of the notarization. If you lose the identifier, you can use the Get Previous Submissions endpoint to get a list of the 100 most recent submissions associated with your team.

Upload your software
Use the attributes object that appears in the response from the Submit Software call to upload your software to Amazon S3. A good way to do this is with the boto3 library, provided by Amazon:

import boto3


aws_info = output["data"]["attributes"]
bucket = aws_info["bucket"]
key = aws_info["object"]
sub_id = output["data"]["id"]


s3 = boto3.client(
         "s3",
         aws_access_key_id=aws_info["awsAccessKeyId"],
         aws_secret_access_key=aws_info["awsSecretAccessKey"],
         aws_session_token=aws_info["awsSessionToken"],
         config=Config(s3={"use_accelerate_endpoint": True})
)


resp = s3.upload_file("OvernightTextEditor_11.6.8.zip", bucket, key)
The temporary S3 credentials contained in the attributes object expire 12 hours after you get them from the notary service. If you don’t use the credentials before they expire, you’ll have to restart the upload process. For more information about using this library and accessing Amazon S3, see the documentation on https://aws.amazon.com.

Check the status of your submission
Notarization typically takes just a few minutes after you complete the upload. To verify the status of your submission, you can call the Get Submission Status endpoint by using the following curl command along with the identifier that the service returned when you started the submission:

% curl -v -H "Authorization: Bearer <token>" "https://appstoreconnect.apple.com/notary/v2/submissions/2efe2717-52ef-43a5-96dc-0797e4ca1041"
After notarization completes, use the Get Submission Log endpoint to get a log file that contains any notarization errors or warnings:

% curl -v -H "Authorization: Bearer <token>" "https://appstoreconnect.apple.com/notary/v2/submissions/2EFE2717-52EF-43A5-96DC-0797E4CA1041/logs"
This call returns a URL that you can use to download a log file in JSON format. The log file provides details about the notarization outcome for this particular submission. The log file URL is valid for only a few hours, but you can ask for a new URL later if you need the log again. Always check the log file, even if notarization succeeds, because it might contain warnings that you can fix prior to your next submission.

Web Service Endpoint
Submit Software
Start the process of uploading a new version of your software to the notary service.
Notary API 2.0.0+
URL
POST https://appstoreconnect.apple.com/notary/v2/submissions
HTTP Body
NewSubmissionRequest
Information about a new software submission that you want to make.
Content-Type: application/json
Response Codes
200
NewSubmissionResponse
OK
The service recieved your submission request. The response includes information that you use to upload your software and track the status of your submission.

Content-Type: application/json
Mentioned in
Submitting software for notarization over the web
Discussion
Use this endpoint to tell the notary service about a new software submission that you want to make. Do this when you want to notarize a new version of your software.

You provide an HTTP body that contains a name for the submission, a hash of the software that you plan to submit, and an optional webhook that the service uses to notify you when notarization completes. For the name, use the name of the file that you upload, including the dmg or zip extension. The service responds with temporary security credentials that you use to submit the software to Amazon S3 and a submission identifier that you use to track the submission’s status.

After uploading your software, you can use the identifier to ask the notary service for the status of your submission using the Get Submission Status endpoint. If you lose the identifier, you can get a list of your team’s 100 most recent submissions using the Get Previous Submissions endpoint. After notarization completes, use the Get Submission Log to get details about the outcome of notarization. Do this even if notarization succeeds, because the log might contain warnings that you can fix before your next submission.

Example
Request
Response
{
  "notifications": [
    {
      "channel": "webhook",
      "target": "https://example.com"
    }
  ],
  "sha256": "68d561c564ef61f718e99a81b13bcb52af11b7ac9baf538af3ea0c83326fb6a1",
  "submissionName": "OvernightTextEditor_11.6.8.zip"
} 

Object
NewSubmissionRequest
Data that you provide when starting a submission to the notary service.
Notary API 2.0.0+
object NewSubmissionRequest
Properties
notifications
[NewSubmissionRequest.Notifications]
An optional array of notifications that you want to receive when notarization finishes. Omit this key if you don’t need a notification.
sha256
string
(Required) A cryptographic hash of the software that you want to notarize, computed using Secure Hashing Algorithm 2 (SHA-2) with a 256-bit digest. Supply the hash as a string of 64 hexadecimal digits. You must compute the hash from the exact version of the software that you plan to upload to Amazon S3.
Value: /[A-Fa-f0-9]{64}/
submissionName
string
(Required) The name of the file that you plan to submit. The service includes this name in its responses when you ask for the status of a submission, get a list of previous submissions, or get a log file corresponding to a submission. The file name doesn’t have to be unique among all your submissions, but making it so might help you to distinguish among submissions in service responses.
Attributes
Possible types:
Discussion
Use a structure of this type as the HTTP body when you post to the Submit Software endpoint.

Object
NewSubmissionRequest.Notifications
A notification that the notary service sends you when notarization finishes.
Notary API 2.0.0+
object NewSubmissionRequest.Notifications
Properties
channel
string
The channel that the service uses to notify you when notarization completes. The only supported value for this key is webhook.
target
string
The URL that the notary service accesses when notarization completes.
Attributes
Possible types:
Mentioned in
Submitting software for notarization over the web
Discussion
If you want a notification when notarization completes, include a data structure of this type in the notifications array that’s part of the body when you post to the Submit Software endpoint. Set the value for the channel key to webhook, and provide a URL as the value for the target key. The service indicates when notarization finishes by posting a body to the URL like this:

{
  "payload": "{\"completed_time\":\"2022-06-08T22:04:00.886Z\",\"event\":\"processing-complete\",\"start_time\":\"2022-06-08T22:03:42.801Z\",\"submission_id\":\"c00fe84b-95f2-4890-b7ea-019a7f546abd\",\"team_id\":\"JA62H4Q78D\"}",
  "signature": "MEUCIEqr...",
  "cert_chain": "LS0tLS1CR..."
}


The value for the payload key indicates when the operation starts and completes, as well as the submission ID and your Team ID. The submission ID matches the value that you receive in response to the Submit Software call. You can use the signature and cert_chain fields to verify the authenticity of the message against the Apple Inc. Root certificate that you can download from the Apple PKI site. If you need the certificate repeatedly, store a copy of the certificate on your server rather than downloading it every time you need it.

Object
NewSubmissionResponse
The notary service’s response to a software submission.
Notary API 2.0.0+
object NewSubmissionResponse
Properties
data
NewSubmissionResponse.Data
Data that describes the result of the submission request.
meta
NewSubmissionResponse.Meta
An empty object that you can ignore.
Attributes
Possible types:
Discussion
You receive a structure of this type in response to a call to the Submit Software endpoint. Use the temporary security credentials this response contains to make a call to Amazon S3 to upload your software.

Topics
Objects
object NewSubmissionResponse.Data
Information that the notary service provides for uploading your software for notarization and tracking the submission.
object NewSubmissionResponse.Meta
An empty object.

Object
NewSubmissionResponse.Data
Information that the notary service provides for uploading your software for notarization and tracking the submission.
Notary API 2.0.0+
object NewSubmissionResponse.Data
Properties
attributes
NewSubmissionResponse.Data.Attributes
Information that you use to upload your software to Amazon S3.
id
string
A unique identifier for this submission. Use this value to track the status of your submission. For example, you use it as the submissionID parameter in the Get Submission Status call, or to match against the id field in the response from the Get Previous Submissions call.
type
string
The resource type.
Attributes
Possible types:

NewSubmissionResponse.Data.Attributes
Information that you use to upload your software for notarization.
Notary API 2.0.0+
object NewSubmissionResponse.Data.Attributes
Properties
awsAccessKeyId
string
An access key that you use in a call to Amazon S3.
awsSecretAccessKey
string
A secret key that you use in a call to Amazon S3.
awsSessionToken
string
A session token that you use in a call to Amazon S3.
bucket
string
The Amazon S3 bucket that you upload your software into.
object
string
The object key that identifies your software upload within the bucket.
Attributes
Possible types:
Discussion
Use the temporary security credentials in this object, along with the bucket and object, to upload your software to Amazon S3. A good way to do this is with the boto3 library provided by Amazon, as described in Submitting software for notarization over the web. Be sure to use the S3 credentials before they expire, which happens 12 hours after you receive them.

For more information about using this library and accessing Amazon S3, see the documentation on https://aws.amazon.com.

Object
NewSubmissionResponse.Meta
An empty object.
Notary API 2.0.0+
object NewSubmissionResponse.Meta
Attributes
Possible types:
Discussion
This object is reserved for future use.

See Also
Objects

Get Submission Status
Fetch the status of a software notarization submission.
Notary API 2.0.0+
URL
GET https://appstoreconnect.apple.com/notary/v2/submissions/{submissionId}
Path Parameters
submissionId
uuid
(Required) The identifier that you receive from the notary service when you post to Submit Software to start a new submission.
Value: /[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/
Response Codes
200
SubmissionResponse
OK
The status request succeeded. The response contains the status.

Content-Type: application/json
403
ErrorResponse
Forbidden
An authentication failure occurred.

Content-Type: application/json
404
ErrorResponse
Not Found
The specified identifier can’t be found.

Content-Type: application/json
Mentioned in
Submitting software for notarization over the web
Discussion
Use this endpoint to fetch the status of a submission request. Form the URL for the call using the identifier that you receive in the id field of the response to the Submit Software endpoint. If you lose the identifier, you can get a list of the most recent 100 submissions by calling the Get Previous Submissions endpoint.

Along with the status of the request, the response indicates the date that you initiated the request and the software name that you provided at that time.

Example
Request
Response
https://appstoreconnect.apple.com/notary/v2/submissions/2efe2717-52ef-43a5-96dc-0797e4ca1041
See Also
Notarization results
object SubmissionResponse
The notary service’s response to a request for the status of a submission.
Get Submission Log
Fetch details about a single completed notarization.
object SubmissionLogURLResponse
The notary service’s response to a request for the log information about a completed submission.

Example
Request
Response
{
  "data": {
    "attributes": {
      "createdDate": "2022-06-08T01:38:09.498Z",
      "name": "OvernightTextEditor_11.6.8.zip",
      "status": "Accepted"
    },
    "id": "2efe2717-52ef-43a5-96dc-0797e4ca1041",
    "type": "submissions"
  },
  "meta": {
  }
} 

Object
SubmissionResponse
The notary service’s response to a request for the status of a submission.
Notary API 2.0.0+
object SubmissionResponse
Properties
data
SubmissionResponse.Data
Data that describes the status of the submission request.
meta
SubmissionResponse.Meta
An empty object that you can ignore.
Attributes
Possible types:
Discussion
You receive a structure of this type in response to a call to the Get Submission Status endpoint.

Topics
Objects
object SubmissionResponse.Data
Information that the service provides about the status of a notarization submission.
object SubmissionResponse.Meta
An empty object.

Object
SubmissionResponse.Data
Information that the service provides about the status of a notarization submission.
Notary API 2.0.0+
object SubmissionResponse.Data
Properties
attributes
SubmissionResponse.Data.Attributes
Information about the status of a submission.
id
string
The unique identifier for this submission. This value matches the value that you provided as a path parameter to the Get Submission Status call that elicited this response.
type
string
The resource type.
Attributes
Possible types:
Topics
Objects
object SubmissionResponse.Data.Attributes
Information about the status of a submission.

Object
SubmissionResponse.Data.Attributes
Information about the status of a submission.
Notary API 2.0.0+
object SubmissionResponse.Data.Attributes
Properties
createdDate
string
The date that you started the submission process, given in ISO 8601 format, like 2022-06-08T01:38:09.498Z.
name
string
The name that you specified in the submissionName field of the Submit Software call when you started the submission.
status
string
The status of the submission. The associated string contains one of the following: Accepted, In Progress, Invalid, or Rejected.
Attributes
Possible types:

Web Service Endpoint
Get Submission Log
Fetch details about a single completed notarization.
Notary API 2.0.0+
URL
GET https://appstoreconnect.apple.com/notary/v2/submissions/{submissionId}/logs
Path Parameters
submissionId
string
(Required) The identifier that you receive from the notary service when you post to Submit Software to start a new submission.
Response Codes
200
SubmissionLogURLResponse
OK
The request succeeded. The response contains a URL that you access to download the log information.

Content-Type: application/json
403
ErrorResponse
Forbidden
An authentication failure occurred.

Content-Type: application/json
404
ErrorResponse
Not Found
No data was found for this team.

Content-Type: application/json
Mentioned in
Submitting software for notarization over the web
Discussion
Use this endpoint to get a URL that you can download a log file from that enumerates any issues found by the notary service. The URL that you receive is temporary, so be sure to use it to immediately fetch the log. If you need the log again in the future, ask for the URL again.

The log file that you download contains JSON-formatted data, and might include both errors and warnings. For information about how to deal with common notarization problems, see Resolving common notarization issues.

Example
Request
Response
https://appstoreconnect.apple.com/notary/v2/submissions/2EFE2717-52EF-43A5-96DC-0797E4CA1041/logs

{
  "data": {
    "attributes": {
      "developerLogUrl": "https://..."
    },
    "id": "2efe2717-52ef-43a5-96dc-0797e4ca1041",
    "type": "submissionsLog"
  },
  "meta": {
  }
} 

Object
SubmissionLogURLResponse
The notary service’s response to a request for the log information about a completed submission.
Notary API 2.0.0+
object SubmissionLogURLResponse
Properties
data
SubmissionLogURLResponse.Data
Data that indicates how to get the log information for a particular submission.
meta
SubmissionLogURLResponse.Meta
An empty object that you can ignore.
Attributes
Possible types:
Discussion
You receive a structure of this type in response to a call to the Get Submission Log endpoint.

Topics
Objects
object SubmissionLogURLResponse.Data
Data that indicates how to get the log information for a particular submission.
object SubmissionLogURLResponse.Meta
An empty object.

Object
SubmissionLogURLResponse.Data
Data that indicates how to get the log information for a particular submission.
Notary API 2.0.0+
object SubmissionLogURLResponse.Data
Properties
attributes
SubmissionLogURLResponse.Data.Attributes
Information about the log associated with the submission.
id
string
The unique identifier for this submission. This value matches the value that you provided as a path parameter to the Get Submission Log call that elicited this response.
type
string
The resource type.
Attributes
Possible types:

Object
SubmissionLogURLResponse.Data.Attributes
Information about the log associated with the submission.
Notary API 2.0.0+
object SubmissionLogURLResponse.Data.Attributes
Properties
developerLogUrl
string
The URL that you use to download the logs for a submission. The URL serves a JSON-encoded file that contains the log information. The URL is valid for only a few hours. If you need the log again later, ask for the URL again by making another call to the Get Submission Log endpoint.

Object
SubmissionLogURLResponse.Meta
An empty object.
Notary API 2.0.0+
object SubmissionLogURLResponse.Meta
Attributes
Possible types:
Discussion
This object is reserved for future use.

See Also
Objects
object SubmissionLogURLResponse.Data
Data that indicates how to get the log information for a particular submission.
Web Service Endpoint
Get Previous Submissions
Fetch a list of your team’s previous notarization submissions.
Notary API 2.0.0+
URL
GET https://appstoreconnect.apple.com/notary/v2/submissions
Response Codes
200
SubmissionListResponse
OK
The submission list request succeeded. The response contains a list of recent submissions, truncated to the 100 most recent.

Content-Type: application/json
403
ErrorResponse
Forbidden
An authentication failure occurred.

Content-Type: application/json
404
ErrorResponse
Not Found
No data was found for this team.

Content-Type: application/json
Mentioned in
Submitting software for notarization over the web
Discussion
Use this endpoint to get the list of submissions associated with your team. The response holds an array of values that include the unique identifier for the submission, the date you initiated the submission, the name of the associated software, and the status of the submission. The response returns information about only the 100 most recent submissions.

If you need information about just one submission, and you have the associated identifier, use Get Submission Status instead.

Example
Request
Response
https://appstoreconnect.apple.com/notary/v2/submissions

Example
Request
Response
{
  "data": [
    {
      "attributes": {
        "createdDate": "2021-04-29T01:38:09.498Z",
        "name": "OvernightTextEditor_11.6.8.zip",
        "status": "Accepted"
      },
      "id": "2efe2717-52ef-43a5-96dc-0797e4ca1041",
      "type": "submissions"
    },
    {
      "attributes": {
        "createdDate": "2021-04-23T17:44:54.761Z",
        "name": "OvernightTextEditor_11.6.7.zip",
        "status": "Accepted"
      },
      "id": "cf0c235a-dad2-4c24-96eb-c876d4cb3a2d",
      "type": "submissions"
    },
    {
      "attributes": {
        "createdDate": "2021-04-19T16:56:17.839Z",
        "name": "OvernightTextEditor_11.6.7.zip",
        "status": "Invalid"
      },
      "id": "38ce81cc-0bf7-454b-91ef-3f7395bf297b",
      "type": "submissions"
    }
  ],
  "meta": {
  }
} 

Object
SubmissionListResponse
The notary service’s response to a request for information about your team’s previous submissions.
Notary API 2.0.0+
object SubmissionListResponse
Properties
data
[SubmissionListResponse.Data]
An array of objects, each of which describes one of your team’s previous submissions.
meta
SubmissionListResponse.Meta
An empty object that you can ignore.
Attributes
Possible types:
Discussion
You receive a structure of this type in response to a call to the Get Previous Submissions endpoint. The list includes only the 100 most recent submissions.

Topics
Objects
object SubmissionListResponse.Data
Data that describes one of your team’s previous submissions.
object SubmissionListResponse.Meta
An empty object.

Object
SubmissionListResponse.Data
Data that describes one of your team’s previous submissions.
Notary API 2.0.0+
object SubmissionListResponse.Data
Properties
attributes
SubmissionListResponse.Data.Attributes
Information about a particular submission.
id
string
The unique identifier for a submission. This value matches the value that you received in the id field that appeared in the response to the Submit Software call that you used to start the submission.
type
string
The resource type.
Attributes
Possible types:

Object
ErrorResponse
The notary service’s response when an error occurs.
Notary API 2.0.0+
object ErrorResponse
Properties
description
string
A string that describes the reason for the error.
labels
[string]
Additional information about the error.
name
string
The name of the error.
Attributes
Possible types:

