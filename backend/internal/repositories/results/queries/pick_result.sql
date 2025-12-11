-- name: PickResult :one
WITH votes_per_game AS (
    SELECT
        g.id,
        COUNT(v.id) AS votes
    FROM games g
    LEFT JOIN votes v ON v.game_id = g.id
    WHERE g.room_id = $1
    GROUP BY g.id
),
sum_votes AS (
    SELECT SUM(votes) AS total_votes
    FROM votes_per_game
),
ordered AS (
    SELECT
        id,
        votes,
        SUM(votes) OVER (ORDER BY id) AS cum_votes
    FROM votes_per_game
),
rand AS (
    SELECT random() * (SELECT total_votes FROM sum_votes) AS r
)
SELECT id
FROM ordered, rand
WHERE cum_votes >= r
ORDER BY cum_votes
LIMIT 1;