package handlers

import (
	"Void/internal/middleware"
	"Void/pkg/session"
	"log"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"

	"Void/internal/db"
	"Void/internal/models"
)

// RegisterAuthRoutes registers all authentication-related routes, including registration, login, and logout handlers.
func RegisterAuthRoutes(app *fiber.App) {
	app.Get("/register", middleware.GetUserFromSession, ShowRegister)
	app.Post("/register", middleware.GetUserFromSession, Register)
	app.Get("/login", middleware.GetUserFromSession, ShowLogin)
	app.Post("/login", middleware.GetUserFromSession, Login)
	app.Get("/logout", middleware.GetUserFromSession, Logout)
}

// ShowRegister renders the registration page using the "register" template and the "layouts/main" layout.
func ShowRegister(c *fiber.Ctx) error {

	return c.Render("register", fiber.Map{
		"UserID": nil,
	}, "layouts/main")
}

// Register handles user registration by creating a new user with hashed password and saving it to the database.
// It validates required input fields (username, email, password) and redirects to the login page upon success.
// Returns an error or appropriate HTTP response in case of failure or missing fields.
func Register(c *fiber.Ctx) error {
	username := c.FormValue("username")
	email := c.FormValue("email")
	password := c.FormValue("password")

	if username == "" || email == "" || password == "" {
		return c.SendString("All fields are required")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).SendString("Error processing password")
	}

	user := models.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}
	db.DB.Create(&user)

	return c.Redirect("/login")
}

// ShowLogin renders the login page using the "login" template and the "layouts/main" layout.
func ShowLogin(c *fiber.Ctx) error {
	log.Println("Rendering login template")
	return c.Render("login", fiber.Map{
		"UserID": nil,
	}, "layouts/main")
}

// Login handles user authentication by verifying their email and password and establishes a session upon success.
// It fetches the user from the database using the provided email and validates the password using bcrypt.
// On successful authentication, the user's ID is saved in the session, and the user is redirected to the homepage.
// Returns an error response in case of missing user, invalid credentials, or session handling issues.
func Login(c *fiber.Ctx) error {
	log.Println("Login handler reached")
	email := c.FormValue("email")
	password := c.FormValue("password")

	var user models.User
	result := db.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return c.SendString("User not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return c.SendString("Invalid credentials")
	}

	// Set the user ID in the session using our session package helper.
	if err := session.SetUserID(c, user.ID); err != nil {
		return c.Status(500).SendString("Failed to save session")
	}

	return c.Redirect("/")
}

// Logout terminates the user's session and redirects them to the login page. Returns an error on failure.
func Logout(c *fiber.Ctx) error {
	if err := session.DestroySession(c); err != nil {
		return c.Status(500).SendString("Failed to Log out.")
	}
	return c.Redirect("/login")
}
