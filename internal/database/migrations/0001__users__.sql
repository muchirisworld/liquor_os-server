-- +goose up

CREATE TABLE "user" (
	"id" varchar PRIMARY KEY NOT NULL,
	"name" text NOT NULL,
	"email" email_citext NOT NULL UNIQUE,
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