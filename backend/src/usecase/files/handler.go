package files

import (
	"backend/src/internal/api/message"
	"context"
	"log"
	"net/http"
)

type service interface {
	Upload(ctx context.Context, r *http.Request) ([]MetaData, error)
	FindAllMetadata(ctx context.Context, request GetAllMetadataRequest) ([]MetaDataResponse, error)
	FindMetadata(ctx context.Context, request FindMetadataRequest) ([]MetaData, error)
	Delete(ctx context.Context, request DeleteRequest) error
	MoveToRubbish(ctx context.Context, request DeleteRequest) error
}

type Handler struct {
	svc service
}

// Constructor
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// Upload Handler method for file multipart form upload.
func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	svc := h.svc
	var newMetadata []MetaData

	r.Body = http.MaxBytesReader(w, r.Body, 100<<20)

	newMetadata, err := svc.Upload(r.Context(), r)
	if err != nil {
		log.Printf("unable to save uploaded file: %v", err)
		http.Error(w, "could not save file", http.StatusInternalServerError)
	}

	log.Printf("Metadata: %v", newMetadata)
	if err := message.Response(w, "uploaded", newMetadata); err != nil {
		log.Print(err)
		return
	}
}

// GetAll method for retrieving all file metadata.
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	svc := h.svc

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

	files, err := svc.FindAllMetadata(r.Context(), request)
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

// FindMetadata method for searching files by metadata contents.
func (h *Handler) FindMetadata(w http.ResponseWriter, r *http.Request) {
	svc := h.svc
	var request FindMetadataRequest

	if err := request.Bind(r); err != nil {
		log.Printf("unable to bind raw request to FindMetadataRequest, %v", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
	}

	files, err := svc.FindMetadata(r.Context(), request)
	if err != nil {
		log.Printf("couldn't get all user's files: %v", err)
		http.Error(w, "unable to get user's files", http.StatusInternalServerError)
	}

	if err = message.Response(w, "fetched user's files", files); err != nil {
		log.Printf("unable to return response, %v", err)
		return
	}
}

func (h *Handler) Download(w http.ResponseWriter, r *http.Request) {
	panic("not implemented.")
}

func (h *Handler) UpdateMetadata(w http.ResponseWriter, r *http.Request) {
	panic("not implemented.")
}

// Delete handler method for the permanent removal of files and their respective metadata in DB.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	var request DeleteRequest

	if err := request.Bind(r); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if err := h.svc.Delete(r.Context(), request); err != nil {
		http.Error(w, "unable to delete file", http.StatusExpectationFailed)
		return
	}

	if err := message.Response(w, "deleted", nil); err != nil {
		log.Printf("unable to return response: %v", err)
		return
	}
}

// TempDelete handler method for the flagging of data to be removed after
// specified timelapse.
func (h *Handler) TempDelete(w http.ResponseWriter, r *http.Request) {
	var request DeleteRequest

	if err := request.Bind(r); err != nil {
		log.Printf("unable to bind request to DeleteRequest %v", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
	}

	if err := h.svc.MoveToRubbish(r.Context(), request); err != nil {
		log.Printf("unable to delete file metadata, %v", err)
		http.Error(w, "unable to delete file", http.StatusExpectationFailed)
	}
}
