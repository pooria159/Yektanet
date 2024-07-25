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
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables.")
	}

	host := os.Getenv("HOST")
	port := os.Getenv("DBPORT")
	user := os.Getenv("USER")
	password := os.Getenv("PASSWORD")
	dbname := os.Getenv("DBNAME")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	fmt.Println(psqlInfo)

	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	DB = db
}

func CreateDB() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables.")
	}
	
	host := os.Getenv("HOST")
	port := os.Getenv("DBPORT")
	user := os.Getenv("USER")
	password := os.Getenv("PASSWORD")
	dbname := os.Getenv("DBNAME")

	dsnWithoutDB := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable",
		host, port, user, password)

	fmt.Println(dsnWithoutDB)
	db, err := gorm.Open(postgres.Open(dsnWithoutDB), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get DB instance: %v", err)
	}
	defer sqlDB.Close()

	_, err = sqlDB.Exec(fmt.Sprintf("CREATE DATABASE %s;", dbname))
	if err != nil {
		fmt.Printf("Database %s fail to create.\n", dbname)
	} else {
		fmt.Printf("Database %s created successfully.\n", dbname)
	}
}

func Ping() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance from GORM: %w", err)
	}
	err = sqlDB.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}
	fmt.Println("Database connection is alive")
	return nil
}


func Migrate(models ...interface{}) {
	for _, model := range models {
		err := DB.AutoMigrate(model)
		if err != nil {
			log.Fatal("failed to migrate database: ", err)
		}
	}
}
