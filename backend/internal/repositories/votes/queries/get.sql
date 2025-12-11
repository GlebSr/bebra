-- name: Get :one
SELECT *
FROM votes
WHERE id = $1;