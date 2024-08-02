package config

import (
	"bank_soal/models"
	"log"

	"gorm.io/gorm"
)

func DBMigration(db *gorm.DB) {
	// Lakukan migrasi pada model
	err := db.AutoMigrate(
		&models.Users{},
		&models.Soals{},
		&models.Category{},
	)

	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}
