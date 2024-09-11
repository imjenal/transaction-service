-- name: CreateTransaction :one
INSERT INTO public.transactions (account_id, amount, operation_type_id, balance, event_date)
VALUES ($1, $2, $3, $4, NOW())
RETURNING uuid, serial_id, account_id, amount, operation_type_id, event_date, balance, updated_at;

-- name: GetTransactionDetailsByTransactionId :one
SELECT uuid, serial_id, account_id, amount, operation_type_id, event_date, balance, updated_at
FROM public.transactions
WHERE uuid = $1;

-- name: GetNegativeBalanceTransactionsByAccountID :many
SELECT uuid, account_id, operation_type_id, amount, balance, event_date FROM public.transactions
WHERE  account_id = $1 AND balance < 0
ORDER BY event_date;

-- name: UpdateTransactionBalances :exec
UPDATE public.transactions SET balance = $2 WHERE uuid = $1;