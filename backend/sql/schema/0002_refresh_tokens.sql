-- +goose up
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

-- +goose down
DROP TABLE refresh_tokens;