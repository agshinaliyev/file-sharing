package service

import (
	"crypto/rand"
	"file-sharing/model/entity"
	"file-sharing/repository"
	"fmt"
	"time"
)

type SharingService interface {
	CreateShareLink(filePath string, expiry time.Duration) (string, error)
	ValidateAndGetURL(filePath, token string) (string, error)
}

type sharingService struct {
	linkRepo repository.LinkRepository
	minioSvc MinioService
}

func NewSharingService(linkRepo repository.LinkRepository, minioSvc MinioService) SharingService {
	return &sharingService{
		linkRepo: linkRepo,
		minioSvc: minioSvc,
	}
}

func (s *sharingService) CreateShareLink(filePath string, expiry time.Duration) (string, error) {
	token := generateToken()
	err := s.linkRepo.GenerateURL(&entity.Link{
		FilePath:  filePath,
		Token:     token,
		ExpiresAt: time.Now().Add(expiry),
	})
	if err != nil {
		return "", fmt.Errorf("failed to create share link: %w", err)
	}
	return fmt.Sprintf("/share/%s?token=%s", filePath, token), nil
}

func (s *sharingService) ValidateAndGetURL(filePath, token string) (string, error) {
	if !s.linkRepo.IsValid(filePath, token) {
		return "", fmt.Errorf("invalid or expired link")
	}
	return s.minioSvc.GeneratePresignedURL(filePath, 15*time.Minute)
}

func generateToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
