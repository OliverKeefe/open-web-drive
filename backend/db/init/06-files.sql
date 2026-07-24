-- FILE (stable identity)
CREATE TABLE files (
    id UUID PRIMARY KEY NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
COMMENT ON TABLE files IS 'Stable file identity. Immutable across renames, moves, and versions.';

-- FILE METADATA (versioned file attributes)
CREATE TABLE file_metadata (
    id UUID NOT NULL PRIMARY KEY,
    file_id UUID NOT NULL REFERENCES files(id) ON DELETE CASCADE,
    version INTEGER NOT NULL DEFAULT 1,
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    file_name VARCHAR(255) NOT NULL,
    path VARCHAR(255) NOT NULL,
    relative_path VARCHAR(255) NOT NULL,
    size BIGINT NOT NULL,
    file_type VARCHAR(12) NOT NULL,
    hash CHAR(64) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    modified_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    uploaded_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT uq_file_metadata_version UNIQUE (file_id, version)
);
COMMENT ON TABLE file_metadata IS 'Versioned file attributes. Each upload creates a new row with version + 1.';

-- INDEX FOR FILE METADATA PAGINATION
CREATE INDEX idx_file_owner_modified
    ON file_metadata (owner_id, modified_at DESC, id DESC);
