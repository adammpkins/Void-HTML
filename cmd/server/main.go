package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"

	"Void/internal/db"
	"Void/internal/handlers"
	"Void/internal/models"
	"Void/internal/services/notifications"
	"Void/pkg/rabbitmq"
	"Void/pkg/session"
)

// In main.go
func main() {
	db.InitDB()

	session.InitStore()

	// Initialize the RabbitMQ connection.
	if err := rabbitmq.Init("amqp://guest:guest@localhost:5672/"); err != nil {
		log.Fatalf("Failed to initialize RabbitMQ: %v", err)
	}
	log.Println("RabbitMQ connected successfully.")

	// Add this before engine initialization
	templates, err := os.ReadDir("./web/templates")
	if err != nil {
		log.Fatalf("Failed to read templates directory: %v", err)
	}
	for _, template := range templates {
		log.Printf("Found template: %s", template.Name())
	}

	engine := html.New("./web/templates", ".html")
	engine.AddFunc("formatDate", func(t time.Time) string {
		return t.Format("Jan 2, 2006 at 3:04pm")
	})
	engine.Debug(true)

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	go func() {
		msgs, err := rabbitmq.ConsumeNotifications()
		if err != nil {
			log.Fatalf("Failed to start consumer: %v", err)
		}
		for d := range msgs {
			log.Printf("Received notification message: %s", d.Body)

			var event models.ShoutCreatedEvent
			if err := json.Unmarshal(d.Body, &event); err != nil {
				log.Printf("Error unmarshalling event: %v", err)
				d.Nack(false, true) // Bad message, requeue it
				continue
			}

			// Process the event (no error check needed)
			notifications.SendNewShoutNotifications(event)

			// Successfully processed, acknowledge the message
			d.Ack(false)
		}
	}()

	app.Use(func(c *fiber.Ctx) error {
		log.Printf("Request received: %s %s", c.Method(), c.Path())
		return c.Next()
	})

	app.Static("/static", "./web/static")
	log.Println("Static routes registered")

	handlers.RegisterAuthRoutes(app)
	log.Println("Auth routes registered")

	handlers.RegisterUserRoutes(app)
	log.Println("User routes registered")

	handlers.RegisterVoidRoutes(app)
	log.Println("Void routes registered")

	handlers.RegisterNotificationRoutes(app)
	log.Println("Notification routes registered")

	app.Get("/test", func(c *fiber.Ctx) error {
		log.Println("Test route hit")
		return c.SendString("Test route working")
	})

	log.Fatal(app.Listen(":3000"))
}
