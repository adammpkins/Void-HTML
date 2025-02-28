package handlers

import (
	"Void/internal/events"
	"log"

	"Void/internal/db"
	"Void/internal/middleware"
	"Void/internal/models"
	"github.com/gofiber/fiber/v2"
)

// RegisterVoidRoutes configures routes for both global and authenticated functionalities for the app.
func RegisterVoidRoutes(app *fiber.App) {
	// Register these FIRST - before the auth group
	app.Get("/echo-chamber", middleware.GetUserFromSession, GetGlobalFeed)
	app.Get("/global/shout/:id", middleware.GetUserFromSession, GetGlobalShout)
	app.Post("/global/shout/:id/echo", middleware.GetUserFromSession, CreateGlobalEcho)
	// Then register the auth group
	authGroup := app.Group("/", middleware.GetUserFromSession, middleware.RequireLogin)
	authGroup.Get("/", GetShouts)
	authGroup.Post("/shout", CreateShout)
	authGroup.Get("/shout/:id", GetShout)
	authGroup.Post("/shout/:id/echo", CreateEcho)
	authGroup.Get("/shout/:id/edit", EditShoutForm)
	authGroup.Post("/shout/:id/update", UpdateShout)
	authGroup.Post("/shout/:id/delete", DeleteShout)
}

// GetShouts retrieves all shouts for the logged-in user.
func GetShouts(c *fiber.Ctx) error {
	log.Println("=== GetShouts handler started ===")

	uid := c.Locals("UserID").(uint)
	var user models.User
	if err := db.DB.First(&user, uid).Error; err != nil {
		return c.Status(404).SendString("User not found")
	}

	log.Printf("User from session: %v", user)

	log.Printf("UserID from session: %v", uid)
	var shouts []models.Shout
	result := db.DB.Preload("Echoes").Preload("User").Where("user_id = ?", uid).Find(&shouts)
	if result.Error != nil {
		log.Printf("Database error: %v", result.Error)
		return c.Status(500).SendString("Database error")
	}

	log.Printf("Found %d shouts for user %d", len(shouts), uid)
	log.Println("Attempting to render index template")

	var count int64
	db.DB.Model(&models.Notification{}).Where("user_id = ? AND read = ?", uid, false).Count(&count)
	return c.Render("index", fiber.Map{
		"Shouts":            shouts,
		"UserID":            uid,
		"User":              user,
		"NotificationCount": count,
	}, "layouts/main")
}

// CreateShout processes a request to create a new shout and persist it in the database, while managing related operations.
// Retrieves the user ID from the session and validates the input shout content.
// Saves the new shout to the database and populates its related user data for enriching event publication.
// Converts the shout into an event payload and publishes the event using the configured messaging mechanism.
// Redirects the user after successful creation or handles errors appropriately during the process.
func CreateShout(c *fiber.Ctx) error {
	uid := c.Locals("UserID").(uint)

	content := c.FormValue("content")
	if content == "" {
		return c.Redirect("/")
	}

	shout := models.Shout{
		Content: content,
		UserID:  uid,
	}

	// Save the shout to the database.
	if err := db.DB.Create(&shout).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// Load the associated user (to fill in Username, etc.)
	if err := db.DB.First(&shout.User, uid).Error; err != nil {
		log.Printf("Failed to load user: %v", err)
		// Continue even if this fails.
	}

	// Convert the shout to its event payload.
	eventPayload := shout.ToEvent()

	// Publish the event using the interface-based function.
	if err := events.PublishShoutEvent(eventPayload); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Redirect("/")
}

// GetShout retrieves a single shout.
func GetShout(c *fiber.Ctx) error {
	uid := c.Locals("UserID").(uint)

	id := c.Params("id")
	var shout models.Shout
	result := db.DB.Preload("Echoes").First(&shout, id)
	if result.Error != nil {
		return c.SendStatus(404)
	}
	if shout.UserID != uid {
		return c.Status(403).SendString("Access denied")
	}
	var count int64
	db.DB.Model(&models.Notification{}).Where("user_id = ? AND read = ?", uid, false).Count(&count)

	return c.Render("shout", fiber.Map{
		"Shout":             shout,
		"UserID":            uid,
		"NotificationCount": count,
	}, "layouts/main")
}

// CreateEcho handles the creation of echoes for shouts.
func CreateEcho(c *fiber.Ctx) error {
	uid := c.Locals("UserID").(uint)

	id := c.Params("id")
	var shout models.Shout
	if err := db.DB.First(&shout, id).Error; err != nil {
		return c.Status(404).SendString("Shout not found")
	}

	if shout.UserID != uid {
		return c.Status(403).SendString("Unauthorized")
	}

	content := c.FormValue("content")
	if content == "" {
		return c.Redirect("/shout/" + id)
	}

	echo := models.Echo{
		Content: content,
		ShoutID: shout.ID,
	}

	if err := echo.Create(db.DB); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Redirect("/shout/" + id)
}

