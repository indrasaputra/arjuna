-- Create "accounts" table
CREATE TABLE public.accounts (id uuid NOT NULL, user_id uuid NOT NULL, email text NOT NULL, password text NOT NULL, created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, deleted_at timestamp NULL, created_by uuid NOT NULL, updated_by uuid NOT NULL, deleted_by uuid NULL, PRIMARY KEY (id), CONSTRAINT accounts_email_key UNIQUE (email), CONSTRAINT accounts_user_id_key UNIQUE (user_id), CONSTRAINT email_length CHECK (LENGTH(email) <= 255), CONSTRAINT password_length CHECK (LENGTH(email) <= 255));
-- Create index "index_on_accounts_on_email" to table: "accounts"
CREATE INDEX index_on_accounts_on_email ON public.accounts (email);
