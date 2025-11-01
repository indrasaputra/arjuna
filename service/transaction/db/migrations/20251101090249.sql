-- Create "transactions" table
CREATE TABLE public.transactions (id uuid NOT NULL, sender_id uuid NOT NULL, receiver_id uuid NOT NULL, amount numeric NOT NULL, created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, deleted_at timestamp NULL, created_by uuid NOT NULL, updated_by uuid NOT NULL, deleted_by uuid NULL, PRIMARY KEY (id));
