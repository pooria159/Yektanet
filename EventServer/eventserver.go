//package eventserver

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/zsais/go-gin-prometheus"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// CONSTS
var JWT_ENCRYPTION_KEY = []byte("Golangers:Pooria-Mohammad-Roya-Sina") // Encryption key used to sign responses.
const (
	recaptchaSecret = "6LfwfxsqAAAAAOjEjdTLn64TaPePPYRtIzTDVmDI"
)
const requestThreshold = 2
const timeframe = 60 // in

var blacklistedUserAgents = []string{
	"Python",                // Python scripts
	"curl",                  // cURL
	"Postman",               // Postman API client
	"HttpClient",            // Generic HTTP client
	"Java",                  // Java clients
	"Go-http-client",        // Go's default HTTP client
	"Wget",                  // Wget utility
	"php",                   // PHP scripts
	"Ruby",                  // Ruby scripts
	"Node.js",               // Node.js scripts
	"BinGet",                // BinGet utility
	"libwww-perl",           // Perl library
	"Microsoft URL Control", // Microsoft URL Control tool
	"Peach",                 // Peach fuzzing tool
	"pxyscand",              // Proxy scanner
	"PycURL",                // Python binding to libcurl
	"Python-urllib",         // Python urllib library
}

//MODELS

type RequestData struct {
	Count         int
	LastRequestAt time.Time
}

var requestLog = make(map[string]RequestData)

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

type RecaptchaResponse struct {
	Success     bool     `json:"success"`
	ChallengeTs string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	ErrorCodes  []string `json:"error-codes"`
}

//CONTROLLERS

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
func UserAgentBlacklist() gin.HandlerFunc {
	return func(c *gin.Context) {
		userAgent := c.GetHeader("User-Agent")
		if isBlacklisted(userAgent) {
			// Block the request
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
	// Extract event information from url
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

		//		s.callAPI(event)
		if err := s.callAPI(event); err != nil {
			log.Printf("Failed to call API for impression event: %v\n", err)
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "Impression processed"})
}
func (s *EventServer) captchaPage(c *gin.Context) {
	eventInfoToken := c.Query("info")
	c.HTML(http.StatusOK, "captcha.html", gin.H{
		"eventInfoToken": eventInfoToken,
	})
}
func validateCaptcha(c *gin.Context) bool {
	recaptchaResponse := c.PostForm("g-recaptcha-response")
	if recaptchaResponse == "" {
		return false
	}

	resp, err := http.PostForm("https://www.google.com/recaptcha/api/siteverify",
		url.Values{"secret": {recaptchaSecret}, "response": {recaptchaResponse}})
	if err != nil {
		log.Printf("Failed to verify reCAPTCHA: %v", err)
		return false
	}
	defer resp.Body.Close()

	var result RecaptchaResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Failed to parse reCAPTCHA response: %v", err)
		return false
	}

	return result.Success
}

func (s *EventServer) verifyCaptcha(c *gin.Context) {
	eventInfoToken := c.PostForm("info")

	if !validateCaptcha(c) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "CAPTCHA validation failed"})
		return
	}

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
			log.Printf("Failed to call API for impression event: %v\n", err)
		}
	}
	c.Redirect(http.StatusSeeOther, event.AdURL)

}

// handleClick handles the click events
func (s *EventServer) handleClick(c *gin.Context) {
	// Extract event information from url
	eventInfoToken := c.Param("info")
	var event Event
	parsedToken, err := jwt.ParseWithClaims(eventInfoToken, &event, func(t *jwt.Token) (interface{}, error) {
		return JWT_ENCRYPTION_KEY, nil
	})
	if err != nil || !parsedToken.Valid {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid click token"})
	}

	if claims, ok := parsedToken.Claims.(*Event); ok {
		issuedAt := claims.IssuedAt
		currentTime := time.Now().Unix()
		secondsElapsed := currentTime - issuedAt

		const tokenValidityThreshold = 4
		if secondsElapsed < tokenValidityThreshold {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Click Time Invalid"})
			return
		}
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid claims"})
		return
	}

	clientIP := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	key := clientIP + "_" + userAgent

	currentTime := time.Now()
	requestData, exists := requestLog[key]
	if exists {
		if currentTime.Sub(requestData.LastRequestAt).Seconds() <= timeframe {
			requestData.Count++
		} else {
			requestData.Count = 1
		}
		requestData.LastRequestAt = currentTime
	} else {
		requestData = RequestData{
			Count:         1,
			LastRequestAt: currentTime,
		}
	}
	requestLog[key] = requestData
	if requestData.Count > requestThreshold {
		c.Redirect(http.StatusSeeOther, "/captcha?info="+eventInfoToken)
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
	p := ginprometheus.NewPrometheus("eventserver")
	p.Use(router)
	router.Use(UserAgentBlacklist())
	router.Use(CORSMiddleware())

	router.LoadHTMLFiles("captcha.html")

	router.GET("/captcha", s.captchaPage)
	router.POST("/verify-captcha", s.verifyCaptcha)
	router.GET("/impression/:info", s.handleImpression)
	router.GET("/click/:info", s.handleClick)

	return router
}

//main

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
