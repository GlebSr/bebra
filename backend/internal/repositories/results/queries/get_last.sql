-- name: GetLastResult :one
SELECT * FROM random_results
WHERE room_id = $1
ORDER BY created_at DESC
LIMIT 1;