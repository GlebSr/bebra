-- name: Clear :exec
DELETE FROM refresh_tokens WHERE expires_at < NOW();