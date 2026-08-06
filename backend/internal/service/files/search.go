package files

import (
	"backend/internal/api/message"
	"backend/internal/auth"
	"backend/internal/platform"
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type searchRepository interface {
	FindMetadataByID(ctx context.Context, ID uuid.UUID) (FileMetadata, error)
	FindMetadata(ctx context.Context, m FileMetadata) ([]FileMetadata, error)
	FindAllMetadata(ctx context.Context, ownerID uuid.UUID, cursor *MetadataCursor, limit int) ([]FileMetadata, error)
}

type SearchService struct {
	Db                searchRepository
	BlobStorageClient *platform.BlobStorageClient
	BucketURL         string
}

func NewSearchService(db searchRepository, client *platform.BlobStorageClient, bucketUrl string) *SearchService {
	return &SearchService{
		Db:                db,
		BlobStorageClient: client,
		BucketURL:         bucketUrl,
	}
}

func (svc *SearchService) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil || r.ContentLength <= 0 {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	request, err := message.Bind[GetAllMetadataRequest](r)
	if err != nil {
		log.Printf("couldn't map http request to dto: %v", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	files, err := svc.execute(r.Context(), request)
	if err != nil {
		log.Printf("couldn't get all user's files: %v", err)
		http.Error(w, "unable to get user's files", http.StatusInternalServerError)
		return
	}

	if err = message.Response(w, "fetched user's files", files); err != nil {
		log.Printf("unable to return response, %v", err)
		return
	}
}

func (svc *SearchService) execute(ctx context.Context, request GetAllMetadataRequest) ([]FileMetadataResponse, error) {
	userID, ok := auth.UserIDFromCtx(ctx)
	if !ok {
		return nil, errors.New("unable to get userID from context")
	}
	ownerID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user id in JWT claims")
	}

	files, err := svc.Db.FindAllMetadata(ctx, ownerID, request.Cursor, request.Limit)
	if err != nil {
		return nil, err
	}

	response := make([]FileMetadataResponse, len(files))
	for i, file := range files {
		response[i] = FileMetadataResponse(file)
	}

	return response, nil
}

func (svc *SearchService) findMetadata(ctx context.Context, request FindMetadataRequest) ([]FileMetadata, error) {
	model := request.ToModel()
	return svc.Db.FindMetadata(ctx, model)
}
