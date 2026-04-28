package router

import (
	"backend/src/internal/auth"
	"backend/src/internal/db/metadb"
	"backend/src/internal/middleware"
	"backend/src/internal/platform"
	filesvc "backend/src/usecase/files"
	"net/http"
)

var (
	route = func(a *auth.Authenticator, h http.HandlerFunc) http.Handler {
		return middleware.Protect(a, h)
	}
)

func RegisterFileRoutes(
	mux *http.ServeMux,
	a *auth.Authenticator,
	db *metadb.MetadataDatabase,
) {

	s3Client, err := platform.NewS3Client()
	if err != nil {
		panic("Can't connect to S3 bucket.")
	}
	repo := filesvc.NewRepository(db.Pool, s3Client.Client, "temp-buck")
	svc := filesvc.NewService(repo)
	h := filesvc.NewHandler(svc)

	upload := route(a, h.Upload)
	mux.Handle(
		"POST /api/files/upload",
		upload,
	)
	findMetadata := route(a, h.FindMetadata)
	mux.Handle(
		"POST /api/files/find",
		findMetadata,
	)
	getAllMetadata := route(a, h.GetAll)
	mux.Handle(
		"POST /api/files/get-all",
		getAllMetadata,
	)
	deleteFile := route(a, h.Delete)
	mux.Handle(
		"POST /api/files/delete",
		deleteFile,
	)
}
