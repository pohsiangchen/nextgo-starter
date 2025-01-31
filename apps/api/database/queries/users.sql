-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id;

-- name: CreateUser :one
INSERT INTO users (
  email, username, "password"
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: UpdateUser :one
UPDATE users
SET username = coalesce(sqlc.narg('username'), username)
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: UpdateUserPassword :one
UPDATE users
SET "password" = $2
WHERE id = $1
RETURNING *;
