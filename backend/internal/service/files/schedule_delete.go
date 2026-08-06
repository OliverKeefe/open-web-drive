package files

import (
	"backend/internal/auth"
	"backend/internal/platform"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type scheduleDeleteRepository interface {
	MarkForDeletion(ctx context.Context, id uuid.UUID, ownerId uuid.UUID) error
}

type ScheduleDeleteService struct {
	db                scheduleDeleteRepository
	BlobStorageClient *platform.BlobStorageClient
	BucketURL         string
}

type scheduleDeleteRequest struct {
	ID []uuid.UUID `json:"id"`
}

func (svc *ScheduleDeleteService) TempDelete(w http.ResponseWriter, r *http.Request) {
	var request DeleteRequest

	if err := request.Bind(r); err != nil {
	request, err := message.Bind[scheduleDeleteRequest](r)
	if err != nil {
		log.Printf("unable to bind request to DeleteRequest %v", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if err := svc.execute(r.Context(), request); err != nil {
		log.Printf("unable to delete file metadata, %v", err)
		http.Error(w, "unable to delete file", http.StatusExpectationFailed)
		return
	}
}

func (svc *ScheduleDeleteService) execute(ctx context.Context, request DeleteRequest) error {
	userID, ok := auth.UserIDFromCtx(ctx)
	if !ok {
		return errors.New("unable to get userID from context")
	}

	ownerID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("unable to parse userID string to uuid: %w", err)
	}

	err = svc.db.MarkForDeletion(ctx, request.ID, ownerID)
	if err != nil {
		log.Printf("unable to move file or metadata to rubbish bin, %v", err)
		return err
	}

	return nil
}
