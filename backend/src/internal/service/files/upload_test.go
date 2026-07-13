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

type mockBlobStorageClient struct {
	UploadFunc func(ctx context.Context, bucket string, key string, data []byte) error
type mockBlobStorage struct {
	multipartUploadFunc func(ctx context.Context, key string, dataStream io.Reader, opts *blob.WriterOptions) error
}

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

	if got.BlobStorageClient != want.BlobStorageClient {
		t.Errorf("UploadService.BlobStorageClient was not initialized correctly")
	}

	if got.Db.pool == want.Db.pool {
		t.Errorf("UploadService.Db.Pool was not initialized")
	}

}
