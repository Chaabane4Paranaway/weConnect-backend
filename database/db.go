package database

import (
	"fmt"
	"go-backend/models"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	// host     = os.Getenv("DB_HOST")
	// user     = os.Getenv("DB_USER")
	// password = os.Getenv("DB_PASSWORD")
	// dbname   = os.Getenv("DB_NAME")
	// port     = os.Getenv("DB_PORT")
	host     = "localhost"
	user     = "postgres"
	password = "pass"
	dbname   = "mygo_db"
	port     = "5432"
	sslmode  = "disable"
	timezone = "Africa/Ouagadougou"
)

var DB *gorm.DB

func InitDB() {
	defer func() {
		if err := DB.AutoMigrate(&models.User{}, &models.Message{}); err != nil {
			log.Fatal("💥 Failed to migrate tables: ", err)
		} else {
			fmt.Println("✅ User tables migrated successfully")
		}
		fmt.Println("✅ PostgreSQL connected successfully")
	}()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", host, user, password, dbname, port, sslmode, timezone)
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("💥 Failed to connect to database: ", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("💥 Error getting generic DB: ", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
}
