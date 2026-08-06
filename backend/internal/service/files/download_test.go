package files

import (
	"backend/internal/auth"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"gocloud.dev/blob"
)

type mockDownloadRepository struct {
	findMetadataByIDFunc func(ctx context.Context, id uuid.UUID) (FileMetadata, error)
}

func (m *mockDownloadRepository) FindMetadataByID(ctx context.Context, id uuid.UUID) (FileMetadata, error) {
	if m.findMetadataByIDFunc != nil {
		return m.findMetadataByIDFunc(ctx, id)
	}
	return FileMetadata{}, nil
}

type mockDownloadBlob struct {
	downloadFunc func(ctx context.Context, key string, writer io.Writer, opts *blob.ReaderOptions) error
}

func (m *mockDownloadBlob) Download(ctx context.Context, key string, writer io.Writer, opts *blob.ReaderOptions) error {
	if m.downloadFunc != nil {
		return m.downloadFunc(ctx, key, writer, opts)
	}
	return nil
}

func TestDownloadService_Handle_Success(t *testing.T) {
	ownerID := uuid.New()
	metaID := uuid.New()

	repo := &mockDownloadRepository{
		findMetadataByIDFunc: func(ctx context.Context, id uuid.UUID) (FileMetadata, error) {
			if id != metaID {
				t.Errorf("FindMetadataByID called with wrong id: got %v, want %v", id, metaID)
			}
			return FileMetadata{
				ID:       metaID,
				OwnerID:  ownerID,
				FileName: "report.pdf",
				Size:     2048,
				FileType: ".pdf",
			}, nil
		},
	}

	blob := &mockDownloadBlob{
		downloadFunc: func(ctx context.Context, key string, writer io.Writer, opts *blob.ReaderOptions) error {
			expectedKey := ownerID.String() + "/report.pdf"
			if key != expectedKey {
				t.Errorf("Download key = %q; want %q", key, expectedKey)
			}
			_, err := writer.Write([]byte("fake pdf content"))
			return err
		},
	}

	svc := NewDownloadService(repo, blob, "test-bucket")

	body, _ := json.Marshal(DownloadRequest{ID: metaID})
	req := httptest.NewRequest(http.MethodPost, "/api/files/download", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	ctx := context.Background()
	authenticator := &auth.Authenticator{}
	ctx = authenticator.InjectUserID(ctx, ownerID.String())
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()
	svc.Handle(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/octet-stream" {
		t.Errorf("Content-Type = %q; want %q", ct, "application/octet-stream")
	}
	if cd := rec.Header().Get("Content-Disposition"); cd != `attachment; filename="report.pdf"` {
		t.Errorf("Content-Disposition = %q; want %q", cd, `attachment; filename="report.pdf"`)
	}
	if cl := rec.Header().Get("Content-Length"); cl != "2048" {
		t.Errorf("Content-Length = %q; want %q", cl, "2048")
	}
	if rec.Body.String() != "fake pdf content" {
		t.Errorf("body = %q; want %q", rec.Body.String(), "fake pdf content")
	}
}

func TestDownloadService_Handle_NotFound(t *testing.T) {
	repo := &mockDownloadRepository{
		findMetadataByIDFunc: func(ctx context.Context, id uuid.UUID) (FileMetadata, error) {
			return FileMetadata{}, errors.New("no rows")
		},
	}
	blob := &mockDownloadBlob{}

	svc := NewDownloadService(repo, blob, "test-bucket")

	body, _ := json.Marshal(DownloadRequest{ID: uuid.New()})
	req := httptest.NewRequest(http.MethodPost, "/api/files/download", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	ctx := context.Background()
	authenticator := &auth.Authenticator{}
	ctx = authenticator.InjectUserID(ctx, uuid.New().String())
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()
	svc.Handle(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}

func TestDownloadService_Handle_Forbidden(t *testing.T) {
	ownerID := uuid.New()
	otherUser := uuid.New()
	metaID := uuid.New()

	repo := &mockDownloadRepository{
		findMetadataByIDFunc: func(ctx context.Context, id uuid.UUID) (FileMetadata, error) {
			return FileMetadata{
				ID:       metaID,
				OwnerID:  ownerID,
				FileName: "secret.pdf",
				Size:     1024,
			}, nil
		},
	}
	blob := &mockDownloadBlob{}

	svc := NewDownloadService(repo, blob, "test-bucket")

	body, _ := json.Marshal(DownloadRequest{ID: metaID})
	req := httptest.NewRequest(http.MethodPost, "/api/files/download", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	ctx := context.Background()
	authenticator := &auth.Authenticator{}
	ctx = authenticator.InjectUserID(ctx, otherUser.String())
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()
	svc.Handle(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}

func TestDownloadService_Handle_BlobError(t *testing.T) {
	ownerID := uuid.New()
	metaID := uuid.New()

	repo := &mockDownloadRepository{
		findMetadataByIDFunc: func(ctx context.Context, id uuid.UUID) (FileMetadata, error) {
			return FileMetadata{
				ID:       metaID,
				OwnerID:  ownerID,
				FileName: "report.pdf",
				Size:     2048,
			}, nil
		},
	}

	blob := &mockDownloadBlob{
		downloadFunc: func(ctx context.Context, key string, writer io.Writer, opts *blob.ReaderOptions) error {
			return errors.New("blob read failed")
		},
	}

	svc := NewDownloadService(repo, blob, "test-bucket")

	body, _ := json.Marshal(DownloadRequest{ID: metaID})
	req := httptest.NewRequest(http.MethodPost, "/api/files/download", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	ctx := context.Background()
	authenticator := &auth.Authenticator{}
	ctx = authenticator.InjectUserID(ctx, ownerID.String())
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()
	svc.Handle(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}
