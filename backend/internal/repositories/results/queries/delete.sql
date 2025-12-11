-- name: Delete :exec
DELETE FROM random_results
WHERE id = $1;