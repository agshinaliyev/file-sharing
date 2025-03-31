package entity

type User struct {
	ID       uint   `gorm:"primaryKey;unique;autoIncrement"`
	Username string `gorm:"size:64"`
	Email    string `gorm:"unique;size:255"`
	Password string
}

type FileMetadata struct {
	Name      string `json:"name"`
	Size      int64  `json:"size"`
	URL       string `json:"url,omitempty"`
	CreatedAt string `json:"created_at"`
}
