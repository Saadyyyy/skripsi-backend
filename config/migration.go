package config

import (
	"bank_soal/models"
	"log"

	"gorm.io/gorm"
)

func DBMigration(db *gorm.DB) {
	// Migrate Category
	err := db.AutoMigrate(&models.Category{})
	if err != nil {
		log.Fatalf("Failed to migrate Category: %v", err)
	}

	// Migrate Soals
	err = db.AutoMigrate(&models.Soals{})
	if err != nil {
		log.Fatalf("Failed to migrate Soals: %v", err)
	}

	// Migrate Users
	err = db.AutoMigrate(&models.Users{})
	if err != nil {
		log.Fatalf("Failed to migrate Users: %v", err)
	}
}
