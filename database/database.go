package database

import (
	"fmt"

	"github.com/stevenwijaya/finance-tracker/config"
	"github.com/stevenwijaya/finance-tracker/models/transactions"
	"github.com/stevenwijaya/finance-tracker/models/users"
	"github.com/stevenwijaya/finance-tracker/pkg/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	cfg := config.LoadConfig()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Info("Connected to PostgreSQL Success")

	db.AutoMigrate(&users.User{})
	db.AutoMigrate(&transactions.Transaction{})
	DB = db
}
