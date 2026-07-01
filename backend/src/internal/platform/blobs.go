package platform

import (
	"context"
	"gocloud.dev/blob"
	_ "gocloud.dev/blob/s3blob"
	"io"
)

type BlobStorageClient struct {
	Bucket *blob.Bucket
}
