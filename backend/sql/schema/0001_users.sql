-- +goose up
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

-- +goose down
DROP TABLE users;