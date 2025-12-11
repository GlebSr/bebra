-- name: Add :one
INSERT INTO GAMES (
    id, room_id, title
) VALUES (
    $1, $2, $3
) RETURNING *;