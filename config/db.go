package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	_ "github.com/lib/pq"
)

var DB *gorm.DB

func Connect() {
	dsn := "user=youruser dbname=mydb password=1234 sslmode=disable"
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
