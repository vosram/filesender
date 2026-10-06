-- +goose up
CREATE TABLE email_login_tokens(
  token_hash TEXT PRIMARY KEY,
  email TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  expires_at TIMESTAMPTZ NOT NULL
);

-- +goose down
DROP TABLE email_login_tokens;