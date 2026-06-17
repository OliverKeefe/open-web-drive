-- GROUPS IDENTITY COLLECTION & MANAGEMENT
CREATE TABLE groups (
    id UUID PRIMARY KEY NOT NULL,
    name VARCHAR(64) NOT NULL
);

CREATE TABLE group_memberships (
    group_id REFERENCES groups(id) ON DELETE CASCADE,
    user_id REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (group_id, user_id)
);