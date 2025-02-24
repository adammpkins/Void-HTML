package notifications

import (
	"log"

	"Void/internal/db"
	"Void/internal/events"
	"Void/internal/models"
)

// SendNewShoutNotifications sends a notification to every user except the shout's author.
func SendNewShoutNotifications(event events.ShoutEvent) {
	log.Printf("Creating notifications for shout ID: %d", event.GetShoutID())
	var recipients []models.User

	if err := db.DB.
		Where("id != ?", event.GetUserID()).
		Find(&recipients).Error; err != nil {
		log.Printf("Error fetching recipients: %v", err)
		return
	}

	log.Printf("Found %d recipients", len(recipients))

	for _, user := range recipients {
		notification := models.Notification{
			UserID:         user.ID,
			Message:        truncate(event.GetContent(), 50),
			AuthorUsername: event.GetUsername(),
			ShoutID:        event.GetShoutID(),
		}
		if err := db.DB.Create(&notification).Error; err != nil {
			log.Printf("Error creating notification for user %d: %v", user.ID, err)
		} else {
			log.Printf("Created notification ID %d for user %d", notification.ID, user.ID)
		}
	}
}

// truncate shortens a string to a specified length and appends "..." if truncation occurs.
func truncate(s string, n int) string {
	if len(s) > n {
		return s[:n] + "..."
	}
	return s
}
