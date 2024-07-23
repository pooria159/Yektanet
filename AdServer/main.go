package main

import (
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

/* Constants Configuring Functionality of the Server */

var TEST_RAW_RESPONSE = []byte(`[{"Id":1,"Title":"12","ImagePath":"uploads\\treesample.png","BidValue":12,"IsActive":true,"Clicks":0,"Impressions":0,"AdvertiserID":2,"Advertiser":{"Id":0,"Name":"","Credit":0}},{"Id":6,"Title":"144","ImagePath":"media\\treesample.png","BidValue":144,"IsActive":true,"Clicks":0,"Impressions":0,"AdvertiserID":2,"Advertiser":{"Id":0,"Name":"","Credit":0}},{"Id":11,"Title":"test","ImagePath":"media/swoled_20240722144230_2.jpg","BidValue":12,"IsActive":true,"Clicks":0,"Impressions":0,"AdvertiserID":2,"Advertiser":{"Id":0,"Name":"","Credit":0}},{"Id":10,"Title":"first","ImagePath":"media/s.jpg","BidValue":100,"IsActive":true,"Clicks":0,"Impressions":0,"AdvertiserID":2,"Advertiser":{"Id":0,"Name":"","Credit":0}}]`)

const FETCH_PERIOD = 60	// How many seconds to wait between fetching
						// Ads from Panel.
const FETCH_URL = "http://localhost:8080/api/v1/ads/active/"	// Address from which ads are to be fetched.
const EVENT_URL = "http://localhost:7070/"						// Address to which ads are to be sent.
const API_TEMPLATE = "/api/ads"									// URL that will be routed to the getNewAd handler.
const PUBLISHER_ID_PARAM = "publisherID"						// Name of the parameter in URL that specifies publisher's id.

const PRINT_RESPONSE = true // Whether to print allAds after it is fetched.
const USER_TOKEN_SIZE = 30	// User token is a random token attached to the sent click and impression link.

/* User-defined Types and Structs */

type FetchedAd struct {
	Id           int    `json:"Id"`
	Title        string `json:"Title"`
	ImageSource  string `json:"ImagePath"`
	Bid          int    `json:"BidValue"`
	RedirectLink string `json:"RedirectLink"`
}

type ResponseInfo struct {
	Title			string 	`json:"Title"`
	ImagePath		string	`json:"ImagePath"`
	ClickLink		string	`json:"ClickLink"`
	ImpressionLink	string	`json:"ImpressionLink"`
}

/* Global Objects */

var allFetchedAds []FetchedAd		// A slice containing all ads.

/* Functions of the Server */

/* In an infinite loop, waits for `FETCH_PERIOD` minutes
   and then fetches ads from Panel. */
func fetchAds() error {
	for {
		
		client := http.DefaultClient
		req, err := http.NewRequest("GET", FETCH_URL, nil)
		if err != nil {
			log.Print("error in making request:", err)
			time.Sleep((FETCH_PERIOD / 2) *  time.Second)
			continue
		}

		resp, err := client.Do(req)
		
		if err != nil {
			log.Print("error in doing request:", err)
			time.Sleep((FETCH_PERIOD / 2) *  time.Second)
			continue
		}
		
		_, err = io.ReadAll(resp.Body)
		if err != nil {
			log.Print("error in reading response body:", err)
			time.Sleep((FETCH_PERIOD / 2) *  time.Second)
			continue
		}

		// Replacing the returned response with a test respone.
		// TODO: Remove the replacement of response done here.
		json.Unmarshal(TEST_RAW_RESPONSE, &allFetchedAds)
		
		if PRINT_RESPONSE {
			log.Printf("Successful Ad Fetch.\nallAds: %+v\n", allFetchedAds)
		}

		/* Sleep for FETCH_PERIOD seconds. */
		time.Sleep(FETCH_PERIOD * time.Minute)
	}
}

func selectAd() FetchedAd {
	var bestAd FetchedAd
	var maxBid int = 0

	for _, ad := range allFetchedAds {
		if ad.Bid > maxBid {
			maxBid = ad.Bid
			bestAd = ad
		}
	}
	
	return bestAd
}

/* Returns a uniformly random int
   in the interval [a, b). */
func randomInRange(a, b int) int {
	return a + rand.Intn(b - a)
}

/* Generate a random token of given size.
   ASCII code of each character is between '0' and 'z' (inclusive).
   Hence, it does not contain '/'. */
func generateRandomToken(size int) string {
	var builder strings.Builder
	builder.Reset()
	var randomChar byte

	for i := 0; i < size; i++ {
		randomChar = byte(randomInRange('0', 'z' + 1)) // Select a random alphanumeric character.
		builder.WriteByte(randomChar)
	}
	return builder.String()
}


/* Generates link to be sent in the response to publisher which
   in turn will be requested form event server.
   `action` determines the meaning of this link, by specifying
   the situation in which this link is requested from event server.
   Current values are `click` or `impression` for now. */
func generateEventServerLink(action string, selectedAd FetchedAd, requestingPublisherId int) string {
	var builder strings.Builder
	builder.Reset()
	builder.WriteString(EVENT_URL)
	builder.WriteString(action)
	builder.WriteRune('/')
	userToken := generateRandomToken(USER_TOKEN_SIZE)
	builder.WriteString(userToken)
	builder.WriteRune('/')
	builder.WriteString(strconv.Itoa(requestingPublisherId))
	builder.WriteRune('/')
	builder.WriteString(strconv.Itoa(selectedAd.Id))
	builder.WriteRune('/')
	builder.WriteString(selectedAd.RedirectLink)
	return builder.String()
}

/* Makes a Response instance, puts info that is to be sent 
   in it and returns it. */
func generateResponse (selectedAd FetchedAd, requestingPublisherId int) ResponseInfo {
	var response ResponseInfo
	response.Title			= selectedAd.Title
	response.ImagePath		= selectedAd.ImageSource
	response.ClickLink		= generateEventServerLink("click", selectedAd, requestingPublisherId)	
	response.ImpressionLink	= generateEventServerLink("impression", selectedAd, requestingPublisherId)
	return response
}

/* Handels GET requests from publishers requesting
   for a new ad. */
func getNewAd(c *gin.Context) {
	selectedAd := selectAd()
	publisherId, _ := strconv.Atoi(c.Param(PUBLISHER_ID_PARAM))
	response := generateResponse(selectedAd, publisherId)
	c.IndentedJSON(http.StatusOK, response)
}

func main() {
	/* Configure Go's predefined logger. */
	log.SetPrefix("AdServer:")
	log.SetFlags(log.Ltime | log.Ldate)

	/* Run the two main workers: ad-fetcher
	   and query-responser. */
	go fetchAds()
	router := gin.Default()
	router.GET("/new-ad", getNewAd)

	router.Run("localhost:9090")
}
