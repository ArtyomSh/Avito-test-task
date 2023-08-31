package postgres

import (
	"AvitoTesting/internal/config"
	"AvitoTesting/pkg/client/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init(cfg *config.Config) (*gorm.DB, error) {

	userName := cfg.Postrgres.Username
	password := cfg.Postrgres.Password
	database := cfg.Postrgres.Database
	host := cfg.Postrgres.Host
	port := cfg.Postrgres.Port

	dbUri := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, userName, password, database, port)

	db, err := gorm.Open(postgres.Open(dbUri), &gorm.Config{})

	db.AutoMigrate(&models.Segment{})

	return db, err
}
