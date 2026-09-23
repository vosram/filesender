-- +goose up
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
  ON UPDATE CASCADE,
  PRIMARY KEY (file_id, share_id)
);

-- +goose down
DROP TABLE shares_files;