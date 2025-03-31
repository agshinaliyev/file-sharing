package entity

import "time"

type User struct {
	ID       uint   `gorm:"primaryKey;unique;autoIncrement"`
	Username string `gorm:"size:64"`
	Email    string `gorm:"unique;size:255"`
	Password string
}

type Link struct {
	ID        uint   `gorm:"primaryKey"`
	URL       string `gorm:"unique"`
	FilePath  string `gorm:"uniqueIndex:idx_file_token"`
	Token     string `gorm:"uniqueIndex:idx_file_token"`
	ExpiresAt time.Time
	CreatedAt time.Time
}
