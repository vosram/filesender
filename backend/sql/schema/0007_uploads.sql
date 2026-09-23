-- +goose up
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

-- +goose down
DROP TABLE uploads;