package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

var DB *gorm.DB

func Connect() {
	//DSN=host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := os.Getenv("DSN")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}
	DB = db
}

func Ping() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("failed to get database instance from GORM: ", err)
	}
	err = sqlDB.Ping()
	if err != nil {
		log.Fatal("failed to ping database: ", err)
	}
	fmt.Println("Database connection is alive")
}

func Migrate(models ...interface{}) {
	for _, model := range models {
		err := DB.AutoMigrate(model)
		if err != nil {
			log.Fatal("failed to migrate database: ", err)
		}
	}
}
