-- GROUPS IDENTITY COLLECTION & MANAGEMENT
CREATE TABLE groups (
    id UUID PRIMARY KEY NOT NULL,
    name VARCHAR(64) NOT NULL
);

CREATE TABLE group_memberships (
    group_id UUID REFERENCES groups(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT pk_group_memberships PRIMARY KEY (group_id, user_id)
);