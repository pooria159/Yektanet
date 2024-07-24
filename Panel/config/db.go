package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"database/sql"
)

var DB *gorm.DB

func Connect() {
	host := os.Getenv("host")
	port := os.Getenv("port")
	user := os.Getenv("user")
	password := os.Getenv("password")
	dbname := os.Getenv("dbname")

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    	"password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	// dsn := os.Getenv("DSN")
	// db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	db, err := sql.Open("postgres", psqlInfo)
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
