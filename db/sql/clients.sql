-- name: CreateClient :one
INSERT INTO clients(fio, phone) VALUES($1, $2) RETURNING *;

