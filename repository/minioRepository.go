package repository

import (
	"file-sharing/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"golang.org/x/net/context"
	"os"
	"time"
)

type MinioRepository struct {
	client     *minio.Client
	bucketName string
}

func NewMinioRepository(cfg config.MinioConfig) (*MinioRepository, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	// Verify bucket exists (critical fix for 500 errors)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	exists, err := client.BucketExists(ctx, cfg.BucketName)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, os.ErrNotExist // Or create it automatically
	}

	return &MinioRepository{
		client:     client,
		bucketName: cfg.BucketName,
	}, nil
}
func (r *MinioRepository) UploadFile(ctx context.Context, objectName, filePath string) error {
	if _, err := os.Stat(filePath); err != nil {
		return err
	}

	// Upload with timeout
	uploadCtx, cancel := context.WithTimeout(ctx, 100*time.Second)
	defer cancel()

	_, err := r.client.FPutObject(uploadCtx, r.bucketName, objectName, filePath, minio.PutObjectOptions{
		ContentType: "application/octet-stream", // Fix for MIME type issues
	})
	return err
}

func (r *MinioRepository) GetFileURL(ctx context.Context, objectName string) (string, error) {
	url, err := r.client.PresignedGetObject(ctx, r.bucketName, objectName, 24*time.Hour, nil)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}
