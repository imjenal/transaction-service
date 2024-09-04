-- name: UserExists :one
SELECT EXISTS(SELECT 1 FROM public.users WHERE uuid = $1) AS exists;
