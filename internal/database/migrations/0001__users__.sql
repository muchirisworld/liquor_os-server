-- +goose up

CREATE EXTENSION IF NOT EXISTS citext;

-- +goose statementbegin
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'email_citext') THEN
        CREATE DOMAIN email_citext AS citext
        CHECK (
            VALUE ~* '^[A-Z0-9._%+\-]+@[A-Z0-9.\-]+\.[A-Z]{2,63}$'
        );
    END IF;
END$$;

CREATE OR REPLACE FUNCTION trg_set_updated_at()
RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    NEW.updated_at := now();
    RETURN NEW;
END$$;

CREATE OR REPLACE FUNCTION trg_normalize_email()
RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    NEW.email := lower(NEW.email);
    RETURN NEW;
END$$;
-- +goose statementend

CREATE TABLE "user" (
	"id" varchar PRIMARY KEY NOT NULL,
	"name" text NOT NULL,
	"email" varchar NOT NULL UNIQUE,
	"email_verified" boolean DEFAULT false NOT NULL,
	"image" text,
	"created_at" timestamptz DEFAULT now() NOT NULL,
	"updated_at" timestamptz DEFAULT now() NOT NULL
);

-- +goose statementbegin
CREATE TRIGGER trg_users_updated_at
BEFORE UPDATE ON "user"
FOR EACH ROW EXECUTE FUNCTION trg_set_updated_at();

CREATE TRIGGER trg_users_normalize_email
BEFORE INSERT OR UPDATE ON "user"
FOR EACH ROW EXECUTE FUNCTION trg_normalize_email();
-- +goose statementend

-- +goose down
DROP TABLE IF EXISTS "user";
DROP DOMAIN IF EXISTS email_citext;