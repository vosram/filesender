# Database Models

## Role Type ENUM

```sql
CREATE TYPE role_type
AS
ENUM('user', 'admin');
```

## Membership Type ENUM

```sql
CREATE TYPE membership_type
AS
ENUM('basic', 'pro', 'premium');
```

## filetype_type ENUM

```sql
CREATE TYPE filetype_type
AS
ENUM('data', 'doc', 'image', 'video');
```

## Users Model

```sql
CREATE TABLE IF NOT EXISTS users(
  id UUID PRIMARY KEY,
  name TEXT NOT NULL,
  email TEXT UNIQUE NOT NULL,
  membership membership_type NOT NULL DEFAULT 'basic',
  banned_until TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  role role_type NOT NULL DEFAULT 'user',
  last_login TIMESTAMPTZ
);
```

## Refresh_Tokens

Refresh tokens will be 32 byte base64 url safe generated from `crypto/rand`. These refresh_tokens should only be stored in httpOnly cookies.

```sql
CREATE TABLE IF NOT EXISTS refresh_tokens(
  token TEXT PRIMARY KEY,
  user_id UUID NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  expires_at TIMESTAMPTZ NOT NULL,
  FOREIGN KEY (user_id)
  REFERENCES users(id)
  ON DELETE CASCADE
  ON UPDATE CASCADE
);
```

## Email Login Tokens

Email login tokens will be 32 byte base64 url safe generated from `crypto/rand`.

```sql
CREATE TABLE IF NOT EXISTS refresh_tokens(
  token TEXT PRIMARY KEY,
  user_id UUID NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  expires_at TIMESTAMPTZ NOT NULL,
  FOREIGN KEY (user_id)
  REFERENCES users(id)
  ON DELETE CASCADE
  ON UPDATE CASCADE
);
```

## Files Model

```sql
CREATE TABLE IF NOT EXISTS files(
  id UUID PRIMARY KEY,
  filename TEXT NOT NULL,
  size INTEGER NOT NULL,
  mime_type TEXT NOT NULL,
  key TEXT NOT NULL,
  owner_id UUID NOT NULL,
  filetype filetype_type NOT NULL DEFAULT 'data',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  expires_at TIMESTAMPTZ,
  FOREIGN KEY (owner_id)
  REFERENCES users(id)
  ON DELETE CASCADE
  ON UPDATE CASCADE
);
```

## upload_status_type ENUM

```sql
CREATE TYPE upload_status_type
AS
ENUM('initiated', 'in_progress', 'completed', 'aborted');
```

## Uploads Model

Upload intents are created when the client requests presigned urls (`POST /files/upload`). They record the server-generated key, the s3 `upload_id` for multipart uploads, and the owner so the completion step can be validated server-side. No `files` record exists until the upload is confirmed.

```sql
CREATE TABLE IF NOT EXISTS uploads(
  client_id UUID PRIMARY KEY,
  owner_id UUID NOT NULL,
  filename TEXT NOT NULL,
  size INTEGER NOT NULL,
  mime_type TEXT NOT NULL,
  key TEXT NOT NULL,
  upload_id TEXT,
  parts INTEGER NOT NULL DEFAULT 1,
  status upload_status_type NOT NULL DEFAULT 'initiated',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  FOREIGN KEY (owner_id)
  REFERENCES users(id)
  ON DELETE CASCADE
  ON UPDATE CASCADE
);
```

## Shares Model

```sql
CREATE TABLE IF NOT EXISTS shares(
  id UUID PRIMARY KEY,
  name TEXT NOT NULL,
  password TEXT,
  owner_id UUID NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  expires_on TIMESTAMPTZ,
  FOREIGN KEY (owner_id)
  REFERENCES users(id)
  ON DELETE CASCADE
  ON UPDATE CASCADE
);
```

## Shares_Files Model

```sql
CREATE TABLE IF NOT EXISTS shares_files(
  file_id UUID NOT NULL,
  share_id UUID NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  FOREIGN KEY (file_id)
  REFERENCES files(id)
  ON DELETE CASCADE
  ON UPDATE CASCADE,
  FOREIGN KEY (share_id)
  REFERENCES shares(id)
  ON DELETE CASCADE
  ON DELETE CASCADE,
  PRIMARY KEY (file_id, share_id)
);
```
