-- FILE
CREATE TABLE files (
    id UUID PRIMARY KEY NOT NULL,
    permissions_id UUID REFERENCES permissions(id)
);

CREATE TABLE file_permissions (
    file_id UUID REFERENCES files(id) NOT NULL,
    permissions_id UUID REFERENCES permissions(id) NOT NULL,
    PRIMARY KEY (file_id, permissions_id)
);

CREATE TABLE file_groups (
    file_id UUID REFERENCES files(id) ON DELETE CASCADE,
    group_id UUID REFERENCES groups(id) ON DELETE CASCADE,
    PRIMARY KEY (file_id, group_id)
);

-- METADATA
CREATE TABLE metadata (
    id UUID PRIMARY KEY NOT NULL,
    owner UUID REFERENCES users(id) ON DELETE CASCADE,
    file_name VARCHAR(255) NOT NULL,
    path VARCHAR(255) NOT NULL,
    relative_path VARCHAR(255) NOT NULL,
    size BIGINT NOT NULL,
    file_type VARCHAR(12) NOT NULL,
    modified_at TIMESTAMP NOT NULL,
    uploaded_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    hash BYTEA NOT NULL
);

-- FILE <-> METADATA
CREATE TABLE file_metadata (
    id UUID DEFAULT uuidv7() PRIMARY KEY,
    file_id UUID NOT NULL REFERENCES files(id) ON DELETE CASCADE,
    metadata_id UUID NOT NULL REFERENCES metadata(id) ON DELETE CASCADE,
    CONSTRAINT uq_file_metadata UNIQUE (file_id, metadata_id)
);

-- FILE VERSIONS
CREATE TABLE file_versions (
    version_id UUID REFERENCES file_metadata(id) NOT NULL,
    file_id UUID REFERENCES files(id) NOT NULL,
    CONSTRAINT pk_file_version PRIMARY KEY (file_id, version_id)
);

-- FILE <-> USER ACCESS (MtM)
CREATE TABLE file_metadata_access (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    file_metadata_id UUID NOT NULL REFERENCES file_metadata(id) ON DELETE CASCADE,
    CONSTRAINT pk_file_metadata_access PRIMARY KEY (user_id, file_metadata_id)
);
-- FILE <->  GROUP ACCESS (MtM)
CREATE TABLE file_metadata_group_access (
    file_id UUID REFERENCES file_metadata(id) ON DELETE CASCADE,
    group_id UUID REFERENCES groups(id) ON DELETE CASCADE
);

CREATE TABLE file_user_access (
    file_id UUID REFERENCES file_metadata(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE
);

-- INDEX FOR FILE METADATA PAGINATION
CREATE INDEX CONCURRENTLY idx_file_user_modified_desc
    ON file_metadata (user_id, modified_at DESC, id DESC);