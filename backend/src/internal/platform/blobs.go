package platform

import (
	"context"
	"fmt"
	"io"

	"gocloud.dev/blob"
	_ "gocloud.dev/blob/s3blob"
)

type BlobStorageClient struct {
	Bucket *blob.Bucket
}

func (b *BlobStorageClient) Upload(ctx context.Context, key string, data []byte, opts *blob.WriterOptions) error {
	return b.Bucket.WriteAll(ctx, key, data, opts)
func NewBlobStorageClient(ctx context.Context, bucketURL string) (*BlobStorageClient, error) {
	bucket, err := blob.OpenBucket(ctx, bucketURL)
	if err != nil {
		return nil, err
	}

	return &BlobStorageClient{
		Bucket: bucket,
	}, nil
}

}

func (b *BlobStorageClient) Download(ctx context.Context, key string, writer io.Writer, opts *blob.ReaderOptions) error {
	return b.Bucket.Download(ctx, key, writer, opts)
}
