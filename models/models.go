package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique" json:"username" validate:"required,min=3,max=10"`
	Email    string `gorm:"unique" json:"email" validate:"required,email"`
	Password string `json:"-" validate:"required,min=6,max=22"`
	Posts    []Post `json:"posts"`
}

type Post struct {
	gorm.Model
	Text     string `json:"text" validate:"required" `
	ImageUrl string `json:"imageUrl"`
	UserID   uint   `json:"userId"`
	User     User   `json:"user"`
}
