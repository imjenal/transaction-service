DROP TRIGGER IF EXISTS set_updated_at_on_operation_types_update
    ON public.operation_types;

DROP TABLE IF EXISTS public.operation_types;

DROP TYPE IF EXISTS public.transaction_type;