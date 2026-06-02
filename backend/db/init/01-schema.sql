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