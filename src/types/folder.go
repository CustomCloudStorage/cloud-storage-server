package types

import "time"

type Folder struct {
	ID        int        `json:"id" gorm:"primaryKey;column:id"`
	UserID    int        `json:"user_id" gorm:"not null;column:user_id"`
	Name      string     `json:"name" gorm:"not null;column:name"`
	ParentID  *int       `json:"parent_id,omitempty" gorm:"column:parent_id"`
	CreatedAt time.Time  `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"column:deleted_at"`
}
