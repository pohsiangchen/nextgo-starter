-- name: CreateFollow :one
INSERT INTO follows (
  following_user_id, followed_user_id
) VALUES (
  $1, $2
)
RETURNING *;

-- name: DeleteFollow :exec
DELETE FROM follows
WHERE following_user_id = $1 AND followed_user_id = $2;
