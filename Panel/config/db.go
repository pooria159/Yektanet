package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var DB *gorm.DB

func ConnectWithLogger() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	host := os.Getenv("HOST")
	port := os.Getenv("DBPORT")
	user := os.Getenv("USER")
	password := os.Getenv("PASSWORD")
	dbname := os.Getenv("DBNAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)

	// Globally mode

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		panic("failed to connect database")
	}
	DB = db
	//return db
}

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
