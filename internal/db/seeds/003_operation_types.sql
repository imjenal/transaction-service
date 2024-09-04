-- Insert "Normal Purchase"
INSERT INTO public.operation_types (description, amount_behavior)
VALUES ('NORMAL_PURCHASE', 'NEGATIVE')
RETURNING uuid, serial_id, description, amount_behavior, created_at, updated_at;

-- Insert "Purchase with Installments"
INSERT INTO public.operation_types (description,amount_behavior)
VALUES ('PURCHASE_WITH_INSTALLMENTS', 'NEGATIVE')
RETURNING uuid, serial_id, description,amount_behavior, created_at, updated_at;

-- Insert "Withdrawal"
INSERT INTO public.operation_types (description,amount_behavior)
VALUES ('WITHDRAWAL', 'NEGATIVE')
RETURNING uuid, serial_id, description,amount_behavior, created_at, updated_at;

-- Insert "Credit Voucher"
INSERT INTO public.operation_types (description,amount_behavior)
VALUES ('CREDIT_VOUCHER', 'POSITIVE')
RETURNING uuid, serial_id, description,amount_behavior, created_at, updated_at;
