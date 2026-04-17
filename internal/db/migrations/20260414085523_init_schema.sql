-- +goose Up
CREATE TABLE users (
    id uuid PRIMARY KEY DEFAULT uuidv7(),
    name text NOT NULL,
    role text NOT NULL DEFAULT 'guest',
    image text,
    email text UNIQUE,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz,
    CONSTRAINT check_users_role CHECK (role IN ('guest', 'author'))
);

CREATE TABLE account (
    id uuid PRIMARY KEY DEFAULT uuidv7(),
    user_id uuid NOT NULL REFERENCES users(id),
    provider text NOT NULL,
    provider_account_id text NOT NULL,
    username text,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz,
    UNIQUE(provider, provider_account_id),
    CONSTRAINT check_account_provider CHECK (provider IN ('github'))
);

CREATE TABLE guestbook (
    id uuid PRIMARY KEY DEFAULT uuidv7(),
    text text NOT NULL,
    user_id uuid NOT NULL REFERENCES users(id),
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz
);

CREATE TABLE ip_address (
    id uuid PRIMARY KEY DEFAULT uuidv7(),
    ip_address text UNIQUE NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz
);

CREATE TABLE content (
    id uuid PRIMARY KEY DEFAULT uuidv7(),
    slug text NOT NULL UNIQUE,
    kind text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz,
    CONSTRAINT check_content_kind CHECK (kind IN ('blog', 'project'))
);

CREATE TABLE content_meta (
    id uuid PRIMARY KEY DEFAULT uuidv7(),
    content_id uuid NOT NULL REFERENCES content(id),
    ip_address_id uuid NOT NULL REFERENCES ip_address(id),
    views int NOT NULL DEFAULT 0,
    likes int NOT NULL DEFAULT 0,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    UNIQUE(content_id, ip_address_id),
    CONSTRAINT check_content_meta_views CHECK (views >= 0),
    CONSTRAINT check_content_meta_likes CHECK (
        likes BETWEEN 0 AND 5
    )
);

-- +goose Down
DROP TABLE content_meta;

DROP TABLE ip_address;

DROP TABLE content;

DROP TABLE guestbook;

DROP TABLE account;

DROP TABLE users;