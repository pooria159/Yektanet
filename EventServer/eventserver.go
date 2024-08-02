package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
)

// Constants
var JWT_ENCRYPTION_KEY = []byte("Golangers:Pooria-Mohammad-Roya-Sina") // Encryption key used to sign responses.
const kafkaBrokerAddress = "localhost:9092"
const kafkaTopic = "events_topic"

// Event represents an event with user, publisher, ad IDs and URL
type Event struct {
	UserID      string
	PublisherID string
	AdID        string
	AdURL       string
	EventType   string

	jwt.StandardClaims
}

// Value represents the value stored in impression and clicks
type Value struct {
	PublisherID string
	AdID        string
}

// EventServer holds the channels for buffering events and maps for deduplication
type EventServer struct {
	impressions    map[string]Value
	clicks         map[string]Value
	clickchan      chan Event
	impressionchan chan Event
	kafkaWriter    *kafka.Writer // Kafka writer
}

// NewEventServer creates a new EventServer with initialized maps and channel
func NewEventServer() *EventServer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{kafkaBrokerAddress},
		Topic:    kafkaTopic,
		Balancer: &kafka.LeastBytes{},
	})

	return &EventServer{
		impressions:    make(map[string]Value),
		clicks:         make(map[string]Value),
		clickchan:      make(chan Event, 100), // Buffer size of 100
		impressionchan: make(chan Event, 100), // Buffer size of 100
		kafkaWriter:    writer,
	}
}

var blacklistedUserAgents = []string{
	"Python", "curl", "Postman", "HttpClient", "Java", "Go-http-client",
	"Wget", "php", "Ruby", "Node.js", "BinGet", "libwww-perl",
	"Microsoft URL Control", "Peach", "pxyscand", "PycURL", "Python-urllib",
}

func UserAgentBlacklist() gin.HandlerFunc {
	return func(c *gin.Context) {
		userAgent := c.GetHeader("User-Agent")
		if isBlacklisted(userAgent) {
			c.JSON(http.StatusForbidden, gin.H{"message": "Blocked: Disallowed User-Agent"})
			c.Abort()
		} else {
			c.Next() // Continue to the next handler
		}
	}
}

func isBlacklisted(userAgent string) bool {
	for _, ua := range blacklistedUserAgents {
		if strings.Contains(userAgent, ua) {
			return true
		}
	}
	return false
}

// handleImpression handles the impression events
func (s *EventServer) handleImpression(c *gin.Context) {
	eventInfoToken := c.Param("info")
	var event Event
	parsedToken, err := jwt.ParseWithClaims(eventInfoToken, &event, func(t *jwt.Token) (interface{}, error) {
		return JWT_ENCRYPTION_KEY, nil
	})
	if err != nil || !parsedToken.Valid {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid impression token"})
	}

	if _, ok := s.impressions[event.UserID]; !ok {
		value := Value{
			AdID:        event.AdID,
			PublisherID: event.PublisherID,
		}
		s.impressions[event.UserID] = value
		s.impressionchan <- event

		if err := s.callAPI(event); err != nil {
			log.Printf("Failed to call API for impression event: %v\n", err)
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "Impression processed"})
}

// handleClick handles the click events
func (s *EventServer) handleClick(c *gin.Context) {
	eventInfoToken := c.Param("info")
	var event Event
	parsedToken, err := jwt.ParseWithClaims(eventInfoToken, &event, func(t *jwt.Token) (interface{}, error) {
		return JWT_ENCRYPTION_KEY, nil
	})
	if err != nil || !parsedToken.Valid {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid click token"})
		return
	}

	if _, ok := s.clicks[event.UserID]; !ok {
		value := Value{
			AdID:        event.AdID,
			PublisherID: event.PublisherID,
		}
		s.clicks[event.UserID] = value
		s.clickchan <- event

		if err := s.callAPI(event); err != nil {
			log.Printf("Failed to call API for click event: %v\n", err)
		}
	}

	c.Redirect(http.StatusSeeOther, event.AdURL)
}

// callInternalAPI simulates calling an internal API to handle the click
func (s *EventServer) callAPI(event Event) error {
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

// processEvents processes events and sends them to Kafka
func (s *EventServer) processEvents() {
	for {
		select {
		case event := <-s.impressionchan:
			s.sendToKafka(event, "impression")
		case event := <-s.clickchan:
			s.sendToKafka(event, "click")
		}
	}
}

// sendToKafka sends an event to Kafka
func (s *EventServer) sendToKafka(event Event, eventType string) {
	eventData, err := json.Marshal(event)
	if err != nil {
		log.Printf("could not marshal event: %v", err)
		return
	}

	msg := kafka.Message{
		Key:   []byte(event.AdID),
		Value: eventData,
	}

	err = s.kafkaWriter.WriteMessages(context.Background(), msg)
	if err != nil {
		log.Printf("could not send event to Kafka: %v", err)
	} else {
		log.Printf("Sent %s event to Kafka: %s", eventType, event.AdID)
	}
}

// SetupRouter sets up the routes for the EventServer
func (s *EventServer) SetupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(UserAgentBlacklist())

	router.GET("/impression/:info", s.handleImpression)
	router.GET("/click/:info", s.handleClick)
	return router
}

func main() {
	server := NewEventServer()
	router := server.SetupRouter()

	// Start processing events
	go server.processEvents()

	err := router.Run(":8081")
	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}
