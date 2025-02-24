// session/session.go
package session

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

// Store is the global session store.
var Store *session.Store

// InitStore initializes the session store.
// Optionally, you can pass a custom session.Config.
func InitStore(config ...session.Config) {
	if len(config) > 0 {
		Store = session.New(config[0])
	} else {
		Store = session.New() // Uses the default configuration
	}
}

// GetSession retrieves the current session from the context.
func GetSession(c *fiber.Ctx) (*session.Session, error) {
	if Store == nil {
		return nil, errors.New("session store is not initialized")
	}
	sess, err := Store.Get(c)
	if err != nil {
		return nil, err
	}
	return sess, nil
}

// GetUserID retrieves the user ID from the session.
// Returns an error if the session does not have a valid user ID.
func GetUserID(c *fiber.Ctx) (uint, error) {
	sess, err := GetSession(c)
	if err != nil {
		return 0, err
	}

	userIDVal := sess.Get("user_id")
	if userIDVal == nil {
		return 0, errors.New("no user id in session")
	}

	switch v := userIDVal.(type) {
	case uint:
		return v, nil
	case float64:
		return uint(v), nil
	default:
		return 0, errors.New("invalid user id type in session")
	}
}

// SetUserID sets the user ID in the session and saves it.
func SetUserID(c *fiber.Ctx, userID uint) error {
	sess, err := GetSession(c)
	if err != nil {
		return err
	}

	sess.Set("user_id", userID)
	return sess.Save()
}

// DestroySession destroys the current session.
func DestroySession(c *fiber.Ctx) error {
	sess, err := GetSession(c)
	if err != nil {
		return err
	}
	return sess.Destroy()
}
