-- +goose Up

CREATE TABLE members (
    "user_id" varchar NOT NULL,
    "organization_id" varchar NOT NULL,
    "identifier" text,
    "image_url" text,
    "profile_image_url" text,
    "role" text NOT NULL DEFAULT 'member',
    "role_name" text,
    "created_at" timestamptz DEFAULT now() NOT NULL,
    "updated_at" timestamptz DEFAULT now() NOT NULL,

    -- Composite Primary Key
    PRIMARY KEY ("user_id", "organization_id"),

    CONSTRAINT fk_member_user 
        FOREIGN KEY ("user_id") 
        REFERENCES users("id") 
        ON DELETE CASCADE,

    CONSTRAINT fk_member_organization 
        FOREIGN KEY ("organization_id") 
        REFERENCES organizations("id") 
        ON DELETE CASCADE
);

-- +goose StatementBegin
CREATE TRIGGER trg_members_updated_at
    BEFORE UPDATE ON members
    FOR EACH ROW
    EXECUTE FUNCTION trg_set_updated_at();
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS members;
