-- name: Get :one
SELECT * FROM GAMES
WHERE id = $1;