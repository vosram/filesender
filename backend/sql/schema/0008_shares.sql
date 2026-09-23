-- +goose up
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

-- +goose down
DROP TABLE shares;