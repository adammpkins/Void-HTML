package db

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"Void/internal/models" // Updated import path
)

// DB is a global variable that holds the GORM database connection instance.
var DB *gorm.DB

// InitDB initializes the database connection using GORM and performs schema migrations for the given models.
func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("void.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	DB.AutoMigrate(&models.Shout{}, &models.Echo{}, &models.User{}, &models.Notification{})
}
