-- USERS
CREATE TABLE users (
    id UUID PRIMARY KEY NOT NULL,
    name VARCHAR(64) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
COMMENT ON TABLE users IS 'File system users.';
