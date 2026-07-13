package files

import (
	"context"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
	"gocloud.dev/blob"
)

type downloadRepository interface {
	CheckExists(ctx context.Context, ID uuid.UUID) (bool, error)
}

type downloadBlob interface {
	Download(ctx context.Context, key string, writer io.Writer, opts *blob.ReaderOptions) error
}

type DownloadService struct {
	Db                downloadRepository
	BlobStorageClient downloadBlob
	BucketURL         string
}

func NewDownloadService(db downloadRepository, client downloadBlob, bucketUrl string) *DownloadService {
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

func (svc *DownloadService) execute(ctx context.Context, request DownloadRequest) ([]byte, error) {
	//userID, ok := auth.UserIDFromCtx(ctx)
	//if !ok {
	//	return errors.New("unable to get userID from context")
	//}

	//ownerID, err := uuid.Parse(userID)
	//if err != nil {
	//	log.Printf("unable to parse userID string to uuid")
	//}

	chunk, err := svc.download(ctx, request.ID)
	if err != nil {
		log.Printf("could not delete file metadata, %v", err)
		return nil, err
	}

	return chunk, nil
}

func (svc *DownloadService) download(ctx context.Context, keys []uuid.UUID) ([]byte, error) {
	panic("not implemented")
}
