-- name: GetAllRoomGames :many
SELECT * FROM GAMES
WHERE room_id = $1;