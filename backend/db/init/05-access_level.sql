-- ACCESS LEVELS (lookup / enum)
CREATE TABLE access_level (
    id UUID NOT NULL PRIMARY KEY,
    name VARCHAR(16) NOT NULL,
    CONSTRAINT uq_access_level_name UNIQUE (name)
);