// GetGlobalFeed retrieves all global shouts.
func GetGlobalFeed(c *fiber.Ctx) error {
	log.Println("GetGlobalFeed handler started")
	var shouts []models.Shout
	result := db.DB.Preload("Echoes").Preload("User").Order("created_at desc").Find(&shouts)
	if result.Error != nil {
		log.Printf("Database error: %v", result.Error)
		return c.Status(500).SendString("Database error")
	}
	log.Printf("Found %d global shouts", len(shouts))

	uid := c.Locals("UserID").(uint)
	if uid == 0 {
		// If no valid user is found, render with nil UserID
		return c.Render("echo_chamber", fiber.Map{
			"Shouts": shouts,
			"UserID": nil,
		}, "layouts/main")
	}

	var count int64
	db.DB.Model(&models.Notification{}).Where("user_id = ? AND read = ?", uid, false).Count(&count)

	return c.Render("echo_chamber", fiber.Map{
		"Shouts":            shouts,
		"UserID":            uid,
		"NotificationCount": count,
	}, "layouts/main")
}

// GetGlobalShout retrieves a single global shout.
func GetGlobalShout(c *fiber.Ctx) error {
	// Try to get the user ID; if there's an error or uid is zero, treat the user as anonymous.
	uid := c.Locals("UserID").(uint)

	// Get the shout ID from the URL.
	id := c.Params("id")
	var shout models.Shout
	result := db.DB.Preload("Echoes").Preload("User").First(&shout, id)
	if result.Error != nil {
		return c.SendStatus(404)
	}

	// If a valid user is logged in, fetch notification count.
	if uid != 0 {
		var count int64
		db.DB.Model(&models.Notification{}).
			Where("user_id = ? AND read = ?", uid, false).
			Count(&count)
		return c.Render("global_shout", fiber.Map{
			"Shout":             shout,
			"UserID":            uid,
			"NotificationCount": count,
		}, "layouts/main")
	}

	// Otherwise, render without user-specific data.
	return c.Render("global_shout", fiber.Map{
		"Shout":  shout,
		"UserID": nil,
	}, "layouts/main")
}

// CreateGlobalEcho handles the creation of echoes for global shouts.
func CreateGlobalEcho(c *fiber.Ctx) error {
	uid := c.Locals("UserID").(uint)

	if uid == 0 {
		return c.Redirect("/login")
	}

	id := c.Params("id")
	var shout models.Shout
	result := db.DB.First(&shout, id)
	if result.Error != nil {
		return c.SendStatus(404)
	}

	content := c.FormValue("content")
	if content == "" {
		return c.Redirect("/global/shout/" + id)
	}
	echo := models.Echo{Content: content, ShoutID: shout.ID}
	db.DB.Create(&echo)
	return c.Redirect("/global/shout/" + id)
}

// EditShoutForm renders the edit form for a given shout.
func EditShoutForm(c *fiber.Ctx) error {
	uid := c.Locals("UserID").(uint)

	shoutID := c.Params("id")
	var shout models.Shout
	if err := db.DB.First(&shout, shoutID).Error; err != nil {
		return c.Status(404).SendString("Shout not found")
	}

	// Ensure that the shout belongs to the logged-in user.
	if shout.UserID != uid {
		return c.Status(403).SendString("Unauthorized")
	}

	var count int64
	db.DB.Model(&models.Notification{}).Where("user_id = ? AND read = ?", uid, false).Count(&count)

	return c.Render("edit_shout", fiber.Map{
		"Shout":             shout,
		"UserID":            uid,
		"NotificationCount": count,
	}, "layouts/main")
}

// UpdateShout handles the update request.
func UpdateShout(c *fiber.Ctx) error {
	uid := c.Locals("UserID").(uint)

	shoutID := c.Params("id")
	var shout models.Shout
	if err := db.DB.First(&shout, shoutID).Error; err != nil {
		return c.Status(404).SendString("Shout not found")
	}

	if shout.UserID != uid {
		return c.Status(403).SendString("Unauthorized")
	}

	newContent := c.FormValue("content")
	if err := shout.UpdateContent(db.DB, newContent); err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.Redirect("/shout/" + shoutID)
}

// DeleteShout handles deleting a shout.
func DeleteShout(c *fiber.Ctx) error {
	uid := c.Locals("UserID").(uint)

	shoutID := c.Params("id")
	var shout models.Shout
	if err := db.DB.First(&shout, shoutID).Error; err != nil {
		return c.Status(404).SendString("Shout not found")
	}

	// Ensure the shout belongs to the logged-in user
	if shout.UserID != uid {
		return c.Status(403).SendString("Unauthorized")
	}

	if err := db.DB.Delete(&shout).Error; err != nil {
		return c.Status(500).SendString("Failed to delete shout")
	}

	return c.Redirect("/")
}
