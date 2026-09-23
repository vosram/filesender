-- name: CreateUserAndReturnId :one
INSERT INTO users (id, name, email, credential_type)
VALUES (
  gen_random_uuid(),
  $1,
  $2,
  $3
)
RETURNING id;

-- name: UpdateLastLoginByUserID :exec
UPDATE users
SET last_login = NOW()
WHERE id = $1;
