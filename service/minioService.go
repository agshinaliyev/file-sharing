package service

import (
	"context"
	"github.com/minio/minio-go/v7"
	"time"
)

type MinioService interface {
	GeneratePresignedURL(filePath string, expiry time.Duration) (string, error)
	FileExists(filePath string) bool
}

type minioService struct {
	client *minio.Client
	bucket string
}

func NewMinioService(client *minio.Client, bucket string) MinioService {
	return &minioService{
		client: client,
		bucket: bucket,
	}
}

func (s *minioService) GeneratePresignedURL(filePath string, expiry time.Duration) (string, error) {
	url, err := s.client.PresignedGetObject(
		context.Background(),
		s.bucket,
		filePath,
		expiry,
		nil,
	)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}

func (s *minioService) FileExists(filePath string) bool {
	_, err := s.client.StatObject(
		context.Background(),
		s.bucket,
		filePath,
		minio.StatObjectOptions{},
	)
	return err == nil
}
