package models

import (
	"errors"
	"gorm.io/gorm"
)

// Echo represents a reply (echo) to a shout.
type Echo struct {
	gorm.Model
	Content string `gorm:"not null"`
	ShoutID uint   `gorm:"not null"`
}

// Create persists the echo using the provided DB instance.
// It also performs a simple validation to ensure content is not empty.
func (e *Echo) Create(db *gorm.DB) error {
	if e.Content == "" {
		return errors.New("echo content cannot be empty")
	}
	return db.Create(e).Error
}
