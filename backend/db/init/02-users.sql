-- USERS
CREATE TABLE users (
    id UUID PRIMARY KEY NOT NULL,
    name VARCHAR(64)
);
COMMENT ON TABLE users IS 'File system users.'
