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

func (m *mockBlobStorage) MultipartUpload(ctx context.Context, key string, dataStream io.Reader, opts *blob.WriterOptions) error {
	if m.multipartUploadFunc != nil {
		return m.multipartUploadFunc(ctx, key, dataStream, opts)
	}
	return nil
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

func TestUploadService_saveMetadata(t *testing.T) {
	ctx := context.Background()
	metadata := FileMetadata{
		ID:       uuid.New(),
		FileName: "test.txt",
	}

	tests := []struct {
		name        string
		repo        *mockUploadRepository
		wantErr     bool
		errContains string
	}{
		{
			name: "success",
			repo: &mockUploadRepository{
				checkExistsFunc: func(ctx context.Context, ID uuid.UUID) (bool, error) {
					if ID != metadata.ID {
						t.Errorf("CheckExists called with wrong ID: got %v, want %v", ID, metadata.ID)
					}
					return false, nil
				},
				persistMetadataFunc: func(ctx context.Context, md FileMetadata) error {
					if md.ID != metadata.ID {
						t.Errorf("PersistMetadata called with wrong metadata: got ID %v, want %v", md.ID, metadata.ID)
					}
					return nil
				},
			},
			wantErr: false,
		},
		{
			name: "check exists returns error",
			repo: &mockUploadRepository{
				checkExistsFunc: func(ctx context.Context, ID uuid.UUID) (bool, error) {
					return false, errors.New("db connection failed")
				},
			},
			wantErr:     true,
			errContains: "db connection failed",
		},
		{
			name: "metadata already exists",
			repo: &mockUploadRepository{
				checkExistsFunc: func(ctx context.Context, ID uuid.UUID) (bool, error) {
					return true, nil
				},
			},
			wantErr:     true,
			errContains: "file metadata already present in database",
		},
		{
			name: "persist metadata returns error",
			repo: &mockUploadRepository{
				checkExistsFunc: func(ctx context.Context, ID uuid.UUID) (bool, error) {
					return false, nil
				},
				persistMetadataFunc: func(ctx context.Context, md FileMetadata) error {
					return errors.New("insert failed")
				},
			},
			wantErr:     true,
			errContains: "insert failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &UploadService{Db: tt.repo}

			err := svc.saveMetadata(ctx, metadata)

			if (err != nil) != tt.wantErr {
				t.Fatalf("saveMetadata error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && tt.errContains != "" {
				if err.Error() != tt.errContains {
					t.Fatalf("saveMetadata error = %q, want %q", err.Error(), tt.errContains)
				}
			}
		})
	}
}

func TestUploadService_saveFileData(t *testing.T) {
	ctx := context.Background()
	ownerID := uuid.New()
	fileName := "test.txt"

	tests := []struct {
		name        string
		blob        *mockBlobStorage
		wantErr     bool
		errContains string
	}{
		{
			name: "success",
			blob: &mockBlobStorage{
				multipartUploadFunc: func(ctx context.Context, key string, dataStream io.Reader, opts *blob.WriterOptions) error {
					expectedKey := ownerID.String() + "-" + fileName
					if key != expectedKey {
						t.Errorf("MultipartUpload key = %q; want %q", key, expectedKey)
					}
					data, err := io.ReadAll(dataStream)
					if err != nil {
						t.Fatalf("failed to read dataStream: %v", err)
					}
					if string(data) != "file content" {
						t.Errorf("MultipartUpload data = %q; want %q", string(data), "file content")
					}
					return nil
				},
			},
			wantErr: false,
		},
		{
			name: "upload error",
			blob: &mockBlobStorage{
				multipartUploadFunc: func(ctx context.Context, key string, dataStream io.Reader, opts *blob.WriterOptions) error {
					return errors.New("upload failed")
				},
			},
			wantErr:     true,
			errContains: "upload failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &UploadService{BlobStorageClient: tt.blob}
			data := strings.NewReader("file content")

			err := svc.saveFileData(ctx, ownerID, data, fileName)

			if (err != nil) != tt.wantErr {
				t.Fatalf("saveFileData() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && tt.errContains != "" {
				if err.Error() != tt.errContains {
					t.Fatalf("saveFileData() error = %q, want %q", err.Error(), tt.errContains)
				}
			}
		})
	}
}
