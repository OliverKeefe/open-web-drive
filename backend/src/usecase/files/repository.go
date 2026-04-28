package files

import (
	"backend/src/internal/auth"
	"backend/src/internal/db/metadb"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db       metadb.Pool
	s3Client *s3.Client
	bucket   string
}

func NewRepository(db metadb.Pool, s3Client *s3.Client, bucket string) *Repository {
	return &Repository{
		db:       db,
		s3Client: s3Client,
		bucket:   bucket,
	}
}

func (repo *Repository) SaveMetaData(ctx context.Context, meta MetaData) (MetaData, error) {
	const query = `INSERT INTO file_metadata (
                           id,
                           file_name,
                           path, size, 
                           file_type,
                           modified_at,
                           uploaded_at,
                           version,
                           checksum,
                           owner_id) 
					VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);`

	_, err := repo.db.Exec(
		ctx,
		query,
		meta.ID,
		meta.FileName,
		meta.Path,
		meta.Size,
		meta.FileType,
		meta.ModifiedAt,
		meta.UploadedAt,
		meta.Version,
		meta.CheckSum,
		meta.Owner,
	)
	if err != nil {
		return meta, err
	}
	return meta, nil
}

func (repo *Repository) SaveToS3(ctx context.Context, basePath string, rdr io.Reader, filename string) error {
	ownerID, ok := auth.UserIDFromCtx(ctx)
	if !ok {
		return errors.New("unable to get ownerID from context")
	}

	key := fmt.Sprintf("%s/%s", ownerID, filename)

	var partMiBs int64 = 10

	uploader := transfermanager.New(repo.s3Client, func(o *transfermanager.Options) {
		o.PartSizeBytes = partMiBs * 1024 * 1024
		o.Concurrency = 3
	})

	_, err := uploader.UploadObject(ctx, &transfermanager.UploadObjectInput{
		Bucket: aws.String(repo.bucket),
		Key:    aws.String(key),
		Body:   rdr,
	})
	if err != nil {
		log.Printf("S3 transfer manager put error: %v", err)
		return err
	}

	return nil
}

func (repo *Repository) SaveFileData(basePath string, rdr io.Reader, filename string) error {
	log.Printf("SaveFileData called")
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return err
	}

	fileExtension := filepath.Ext(filename)
	base := strings.TrimSuffix(filename, fileExtension)

	tmp, err := os.CreateTemp(basePath, base+"-*"+fileExtension)
	if err != nil {
		return err
	}
	defer tmp.Close()

	if _, err := io.Copy(tmp, rdr); err != nil {
		return err
	}

	return nil
}

func (repo *Repository) FindAllMetadata(
	ctx context.Context,
	req GetAllMetadataRequest,
) ([]MetaData, error) {

	var (
		rows pgx.Rows
		err  error
	)

	if req.Cursor == nil || req.Cursor.ID == uuid.Nil || req.Cursor.ModifiedAt.IsZero() {
		rows, err = repo.db.Query(ctx, `
			SELECT id, file_name, path, size, file_type, modified_at,
				uploaded_at, owner_id, checksum, version 
			FROM file_metadata
			WHERE owner_id = $1
			ORDER BY modified_at DESC, id DESC
			LIMIT $2;
		`, req.UserID, req.Limit)

	} else {
		rows, err = repo.db.Query(ctx, `
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

func (repo *Repository) FindMetadata(ctx context.Context, model MetaData) ([]MetaData, error) {
	var result []MetaData

	query, args := FindMetadataQuery(model)

	rows, err := repo.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
			&model.AccessTo,
			&model.Group,
			&model.Version,
		); err != nil {
			return nil, err
		}

		result = append(result, model)
	}

	return result, nil
}

func (repo *Repository) DeleteMetadata(ctx context.Context, id uuid.UUID, ownerId uuid.UUID) error {
	const query = `DELETE FROM file_metadata WHERE id = $1 AND owner_id = $2;`

	status, err := repo.db.Exec(ctx, query, id, ownerId)
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

func (repo *Repository) DeleteFileData(ctx context.Context, id uuid.UUID, ownerId uuid.UUID) error {
	panic("not implemented")
}

func (repo *Repository) Modify(ctx context.Context, model MetaData) (error, MetaData) {
	panic("not implemented")
}

func (repo *Repository) MarkForDeletion(ctx context.Context, id uuid.UUID, id2 uuid.UUID) error {
	panic("not implemented")
	// Should add flag to relation in metadata to delete x days (depending on user policy
	// default = 30d), then need cleaner func to bulk cleanse metadata db and before this,
	// remove file data from both virtual disc and IPFS node.
}
