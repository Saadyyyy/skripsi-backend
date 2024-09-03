package main

import (
	"bank_soal/config"
	"bank_soal/route"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := config.InitConfig()

	// Initialize GORM DB
	gormDB, err := config.InitDBPostgres(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Perform migrations
	config.DBMigration(gormDB)

	// Retrieve underlying *sql.DB from GORM
	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Fatalf("Failed to get *sql.DB instance from GORM: %v", err)
	}

	// Convert *sql.DB to *sqlx.DB
	sqlxDB := sqlx.NewDb(sqlDB, "postgres")

	// Initialize Echo instance
	e := echo.New()

	// Middleware setup
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Serve static files
	e.Static("/assets", "assets")

	// Register routes with Echo
	route.Register(sqlxDB, e)

	port := fmt.Sprintf(":%d", cfg.SERVERPORT)
	log.Printf("Starting server on port %s ", port)
	e.Logger.Fatal(e.Start(port))
}
