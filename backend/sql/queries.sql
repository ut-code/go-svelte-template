-- name: GetWordByDate :one
SELECT * FROM words WHERE date = $1;

-- name: GetTodayWord :one
SELECT * FROM words WHERE date = CURRENT_DATE;

-- name: ListWords :many
SELECT * FROM words ORDER BY date DESC;

-- name: CreateWord :one
INSERT INTO words (word, date) VALUES ($1, $2) RETURNING *;

-- name: CreateScore :one
INSERT INTO scores (player_name, word_id, attempts, solved)
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetScoresByWord :many
SELECT * FROM scores WHERE word_id = $1 ORDER BY attempts ASC;

-- name: GetScoresByPlayer :many
SELECT s.*, w.word, w.date
FROM scores s
JOIN words w ON s.word_id = w.id
WHERE s.player_name = $1
ORDER BY w.date DESC;
