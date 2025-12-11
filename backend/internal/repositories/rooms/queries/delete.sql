-- name: Delete :exec
DELETE FROM rooms
WHERE id = $1;