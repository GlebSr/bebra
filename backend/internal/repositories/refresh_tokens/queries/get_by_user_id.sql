-- name: GetByUserID :many
SELECT * FROM refresh_tokens
WHERE user_id = $1;