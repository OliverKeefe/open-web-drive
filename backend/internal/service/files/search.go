package files

import (
	"backend/internal/api/message"
	"backend/internal/platform"
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type searchRepository interface {
	FindMetadataByID(ctx context.Context, ID uuid.UUID) (Metadata, error)
	FindMetadata(m MetaData) (string, []any)
	FindAllMetadata(ctx context.Context, req GetAllMetadataRequest) ([]MetaData, error)
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

	var request GetAllMetadataRequest
	err := request.Bind(r)
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

func (svc *SearchService) execute(ctx context.Context, request GetAllMetadataRequest) ([]MetaDataResponse, error) {
	var (
		files    []MetaData
		response []MetaDataResponse
	)

	files, err := svc.Db.FindAllMetadata(ctx, request)
	if err != nil {
		log.Printf("unable to get all files for user: %s, %v ", request.UserID, err)
	}

	for _, file := range files {
		file := file.ToResponse()
		if err != nil {
			log.Printf("unable to map file metadata: %v, to dto: %v", file, err)
		}
		response = append(response, file)
	}

	return response, nil
}

func (svc *SearchService) findMetadata(ctx context.Context, request FindMetadataRequest) ([]MetaData, error) {
	//var (
	//	files []MetaData
	//)
	//
	//model := request.ToModel()
	//files, err := svc.Db.FindMetadata(ctx, model)
	//if err != nil {
	//	log.Printf("unable to get file metadata: %v", err)
	//	return files, err
	//}
	//
	//return files, nil
	panic("not implemented")
}
