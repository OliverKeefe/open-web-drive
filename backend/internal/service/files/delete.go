package files

import (
	"backend/internal/api/message"
	"backend/internal/auth"
	"backend/internal/platform"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type deleteRepository interface {
	DeleteMetadata(ctx context.Context, id uuid.UUID, ownerId uuid.UUID) error
	DeleteFileData(ctx context.Context, id uuid.UUID, ownerId uuid.UUID) error
}

type DeleteService struct {
	Db                deleteRepository
	BlobStorageClient *platform.BlobStorageClient
	BucketURL         string
}

func NewDeleteService(db deleteRepository, client *platform.BlobStorageClient, bucketUrl string) *DeleteService {
	return &DeleteService{
		Db:                db,
		BlobStorageClient: client,
		BucketURL:         bucketUrl,
	}
}

func (svc *DeleteService) Handle(w http.ResponseWriter, r *http.Request) {
	request, err := message.Bind[DeleteRequest](r)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if err := svc.execute(r.Context(), request); err != nil {
		http.Error(w, "unable to delete file", http.StatusExpectationFailed)
		return
	}

	if err := message.Response(w, "deleted", nil); err != nil {
		log.Printf("unable to return response: %v", err)
		return
	}
}

func (svc *DeleteService) execute(ctx context.Context, request DeleteRequest) error {
	userID, ok := auth.UserIDFromCtx(ctx)
	if !ok {
		return errors.New("unable to get userID from context")
	}

	ownerID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("unable to parse userID string to uuid: %w", err)
	}

	err = svc.Db.DeleteFileData(ctx, request.ID, ownerID)
	if err != nil {
		return fmt.Errorf("unable to delete file data, %w", err)
	}

	return nil
}
