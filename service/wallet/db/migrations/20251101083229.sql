-- Create "wallets" table
CREATE TABLE public.wallets (id uuid NOT NULL, user_id uuid NOT NULL, balance numeric NOT NULL, created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, deleted_at timestamp NULL, created_by uuid NOT NULL, updated_by uuid NOT NULL, deleted_by uuid NULL, PRIMARY KEY (id), CONSTRAINT non_negative_balance CHECK (balance >= (0)::numeric));
-- Create index "index_on_wallets_on_id_and_user_id" to table: "wallets"
CREATE INDEX index_on_wallets_on_id_and_user_id ON public.wallets (id, user_id);
