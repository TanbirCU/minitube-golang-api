package models
import "time"

type Video struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `json:"user_id"`
	CategoryID  uint      `json:"category_id"`

	Title       string    `json:"title"`
	Description string    `json:"description"`

	Filename    string    `json:"filename"`
	VideoPath   string    `json:"video_path"`
	Thumbnail   string    `json:"thumbnail"`

	Views       uint64    `json:"views"`

	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	User        User      `json:"user,omitempty"`
	Category    Category  `json:"category,omitempty"`
}
