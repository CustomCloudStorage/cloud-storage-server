package types

import "time"

type Registration struct {
	Email      string    `gorm:"primaryKey;column:email;size:255" json:"email"`
	Password   string    `gorm:"column:password;not null" json:"-"`
	Code       string    `gorm:"column:code;size:6;not null" json:"-"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime" json:"-"`
	LastSentAt time.Time `gorm:"column:last_sent_at;autoUpdateTime" json:"-"`
}
