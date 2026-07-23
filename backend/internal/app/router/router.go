package router

import (
	"backend/internal/auth"
	"backend/internal/database"
	"backend/internal/middleware"
	"backend/internal/platform"
	files2 "backend/internal/service/files"
	"context"
	"net/http"
	"os"
)

var (
	route = func(a *auth.Authenticator, h http.HandlerFunc) http.Handler {
		return middleware.Protect(a, h)
	}
)

func RegisterFileRoutes(mux *http.ServeMux, a *auth.Authenticator, db *database.MetadataDatabase) error {
	bucketURL := os.Getenv("BLOB_BUCKET_URL")
	if bucketURL == "" {
		bucketURL = "s3://temp-buck"
	}

	client, err := platform.NewBlobStorageClient(context.Background(), bucketURL)
	if err != nil {
		panic(err)
	}

	repository := files2.NewFileRepository(db.Pool)

	uploadSvc := files2.NewUploadService(repository, client, bucketURL)
	downloadSvc := files2.NewDownloadService(repository, client, bucketURL)
	deleteSvc := files2.NewDeleteService(repository, client, bucketURL)
	updateSvc := files2.NewUpdateMetadataService(repository, client, bucketURL)
	searchSvc := files2.NewSearchService(repository, client, bucketURL)

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
		"PUT /api/files/metadata",
		updateMetadataRoute,
	)
	searchRoute := route(a, searchSvc.Handle)
	mux.Handle(
		"POST /api/files/get-all",
		searchRoute,
	)

	return nil
}
