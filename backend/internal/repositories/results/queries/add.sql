-- name: Add :one
INSERT INTO random_results (
    id, room_id, game_id, chosen_by
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;