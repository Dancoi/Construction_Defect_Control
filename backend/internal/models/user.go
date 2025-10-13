package models

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"size:255" json:"name"`
	Email        string    `gorm:"size:255;uniqueIndex" json:"email"`
	PasswordHash string    `gorm:"size:512" json:"-"`
	Role         string    `gorm:"size:50" json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
