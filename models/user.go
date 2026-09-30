package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:191" json:"name"`
	Email     string    `gorm:"uniqueIndex;size:191" json:"email"`
	Password  string    `gorm:"size:255" json:"-"`
	Avatar    string    `gorm:"size:500" json:"avatar"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}