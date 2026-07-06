package files

import (
	"backend/src/internal/api/message"
	"backend/src/internal/auth"
	"backend/src/internal/platform"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type DeleteService struct {
	Db                DB
	BlobStorageClient *platform.BlobStorageClient
	BucketURL         string
}

func NewDeleteService(db DB, client *platform.BlobStorageClient, bucketUrl string) *DeleteService {
	return &DeleteService{
		Db:                db,
		BlobStorageClient: client,
		BucketURL:         bucketUrl,
	}
}

func (svc *DeleteService) Handle(w http.ResponseWriter, r *http.Request) {
	var request DeleteRequest

	if err := request.Bind(r); err != nil {
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
		log.Printf("unable to parse userID string to uuid")
	}

	err = svc.Db.DeleteFileData(ctx, request.ID, ownerID)
	if err != nil {
		return fmt.Errorf("unable to delete file data, %v", err)
	}

	err = svc.Db.DeleteMetadata(ctx, request.ID, ownerID)
	if err != nil {
		log.Printf("could not delete file metadata, %v", err)
		return err
	}

	return nil
}
