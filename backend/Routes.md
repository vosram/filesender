# Routes

## User related routes

These user routes are mostly perform CRUD operations on the currently logged in user or be used by the admin to make changes.

- [ ] `GET /api/users/me`
- [ ] `PUT /api/users/me`
- [ ] `DELETE /api/users/me`

### `GET /api/users/me`

- Auth Requred: True
- Request Body: False
- Response Body: JSON

This endpoint gets the currently logged in users data. which is really just name and email. `credential` is what login method was used to create the account. This is important because `PUT /api/users/me` can only change the email if the user's credential value is `magic_email`.

```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "credential": "email"
}
```

### `PUT /api/users/me`

- Auth Required: True
- Request Body: JSON
- Response Body: JSON

This endpoint will update the name of the user and if the email is different, then it will launch the email change workflow. This workflow will only start if the user's credential is `magic_email`.

Email change workflow would look like this:

```text
PUT /api/users/me
  |
  |
  V
Verification email send to current email
link with token /email-change-verify?token=<token>
  |
  |
  V
User clicks on verification email
GET /email-change-verify?token=<token>
  |
  |
  V
FE sends token to BE
POST /api/auth/email/change-verify
{
  "verifyToken": "d953144ddab8e0337db17ea6137733f75c8a7435cb0e785d43c8cf166f2318af"
}
FE lets user know a new email has been sent to the new email to finalize verification
  |
  |
  V
User clicks link sent to new email
GET /email-change-confirmation?token=<token>
  |
  |
  V
FE sends token to BE
POST /api/auth/email/change-confirmation
{
  "confirmToken": "fe36208ef7fc352ea9ca78610960f3d8733d34205f12562b2a06b4c3cc97a3c1"
}
BE verifies the token and updates the user's email.
on success, FE will redirect to FE /account
```

## Auth Routes

- [x] `POST /api/auth/email/login`
- [x] `POST /api/auth/email/verify`
- [ ] `POST /api/auth/email/change-verify`
- [ ] `POST /api/auth/email/change-confirmation`
- [x] `GET /api/auth/google/login`
- [x] `GET /api/auth/google/callback`
- [x] `GET /api/auth/github/login`
- [x] `GET /api/auth/github/callback`
- [x] `POST /api/auth/refresh`
- [x] `GET /api/auth/me`

### JWT claims

This is what claims should be set in JWT access tokens:

```json
{
  "sub": "86a6b2b5-7e04-4d89-885e-48cdbd9da98e",
  "iss": "FileSender",
  "role": "user",
  "membership": "pro",
  "iat": 1516239022,
  "exp": 1516249022
}
```

### `POST /api/auth/email/login

- Auth Required: False
- Request Body: JSON
- Response Body: JSON only on error

This endpoint is used to send a login token to the users email. If the user doesn't have an account registered it will create an account. The login token should be a 32-byte base64url token generated from `crypto/rand` and saved in the DB table `email_login_tokens`. They should look something like `xpzNeE4DIQwRQiDAVoyxnV9qeHWwAb-P7c8uY6RBbUw`. The email sent to the user should add a link like `https://domain.com/auth/email/authenticate?token=<token>`. That will be a FE route that should fire off another call to the BE at `POST /api/auth/verify`.

**Request Body:**

```json
{
  "email": "name@email.com"
}
```

### `POST /auth/email/verify`

- Auth Required: False
- Request Body: JSON
- Response Body: JSON

This endpoint validates the token with the DB table `email_login_tokens`. If valid, an httpOnly cookie with a refresh token and the accessToken in a JSON response should be returned.

**Request Body**

```json
{
  "token": "P9Awma1P6STMOn8zQ8p9jwyjkIIbihtrKt1jLVkzpnw"
}
```

**Response Cookies**

- name: refresh_token
- secure: True
- httpOnly: True
- sameSite: 'Strict'
- Path: '/api/auth'

**Response Body**

