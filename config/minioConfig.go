package config

import "os"

type MinioConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	BucketName      string
}

func LoadMinioConfig() MinioConfig {
	return MinioConfig{
		Endpoint:        os.Getenv("MINIO_ENDPOINT"),
		AccessKeyID:     os.Getenv("MINIO_ACCESS_KEY"),
		SecretAccessKey: os.Getenv("MINIO_SECRET_KEY"),
		UseSSL:          os.Getenv("MINIO_USE_SSL") == "true",
		BucketName:      os.Getenv("MINIO_BUCKET_NAME"),
	}

}
