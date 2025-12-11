-- name: Add :one
INSERT INTO votes (
    id, room_id, game_id, user_id
) VALUES (
    $1, $2, $3, $4
) RETURNING *;