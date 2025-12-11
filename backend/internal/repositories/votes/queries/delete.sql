-- name: Delete :exec
DELETE FROM votes
WHERE id = $1;