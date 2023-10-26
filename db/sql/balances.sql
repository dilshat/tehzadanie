
-- name: CreateClientBalance :one
INSERT INTO balances(client_id, balance) VALUES($1, $2) RETURNING *;

-- name: Deposit :one
UPDATE balances SET balance = balance + $1 WHERE client_id = $2 RETURNING balance;

-- name: Withdraw :execrows
UPDATE balances SET balance = balance - $1 WHERE client_id = $2 AND balance - $1 >= 0;

-- name: GetClientBalance :one
SELECT balance FROM balances WHERE client_id = $1 LIMIT 1;