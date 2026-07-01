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

func (b *BlobStorageClient) Upload(ctx context.Context, key string, data []byte, opts *blob.WriterOptions) error {
	return b.Bucket.WriteAll(ctx, key, data, opts)
}

func (b *BlobStorageClient) Download(ctx context.Context, key string, writer io.Writer, opts *blob.ReaderOptions) error {
	return b.Bucket.Download(ctx, key, writer, opts)
}
