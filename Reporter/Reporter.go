package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/segmentio/kafka-go"
)

var db *gorm.DB

type Event struct {
	gorm.Model
	EventType    string    `json:"type" gorm:"column:event_type"`
	AdID         string    `json:"ad_id" gorm:"column:ad_id"`
	AdvertiserID string    `json:"advertiser_id" gorm:"column:advertiser_id"`
	PublisherID  string    `json:"publisher_id" gorm:"column:publisher_id"`
	Credit       int       `json:"Credit" gorm:"column:credit"`
	Time         time.Time `json:"Time" gorm:"column:time"`
}

func setupKafkaReader() *kafka.Reader {
	brokerAddress := "localhost:9092" //temp
	topic := "events_topic"
	groupID := "reporter_group"

	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{brokerAddress},
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
}

func consumeEvents(reader *kafka.Reader) {
	defer reader.Close()

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("could not read message: %v", err)
			continue
		}
		processEvent(msg.Value)
	}
}

func processEvent(eventData []byte) {
	// Unmarshal the event data
	event := &Event{}
	if err := json.Unmarshal(eventData, event); err != nil {
		log.Printf("could not unmarshal event: %v", err)
		return
	}

	if err := insertEventIntoDB(event); err != nil {
		log.Printf("could not insert event into DB: %v", err)
		return
	}
	// Log successful insertion
	log.Printf("Inserted event into DB: %v, AdID: %v, EventType: %v", event.Time, event.AdID, event.EventType)

}

func insertEventIntoDB(event *Event) error {
	// db, err := sql.Open("postgres", "  ")
	// if err != nil {
	// 	return err
	// }
	// defer db.Close()
	// query := `
	//     INSERT INTO Events (ad_id, event_type, time, publisher_id, credit, advertiser_id)
	//     VALUES ($1, $2, $3, $4, $5, $6)
	// `
	// _, err = db.Exec(query, event.AdID, event.EventType, event.Time, event.PublisherID, event.Credit, event.AdvertiserID)
	// return err

	// GORM: Insert the event into the database
	result := db.Create(event)
	return result.Error

}

func main() {
	// GORM: Initialize the database connection
	var err error
	db, err = gorm.Open(postgres.Open(""), &gorm.Config{}) //??
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// GORM: Auto Migrate the schema
	err = db.AutoMigrate(&Event{})
	if err != nil {
		log.Fatalf("failed to auto migrate: %v", err)
	}

	reader := setupKafkaReader()
	/* Run the two main workers:
	 event-consumer and api-handler. */
	go consumeEvents(reader)
	setupAndRunAPIRouter()
}
