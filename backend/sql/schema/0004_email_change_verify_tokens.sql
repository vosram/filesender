-- +goose up
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

-- +goose down
DROP TABLE email_change_verify_tokens;