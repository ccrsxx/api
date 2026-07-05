-- +goose Up
CREATE TABLE contact (
    id uuid PRIMARY KEY DEFAULT uuidv7(),
    name text NOT NULL,
    email text NOT NULL,
    message text NOT NULL,
    delivered_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz
);

-- +goose Down
DROP TABLE contact;