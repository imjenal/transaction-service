DROP TRIGGER IF EXISTS set_updated_at_on_transactions_update
    ON public.transactions;

DROP TABLE IF EXISTS public.transactions;