package models

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"unique" json:"username" validate:"required,min=3,max=10"`
	Email     string    `gorm:"unique" json:"email" validate:"required,email"`
	Password  string    `json:"-" validate:"required,min=6,max=22"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
