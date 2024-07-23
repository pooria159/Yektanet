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
<<<<<<< HEAD
	dsn := "user=postgres dbname=postgres password=Pooria1381 sslmode=disable"
=======
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := os.Getenv("DSN")
>>>>>>> 1fc507389e2e24ce35e00b2ddf9f50f78acdef95
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
