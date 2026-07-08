package files

import (
	"backend/src/internal/platform"
	"context"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
	"gocloud.dev/blob/memblob"
)

type mockBlobStorageClient struct {
	UploadFunc func(ctx context.Context, bucket string, key string, data []byte) error
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

func TestNewUploadService(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Errorf("unable to mock db, %v", err)
	}
	defer mockPool.Close()

	memBucket := memblob.OpenBucket(nil)
	defer memBucket.Close()

	mockBlobStorageClient := platform.NewBlobStorageClient(memBucket)
	var mockDb = DB{
		pool: mockPool,
	}

	want := UploadService{
		Db:                mockDb,
		BlobStorageClient: mockBlobStorageClient,
		bucket:            "test-bucket",
	}

	got := NewUploadService(mockDb, mockBlobStorageClient, "test-bucket")

	if got.bucket != want.bucket {
		t.Errorf("UploadService.bucket = %q; want %q", got.bucket, "test-bucket")
	}

	if got.BlobStorageClient != want.BlobStorageClient {
		t.Errorf("UploadService.BlobStorageClient was not initialized correctly")
	}

	if got.Db.pool == want.Db.pool {
		t.Errorf("UploadService.Db.Pool was not initialized")
	}

}
