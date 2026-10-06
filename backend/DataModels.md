# Database Models

## Users Model

Credential Type is a string but the backend will have set values to input. the values will be as follows:

- "magic_email"
- "google"

Membership is a text but the backend will have set values to input. The values will be as follows:

- "basic"
- "pro"
- "premium"

Role are the permissions role a user has. It's just a text field, the backend will have set have fixed roles:

- "user"
- "admin"

```sql
CREATE TABLE users(
  id UUID PRIMARY KEY,
  name TEXT NOT NULL,
  email TEXT UNIQUE NOT NULL,
  credential_type TEXT NOT NULL,
  membership TEXT NOT NULL DEFAULT 'basic',
  banned_until TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  role TEXT NOT NULL DEFAULT 'user',
  last_login TIMESTAMPTZ
);
```

## Refresh_Tokens

Refresh tokens will be 32 byte hexidecimal safe generated from `crypto/rand`. These refresh_tokens should only be stored in httpOnly cookies. the refreshToken saved in DB should be a SHA256 hash of the token

```sql
CREATE TABLE refresh_tokens(
  token_hash TEXT PRIMARY KEY,
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

Email login tokens will be 32 byte hexidecimal safe generated from `crypto/rand`. The token will be added to a url in an email link. The token saved to DB should be a SHA256 hash of this token.

```sql
CREATE TABLE email_login_tokens(
  token_hash TEXT PRIMARY KEY,
  email TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  expires_at TIMESTAMPTZ NOT NULL
);
```

## Email Change Verify Tokens

```sql
CREATE TABLE email_change_verify_tokens(
  token_hash TEXT PRIMARY KEY,
  user_id UUID NOT NULL,
  new_email TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  expires_at TIMESTAMPTZ NOT NULL,
  FOREIGN KEY (user_id)
  REFERENCES users(id)
  ON DELETE CASCADE
  ON UPDATE CASCADE
);
```

## Email Change Confirmation Tokens

```sql
CREATE TABLE email_change_confirmation_tokens(
  token_hash TEXT PRIMARY KEY,
  user_id UUID NOT NULL,
  new_email TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  expires_at TIMESTAMPTZ NOT NULL,
  FOREIGN KEY (user_id)
  REFERENCES users(id)
  ON DELETE CASCADE
  ON UPDATE CASCADE
);
```

## Files Model

`filetype` is just a regular text field. The backend will have set fixed values such as:

- "data"
- "image"
- "video"
- "document"
- "archive"

```sql
CREATE TABLE files(
  id UUID PRIMARY KEY,
  filename TEXT NOT NULL,
  size INTEGER NOT NULL,
  mime_type TEXT NOT NULL,
  key TEXT NOT NULL,
  owner_id UUID NOT NULL,
  file_type TEXT NOT NULL DEFAULT 'data',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  expires_at TIMESTAMPTZ,
  FOREIGN KEY (owner_id)
  REFERENCES users(id)
  ON DELETE CASCADE
  ON UPDATE CASCADE
);
```

## Uploads Model

Upload intents are created when the client requests presigned urls (`POST /files/upload`). They record the server-generated key, the s3 `upload_id` for multipart uploads, and the owner so the completion step can be validated server-side. No `files` record exists until the upload is confirmed.

`status` is a regular text field and the backend will have fixed values such as:

- "initiated"
- "in_progress"
- "completed"
- "aborted"

```sql
CREATE TABLE uploads(
  client_id UUID PRIMARY KEY,
  owner_id UUID NOT NULL,
  file_name TEXT NOT NULL,
  size BIGINT NOT NULL,
  mime_type TEXT NOT NULL,
  key TEXT NOT NULL,
  upload_id TEXT,
  parts INTEGER NOT NULL DEFAULT 1,
  status TEXT NOT NULL DEFAULT 'initiated',
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
CREATE TABLE shares(
  id UUID PRIMARY KEY,
  name TEXT NOT NULL,
  password_hash TEXT,
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
CREATE TABLE shares_files(
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
