package router

import (
	"backend/src/internal/auth"
	"backend/src/internal/database"
	"backend/src/internal/middleware"
	"backend/src/internal/platform"
	"backend/src/internal/service/files"
	"context"
	"net/http"
)

var (
	route = func(a *auth.Authenticator, h http.HandlerFunc) http.Handler {
		return middleware.Protect(a, h)
	}
)

func RegisterFileRoutes(mux *http.ServeMux, a *auth.Authenticator, db *database.MetadataDatabase) error {
	bucketUrl := "test-bucket"

	client, err := platform.NewBlobStorageClient(context.Background(), "test-bucket")
	if err != nil {
		panic(err)
	}

	repository := files.NewFileRepository(db.Pool)

	uploadSvc := files.NewUploadService(repository, client, bucketUrl)
	downloadSvc := files.NewDownloadService(repository, client, bucketUrl)
	deleteSvc := files.NewDeleteService(repository, client, bucketUrl)
	updateSvc := files.NewUpdateMetadataService(repository, client, bucketUrl)

	uploadEndpoint := route(a, uploadSvc.Handle)
	mux.Handle(
		"POST /api/files/upload",
		uploadEndpoint,
	)
	downloadRoute := route(a, downloadSvc.Handle)
	mux.Handle(
		"POST /api/files/download",
		downloadRoute,
	)
	deleteRoute := route(a, deleteSvc.Handle)
	mux.Handle(
		"DELETE /api/files/delete",
		deleteRoute,
	)
	updateMetadataRoute := route(a, updateSvc.Handle)
	mux.Handle(
		"POST /api/files/delete",
		updateMetadataRoute,
	)

	return nil
}
