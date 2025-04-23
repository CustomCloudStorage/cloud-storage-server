package types

import (
	"time"

	"github.com/google/uuid"
)

type UploadSession struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;column:id"`
	UserID     int       `json:"user_id" gorm:"not null;column:user_id"`
	FolderID   *int      `json:"folder_id,omitempty" gorm:"column:folder_id"`
	Name       string    `json:"name" gorm:"not null;column:name"`
	Extension  string    `json:"extension" gorm:"not null;column:extension"`
	TotalParts int       `json:"total_parts" gorm:"not null;column:total_parts"`
	TotalSize  int64     `json:"total_size" gorm:"not null;column:total_size"`
	CreatedAt  time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
}

type UploadPart struct {
	SessionID  uuid.UUID `json:"session_id" gorm:"type:uuid;primaryKey;column:session_id"`
	PartNumber int       `json:"part_number" gorm:"primaryKey;column:part_number"`
	Size       int64     `json:"size" gorm:"not null;column:size"`
	UploadedAt time.Time `json:"uploaded_at" gorm:"column:uploaded_at;autoCreateTime"`
}
