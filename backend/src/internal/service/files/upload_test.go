package files

import (
	"backend/src/internal/auth"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v4"
	"gocloud.dev/blob"
)

type mockBlobStorage struct {
	multipartUploadFunc func(ctx context.Context, key string, dataStream io.Reader, opts *blob.WriterOptions) error
}

type mockUploadRepository struct {
	checkExistsFunc     func(ctx context.Context, ID uuid.UUID) (bool, error)
	persistMetadataFunc func(ctx context.Context, metadata FileMetadata) error
}

func (m *mockUploadRepository) CheckExists(ctx context.Context, ID uuid.UUID) (bool, error) {
	if m.checkExistsFunc != nil {
		return m.checkExistsFunc(ctx, ID)
	}
	return false, nil
}

func (m *mockUploadRepository) PersistMetadata(ctx context.Context, metadata FileMetadata) error {
	if m.persistMetadataFunc != nil {
		return m.persistMetadataFunc(ctx, metadata)
	}
	return nil
}

var (
	dummyMultipart = multipartMetadata{
		Path:             "parent/child1/child2/file.txt",
		RelativePath:     "/file.txt",
		LastModified:     0,
		LastModifiedDate: "",
		UploadedAt:       0,
		Size:             0,
		FileType:         "",
		ID:               "",
		OwnerID:          "",
		FileName:         "",
		CreatedAt:        0,
		Permissions:      nil,
	}
)

func Test_NewUploadService(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("unable to mock db, %v", err)
	}
	defer mockPool.Close()

	blobClient := &mockBlobStorage{}
	repo := FileRepository{Pool: mockPool}

	got := NewUploadService(&repo, blobClient, "test-bucket")

	if got == nil {
		t.Fatal("NewUploadService returned nil")
	}

	if got.BucketURL != "test-bucket" {
		t.Errorf("UploadService.BucketURL = %q; want %q", got.BucketURL, "test-bucket")
	}

	if got.BlobStorageClient != blobClient {
		t.Errorf("UploadService.BlobStorageClient was not set correctly")
	}

	if got.Db != &repo {
		t.Errorf("UploadService.Db was not set correctly")
	}
}

func TestUploadHandler_InvalidUploadRequest(t *testing.T) {
	tests := []struct {
		name           string
		uploadRequest  func(r *http.Request)
		expectedStatus int
	}{
		{
			name: "Missing Content-Type Header",
			uploadRequest: func(r *http.Request) {
				r.Header.Del("Content-Type")
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Incorrect Content-Type Header",
			uploadRequest: func(r *http.Request) {
				r.Header.Set("Content-Type", "application/json")
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Nil Request Body",
			uploadRequest: func(r *http.Request) {
				r.Header.Set("Content-Type", "multipart/form-data")
				r.Body = http.NoBody
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			metadataBytes, _ := json.Marshal(dummyMultipart)
			_ = writer.WriteField("metadata-", string(metadataBytes))

			filePart, _ := writer.CreateFormFile("filedata-", "file.txt")
			_, _ = filePart.Write([]byte("dummy file data"))

			writer.Close()

			r := httptest.NewRequest(http.MethodPost, "/api/files/upload", body)
			r.Header.Set("Content-Type", writer.FormDataContentType())

			tt.uploadRequest(r)

			pool, err := pgxmock.NewPool()
			if err != nil {
				t.Fatalf("unable to mock db, %v", err)
			}
			defer pool.Close()

			ctx := context.Background()
			blobClient := &mockBlobStorage{}
			repo := NewFileRepository(pool)
			uploadSvc := NewUploadService(repo, blobClient, "test-bucket")

			userID := uuid.New().String()
			authenticator := &auth.Authenticator{}
			ctx = authenticator.InjectUserID(ctx, userID)
			ctx = authenticator.InjectClaims(ctx, map[string]interface{}{
				"scopes": "permissions:files:write",
			})
			r = r.WithContext(ctx)

			recorder := httptest.NewRecorder()

			uploadSvc.Handle(recorder, r)

			if recorder.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, recorder.Code)
			}
		})
	}
}

}
