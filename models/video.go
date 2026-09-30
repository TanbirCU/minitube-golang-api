package models
import "time"

type Video struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `json:"user_id"`
	CategoryID  uint      `json:"category_id"`

	Title       string    `gorm:"size:500" json:"title"`
	Description string    `gorm:"size:5000" json:"description"`

	Filename    string    `gorm:"size:255" json:"filename"`
	VideoPath   string    `gorm:"size:500" json:"video_path"`
	Thumbnail   string    `gorm:"size:500" json:"thumbnail"`

	Views       uint64    `json:"views"`

	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	User        User      `json:"user,omitempty"`
	Category    Category  `json:"category,omitempty"`
}