```json
{
  "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiI3ZGU1Mjg3MC01NzU0LTQzMGQtYTkyMi05M2E5M2RkMDlhODAiLCJyb2xlIjoidXNlciIsImlhdCI6MTUxNjIzOTAyMiwiZXhwIjoxNTE2MjQ5MDIyfQ.-YCNhR9TTejFxtJUJaRVG_rLad79nb6j9FHIrXYHdhw"
}
```

### `POST /api/auth/email/change-verify`

- Auth Required: true
- Request Body: JSON
- Response Body: None

This endpoint is to verify that the currently logged in user is authorizing a change in the email login they're using.

**Request Body:**

```json
{
  "verifyToken": "39dd50ac61931c31328c0fc08fe22aed935a735364316789e62cf10cc534cf50"
}
```

That token is hashed to sha256 and a DB lookup to table `email_change_verify_tokens` is performed to find that token. That DB record will have a `new_email` field where our BE will create a confirmation token and then that raw token to the new email. The raw token will be hashed and saved to DB table `email_change_confirmation_tokens`.

### `POST /api/auth/email/change-confirmation`

- Auth Required: true
- Request Body: JSON
- Response Body: None

**Request Body:**

```json
{
  "confirmToken": "e9574e55a57265ec5b77c56f3c0c1d184107c1cf3046d8f9e6d8d8ac5f8762a0"
}
```

This endpoint is to confirm the new email address is working. the token is hashed to sha256 and a DB lookup on table `email_change_confirmation_tokens` is performed. If valid, the user's account email is changed to that of the new email. At this point we know the user has access to this new email.

### Google Oauth routes

These routes should be configured with Go's [Oauth2](https://pkg.go.dev/golang.org/x/oauth2) package. The routes that should be setup are:

- `GET /api/auth/google`
- `GET /api/auth/google/callback`

