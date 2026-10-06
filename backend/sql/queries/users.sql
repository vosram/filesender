-- name: FindUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: FindUserById :one
SELECT * FROM users
WHERE id = $1;

-- name: CreateUserAndReturnUser :one
INSERT INTO users (id, name, email, credential_type, last_login)
VALUES (
  gen_random_uuid(),
  $1,
  $2,
  $3,
  now()
)
RETURNING *;

-- name: UpdateLastLoginByUserID :exec
UPDATE users
SET last_login = NOW()
WHERE id = $1;
