package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	server := NewEventServer()
	return server.SetupRouter()
}

// TestHandleImpression tests the handleImpression handler
func TestHandleImpression(t *testing.T) {
	router := setupRouter()

	// Test case: Missing required parameters
	req, _ := http.NewRequest("POST", "/impression", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Missing required query parameters")

	// Test case: Valid request
	req, _ = http.NewRequest("POST", "/impression?user_id=user1&publisher_id=pub1&ad_id=ad1", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Impression processed")
}

// TestHandleClick tests the handleClick handler
func TestHandleClick(t *testing.T) {
	router := setupRouter()

	// Test case: Missing required parameters
	req, _ := http.NewRequest("POST", "/click", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Missing required query parameters")

	// Test case: Valid request
	req, _ = http.NewRequest("POST", "/click?user_id=user1&publisher_id=pub1&ad_id=ad1&ad_url=http://example.com", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusSeeOther, w.Code)
	assert.Equal(t, "http://example.com", w.Header().Get("Location"))
}
