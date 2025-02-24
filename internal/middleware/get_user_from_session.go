package middleware

import (
	"Void/pkg/session"
	"github.com/gofiber/fiber/v2"
)

func GetUserFromSession(c *fiber.Ctx) error {
	uid, _ := session.GetUserID(c)

	c.Locals("UserID", uid)
	return c.Next()
}
