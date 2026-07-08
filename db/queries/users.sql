-- name: CreateUser :one
INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: CreateSession :one
INSERT INTO sessions (id, user_id, refresh_token, expires_at) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetSession :one
SELECT * FROM sessions WHERE id = $1;
