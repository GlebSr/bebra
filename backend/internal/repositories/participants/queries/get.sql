-- name: Get :one
SELECT * FROM room_participants
WHERE room_id = $1 AND user_id = $2;