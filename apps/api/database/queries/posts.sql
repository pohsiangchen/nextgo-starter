-- name: GetPost :one
SELECT * FROM posts
WHERE id = $1 LIMIT 1;

-- name: ListFeedsByUserId :many
SELECT u.username, sqlc.embed(p), COUNT(c.id) AS comments_count
FROM posts p 
LEFT JOIN "comments" c ON c.post_id = p.id 
LEFT JOIN users u ON u.id = p.user_id 
JOIN follows f ON f.followed_user_id = p.user_id
WHERE f.following_user_id = sqlc.arg('following_user_id') AND p.id < sqlc.arg('last_posts_id')
GROUP BY u.id, p.id
ORDER BY p.id DESC
LIMIT sqlc.arg('limit');

-- name: CreatePost :one
INSERT INTO posts (
  title, content, user_id
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: DeletePost :exec
DELETE FROM posts
WHERE id = $1;

-- name: UpdatePost :one
UPDATE posts
SET title = coalesce(sqlc.narg('title'), title),
content = coalesce(sqlc.narg('content'), content)
WHERE id = sqlc.arg('id')
RETURNING *;
