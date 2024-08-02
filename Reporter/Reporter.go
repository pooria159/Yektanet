package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/robfig/cron"
	"github.com/segmentio/kafka-go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

type AggregatedData struct {
	gorm.Model
	AdID        int       `gorm:"column:ad_id"`
	Clicks      int       `gorm:"column:clicks"`
	Impressions int       `gorm:"column:impressions"`
	Credit      int       `gorm:"column:credit"`
	Time        time.Time `gorm:"column:time"`
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

func aggregateData() {
	var results []struct {
		AdID        int
		Clicks      int
		Impressions int
		Credit      int
		Time        time.Time
	}
	// Query to aggregate data using GORM
	db.Table("events").
		Select("ad_id, " +
			"SUM(CASE WHEN event_type = 'click' THEN 1 ELSE 0 END) as clicks, " +
			"SUM(CASE WHEN event_type = 'impression' THEN 1 ELSE 0 END) as impressions, " +
			"SUM(credit) as credit, " +
			"DATE_TRUNC('hour', time) as time").
		Group("ad_id, DATE_TRUNC('hour', time)").
		Order("time").
		Scan(&results)

	// Insert aggregated data into the new table
	for _, result := range results {
		aggregatedData := AggregatedData{
			AdID:        result.AdID,
			Clicks:      result.Clicks,
			Impressions: result.Impressions,
			Credit:      result.Credit,
			Time:        result.Time,
		}
		db.Create(&aggregatedData)
	}
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

	
	// Set up and start the cron job
	c := cron.New()
	err = c.AddFunc("@hourly", aggregateData)
	if err != nil {
		log.Fatalf("failed to add cron job: %v", err)
	}
	c.Start()

	// Set up Kafka reader
	reader := setupKafkaReader()
	go consumeEvents(reader)
	setupAndRunAPIRouter()
}
