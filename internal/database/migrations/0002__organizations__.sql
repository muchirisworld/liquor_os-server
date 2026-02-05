-- +goose Up

CREATE TABLE organizations (
    "id" varchar PRIMARY KEY NOT NULL,
    "name" text NOT NULL,
    "slug" text NOT NULL UNIQUE,
    "created_by" varchar NOT NULL,
    "logo_url" text,
    "image_url" text,
    "metadata" jsonb DEFAULT '{}'::jsonb,
    "max_allowed_memberships" INTEGER NOT NULL DEFAULT 5,
    "created_at" timestamptz DEFAULT now() NOT NULL,
    "updated_at" timestamptz DEFAULT now() NOT NULL,
    
    CONSTRAINT "fk_organization_creator"
        FOREIGN KEY("created_by")
        REFERENCES users("id")
        ON DELETE CASCADE
);

-- +goose StatementBegin
CREATE TRIGGER trg_organizations_updated_at
    BEFORE UPDATE ON organizations
    FOR EACH ROW
    EXECUTE FUNCTION trg_set_updated_at();
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS organizations;
