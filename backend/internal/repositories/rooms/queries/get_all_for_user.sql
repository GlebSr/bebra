-- name: GetAllForUser :many
SELECT r.*
FROM rooms r
JOIN room_participants rp ON rp.room_id = r.id
WHERE rp.user_id = $1;