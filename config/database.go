package config

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDBPostgres(cfg *AppConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Jakarta",
		cfg.DBHOST, cfg.DBUSERNAME, cfg.DBPASSWORD, cfg.DBNAME, cfg.DBPORT)
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		panic(err)
	}
	// Membuat instance router Echo
	router := echo.New()

	// Middleware untuk log
	router.Use(middleware.Logger())

	return db, nil

}
