\connect appdb;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- USERS
CREATE TABLE users (
    id UUID PRIMARY KEY NOT NULL,
    name VARCHAR(64)
);

-- TENANTS
CREATE TABLE orgs (
    id UUID PRIMARY KEY NOT NULL,
    name VARCHAR(64) NOT NULL
);

-- GROUPS IDENTITY COLLECTION & MANAGEMENT
CREATE TABLE groups (
    id UUID PRIMARY KEY NOT NULL,
    name VARCHAR(64) NOT NULL
);

CREATE TABLE group_permissions (
    group_id REFERENCES groups(id) ON DELETE CASCADE,
    user_id REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (group_id, user_id)
);

CREATE TABLE roles (
    id UUID PRIMARY KEY NOT NULL,
    name VARCHAR(64) NOT NULL
);

-- BUCKETS
CREATE TABLE storage_bucket (
    id UUID PRIMARY KEY NOT NULL,
    org_id UUID NOT NULL REFERENCES orgs(id) ON DELETE CASCADE,
    name VARCHAR(64) NOT NULL,
    groups NOT NULL REFERENCES storage_bucket_groups
    owner UUID NOT NULL REFERENCES users(id)
);

-- BUCKET GROUP MEMBERSHIP
CREATE TABLE storage_bucket_groups (
    bucket_id UUID NOT NULL REFERENCES bucket(id) ON DELETE CASCADE,
    group_id UUID NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    PRIMARY KEY (bucket_id, group_id)
);

-- FILE
CREATE TABLE files(
    id UUID PRIMARY KEY NOT NULL,
    permissions_id REFERENCES permissions(id)
);

CREATE TABLE permissions(
    id UUID PRIMARY KEY NOT NULL,
    file_id REFERENCES files(id) NOT NULL,
    grantee_type VARCHAR(10) NOT NULL,
    grantee_id UUID NOT NULL,
    access_level CHAR -- May need to be altered again.
);

CREATE TABLE file_permissions(
    file_id REFERENCES files(id) NOT NULL,
    permissions_id REFERENCES permissions(id) NOT NULL,
    PRIMARY KEY (file_id, permissions_id)
);


-- METADATA
CREATE TABLE metadata(
    id UUID PRIMARY KEY NOT NULL DELETE CASCADE,
    file_name VARCHAR(255) NOT NULL DELETE CASCADE,
    path VARCHAR(255) NOT NULL DELETE CASCADE,
    relative_path VARCHAR(255) NOT NULL DELETE CASCADE,
    size BIGINT NOT NULL DELETE CASCADE,
    file_type VARCHAR(12) NOT NULL DELETE CASCADE,
    modified_at TIMESTAMP NOT NULL DELETE CASCADE,
    uploaded_at TIMESTAMP NOT NULL DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DELETE CASCADE,
    version UUID NOT NULL DELETE CASCADE,
    hash BYTEA NOT NULL DELETE CASCADE,
);

-- FILE <-> METADATA
CREATE TABLE file_metadata(
    file_id UUID NOT NULL REFERENCES files(id) ON DELETE CASCADE,
    metadata_id UUID NOT NULL REFERENCES metadata(id) ON DELETE CASCADE,
    CONSTRAINT pk_file_metadata PRIMARY KEY (file_id, metadata_id)
);

-- FILE VERSIONS
CREATE TABLE file_versions(
    version_id UUID REFERENCES file_metadata(id) NOT NULL,
    file_id UUID REFERENCES files(id) NOT NULL,
    CONSTRAINT pk_file_version PRIMARY KEY (file_id, version_id)
);

-- INDEX FOR FILE METADATA PAGINATION
CREATE INDEX CONCURRENTLY idx_file_user_modified_desc
ON file_metadata (user_id, modified_at DESC, id DESC);

-- FILE <-> USER ACCESS (MtM)
CREATE TABLE file_metadata_access(
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    file_id UUID NOT NULL REFERENCES file_metadata(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, file_id)
);

-- FILE <-> GROUP ACCESS (MtM)
CREATE TABLE file_metadata_group_access(
    file_id UUID REFERENCES file_metadata(id) ON DELETE CASCADE,
    group_id UUID REFERENCES groups(id) ON DELETE CASCADE
);

CREATE TABLE file_user_access(
    file_id UUID REFERENCES file_metadata(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE
);