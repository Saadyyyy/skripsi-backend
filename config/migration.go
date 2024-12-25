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
	err2 := db.AutoMigrate(&models.Soals{})
	if err2 != nil {
		log.Fatalf("Failed to migrate Soals: %v", err2)
	}

	// Migrate Users
	err3 := db.AutoMigrate(&models.Users{})
	if err3 != nil {
		log.Fatalf("Failed to migrate Users: %v", err3)
	}

	// Migrate Rangkings
	err4 := db.AutoMigrate(&models.Rangking{})
	if err4 != nil {
		log.Fatalf("Failed to migrate Rangking: %v", err4)
	}

	// Migrate Soal Exsample
	err5 := db.AutoMigrate(&models.SoalExsample{})
	if err5 != nil {
		log.Fatalf("Failed to migrate Rangking: %v", err5)
	}
}
