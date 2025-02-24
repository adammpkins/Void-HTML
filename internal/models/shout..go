package models

import (
	"encoding/json"
	"errors"
	"log"

	"Void/pkg/rabbitmq"

	"gorm.io/gorm"
)

// Shout represents a post in Void.
type Shout struct {
	gorm.Model
	Content string `gorm:"not null"`
	UserID  uint   `gorm:"not null"`
	User    User   // Association to the user.
	Echoes  []Echo `gorm:"foreignKey:ShoutID"`
}

// ShoutCreatedEvent is the event payload for when a shout is created.
// It lives in the models package because it's closely related to the domain.
type ShoutCreatedEvent struct {
	ShoutID  uint   `json:"shout_id"`
	Content  string `json:"content"`
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
}

// ToEvent converts a Shout to a ShoutCreatedEvent.
func (s *Shout) ToEvent() ShoutCreatedEvent {
	return ShoutCreatedEvent{
		ShoutID:  s.ID,
		Content:  s.Content,
		UserID:   s.UserID,
		Username: s.User.Username,
	}
}

// The following methods implement the interface defined in the events package.
// This way, the events package can depend solely on an interface rather than a concrete type.

// GetShoutID returns the unique identifier of the shout associated with this event.
func (e ShoutCreatedEvent) GetShoutID() uint {
	return e.ShoutID
}

// GetContent retrieves the content of the shout associated with the ShoutCreatedEvent.
func (e ShoutCreatedEvent) GetContent() string {
	return e.Content
}

// GetUserID retrieves the unique identifier of the user associated with the ShoutCreatedEvent.
func (e ShoutCreatedEvent) GetUserID() uint {
	return e.UserID
}

// GetUsername retrieves the username of the user associated with the ShoutCreatedEvent.
func (e ShoutCreatedEvent) GetUsername() string {
	return e.Username
}

// Create persists the shout using the provided DB instance and publishes a notification event.
func (s *Shout) Create(db *gorm.DB) error {
	// Save the shout to the database using the injected DB.
	if err := db.Create(s).Error; err != nil {
		return err
	}

	// Load the associated user record so that s.User is populated.
	if err := db.First(&s.User, s.UserID).Error; err != nil {
		log.Printf("Failed to load user: %v", err)
		// Continue even if loading the user fails.
	}

	// Create an event payload using the shared type.
	eventPayload := ShoutCreatedEvent{
		ShoutID:  s.ID,
		Content:  s.Content,
		UserID:   s.UserID,
		Username: s.User.Username,
	}

	// Marshal the event payload to JSON.
	msg, err := json.Marshal(eventPayload)
	if err != nil {
		log.Printf("Failed to marshal event payload: %v", err)
		return err
	}

	// Publish the notification event.
	if err := rabbitmq.PublishNotification(msg); err != nil {
		log.Printf("Failed to publish notification: %v", err)
		return err
	}

	log.Printf("Published notification for shout ID: %d", s.ID)
	return nil
}

// UpdateContent updates the shout's content with some basic validation.
func (s *Shout) UpdateContent(db *gorm.DB, newContent string) error {
	if newContent == "" {
		return errors.New("new content cannot be empty")
	}
	s.Content = newContent
	return db.Save(s).Error
}

// Delete removes the shout from the database.
func (s *Shout) Delete(db *gorm.DB) error {
	return db.Delete(s).Error
}
