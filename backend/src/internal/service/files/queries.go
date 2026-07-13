package files

import (
	query "backend/src/internal/api/query"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Pool interface {
	Begin(context.Context) (pgx.Tx, error)
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
}

type FileRepository struct {
	Pool Pool
}

func NewFileRepository(pool Pool) *FileRepository {
	return &FileRepository{Pool: pool}
}

func FindMetadata(m MetaData) (string, []any) {
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
	panic("not implemented")
}

func (db *FileRepository) FindMetadataByID(ctx context.Context, ID uuid.UUID) (Metadata, error) {
	panic("not implemented")
}

func (db *FileRepository) CheckExists(ctx context.Context, ID uuid.UUID) (bool, error) {
	const sql = `SELECT id FROM metadata WHERE id = ($1)`
	status, err := db.Pool.Exec(ctx, sql, ID)
	if err != nil {
		slog.Error("failed existence check", "error", err, "status", status)
		return false, err
	}
	return true, nil
}

func (db *FileRepository) PersistMetadata(ctx context.Context, metadata FileMetadata) error {
	const q = `INSERT INTO file_metadata (id, file_name, path, size, file_type, 
                           modified_at, uploaded_at, version, checksum, owner_id) 
				   VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);`

	_, err := db.Pool.Exec(
		ctx,
		q,
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

func (db *FileRepository) DeleteFileData(ctx context.Context, id uuid.UUID, ownerId uuid.UUID) error {
	panic("not implemented")
}

func (db *FileRepository) MarkForDeletion(ctx context.Context, id uuid.UUID, ownerId uuid.UUID) error {
	panic("not implemented")
	// Should add flag to relation in metadata to delete x days (depending on user policy
	// default = 30d), then need cleaner func to bulk cleanse metadata db and before this,
	// remove file data from both virtual disc and IPFS node.
}

func (db *FileRepository) FindAllMetadata(ctx context.Context, req GetAllMetadataRequest) ([]MetaData, error) {

	var (
		rows pgx.Rows
		err  error
	)

	if req.Cursor == nil || req.Cursor.ID == uuid.Nil || req.Cursor.ModifiedAt.IsZero() {
		rows, err = db.Pool.Query(ctx, `
			SELECT id, file_name, path, size, file_type, modified_at,
				uploaded_at, owner_id, checksum, version 
			FROM file_metadata
			WHERE owner_id = $1
			ORDER BY modified_at DESC, id DESC
			LIMIT $2;
		`, req.UserID, req.Limit)

	} else {
		rows, err = db.Pool.Query(ctx, `
			SELECT id, file_name, path, size, file_type, modified_at,
				uploaded_at, owner_id, checksum, version
			FROM file_metadata
			WHERE owner_id = $1
				AND (modified_at, id) < ($2, $3)
			ORDER BY modified_at DESC, id DESC
			LIMIT $4;
		`, req.UserID, req.Cursor.ModifiedAt, req.Cursor.ID, req.Limit)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]MetaData, 0, req.Limit)

	for rows.Next() {
		var model MetaData
		if err := rows.Scan(
			&model.ID,
			&model.FileName,
			&model.Path,
			&model.Size,
			&model.FileType,
			&model.ModifiedAt,
			&model.UploadedAt,
			&model.Owner,
			//&model.AccessTo,
			//&model.Group,
			&model.CheckSum,
			&model.Version,
		); err != nil {
			return nil, err
		}

		result = append(result, model)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (db *FileRepository) DeleteMetadata(ctx context.Context, id uuid.UUID, ownerId uuid.UUID) error {
	const q = `DELETE FROM file_metadata WHERE id = $1 AND owner_id = $2;`

	status, err := db.Pool.Exec(ctx, q, id, ownerId)
	if err != nil {
		return fmt.Errorf(
			"status: %s, could not delete file metadata, %w",
			status,
			err,
		)
	}

	rows := status.RowsAffected()

	if rows == 0 {
		return errors.New("no record found")
	}

	return nil
}
