-- name: CreateTransaction :one
INSERT INTO public.transactions (account_id, amount, operation_type_id, event_date)
VALUES ($1, $2, $3, NOW())
RETURNING uuid, serial_id, account_id, amount, operation_type_id, event_date, updated_at;

-- name: GetTransactionDetailsByTransactionId :one
SELECT uuid, serial_id, account_id, amount, operation_type_id, event_date, updated_at
FROM public.transactions
WHERE uuid = $1;