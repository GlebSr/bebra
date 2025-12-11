-- name: Add :one
INSERT INTO room_participants (
    id, room_id, user_id, role
) VALUES (
    $1, $2, $3, $4
) RETURNING *;