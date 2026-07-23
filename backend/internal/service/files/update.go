package files

import (
	"backend/internal/platform"
	"context"
	"net/http"

	"github.com/google/uuid"
)

type updateMetadataRepository interface {
	CheckExists(ctx context.Context, fileID uuid.UUID, version int) (bool, error)
}

type UpdateMetadataService struct {
	Db                updateMetadataRepository
	BlobStorageClient *platform.BlobStorageClient
	BucketURL         string
}

func NewUpdateMetadataService(db updateMetadataRepository, client *platform.BlobStorageClient, bucketUrl string) *UpdateMetadataService {
	return &UpdateMetadataService{
		Db:                db,
		BlobStorageClient: client,
		BucketURL:         bucketUrl,
	}
}

func (svc *UpdateMetadataService) Handle(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

type updateMetadataRequest struct {
}

func (svc *UpdateMetadataService) execute(req updateMetadataRequest) error {
	panic("not implemented")
}
