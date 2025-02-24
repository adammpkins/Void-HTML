package models

import "gorm.io/gorm"

// Notification represents a notification for a user.
type Notification struct {
	gorm.Model
	UserID         uint   `gorm:"not null"` // The recipient's ID
	Message        string `gorm:"not null"` // The plain text message
	AuthorUsername string // The username of the shout's author
	ShoutID        uint   // Optional: the ID of the related shout
	Read           bool   `gorm:"default:false"`
}
