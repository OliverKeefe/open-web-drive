CREATE TABLE access_levels (
    id UUID PRIMARY KEY NOT NULL,
    name VARCHAR(16) NOT NULL
);

CREATE TABLE permissions (
    id UUID PRIMARY KEY NOT NULL,
    file_id REFERENCES files(id) NOT NULL,
    grantee_type VARCHAR(10) NOT NULL,
    grantee_id UUID NOT NULL,
    access_level REFERENCES permission_access_levels(id) NOT NULL -- May need to be altered again.
);

CREATE TABLE permission_access_levels (
    permission_id REFERENCES permissions(id),
    access_level_id REFERENCES access_levels(id),
    PRIMARY KEY (permission_id, access_level_id)
);