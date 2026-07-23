package files

import (
	"context"
	"errors"
	"fmt"
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

// PersistFile inserts a new row into the files table (stable file identity).
func (db *FileRepository) PersistFile(ctx context.Context, fileID uuid.UUID) error {
	const q = `INSERT INTO files (id, created_at) VALUES ($1, now());`
	_, err := db.Pool.Exec(ctx, q, fileID)
	return err
}

// GetNextVersion returns the next version number for a given file_id.
func (db *FileRepository) GetNextVersion(ctx context.Context, fileID uuid.UUID) (int, error) {
	const q = `SELECT COALESCE(MAX(version), 0) + 1 FROM file_metadata WHERE file_id = $1;`
	var next int
	err := db.Pool.QueryRow(ctx, q, fileID).Scan(&next)
	return next, err
}

// UpsertUser inserts a user if they don't already exist.
// TODO: extract preferred_username from JWT claims when auth layer supports it.
func (db *FileRepository) UpsertUser(ctx context.Context, userID uuid.UUID, name string) error {
	const q = `INSERT INTO users (id, name) VALUES ($1, $2)
		ON CONFLICT (id) DO NOTHING;`
	_, err := db.Pool.Exec(ctx, q, userID, name)
	return err
}

// FileExists checks if a files row exists for the given file ID.
func (db *FileRepository) FileExists(ctx context.Context, fileID uuid.UUID) (bool, error) {
	const q = `SELECT id FROM files WHERE id = $1;`
	var id uuid.UUID
	err := db.Pool.QueryRow(ctx, q, fileID).Scan(&id)
	if err == pgx.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// CheckExists checks if a metadata row already exists for the given (file_id, version).
func (db *FileRepository) CheckExists(ctx context.Context, fileID uuid.UUID, version int) (bool, error) {
	const q = `SELECT id FROM file_metadata WHERE file_id = $1 AND version = $2;`
	var id uuid.UUID
	err := db.Pool.QueryRow(ctx, q, fileID, version).Scan(&id)
	if err == pgx.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// PersistMetadata inserts a new file_metadata row with all versioned attributes.
func (db *FileRepository) PersistMetadata(ctx context.Context, metadata FileMetadata) error {
	const q = `INSERT INTO file_metadata
		(id, file_id, version, owner_id, file_name, path, relative_path,
		 size, file_type, hash, modified_at, uploaded_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13);`

	_, err := db.Pool.Exec(
		ctx,
		q,
		metadata.ID,
		metadata.FileID,
		metadata.Version,
		metadata.OwnerID,
		metadata.FileName,
		metadata.Path,
		metadata.RelativePath,
		metadata.Size,
		metadata.FileType,
		metadata.Hash,
		metadata.ModifiedAt,
		metadata.UploadedAt,
		metadata.CreatedAt,
	)
	return err
}

// FindAllMetadata returns all file_metadata rows for a user, with cursor-based pagination.
func (db *FileRepository) FindAllMetadata(ctx context.Context, req GetAllMetadataRequest) ([]FileMetadata, error) {
	var (
		rows pgx.Rows
		err  error
	)

	if req.Cursor == nil || req.Cursor.ID == uuid.Nil || req.Cursor.ModifiedAt.IsZero() {
		rows, err = db.Pool.Query(ctx, `
			SELECT id, file_id, file_name, path, relative_path, size, file_type,
				owner_id, version, hash, created_at, modified_at, uploaded_at
			FROM file_metadata
			WHERE owner_id = $1
			ORDER BY modified_at DESC, id DESC
			LIMIT $2;
		`, req.UserID, req.Limit)
	} else {
		rows, err = db.Pool.Query(ctx, `
			SELECT id, file_id, file_name, path, relative_path, size, file_type,
				owner_id, version, hash, created_at, modified_at, uploaded_at
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

	result := make([]FileMetadata, 0, req.Limit)

	for rows.Next() {
		var model FileMetadata
		if err := rows.Scan(
			&model.ID,
			&model.FileID,
			&model.FileName,
			&model.Path,
			&model.RelativePath,
			&model.Size,
			&model.FileType,
			&model.OwnerID,
			&model.Version,
			&model.Hash,
			&model.CreatedAt,
			&model.ModifiedAt,
			&model.UploadedAt,
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

// FindMetadata builds a dynamic query for filtering file_metadata rows.
// Only non-zero fields in m are used as WHERE conditions.
func (db *FileRepository) FindMetadata(ctx context.Context, m FileMetadata) ([]FileMetadata, error) {
	const baseQuery = `SELECT id, file_id, file_name, path, relative_path, size, file_type,
		owner_id, version, hash, created_at, modified_at, uploaded_at
	FROM file_metadata`

	var (
		clauses []string
		args    []any
		argN    int
	)

	addClause := func(col string, val any) {
		argN++
		clauses = append(clauses, fmt.Sprintf("%s = $%d", col, argN))
		args = append(args, val)
	}

	if m.ID != uuid.Nil {
		addClause("id", m.ID)
	}
	if m.FileID != uuid.Nil {
		addClause("file_id", m.FileID)
	}
	if m.FileName != "" {
		addClause("file_name", m.FileName)
	}
	if m.Path != "" {
		addClause("path", m.Path)
	}
	if m.Size != 0 {
		addClause("size", m.Size)
	}
	if m.FileType != "" {
		addClause("file_type", m.FileType)
	}
	if !m.ModifiedAt.IsZero() {
		addClause("modified_at", m.ModifiedAt)
	}
	if !m.UploadedAt.IsZero() {
		addClause("uploaded_at", m.UploadedAt)
	}
	if m.OwnerID != uuid.Nil {
		addClause("owner_id", m.OwnerID)
	}
	if m.Version != 0 {
		addClause("version", m.Version)
	}
	if m.Hash != "" {
		addClause("hash", m.Hash)
	}

	query := baseQuery
	if len(clauses) > 0 {
		query = fmt.Sprintf("%s WHERE %s", baseQuery, strings.Join(clauses, " AND "))
	}

	rows, err := db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query file_metadata: %w", err)
	}
	defer rows.Close()

	var result []FileMetadata
	for rows.Next() {
		var model FileMetadata
		if err := rows.Scan(
			&model.ID,
			&model.FileID,
			&model.FileName,
			&model.Path,
			&model.RelativePath,
			&model.Size,
			&model.FileType,
			&model.OwnerID,
			&model.Version,
			&model.Hash,
			&model.CreatedAt,
			&model.ModifiedAt,
			&model.UploadedAt,
		); err != nil {
			return nil, fmt.Errorf("scan file_metadata: %w", err)
		}
		result = append(result, model)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

// FindMetadataByID returns a single file_metadata row by its primary key.
func (db *FileRepository) FindMetadataByID(ctx context.Context, ID uuid.UUID) (FileMetadata, error) {
	const q = `SELECT id, file_id, file_name, path, relative_path, size, file_type,
		owner_id, version, hash, created_at, modified_at, uploaded_at
	FROM file_metadata
	WHERE id = $1;`

	var model FileMetadata
	err := db.Pool.QueryRow(ctx, q, ID).Scan(
		&model.ID,
		&model.FileID,
		&model.FileName,
		&model.Path,
		&model.RelativePath,
		&model.Size,
		&model.FileType,
		&model.OwnerID,
		&model.Version,
		&model.Hash,
		&model.CreatedAt,
		&model.ModifiedAt,
		&model.UploadedAt,
	)
	if err != nil {
		return FileMetadata{}, err
	}
	return model, nil
}

// DeleteMetadata deletes a file_metadata row by id and owner_id.
func (db *FileRepository) DeleteMetadata(ctx context.Context, id uuid.UUID, ownerID uuid.UUID) error {
	const q = `DELETE FROM file_metadata WHERE id = $1 AND owner_id = $2;`

	status, err := db.Pool.Exec(ctx, q, id, ownerID)
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

// DeleteFileData deletes the blob from storage and removes the files row if no
// metadata versions remain.
func (db *FileRepository) DeleteFileData(ctx context.Context, id uuid.UUID, ownerID uuid.UUID) error {
	// Query the metadata row to get the blob key components.
	var (
		fileID   uuid.UUID
		fileName string
	)
	const lookup = `SELECT file_id, file_name FROM file_metadata WHERE id = $1 AND owner_id = $2;`
	err := db.Pool.QueryRow(ctx, lookup, id, ownerID).Scan(&fileID, &fileName)
	if err == pgx.ErrNoRows {
		return errors.New("no record found")
	}
	if err != nil {
		return fmt.Errorf("could not look up file metadata: %w", err)
	}

	// Delete all metadata versions for this file.
	const delMeta = `DELETE FROM file_metadata WHERE file_id = $1;`
	if _, err := db.Pool.Exec(ctx, delMeta, fileID); err != nil {
		return fmt.Errorf("could not delete file metadata: %w", err)
	}

	// Delete the file identity row.
	const delFile = `DELETE FROM files WHERE id = $1;`
	status, err := db.Pool.Exec(ctx, delFile, fileID)
	if err != nil {
		return fmt.Errorf("could not delete file: %w", err)
	}
	if status.RowsAffected() == 0 {
		return errors.New("no file record found")
	}

	return nil
}

func (db *FileRepository) MarkForDeletion(ctx context.Context, id uuid.UUID, ownerID uuid.UUID) error {
	panic("not implemented")
}
