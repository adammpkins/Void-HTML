package models

import "gorm.io/gorm"

// Notification represents a notification for a user.
type Notification struct {
	gorm.Model
	UserID         uint   `gorm:"not null"` // The recipient's ID
	Message        string `gorm:"not null"` // The plain text message
	AuthorUsername string // The username of the shout's author
	AuthorAvatar   string // New field for the avatar URL
	ShoutID        uint   // Optional: the ID of the related shout
	Read           bool   `gorm:"default:false"`
}
