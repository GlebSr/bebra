-- name: GetAllParticipants :many
SELECT u.id, u.name, rp.role
FROM room_participants rp
JOIN users u ON rp.user_id = u.id
WHERE rp.room_id = $1;