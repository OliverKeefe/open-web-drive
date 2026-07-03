package files

import (
	"backend/src/internal/api/message"
	"backend/src/internal/auth"
	"backend/src/internal/database"
	"backend/src/internal/platform"
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"io"
	"log"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type UploadService struct {
	repository *FileRepository
	s3Client   *platform.S3Client
	provider   *platform.Provider
	bucket     *string
}

func NewUploadService(client *s3.Client, bucket string) *UploadService {
	return &UploadService{
		s3Client: client,
		bucket:   bucket,
	}
}

func (svc *UploadService) Handle(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromCtx(r.Context())
	if !ok {
		http.Error(w, "upload failed", http.StatusUnauthorized)
		return
	}

	hasClaim, err := auth.HasClaim(r.Context(), uuid.MustParse(userID), "permissions:files:write")
	if err != nil {
		slog.Error("upload failed", "error", err)
	}

	if !hasClaim {
		http.Error(w, "upload failed", http.StatusUnauthorized)
		return
	}

	// FYI - 5<<20 is a bitshift operation (5*2^20 = 5,242,880 Bytes or 5 MegaBytes)
	r.Body = http.MaxBytesReader(w, r.Body, 5<<20)
	if err := svc.execute(r); err != nil {
		slog.Error("upload failed", "error", err)
		http.Error(w, "upload failed", http.StatusBadRequest)
	}

	if err := message.Response(w, "upload succeeded", nil); err != nil {
		slog.Error("upload failed", "error", err)
		return
	}

}

// The raw metadata part of a multipart upload.
type multipartMetadata struct {
	Path             string            `json:"path"`
	RelativePath     string            `json:"relativePath"`
	LastModified     int64             `json:"lastModified"`
	LastModifiedDate string            `json:"lastModifiedDate"`
	UploadedAt       int64             `json:"uploadedAt"`
	Size             uint64            `json:"size"`
	FileType         string            `json:"fileType"`
	ID               string            `json:"id"`
	OwnerID          string            `json:"owner_id"`
	FileName         string            `json:"file_name"`
	CreatedAt        int64             `json:"created_at"`
	Permissions      []FilePermissions `json:"file_permissions"`
}

func (svc *UploadService) execute(r *http.Request) error {
	mr, err := r.MultipartReader()
	if err != nil {
		return err
	}

	ctx := r.Context()

	metadataByID := make(map[string]FileMetadata)

	userID, ok := auth.UserIDFromCtx(ctx)
	if !ok {
		return errors.New("unable to get UserID from context")
	}

	claims, ok := auth.ClaimsFromCtx(r.Context())
	if !ok {
		return errors.New("unable to get claims from context.")
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

		var (
			decodedRequest multipartMetadata
		)
		if err := json.NewDecoder(part).Decode(&decodedRequest); err != nil {
			return err
		}

		switch {
		case strings.HasPrefix(name, "metadata-"):
			idAsStr := strings.TrimPrefix(name, "metadata-")

			versionID, err := uuid.NewV7()
			if err != nil {
				return err
			}

			metadataByID[idAsStr] = FileMetadata{
				ID:           uuid.MustParse(idAsStr),
				OwnerID:      uuid.MustParse(userID),
				FileName:     decodedRequest.FileName,
				Path:         decodedRequest.Path,
				RelativePath: decodedRequest.RelativePath,
				Size:         decodedRequest.Size,
				FileType:     decodedRequest.FileType,
				ModifiedAt:   time.UnixMilli(decodedRequest.LastModified),
				UploadedAt:   time.UnixMilli(decodedRequest.UploadedAt),
				CreatedAt:    time.UnixMilli(decodedRequest.CreatedAt),
				Version:      versionID,
				Permissions:  decodedRequest.Permissions,
			}

		case strings.HasPrefix(name, "file-"):
			idToStr := strings.TrimPrefix(name, "file-")

			data, err := io.ReadAll(part)
			if err != nil {
				return nil, err
			}

			err = svc.saveFileData(
				r.Context(),
				uuid.MustParse(decodedRequest.OwnerID),
				bytes.NewReader(data),
				part.FileName(),
			)
			md := metadataByID[idToStr]
			md.Hash = sha256.Sum256(data)
			metadataByID[idToStr] = md
		}
	}
	var newMetadata []FileMetadata
	for _, md := range metadataByID {
		newMd, err := svc.saveMetadata(r.Context(), md)
		if err != nil {
			return err
		}
		newMetadata = append(newMetadata, newMd)
	}

	return nil
}

func (svc *UploadService) saveFileData(ctx context.Context, ownerID uuid.UUID, r bytes.Reader, fileName string) error {
	err := svc.s3Client.SaveToS3(ctx, ownerID, r, &svc.bucket, fileName)
	if err != nil {
		return err
	}
}

func (svc *UploadService) saveMetadata(ctx context.Context, metadata FileMetadata) error {
	exists, err := svc.Db.CheckExists(ctx, metadata.ID)
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
