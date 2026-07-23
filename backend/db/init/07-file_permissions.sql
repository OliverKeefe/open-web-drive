-- FILE PERMISSIONS (per-grantee access grants)
CREATE TABLE file_permissions (
    id UUID NOT NULL PRIMARY KEY,
    file_id UUID NOT NULL REFERENCES files(id) ON DELETE CASCADE,
    grantee_type VARCHAR(10) NOT NULL,
    grantee_id UUID NOT NULL,
    access_level_id UUID NOT NULL REFERENCES access_level(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT chk_grantee_type CHECK (grantee_type IN ('user', 'group')),
    CONSTRAINT uq_file_permissions UNIQUE (file_id, grantee_type, grantee_id, access_level_id)
);
