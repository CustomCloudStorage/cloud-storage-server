package types

import (
	"io"
	"time"
)

type File struct {
	ID           int        `json:"id" gorm:"primaryKey;column:id"`
	UserID       int        `json:"user_id" gorm:"not null;column:user_id"`
	FolderID     *int       `json:"folder_id,omitempty" gorm:"column:folder_id"`
	Name         string     `json:"name" gorm:"not null;column:name"`
	Extension    string     `json:"extension" gorm:"not null;column:extension"`
	Size         int64      `json:"size" gorm:"not null;column:size"`
	PhysicalName string     `json:"physical_name" gorm:"not null;column:physical_name"`
	CreatedAt    time.Time  `json:"created_at" gorm:"column:created_at"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty" gorm:"column:deleted_at"`
}

type DownloadedFile struct {
	Reader      io.ReadSeeker
	FileName    string
	ContentType string
	FileSize    int64
	ModTime     time.Time
}
