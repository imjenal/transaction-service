ALTER TABLE public.transactions
    ADD COLUMN balance FLOAT DEFAULT 0.0 NOT NULL;