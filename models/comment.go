package models

import "time"

type Comment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	VideoID   uint      `gorm:"index" json:"video_id"`
	UserID    uint      `json:"user_id"`
	Text      string    `gorm:"type:text" json:"comment"`
	Likes     uint      `gorm:"default:0" json:"likes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	User      User      `json:"user,omitempty"`
}
