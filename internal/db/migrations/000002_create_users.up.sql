CREATE TABLE IF NOT EXISTS public.users
(
    uuid                              UUID PRIMARY KEY         NOT NULL DEFAULT gen_random_uuid(),
    serial_id                         BIGSERIAL UNIQUE         NOT NULL,
    first_name                  VARCHAR(255),
    middle_name       VARCHAR(255),
    last_name         VARCHAR(255),
    phone_number      VARCHAR(50) UNIQUE                      NOT NULL,
    email             VARCHAR(255) UNIQUE,
    created_at                        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at                        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TRIGGER set_updated_at_on_users_update
    BEFORE UPDATE
    ON public.users
    FOR EACH ROW
EXECUTE PROCEDURE set_updated_at();
