-- BUCKETS
CREATE TABLE storage_bucket (
    id UUID PRIMARY KEY NOT NULL,
    org_id UUID NOT NULL REFERENCES orgs(id) ON DELETE CASCADE,
    name VARCHAR(64) NOT NULL
    --groups NOT NULL REFERENCES storage_bucket_groups,
    --owner UUID NOT NULL REFERENCES users(id)
);

-- BUCKET GROUP MEMBERSHIP
--CREATE TABLE storage_bucket_groups (
--    bucket_id UUID NOT NULL REFERENCES bucket(id) ON DELETE CASCADE,
--    group_id UUID NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
--    PRIMARY KEY (bucket_id, group_id)
--);
