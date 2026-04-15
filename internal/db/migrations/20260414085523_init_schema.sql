-- +goose Up
CREATE TABLE users (
    id uuid PRIMARY KEY DEFAULT uuidv7(),
    name text NOT NULL,
    role text NOT NULL DEFAULT 'guest',
    image text,
    email text UNIQUE,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz,
    CONSTRAINT valid_role CHECK (role IN ('guest', 'author'))
);
CREATE TABLE account (
    id uuid PRIMARY KEY DEFAULT uuidv7(),
    user_id uuid NOT NULL REFERENCES users(id),
    provider text NOT NULL,
    provider_account_id text NOT NULL,
    username text,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz,
    UNIQUE(provider, provider_account_id),
    CONSTRAINT valid_provider CHECK (provider IN ('github'))
);
CREATE TABLE guestbook (
    id uuid PRIMARY KEY DEFAULT uuidv7(),
    text text NOT NULL,
    user_id uuid NOT NULL REFERENCES users(id),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz
);
CREATE TABLE content (
    id uuid PRIMARY KEY DEFAULT uuidv7(),
    slug text NOT NULL UNIQUE,
    type text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz,
    CONSTRAINT valid_content_type CHECK (type IN ('blog', 'project'))
);
CREATE TABLE ip_address (
    id uuid PRIMARY KEY DEFAULT uuidv7(),
    ip_address text UNIQUE NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz
);
CREATE TABLE content_view (
    id uuid PRIMARY KEY DEFAULT uuidv7(),
    content_id uuid NOT NULL REFERENCES content(id),
    ip_address_id uuid NOT NULL REFERENCES ip_address(id),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz,
    UNIQUE(content_id, ip_address_id)
);
CREATE TABLE content_like (
    id uuid PRIMARY KEY DEFAULT uuidv7(),
    content_id uuid NOT NULL REFERENCES content(id),
    ip_address_id uuid NOT NULL REFERENCES ip_address(id),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz
);
-- +goose Down
DROP TABLE content_like;
DROP TABLE content_view;
DROP TABLE ip_address;
DROP TABLE content;
DROP TABLE guestbook;
DROP TABLE account;
DROP TABLE users;