package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToDB() (*gorm.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables.")
	}

	host := os.Getenv("HOST")
	port := os.Getenv("DBPORT")
	// TODO: Fix the overwriting of environment variable 'USER'
	//user := os.Getenv("USER")
	user := "postgres"
	password := os.Getenv("PASSWORD")
	dbname := os.Getenv("DBNAME")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	fmt.Println(psqlInfo)

	return gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
}

func CreateDBIfNotExists() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables.")
	}

	host := os.Getenv("HOST")
	port := os.Getenv("DBPORT")
	// TODO: Fix the overwriting of environment variable 'USER'
	//user := os.Getenv("USER")
	user := "postgres"
	password := os.Getenv("PASSWORD")
	dbname := os.Getenv("DBNAME")

	fmt.Println(host)
	fmt.Println(port)
	fmt.Println(user)
	fmt.Println(password)
	fmt.Println(dbname)

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
