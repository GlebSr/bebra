-- name: GetAllResults :many
SELECT * FROM random_results
WHERE room_id = $1;