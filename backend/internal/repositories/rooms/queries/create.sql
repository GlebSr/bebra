-- name: Create :one
INSERT INTO rooms (
    id, name, owner_id
) VALUES (
    $1, $2, $3
)
RETURNING *;