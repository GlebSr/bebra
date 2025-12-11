-- name: Delete :exec
DELETE FROM GAMES
WHERE id = $1;