-- name: GetComment :one
SELECT * FROM comments
WHERE id = $1 LIMIT 1;

-- name: ListCommentsByPostId :many
SELECT users.username, sqlc.embed(comments)
FROM comments
JOIN users ON users.id = comments.user_id
WHERE comments.post_id = sqlc.arg('post_id') AND comments.id < sqlc.arg('last_comments_id')
ORDER BY comments.id DESC
LIMIT sqlc.arg('limit');

-- name: CreateComment :one
INSERT INTO comments (
  content, user_id, post_id
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: DeleteComment :exec
DELETE FROM comments
WHERE id = $1;

-- name: UpdateComment :one
UPDATE comments
SET content = coalesce(sqlc.narg('content'), content)
WHERE id = sqlc.arg('id')
RETURNING *;
