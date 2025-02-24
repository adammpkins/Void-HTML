package rabbitmq

import (
	"github.com/streadway/amqp"
)

var (
	// Conn represents the established connection to the RabbitMQ server.
	Conn *amqp.Connection

	// Channel represents an active AMQP channel used for communication with the RabbitMQ server.
	Channel *amqp.Channel
)

// Init establishes a connection to RabbitMQ and declares the queue.
func Init(amqpURL string) error {
	var err error
	Conn, err = amqp.Dial(amqpURL)
	if err != nil {
		return err
	}
	Channel, err = Conn.Channel()
	if err != nil {
		return err
	}
	// Declare a durable queue for notifications.
	_, err = Channel.QueueDeclare(
		"shout_notifications", // name
		true,                  // durable
		false,                 // delete when unused
		false,                 // exclusive
		false,                 // no-wait
		nil,                   // arguments
	)
	return err
}

// PublishNotification sends a message to the shout_notifications queue.
func PublishNotification(message []byte) error {
	return Channel.Publish(
		"",                    // exchange
		"shout_notifications", // routing key
		false,                 // mandatory
		false,                 // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
}

// ConsumeNotifications sets up a consumer to listen on the shout_notifications queue.
func ConsumeNotifications() (<-chan amqp.Delivery, error) {
	return Channel.Consume(
		"shout_notifications", // queue
		"",                    // consumer
		false,                 // auto-ack
		false,                 // exclusive
		false,                 // no-local
		false,                 // no-wait
		nil,                   // args
	)
}
