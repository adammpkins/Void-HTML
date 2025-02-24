package handlers

import (
	"Void/internal/db"
	"Void/internal/middleware"
	"Void/internal/models"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
)

// RegisterUserRoutes registers user-related routes to the provided fiber app instance.
func RegisterUserRoutes(app *fiber.App) {
	app.Get("/users/:username", middleware.GetUserFromSession, GetProfile)
	app.Get("/profile/edit", middleware.GetUserFromSession, ShowEditProfile)
	app.Post("/profile/edit", middleware.GetUserFromSession, UpdateProfile)
}

// GetProfile handles HTTP GET requests to retrieve a user profile based on the provided username parameter.
func GetProfile(c *fiber.Ctx) error {
	//get the user and their shouts
	uid := c.Locals("UserID").(uint)
	log.Println("Getting profile for user: ", c.Params("username"))
	username := c.Params("username")

	var user models.User
	result := db.DB.First(&user, "username = ?", username)
	if result.Error != nil {
		return c.SendString("User not found")
	}

	var shouts []models.Shout
	if err := db.DB.Where("user_id = ?", user.ID).Order("created_at desc").Find(&shouts).Error; err != nil {
		log.Printf("Error fetching shouts for user %s: %v", username, err)
	}

	//render the profile, no authentication should be required.
	return c.Render("profile", fiber.Map{
		"User":   user,
		"UserID": uid,
		"Shouts": shouts,
	}, "layouts/main")

}

// ShowEditProfile renders the edit profile form.
func ShowEditProfile(c *fiber.Ctx) error {
	uid := c.Locals("UserID").(uint)

	var user models.User
	if err := db.DB.First(&user, uid).Error; err != nil {
		return c.SendString("User not found")
	}
	return c.Render("edit_profile", fiber.Map{
		"User":   user,
		"UserID": uid,
	}, "layouts/main")
}

// UpdateProfile processes the form submission to update avatar and bio.
func UpdateProfile(c *fiber.Ctx) error {
	uid := c.Locals("UserID").(uint)
	var user models.User
	if err := db.DB.First(&user, uid).Error; err != nil {
		return c.SendString("User not found")
	}

	// Update bio
	bio := c.FormValue("bio")
	user.Bio = bio

	// Process avatar file upload if provided.
	file, err := c.FormFile("avatar")
	if err == nil { // file was uploaded
		// Save the uploaded file temporarily.
		tempPath := "./tmp/" + file.Filename
		if err := c.SaveFile(file, tempPath); err != nil {
			return c.Status(500).SendString("Error saving avatar file")
		}

		// Open the saved file.
		srcImage, err := imaging.Open(tempPath)
		if err != nil {
			os.Remove(tempPath)
			log.Printf("imaging.Open error: %v", err)
			return c.Status(500).SendString("Error opening uploaded image")
		}
		// Resize image to a fixed size (e.g., 200x200 pixels).
		resizedImage := imaging.Resize(srcImage, 200, 200, imaging.Lanczos)

		// Generate a unique filename.
		avatarFilename := fmt.Sprintf("%d_%s", uid, file.Filename)
		avatarPath := "./web/static/uploads/avatars/" + avatarFilename

		// Save the processed image.
		if err := imaging.Save(resizedImage, avatarPath); err != nil {
			os.Remove(tempPath)
			return c.Status(500).SendString("Error saving processed avatar image")
		}

		// Clean up temporary file.
		os.Remove(tempPath)

		// Update user's Avatar field.
		user.Avatar = "/static/uploads/avatars/" + avatarFilename
	}

	// Save updated user info.
	if err := db.DB.Save(&user).Error; err != nil {
		return c.Status(500).SendString("Error updating profile")
	}

	return c.Redirect("/users/" + user.Username)
}
