package files

import (
	query "backend/src/internal/api/query"
	"backend/src/internal/database"
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/google/uuid"
)

type FileRepository struct {
	pool database.Pool
}

func FindMetadataQuery(m MetaData) (string, []any) {
	const baseQuery = `SELECT id, file_name, path, size, file_type, modified_at, 
       	uploaded_at, version, checksum, owner_id 
		FROM file_metadata
	`

	var b query.Builder

	b.Equal("id", m.ID)
	b.Equal("file_name", m.FileName)
	b.Equal("path", m.Path)
	b.Equal("size", m.Size)
	b.Equal("file_type", m.FileType)
	b.Equal("modified_at", m.ModifiedAt)
	b.Equal("uploaded_at", m.UploadedAt)
	b.Equal("owner_id", m.Owner)
	b.Equal("version", m.Version)
	if len(m.Group) > 0 {
		const selectGroups = `
			EXISTS (
				SELECT * FROM file_metadata_group_access fga 
				WHERE fga.file_id = file_metadata.id
					AND fga.user_id = ANY (?) 
			)`
		b.Raw(selectGroups, m.Group)
	}
	if len(m.AccessTo) > 0 {
		const accessTo = `
			EXISTS (
				SELECT * FROM file_user_access fua
				WHERE fua.file_id = file_medata.id
					AND fua.user_id = ANY (?)
			)`
		b.Raw(accessTo, m.AccessTo)
	}

	sql := fmt.Sprintf("%s WHERE %s", baseQuery, strings.Join(b.Clauses, " AND "))
	args := b.Args

	return sql, args
}

func (db *FileRepository) Persist() {

}

func (db *FileRepository) FindMetadataByID(ctx context.Context, ID uuid.UUID) (Metadata, error) {
	panic("not implemented")
}

func (db *FileRepository) CheckExists(ctx context.Context, ID uuid.UUID) (bool, error) {
	const sql = `SELECT id FROM metadata WHERE id = ($1)`
	status, err := db.pool.Exec(ctx, sql, ID)
	if err != nil {
		slog.Error("failed existence check", "error", "status:", status, err)
		return false, err
	}
	return true, nil
}

func (db *FileRepository) PersistMetadata(ctx context.Context, metadata FileMetadata) error {
	const query = `INSERT INTO file_metadata (id, file_name, path, size, file_type, 
                           modified_at, uploaded_at, version, checksum, owner_id) 
				   VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);`

	_, err := db.pool.Exec(
		ctx,
		query,
		metadata.ID,
		metadata.OwnerID,
		metadata.FileName,
		metadata.Path,
		metadata.RelativePath,
		metadata.Size,
		metadata.FileType,
		metadata.ModifiedAt,
		metadata.UploadedAt,
		metadata.CreatedAt,
		metadata.Version,
		metadata.Hash,
		metadata.Permissions,
	)
	if err != nil {
		return err
	}

	return nil
}
