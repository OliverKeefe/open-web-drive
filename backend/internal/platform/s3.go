package platform

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type S3Client struct {
	Client *s3.Client
}

func NewS3Client() (S3Client, error) {
	awsEndpoint := "http://localhost:4566"
	awsRegion := "us-east-1"

	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(awsRegion),
	)
	if err != nil {
		log.Fatalf("cannot load the aws configs: %s", err)
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
		o.BaseEndpoint = aws.String(awsEndpoint)
	})

	out, err := client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		panic(err)
	}

	for _, bucket := range out.Buckets {
		log.Printf("Bucket Name: %v", &bucket.Name)
	}

	return S3Client{
		Client: client,
	}, nil
}

func (s3 *S3Client) SaveToS3(ctx context.Context, bucket string, reader io.Reader, ownerID uuid.UUID, fileName string) error {
	key := fmt.Sprintf("%s/%s", ownerID, fileName)

	var partMiBs int64 = 5

	uploader := transfermanager.New(s3.Client, func(opts *transfermanager.Options) {
		opts.PartSizeBytes = partMiBs * 1024 * 1024
		opts.Concurrency = 3
	})

	_, err := uploader.UploadObject(ctx, &transfermanager.UploadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   reader,
	})
	if err != nil {
		log.Printf("s3 transfer manager put error: %v", err)
		return err
	}

	return nil

}
