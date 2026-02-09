package db

import (
	"fmt"
	"log"
	"time"

	"github.com/slickip/Subscription-service/internal/config"
	"github.com/slickip/Subscription-service/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
		cfg.TimeZone,
	)

	var db *gorm.DB
	var err error

	for i := 0; i < 5; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("Database not ready yet (attempt %d/5): %v", i+1, err)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		log.Fatalf("failed to connect database after retries: %v", err)
	}

	if err := db.AutoMigrate(
		&models.Subscription{},
	); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	log.Println("Database connected and migrated successfully")
	return db
}
