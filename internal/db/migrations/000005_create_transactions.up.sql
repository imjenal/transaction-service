CREATE TABLE IF NOT EXISTS public.transactions
(
    uuid              UUID PRIMARY KEY         NOT NULL DEFAULT gen_random_uuid(),
    serial_id         BIGSERIAL UNIQUE         NOT NULL,
    account_id        UUID                     NOT NULL REFERENCES public.accounts (uuid),
    amount            FLOAT                    NOT NULL,
    operation_type_id BIGSERIAL                     NOT NULL REFERENCES public.operation_types (serial_id),
    event_date        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TRIGGER set_updated_at_on_transactions_update
    BEFORE UPDATE
    ON public.transactions
    FOR EACH ROW
EXECUTE PROCEDURE set_updated_at();
