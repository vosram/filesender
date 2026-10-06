-- name: SaveEmailLoginToken :exec
INSERT INTO email_login_tokens (token_hash, email, expires_at)
VALUES ($1, $2, $3);

-- name: ConsumeEmailLoginToken :one
DELETE FROM email_login_tokens
WHERE token_hash = $1
RETURNING *;