\connect appdb;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- USERS
CREATE TABLE users (
    id UUID PRIMARY KEY NOT NULL,
    name VARCHAR(64)
);

-- TENANTS
CREATE TABLE tenant (
    id UUID PRIMARY KEY NOT NULL,
    name VARCHAR(64) NOT NULL
);

-- GROUPS
CREATE TABLE groups (
    id UUID PRIMARY KEY NOT NULL,
    name VARCHAR(64) NOT NULL
);

-- BUCKETS
CREATE TABLE bucket (
    id UUID PRIMARY KEY NOT NULL,
    tenant_id UUID NOT NULL REFERENCES tenant(id) ON DELETE CASCADE,
    name VARCHAR(64) NOT NULL,
    size BIGINT NOT NULL,
    creator UUID NOT NULL REFERENCES users(id),
    owner UUID NOT NULL REFERENCES users(id)
);

-- BUCKET GROUP MEMBERSHIP
CREATE TABLE bucket_groups (
    bucket_id UUID NOT NULL REFERENCES bucket(id) ON DELETE CASCADE,
    group_id UUID NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    PRIMARY KEY (bucket_id, group_id)
);

-- ACCESS TO FILE
CREATE TABLE access_to(
    access_group_id UUID REFERENCES groups(id) ON DELETE CASCADE,
    access_user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (access_group_id, access_user_id)
);

-- FILE
CREATE TABLE files(
    id UUID PRIMARY KEY NOT NULL,
    owner_id UUID REFERENCES users(id) NOT NULL,
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

-- FILE METADATA
CREATE TABLE file_metadata(
    id UUID PRIMARY KEY NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    path VARCHAR(255) NOT NULL,
    size BIGINT NOT NULL,
    file_type VARCHAR(12),
    modified_at TIMESTAMP NOT NULL,
    uploaded_at TIMESTAMP NOT NULL,
    version UUID NOT NULL,
    hash BYTEA,
);

-- FILE VERSIONS
CREATE TABLE file_versions(
    version_id UUID REFERENCES file_metadata(id) NOT NULL,
    file_id UUID REFERENCES files(id) NOT NULL,
    CONSTRAINT pk_file_version PRIMARY KEY (file_id, version_id)
);

-- INDEX FOR FILE METADATA PAGINATION
CREATE INDEX CONCURRENTLY idx_file_owner_modified_desc
ON file_metadata (owner_id, modified_at DESC, id DESC);

-- DIRECTORY METADATA
CREATE TABLE dir_metadata(
    id UUID PRIMARY KEY NOT NULL,
    dir_name VARCHAR(255) NOT NULL,
    path VARCHAR(255) NOT NULL,
    modified_at TIMESTAMP NOT NULL,
    uploaded_at TIMESTAMP NOT NULL,
    owner UUID NOT NULL REFERENCES users(id)
);

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