package files

import (
	"backend/src/internal/auth"
	"backend/src/internal/platform"
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type DownloadService struct {
	Db                DB
	BlobStorageClient *platform.BlobStorageClient
	BucketURL         string
}

func NewDownloadService(db DB, client *platform.BlobStorageClient, bucketUrl string) *DownloadService {
	return &DownloadService{
		Db:                db,
		BlobStorageClient: client,
		BucketURL:         bucketUrl,
	}
}

type DownloadRequest struct {
	ID []uuid.UUID `json:"id"`
}

func (svc *DownloadService) Handle(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

func (svc *DownloadService) execute(ctx context.Context, request DownloadRequest) error {
	userID, ok := auth.UserIDFromCtx(ctx)
	if !ok {
		return errors.New("unable to get userID from context")
	}

	ownerID, err := uuid.Parse(userID)
	if err != nil {
		log.Printf("unable to parse userID string to uuid")
	}

	chunk, err := download(ctx, request.ID)
	if err != nil {
		log.Printf("could not delete file metadata, %v", err)
		return err
	}

	return chunk, nil
}

func download(ctx context.Context, key uuid.UUID) ([]byte, error) {
	panic("not implemented")
}