In order the get the [auth code url](https://pkg.go.dev/golang.org/x/oauth2#Config.AuthCodeURL), you need a `state` variable which can be a 16 byte random opaque value that should be saved in an httpOnly cookie before redirecting the client to the auth url, and should be checked with the returned state when `/api/auth/google/callback` is hit during the oauth redirect.

The response to the callback endpoint should have a `refresh_token` sent via httpOnly cookie. A redirect should be sent to `/oauth-callback`. This `/oauth-callback` is a FE route that will hit `POST /api/auth/refresh` to get an access token from the `refresh_token` in the httpOnly cookie.

### Github Oauth routes

- `GET /api/auth/github`
- `GET /api/auth/github/callback`

### `POST /api/auth/refresh`

- Auth Required: True
- Request Body: False
- Response Body: JSON

This endpoint reads the `refreshToken` in the httpOnly cookie and validates it with the records in `refresh_tokens` table in DB. This route should also implement token rotation. If the refresh token is valid, a new token should be created and saved in the `refresh_tokens` table. The tokens saved in the DB should be hashed with SHA256 while the token is sent back via httpOnly cookie should be the unhashed token. This ensures that if the refreshTokens in the DB are compomised, they can't be immidiately used by the hackers.

**Response Cookie**

- Name: refresh_token
- Value: `<token>`
- httpOnly: true
- secure: true
- path: "/api/auth"
- SameSite: "lax"

**Response Body**

```json
{
  "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiI4NmE2YjJiNS03ZTA0LTRkODktODg1ZS00OGNkYmQ5ZGE5OGUiLCJpc3MiOiJGaWxlU2VuZGVyIiwicm9sZSI6InVzZXIiLCJpYXQiOjE1MTYyMzkwMjIsImV4cCI6MTUxNjI0OTAyMn0._SojRglLGVpQ50Xv8yhCEHlIkyBrmKzVVmJMFpq2SCY"
}
```

### `GET /api/auth/me`

- Auth Required: True
- Request Body: False
- Response Body: JSON

This endpoint is just a simple auth check primarily used for debugging and real time auth check for production.

**Response Body**

```json
{
  "name": "John Doe",
  "role": "user"
}
```

## Files related routes

- [ ] `GET /api/files`
- [ ] `GET /api/files/:fileId`
- [ ] `POST /api/files/upload`
- [ ] `POST /api/files/upload/:clientId/complete`
- [ ] `POST /api/files/upload/:clientId/abort`
- [ ] `PATCH /api/files/:fileId`
- [ ] `DELETE /api/files/:fileId`

### `GET /files`

- Auth Required: True
- Query Params:
  - sort: 'name' | 'size' | 'type' | 'uploaded' | 'expires'
  - filter: (multiple keys) '?filter=document&filter=data&filter=image'
  - page: integer
- Response Body: JSON

This endpoint retrieves a list of all files in the users files pool. The way it will collect all filter arguments is through [c.QueryArray](https://pkg.go.dev/github.com/gin-gonic/gin#Context.QueryArray). If no filter arguments are present, all files should be returned. If no page query param if present, it defaults to 1. Files per page are limited to 20 by default.

Response:

```json
{
  "totalPages": 5,
  "page": 1,
  "files": [
    {
      "id": "82b2b9c2-6b26-45db-ae80-b1c2cd18292b",
      "name": "report_q2_draft.pdf",
      "size": "2.4 MB",
      "type": "Doc",
      "uploadedOn": "2026-06-15T09:30:00.0000z",
      "expiresOn": "2026-06-15T09:30:00.0000z"
    },
    {
      "id": "65c9e2a3-69d7-425b-b8f0-81d5b21f430c",
      "name": "july-sales-report.xls",
      "size": "3.1 MB",
      "type": "Doc",
      "uploadedOn": "2026-06-15T09:30:00.0000z",
      "expiresOn": "2026-06-15T09:30:00.0000z"
    },
    {
      "id": "a24d01c8-deb7-4cad-a2c6-289af16fbee0",
      "name": "2025-family-group-picture.jpg",
      "size": "7.2 MB",
      "type": "Image",
      "uploadedOn": "2026-06-15T09:30:00.0000z",
      "expiresOn": null
    }
  ]
}
```

### `GET /files/:fileId`

- Auth Required: True
- Request Body: False
- Response Body: JSON

This endpoint gets all the available data for a file, including a presigned download link.

**JSON Response Body**:

```json
{
  "file": {
    "id": "82b2b9c2-6b26-45db-ae80-b1c2cd18292b",
    "name": "report_q2_draft.pdf",
    "size": "2.4 MB",
    "type": "Doc",
    "uploadedOn": "2026-06-15T09:30:00.0000z",
    "expiresOn": "2026-06-15T09:30:00.0000z",
    "url": "https://aws.com/key?sig=YmY5ZWM0YzE2Njg2NzljNjQxYzk4NzU5YzNlZmJhMzBiNjhhZGIyZDc2MTVjODdiYT"
  }
}
```

### `POST /files/upload`

- Auth Required: True
- Request Body: JSON
- Response Body: JSON

This endpoint initiates an upload for one or more files. It generates s3 presigned urls so the client can upload files directly to s3 instead of passing all that data to the server. The file size is used to validate whether the file is allowed. Each file has a client-generated uuid (`clientId`) so the correct server-side presigned url is returned for the correct file — the filename is not a reliable identifier. The final stored key is generated by the backend.

The server decides whether a file is a single-part or multipart upload based on `size >= SINGLE_PART_THRESHOLD` (default 100 MB). For multipart files, the server calls `CreateMultipartUpload` to obtain an `UploadId`, then presigns one url per part using the server-computed part size (`MULTIPART_PART_SIZE`, default 64 MB, clamped to S3's 5 MB minimum and 10,000 part maximum). Each part url includes the `partNumber` and `uploadId` in its signature.

For every file, the server inserts an `uploads` row (`status = 'initiated'`) that records the generated key, `uploadId`, part count, size, and owner. This row is used later to validate the completion request so the client cannot forge a file key, id, or size. No `files` record is created until the client confirms the upload completed.

**Constants**

- `SINGLE_PART_THRESHOLD` = 100 MB
- `MULTIPART_PART_SIZE` = 64 MB
- `Presigned_url_TTL` = 15 minutes

**S3 Bucket Lifecycle**

It's important to create this rule on the AWS S3 console as it prevents failed multipart to take up s3 storage and increase your costs. Make sure to set the following lifecycle rule on your dev and production bucket.

S3 bucket lifecycle rule: abort incomplete multipart uploads after 7 days.

You can do so by the following:

1. Go to s3 console.
2. Go to your chosen bucket.
3. Go to Management
4. In Lifecycle Configuration, click "Create lifecycle rule"
5. Set Lifecycle rule name, choose rule scope: "Apply to all objects in the bucket"
6. Lifecycle rule Actions: Check "Delete expired object... or incomplete multipart uploads"
7. Check incomplete nultipart uploads > "Delete incomplete multipart uploads"
8. Set Number of days to what you feel is appropriate. 3 days is a good default
9. Finally click "Create rule" orange button

A cleanup job should also remove stale `uploads` rows.

**Request Body**

```JSON
{
  "files": [
    {
      "clientId": "1999c359-77ed-4ad3-bd73-a42c3c89f9c0",
      "name": "somename.XIFF",
      "contentType": "image/jpeg",
      "size": 820000000
    },
    {
      "clientId": "56f85dae-9eb9-4b8e-97ef-8647b826ccab",
      "name": "july-report.xls",
      "contentType": "application/vnd.ms-excel",
      "size": 5300000
    },
    {
      "clientId": "e7b0f2c1-1d2c-4c3a-9b2e-8f3a5d6c7e8f",
      "name": "group-picture.jpg",
      "contentType": "image/jpeg",
      "size": 8100000
    }
  ]
}
```

**Response Body**

```json
{
  "files": [
    {
      "clientId": "1999c359-77ed-4ad3-bd73-a42c3c89f9c0",
      "key": "7e3515f8-2428-4071-af4f-7616fa2526ec.XIFF",
      "type": "multipart",
      "uploadId": "d30cdafad3675935636ccc51f679ae9f",
      "partSize": 64000000,
      "parts": [
        {
          "partNumber": 1,
          "url": "<presigned url>"
        },
        {
          "partNumber": 2,
          "url": "<presigned url>"
        }
      ]
    },
    {
      "clientId": "56f85dae-9eb9-4b8e-97ef-8647b826ccab",
      "key": "9a2c4b1d-5f6e-4a7b-8c3d-2e1f0a9b8c7d.xls",
      "type": "singlepart",
      "url": "<presigned url>"
    },
    {
      "clientId": "e7b0f2c1-1d2c-4c3a-9b2e-8f3a5d6c7e8f",
      "key": "3f8a2b1c-9d4e-4b6a-8c1e-7a5f3d2b9c1e.jpg",
      "type": "singlepart",
      "url": "<presigned url>"
    }
  ]
}
```

### `POST /files/upload/:clientId/complete`

- Auth Required: True
- Request Body: JSON for Multipart uploads

This endpoint confirms that a file finished uploading and creates the `files` record. The server looks up the `uploads` row by `clientId` and verifies it belongs to the authenticated user, which prevents the client from forging the file key, id, or size.

For a single-part upload, the server verifies the object exists in s3 (e.g. `HeadObject`) and that its size matches the recorded size, then creates the `files` row. For a multipart upload, the client supplies the list of uploaded parts and the server calls `CompleteMultipartUpload` with that list, then verifies existence and creates the `files` row. In both cases the `uploads` row is set to `status = 'completed'`.

**Request Body** (multipart only; single-part sends an empty body)

```json
{
  "parts": [
    {
      "partNumber": 1,
      "etag": "a8f2c1d9e0b4f6a7c8d3e5b1f0a9c2e4"
    },
    {
      "partNumber": 2,
      "etag": "b7e1c0d8f9a4b5c6d2e4f0a1b9c8d7e3"
    }
  ]
}
```

### `POST /files/upload/:clientId/abort`

- Auth Required: True

This endpoint cancels an in-progress upload. If the `uploads` row has an `upload_id`, the server calls `AbortMultipartUpload` and sets `status = 'aborted'`. Used when the user cancels an upload or it fails partway through.

### `PATCH /files/:fileId`

- Auth Required: True
- Request Body: JSON

This endpoint is to update fields like filename and expiration date.

```JSON
{
  "filename": "new name.jpg",
  "expiresAt": "2026-08-25T16:51:47.155Z"
}
```

### `DELETE /files/:fileId`

- Auth Required: True
- Request Body: False

## Share related routes

- `GET /api/shares`
- `GET /api/shares/:shareId`
- `POST /api/shares/:shareId/unlock`
- `POST /api/shares`
- `PUT /api/shares/:shareId`
- `DELETE /api/shares/:shareId`

### `GET /shares`

- Auth Required: True
- Request Body: False
- Query Params:
  - page: integer. defaults to 1 when not present
  - sort: string: 'name' | 'createdAt' | 'expiresOn'
- Response Body: True

This endpoint gets all of the shares info for currently logged in user. There should be a limit of 20 shares per page returned

**Response Body:**

```json
{
  "page": 1,
  "totalPages": 1,
  "shares": [
    {
      "id": "a169c76a-fd18-4c1a-a8ca-66abc0e6ae63",
      "name": "Pheonix Project",
      "protected": "Password",
      "createdAt": "2026/09/02 10:18 AM",
      "expiresAt": "2026/10/02 10:18 AM"
    },
    {
      "id": "2a833076-5b41-4e3f-b548-7e4349dacf25",
      "name": "Ocelot Project",
      "protected": "None",
      "createdAt": "2026/09/02 10:19 AM",
      "expiresAt": null
    },
    {
      "id": "3d8fd519-1fb8-4e1f-9be2-cdf258920221",
      "name": "Capybara Project",
      "protected": "Password",
      "createdAt": "2026/09/02 10:20 AM",
      "expiresAt": "2026/09/26 10:20 AM"
    }
  ]
}
```

### `GET /shares/:shareId`

- Auth Required: Only ifpassword protected
- Response Body: JSON

This endpoint is multi-faceted. If the share is not password protected, it returns the files with presigned s3 download urls. No auth header required. However, is the share is password protected, it requires a jwt token passed into the header that either `sub` is the owner of the share or has `shareId` set to the `:shareId`. This ensures that public invidiuals with no account can still download password protected shares.

The front end will show a "Password required" message is the share is protected. the endpoint `POST /shares/:shareId/unlock` should be submitted with the password. the returned jwt will be used by the frontend to insert the `Authorization` header when calling `GET /shares/:shareId` again to get the file urls.

File download urls should expire after 15 minutes

**Response Body**:

```json
{
  "name": "Phoenix Project",
  "expiresOn": "2026-09-01T16:26:37.948Z",
  "files": [
    {
      "filename": "filename.jpg",
      "size": "7.8 Mb",
      "filetype": "Image",
      "url": "<presigned download url>"
    },
    {
      "filename": "final-report.xls",
      "size": "3.1 Mb",
      "filetype": "Doc",
      "url": "<presigned download url>"
    }
  ]
}
```

### `POST /shares/:shareId/unlock`

- Auth Required: False
- Request Body: JSON
- Response Body: JSON

This endpoint will be used to unlock password protected shares. It takes in a password in the request body and returns a jwt that grants access to that share for 15 minutes.

The password is checked against the hashed password in the `shares` DB table. If it passes then the jwt access token is issued. The JWT should have a field called `shareId` that has the shareId to verify later they have permission for that specific share. Inspect the jwt token below for an example, you can use [jwt.io](https://www.jwt.io) to inspect it.

**Request Body:**

```json
{
  "password": "secret-password"
}
```

**Response Body:**

```json
{
  "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIiLCJpc3MiOiJGaWxlU2VuZGVyIiwic2hhcmVJZCI6ImQ5Y2NkYmE4LTM0ZTMtNDcwMi1hYzAzLTEwNmM2NjE3MDkzNSIsImlhdCI6MTUxNjIzOTAyMiwiZXhwIjoxNzg4MjgxMDYzfQ.DaG6-xW0z3DtCXqmVvZK_dbIXOwvDFJvo0Lmg1oFisA"
}
```

### `POST /shares`

- Auth Required: True
- Request Body: JSON
  - name: `string`. Share name
  - password: `string`. Omit if no password will be required
  - expires_on: `string`. Datetime ISO string. Omit if share will be permanent
  - files: `[]string`. list of file IDs from logged in users file pool.
- Response Body: JSON only on error

This endpoint creates a new share by the current user. This request body must include the IDs of the files the share will have. Files will be checked to see if they exist, if they do, they'll be linked to a share via the `shares_files` DB table. The password passed in will create a argon2id hash that will be saved in `shares` DB table. Passwords should be at least 8 chars in length.

**Request Body**

```json
{
  "name": "Phoenix Project",
  "password": "mystrongpassword",
  "expiresOn": "2026-08-31T16:10:56.923Z",
  "files": [
    "1fdce82d-3aaa-4239-ba85-015474662708",
    "712d27eb-a72c-4616-b0cd-da5b6642215e",
    "dada4172-8b1c-4084-9414-eb67539b274c",
    "bbb1e2aa-263d-430d-aa5b-d4c80bdc23aa"
  ]
}
```

On successful creation, only a successful status code will be returned: `201`, no body.

### `PUT /shares/:shareId`

- Auth Required: True
- Request Body: JSON
  - name: `string`. Omit if no update to name
  - expiresOn: `string`. Datetime string ISO. Omit if no update.
  - password: `string`. Omit if no change. if set, it will update or add the new password.
  - files: string[]. List of file ids to be in share.
- Response Body: JSON only on error

Pretty much what it looks like. One detail though is that the `files` field on this endpoint are all the ids of files to be in the share. That means that we will need to do 2 passes.

1. remove any files not found in this list
2. add any files in this list that are not present in the DB records

**Response Body:**

```JSON
{
  "name": "The share name",
  "expiresOn": "2026-09-02T15:38:13.642Z",
  "password": "",
  "files": [
    "13092e7f-3ddb-4839-b33b-967b3cc5bd09",
    "5cc7da9c-21ea-446d-8339-7ea92d43058a",
    "ab2e76bb-9cb5-4828-95d9-a246f474ed5c",
    "2d3da56b-c760-47dd-9b35-641621580457"
  ]
}
```

### `DELETE /shares/:shareId`

- Auth Required: True
- Request Body: False
- Response Body: False

This endpoint deletes the share if it's owned by the logged in user or is admin. Since `shares_files` DB table has cascade delete when the file or share is deleted, we only need to delete the share outright.

## Membership routes

- [x] `POST /api/membership`

### `POST /api/membership`

- Auth Required: True
- Request Body: JSON
- Response Body: False

This endpoint changes the membership of a user. This is done by reading the JSON payload with the plan selected and the secret password to be accepted. The password should be a argon2id hash stored as an ENV variable. A special program should be used to create these hashed that will be used as ENV variables. When a user sends the cleartext passwords, the argon2id hash will be compared to the hash of the cleartext password from the request body.

**JSON Body Example**

```json
{
  "membership": "pro",
  "password": "somesecretpassword"
}
```
