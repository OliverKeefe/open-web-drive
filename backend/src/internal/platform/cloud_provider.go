package platform

import "github.com/aws/aws-sdk-go-v2/service/s3"

type Provider struct {
	ObjectStorageClient *s3.Client
}

// TODO: build wrapper for each provider.
type ObjectStorageProvider interface{}
