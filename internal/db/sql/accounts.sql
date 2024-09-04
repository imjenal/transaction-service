-- name: CreateAccount :one
INSERT INTO public.accounts (document_number, current_balance, user_id)
VALUES ($1, $2, $3)
RETURNING uuid, serial_id, document_number, current_balance, user_id, created_at, updated_at;

-- name: GetAccountDetailsByUUID :one
SELECT uuid, serial_id, document_number, current_balance, user_id, created_at, updated_at
FROM public.accounts
WHERE uuid = $1;

-- name: AccountExists :one
SELECT EXISTS(SELECT 1 FROM public.accounts WHERE uuid = $1) AS exists;