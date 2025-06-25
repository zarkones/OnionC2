package db

import (
	"api/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var ORM *gorm.DB

// Init initializes the database.
func Init(dbName string) error {
	// Use in-memory database if dbName is empty.
	if len(dbName) == 0 {
		dbName = ":memory:"
	}

	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		return err
	}

	if err := db.AutoMigrate(
		&models.Agent{},
		&models.Message{},
		&models.File{},
		&models.Permission{},
		&models.Operator{},
	); err != nil {
		return err
	}

	ORM = db

	return nil
}
