CREATE TABLE access_levels (
    id UUID DEFAULT uuidv7() PRIMARY KEY,
    name VARCHAR(16) NOT NULL
);

CREATE TABLE permissions (
    id UUID PRIMARY KEY NOT NULL,
    grantee_type VARCHAR(10) NOT NULL,
    grantee_id UUID NOT NULL
);

CREATE TABLE permission_access_levels (
    permission_id UUID REFERENCES permissions(id),
    access_level_id UUID REFERENCES access_levels(id),
    CONSTRAINT pk_permission_access_levels PRIMARY KEY (permission_id, access_level_id)
);