package files

import (
	"backend/internal/api/message"
	"backend/internal/auth"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/google/uuid"
	"gocloud.dev/blob"
)

type uploadRepository interface {
	PersistFile(ctx context.Context, fileID uuid.UUID) error
	PersistMetadata(ctx context.Context, metadata FileMetadata) error
	GetNextVersion(ctx context.Context, fileID uuid.UUID) (int, error)
	FileExists(ctx context.Context, fileID uuid.UUID) (bool, error)
	CheckExists(ctx context.Context, fileID uuid.UUID, version int) (bool, error)
	UpsertUser(ctx context.Context, userID uuid.UUID, name string) error
}

type blobStorage interface {
	MultipartUpload(ctx context.Context, key string, dataStream io.Reader, opts *blob.WriterOptions) error
}

type pendingFile struct {
	metadata   FileMetadata
	fileExists bool
}

type UploadService struct {
	Db                uploadRepository
	BlobStorageClient blobStorage
	BucketURL         string
}

func NewUploadService(db uploadRepository, client blobStorage, bucketUrl string) *UploadService {
	return &UploadService{
		Db:                db,
		BlobStorageClient: client,
		BucketURL:         bucketUrl,
	}
}

func (svc *UploadService) Handle(w http.ResponseWriter, r *http.Request) {
	//userID, ok := auth.UserIDFromCtx(r.Context())
	//if !ok {
	//	http.Error(w, "upload failed", http.StatusUnauthorized)
	//	return
	//}

	//hasClaim, err := auth.HasClaim(r.Context(), uuid.MustParse(userID), "")
	//if err != nil {
	//	slog.Error("upload failed", "error", err)
	//}
	//
	//if !hasClaim {
	//	http.Error(w, "upload failed", http.StatusUnauthorized)
	//	return
	//}

	// FYI - 5<<20 is a bitshift operation (5*2^20 = 5,242,880 Bytes or 5 MegaBytes)
	r.Body = http.MaxBytesReader(w, r.Body, 5<<20)
	if err := svc.execute(r); err != nil {
		slog.Error("upload failed", "error", err)
		http.Error(w, "upload failed", http.StatusBadRequest)
		return
	}

	if err := message.Response(w, "upload succeeded", nil); err != nil {
		slog.Error("upload failed", "error", err)
		return
	}
}

// multipartMetadata is the raw metadata part of a multipart upload.
// TODO: reformat json fields, currently match frontend but not consistent with Go idioms.
type multipartMetadata struct {
	Path             string `json:"path"`
	RelativePath     string `json:"relativePath"`
	LastModified     int64  `json:"lastModified"`
	LastModifiedDate string `json:"lastModifiedDate"`
	UploadedAt       int64  `json:"uploadedAt"`
	Size             uint64 `json:"size"`
	FileType         string `json:"fileType"`
	ID               string `json:"id"`
	FileName         string `json:"file_name"`
	CreatedAt        int64  `json:"created_at"`
}

func (svc *UploadService) execute(r *http.Request) error {
	mr, err := r.MultipartReader()
	if err != nil {
		return err
	}

	ctx := r.Context()

	pendingByID := make(map[string]*pendingFile)

	userID, ok := auth.UserIDFromCtx(ctx)
	if !ok {
		return errors.New("unable to get UserID from context")
	}

	parsedUserID := uuid.MustParse(userID)
	if err := svc.Db.UpsertUser(ctx, parsedUserID, userID); err != nil {
		return fmt.Errorf("failed to upsert user: %w", err)
	}

	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		name := part.FormName()

		switch {
		case strings.HasPrefix(name, "metadata-"):
			var decodedRequest multipartMetadata
			if err := json.NewDecoder(part).Decode(&decodedRequest); err != nil {
				return err
			}

			idAsStr := strings.TrimPrefix(name, "metadata-")
			fileID := uuid.MustParse(idAsStr)

			// Check if the file identity row already exists.
			fileExists, err := svc.Db.FileExists(ctx, fileID)
			if err != nil {
				return fmt.Errorf("failed to check file existence: %w", err)
			}

			// Get the next version number.
			nextVersion, err := svc.Db.GetNextVersion(ctx, fileID)
			if err != nil {
				return fmt.Errorf("failed to get next version: %w", err)
			}

			metaID, err := uuid.NewV7()
			if err != nil {
				return fmt.Errorf("failed to generate metadata id: %w", err)
			}

			pendingByID[idAsStr] = &pendingFile{
				metadata: FileMetadata{
					ID:           metaID,
					FileID:       fileID,
					Version:      nextVersion,
					OwnerID:      uuid.MustParse(userID),
					FileName:     decodedRequest.FileName,
					Path:         decodedRequest.Path,
					RelativePath: decodedRequest.RelativePath,
					Size:         int64(decodedRequest.Size),
					FileType:     path.Ext(decodedRequest.FileName),
					ModifiedAt:   time.UnixMilli(decodedRequest.LastModified),
					UploadedAt:   time.UnixMilli(decodedRequest.UploadedAt),
					CreatedAt:    time.UnixMilli(decodedRequest.CreatedAt),
				},
				fileExists: fileExists,
			}

		case strings.HasPrefix(name, "filedata-"):
			idToStr := strings.TrimPrefix(name, "filedata-")

			data, err := io.ReadAll(part)
			if err != nil {
				return err
			}

			pf, ok := pendingByID[idToStr]
			if !ok {
				return fmt.Errorf("filedata part received before metadata part for %s", idToStr)
			}

			ownerID := uuid.MustParse(userID)
			if err := svc.saveFileData(ctx, ownerID, bytes.NewReader(data), part.FileName()); err != nil {
				return err
			}

			hash := sha256.Sum256(data)
			pf.metadata.Hash = hex.EncodeToString(hash[:])
		}
	}

	for _, pf := range pendingByID {
		if err := svc.saveFile(ctx, pf); err != nil {
			return err
		}
		if err := svc.saveMetadata(ctx, pf.metadata); err != nil {
			return err
		}
	}

	return nil
}

func (svc *UploadService) saveFile(ctx context.Context, pf *pendingFile) error {
	if pf.fileExists {
		return nil
	}
	return svc.Db.PersistFile(ctx, pf.metadata.FileID)
}

func (svc *UploadService) saveFileData(ctx context.Context, ownerID uuid.UUID, r io.Reader, fileName string) error {
	key := fmt.Sprintf("%s/%s", ownerID, url.PathEscape(fileName))

	contentType := mime.TypeByExtension(filepath.Ext(fileName))
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	opts := &blob.WriterOptions{ContentType: contentType}

	err := svc.BlobStorageClient.MultipartUpload(ctx, key, r, opts)
	if err != nil {
		return err
	}

	return nil
}

func (svc *UploadService) saveMetadata(ctx context.Context, metadata FileMetadata) error {
	exists, err := svc.Db.CheckExists(ctx, metadata.FileID, metadata.Version)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("file metadata already present in database")
	}

	if err = svc.Db.PersistMetadata(ctx, metadata); err != nil {
		return err
	}

	return nil
}
