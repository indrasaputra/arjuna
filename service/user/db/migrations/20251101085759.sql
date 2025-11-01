-- Create enum type "user_outbox_status"
CREATE TYPE public.user_outbox_status AS ENUM ('READY', 'PROCESSED', 'DELIVERED', 'FAILED');
-- Create "users" table
CREATE TABLE public.users (id uuid NOT NULL, name text NOT NULL, created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, deleted_at timestamp NULL, created_by uuid NOT NULL, updated_by uuid NOT NULL, deleted_by uuid NULL, PRIMARY KEY (id), CONSTRAINT name_length CHECK (LENGTH(name) <= 100));
-- Create index "index_on_users_on_id" to table: "users"
CREATE INDEX index_on_users_on_id ON public.users (id);
-- Create "users_outbox" table
CREATE TABLE public.users_outbox (id uuid NOT NULL, payload jsonb NOT NULL, status public.user_outbox_status NOT NULL DEFAULT 'READY', created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, deleted_at timestamp NULL, created_by uuid NOT NULL, updated_by uuid NOT NULL, deleted_by uuid NULL, PRIMARY KEY (id));
-- Create index "index_on_users_outbox_on_status_and_created_at" to table: "users_outbox"
CREATE INDEX index_on_users_outbox_on_status_and_created_at ON public.users_outbox (status, created_at);
