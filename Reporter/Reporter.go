package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/robfig/cron"
	"github.com/segmentio/kafka-go"
	"gorm.io/gorm"
)

var db *gorm.DB

/* Constants Configuring Reporter's Functionality. */
const BROKER_ADDRESS = "95.217.125.140:29092"
const TOPIC = "test"
const GROUP_ID = "reporter_group"

type Event struct {
	gorm.Model
	EventType string `json:"EventType" gorm:"column:event_type"`
	AdID      string `json:"AdID" gorm:"column:ad_id"`
	//	AdvertiserID string    `json:"advertiser_id" gorm:"column:advertiser_id"`
	PublisherID string `json:"PublisherID" gorm:"column:publisher_id"`
	//	Credit       int       `json:"Credit" gorm:"column:credit"`
	Time int64 `json:"Time" gorm:"column:time"`
}

type AggregatedData struct {
	gorm.Model
	AdID        string `gorm:"column:ad_id"`
	Clicks      string `gorm:"column:clicks"`
	Impressions string `gorm:"column:impressions"`
	//	Credit      string       `gorm:"column:credit"`
	Time int64 `gorm:"column:time"`
}

// callInternalAPI simulates calling an internal API to handle the click
func callAPI(event Event) error {
	url := fmt.Sprintf("https://panel.lontra.tech/api/v1/ads/%s/event", event.AdID)
	payload := map[string]interface{}{
		"publisher_id": event.PublisherID,
		"event_type":   event.EventType,
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}
	fmt.Println("stat was 200")
	return nil
}

func setupKafkaReader() *kafka.Reader {

	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{BROKER_ADDRESS},
		Topic:    TOPIC,
		GroupID:  GROUP_ID,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
}

func consumeEvents(reader *kafka.Reader) {
	defer reader.Close()

	for {
		fmt.Println("Reading a new message ...")
		msg, err := reader.ReadMessage(context.Background())
		fmt.Printf("New message is read: %v\n", msg)
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

	//Added api call here instead of eventserver
	if err := callAPI(*event); err != nil {
		log.Printf("Failed to call API for an event: %v\n", err)
	}

	if err := insertEventIntoDB(event); err != nil {
		log.Printf("could not insert event into DB: %v", err)
		return
	}

	// Log successful insertion
	log.Printf("Inserted event into DB: %v, AdID: %v, EventType: %v", event.Time, event.AdID, event.EventType)

}

func insertEventIntoDB(event *Event) error {

	// GORM: Insert the event into the database
	result := db.Create(event)
	return result.Error

}

func aggregateData() {
	var results []struct {
		AdID        string
		Clicks      string
		Impressions string
		//	Credit      string
		Time int64
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
			//	Credit:      result.Credit,
			Time: result.Time,
		}
		db.Create(&aggregatedData)
	}
}

func main() {
	// GORM: Initialize the database connection
	var err error
	CreateDBIfNotExists()

	db, err = ConnectToDB()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	// GORM: Auto Migrate the schema
	err = db.AutoMigrate(&Event{})

	if err != nil {
		log.Fatalf("failed to auto migrate: %v", err)
	}
	err = db.AutoMigrate(&AggregatedData{})

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
	fmt.Println("Setup successfully!")
	consumeEvents(reader)
	//setupAndRunAPIRouter()
}
