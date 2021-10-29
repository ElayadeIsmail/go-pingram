package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"unique" json:"username" validate:"required,min=3,max=10"`
	Email     string         `gorm:"unique" json:"email" validate:"required,email"`
	Password  string         `json:"-" validate:"required,min=6,max=22"`
	Posts     []Post         `json:"posts"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}

type Post struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Text      string         `json:"text" validate:"required" form:"text"`
	ImageUrl  string         `json:"imageUrl"`
	UserID    uint           `json:"userId"`
	User      User           `json:"user"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}
