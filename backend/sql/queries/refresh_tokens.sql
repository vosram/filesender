-- name: SaveRefreshToken :exec
INSERT INTO refresh_tokens (token_hash, user_id, expires_at)
VALUES ($1, $2, $3);

-- name: ConsumeRefreshToken :one
DELETE FROM refresh_tokens
WHERE token_hash = $1
RETURNING *;