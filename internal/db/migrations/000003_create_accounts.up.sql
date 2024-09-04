CREATE TABLE IF NOT EXISTS public.accounts
(
    uuid                              UUID PRIMARY KEY         NOT NULL DEFAULT gen_random_uuid(),
    serial_id                         BIGSERIAL UNIQUE         NOT NULL,
    document_number                  VARCHAR(255)             NOT NULL,
    current_balance                   FLOAT                    NOT NULL,
    user_id                           UUID                     NOT NULL REFERENCES public.users (uuid),
    created_at                        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at                        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TRIGGER set_updated_at_on_accounts_update
    BEFORE UPDATE
    ON public.accounts
    FOR EACH ROW
EXECUTE PROCEDURE set_updated_at();
