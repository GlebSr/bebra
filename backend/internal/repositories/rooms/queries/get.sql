-- name: Get :one
SELECT *
FROM rooms
WHERE id = $1;