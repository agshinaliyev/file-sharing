package config

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	log "github.com/sirupsen/logrus"
)

func NewMinIOClient() (*minio.Client, error) {
	endpoint := "localhost:9000" // Change if MinIO is hosted elsewhere
	accessKey := "minioadmin"    // Default MinIO credentials (replace in production)
	secretKey := "minioadmin123" // Default MinIO credentials
	useSSL := false
	bucketName := "uploads"

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalf("Failed to create MinIO client: %v", err)
		return nil, err
	}
	exists, error := client.BucketExists(context.Background(), bucketName)
	if error != nil {
		log.Println("Failed to check existing bucket", error)

	}
	if !exists {
		err := client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Println("Failed to create bucket", err)
		}
		log.Infof("Bucket created successfully!")
	}

	log.Println("Connected to MinIO")
	return client, nil
}
