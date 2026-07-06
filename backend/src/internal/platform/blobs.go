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

func NewBlobStorageClient(ctx context.Context, bucketURL string) (*BlobStorageClient, error) {
	bucket, err := blob.OpenBucket(ctx, bucketURL)
	if err != nil {
		return nil, err
	}

	return &BlobStorageClient{
		Bucket: bucket,
	}, nil
}

func (b *BlobStorageClient) MultipartUpload(ctx context.Context, key string, dataStream io.Reader, opts *blob.WriterOptions) error {
	w, err := b.Bucket.NewWriter(ctx, key, opts)
	if err != nil {
		return fmt.Errorf("failed to open cloud writer: %w", err)
	}

	if _, err := io.Copy(w, dataStream); err != nil {
		return fmt.Errorf("multipart upload failed during data stream copy: %w", err)
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("failed to close cloud writer: %w", err)
	}

	return nil
}

func (b *BlobStorageClient) Download(ctx context.Context, key string, writer io.Writer, opts *blob.ReaderOptions) error {
	return b.Bucket.Download(ctx, key, writer, opts)
}
