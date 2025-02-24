package handlers

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"Void/internal/db"
	"Void/internal/middleware"
	"Void/internal/models"
)

func RegisterNotificationRoutes(app *fiber.App) {
	// Group this route so that only logged-in users can access it.
	authGroup := app.Group("/", middleware.GetUserFromSession, middleware.RequireLogin)
	authGroup.Get("/notifications", middleware.GetUserFromSession, GetNotifications)
	authGroup.Get("/notifications/:id/read", middleware.GetUserFromSession, MarkNotificationAsRead)

}

// GetNotifications retrieves notifications for the logged-in user and renders them.
func GetNotifications(c *fiber.Ctx) error {
	// Get session and the user ID
	uid := c.Locals("UserID").(uint)
	if uid != 0 {
		// Query the notifications for this user, ordered by newest first.
		var notifications []models.Notification
		if err := db.DB.
			Where("user_id = ?", uid).
			Order("created_at desc").
			Find(&notifications).Error; err != nil {
			log.Printf("Error fetching notifications: %v", err)
			return c.Status(500).SendString("Error fetching notifications")
		}

		var count int64
		db.DB.Model(&models.Notification{}).Where("user_id = ? AND read = ?", uid, false).Count(&count)

		// Render the notifications view.
		return c.Render("notifications", fiber.Map{
			"Notifications":     notifications,
			"UserID":            uid,
			"NotificationCount": count,
		}, "layouts/main")
	}

	return c.Redirect("/login")
}

// MarkNotificationAsRead marks a specific notification as read for the currently logged-in user.
// It validates the session, ensures the notification exists, and belongs to the user before updating it.
// Redirects the user to the notifications page upon success. Returns an error or HTTP status code on failure.
func MarkNotificationAsRead(c *fiber.Ctx) error {
	// Retrieve the session and user id.
	uid := c.Locals("UserID").(uint)

	// Get the notification id from the URL.
	notifID := c.Params("id")
	var notif models.Notification
	if err := db.DB.First(&notif, notifID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Notification not found")
	}

	// Ensure the notification belongs to the current user.
	if notif.UserID != uid {
		return c.Status(fiber.StatusForbidden).SendString("Access denied")
	}

	// Mark the notification as read.
	notif.Read = true
	if err := db.DB.Save(&notif).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to mark notification as read")
	}

	// Redirect back to the notifications page.
	return c.Redirect("/notifications")
}
