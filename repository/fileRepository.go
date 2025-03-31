package repository

import (
	"file-sharing/model/entity"
	"gorm.io/gorm"
	"time"
)

type LinkRepository interface {
	GenerateURL(link *entity.Link) error
	IsValid(filePath, token string) bool
	CleanExpired() error
}

type linkRepository struct {
	db *gorm.DB
}

func NewLinkRepository(db *gorm.DB) LinkRepository {
	return &linkRepository{db: db}
}

func (r *linkRepository) GenerateURL(link *entity.Link) error {
	return r.db.Create(link).Error
}

func (r *linkRepository) IsValid(filePath, token string) bool {
	var link entity.Link
	err := r.db.
		Where("file_path = ? AND token = ? AND expires_at > ?",
			filePath, token, time.Now()).
		First(&link).
		Error
	return err == nil
}

func (r *linkRepository) CleanExpired() error {
	return r.db.
		Where("expires_at <= ?", time.Now()).
		Delete(&entity.Link{}).
		Error
}
