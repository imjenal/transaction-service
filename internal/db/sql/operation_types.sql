-- name: GetOperationTypeAmountBehavior :one
SELECT amount_behavior FROM public.operation_types WHERE serial_id = $1;