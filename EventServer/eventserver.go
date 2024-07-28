//package eventserver

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/cors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Event represents an event with user, publisher, ad IDs and URL
type Event struct {
	UserID      string
	PublisherID string
	AdID        string
	AdURL       string
	EventType   string
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
}

// NewEventServer creates a new EventServer with initialized maps and channel
func NewEventServer() *EventServer {
	return &EventServer{
		impressions:    make(map[string]Value),
		clicks:         make(map[string]Value),
		clickchan:      make(chan Event, 100), // Buffer size of 100
		impressionchan: make(chan Event, 100), // Buffer size of 100
	}
}

// handleImpression handles the impression events
func (s *EventServer) handleImpression(c *gin.Context) {
	userID := c.Query("user_id")
	publisherID := c.Query("publisher_id")
	adID := c.Query("ad_id")

	if userID == "" || publisherID == "" || adID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required query parameters"})
		return
	}

	event := Event{
		UserID:      userID,
		PublisherID: publisherID,
		AdID:        adID,
		EventType:   "impression",
	}

	if _, ok := s.impressions[event.UserID]; !ok {
		value := Value{
			AdID:        event.AdID,
			PublisherID: event.PublisherID,
		}
		s.impressions[event.UserID] = value
		s.impressionchan <- event

		//		s.callAPI(event)
		if err := s.callAPI(event); err != nil {
			log.Printf("Failed to call API for impression event: %v\n", err)
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "Impression processed"})
}

// handleClick handles the click events
func (s *EventServer) handleClick(c *gin.Context) {
	// Retrieve query parameters
	userID := c.Query("user_id")
	publisherID := c.Query("publisher_id")
	adID := c.Query("ad_id")
	adURL := c.Query("ad_url")

	// Check for missing parameters
	if userID == "" || publisherID == "" || adID == "" || adURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required query parameters"})
		return
	}

	event := Event{
		UserID:      userID,
		PublisherID: publisherID,
		AdID:        adID,
		AdURL:       adURL,
		EventType:   "click",
	}

	if _, ok := s.clicks[event.UserID]; !ok {
		value := Value{
			AdID:        event.AdID,
			PublisherID: event.PublisherID,
		}
		s.clicks[event.UserID] = value
		s.clickchan <- event

		if err := s.callAPI(event); err != nil {
			log.Printf("Failed to call API for impression event: %v\n", err)
		}
	}

	// Redirect to the ad URL
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

func (s *EventServer) processEvents() {
	//TODO Process events e.g., send to Kafka
	// event := <-s.impressionchan
	// go s.sendToKafka(event, "impression")
}

func (s *EventServer) sendToKafka(event Event, eventType string) {
	//TODO
}

// SetupRouter sets up the routes for the EventServer
func (s *EventServer) SetupRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/impression", s.handleImpression)
	router.POST("/click", s.handleClick)
	return router
}

//main

func main() {
	server := NewEventServer()
	router := server.SetupRouter()
	router.Use(cors.Default())

	// Start processing events
	go server.processEvents()

	err := router.Run(":8081")
	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}
