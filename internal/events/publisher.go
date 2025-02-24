package events

import (
	"encoding/json"
	"log"

	"Void/pkg/rabbitmq"
)

// PublishShoutEvent publishes a shout event.
// It depends solely on the ShoutEvent interface.
func PublishShoutEvent(event ShoutEvent) error {
	msg, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to marshal event: %v", err)
		return err
	}
	if err := rabbitmq.PublishNotification(msg); err != nil {
		log.Printf("Failed to publish event: %v", err)
		return err
	}
	log.Printf("Published shout event for shout ID: %d", event.GetShoutID())
	return nil
}
