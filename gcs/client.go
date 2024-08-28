package gcs

import (
	"context"
	"fmt"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

// Gcs holds the Google Cloud Storage client and bucket handle.
type Gcs struct {
	Client *storage.Client
	Bucket *storage.BucketHandle
}

// GoogleCloudStorage initializes a new Google Cloud Storage client and bucket handle
// for the specified context and credentials file.
func GoogleCloudStorage(ctx context.Context) (*Gcs, error) {
	// Create client options with the provided credentials file.
	opts := []option.ClientOption{
		option.WithCredentialsFile(CredentialFile),
	}

	// Create a new storage client.
	client, err := storage.NewClient(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to storage client: %v", err)
	}

	// Initialize the bucket handle using the global BucketName.
	bucket := client.Bucket(BucketName)

	// Return the initialized Gcs struct.
	return &Gcs{
		Client: client,
		Bucket: bucket,
	}, nil
}
