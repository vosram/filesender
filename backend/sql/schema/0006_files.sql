-- +goose up
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

-- +goose down
DROP TABLE files;