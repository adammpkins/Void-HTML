package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// RequireLogin checks if the user is logged in.
func RequireLogin(c *fiber.Ctx) error {
	uid := c.Locals("UserID").(uint)
	if uid == 0 {
		return c.Redirect("/login")
	}
	return c.Next()
}
