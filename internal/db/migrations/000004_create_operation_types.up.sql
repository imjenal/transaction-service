CREATE TYPE public.transaction_type AS ENUM ('NORMAL_PURCHASE', 'WITHDRAWAL', 'CREDIT_VOUCHER', 'PURCHASE_WITH_INSTALLMENTS');
CREATE TYPE public.amount_behavior AS ENUM ('POSITIVE', 'NEGATIVE');

CREATE TABLE IF NOT EXISTS public.operation_types
(
    uuid            UUID PRIMARY KEY         NOT NULL DEFAULT gen_random_uuid(),
    serial_id       BIGSERIAL UNIQUE         NOT NULL,
    description     public.transaction_type  NOT NULL,
    amount_behavior public.amount_behavior   NOT NULL DEFAULT 'NEGATIVE',
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TRIGGER set_updated_at_on_operation_types_update
    BEFORE UPDATE
    ON public.operation_types
    FOR EACH ROW
EXECUTE PROCEDURE set_updated_at();
