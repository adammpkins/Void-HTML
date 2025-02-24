package models

import "gorm.io/gorm"

// User represents a registered user.
type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Avatar   string // Path to avatar file
	Bio      string // Optional user bio
}
