package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type File struct {
	ID         uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID     uint           `gorm:"not null;index"`
	User       User           `gorm:"foreignKey:UserID" json:"-"`
	FileName   string         `gorm:"type:varchar(255);not null" json:"file_name"`
	FilePath   string         `gorm:"type:varchar(255);not null" json:"file_path"`
	UrlPath    string         `gorm:"type:varchar(512);not null" json:"url_path"`
	FileType   string         `gorm:"type:varchar(100);not null" json:"file_type"`
	FileSize   int64          `gorm:"not null" json:"file_size"`
	UploadedAt time.Time      `gorm:"autoCreateTime" json:"uploaded_at"`
	IsDeleted  bool           `gorm:"default:false" json:"is_deleted"`
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"created_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (file *File) TableName() string {
	return "files"
}

func (f *File) BeforeCreate(tx *gorm.DB) error {
	f.ID = uuid.New()
	return nil
}

func (f *File) ToJson() ([]byte, error) {
	return json.Marshal(f)
}
