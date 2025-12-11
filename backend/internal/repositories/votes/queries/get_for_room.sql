-- name: GetForRoom :many
SELECT *
FROM votes
WHERE room_id = $1;