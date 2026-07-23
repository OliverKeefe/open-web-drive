-- FILE <-> USER ACCESS (MtM)
CREATE TABLE file_user_access (
    file_metadata_id UUID NOT NULL REFERENCES file_metadata(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT pk_file_user_access PRIMARY KEY (file_metadata_id, user_id)
);

-- FILE <-> GROUP ACCESS (MtM)
CREATE TABLE file_group_access (
    file_metadata_id UUID NOT NULL REFERENCES file_metadata(id) ON DELETE CASCADE,
    group_id UUID NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    CONSTRAINT pk_file_group_access PRIMARY KEY (file_metadata_id, group_id)
);
