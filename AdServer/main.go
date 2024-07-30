package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

/* Constants Configuring Functionality of the Server */

var TEST_RAW_RESPONSE = []byte(`[{"Id":1,"Title":"12","ImagePath":"uploads\\treesample.png","BidValue":12,"IsActive":true,"Clicks":0,"Impressions":0,"AdvertiserID":2,"Advertiser":{"Id":0,"Name":"","Credit":0}},{"Id":6,"Title":"144","ImagePath":"media\\treesample.png","BidValue":144,"IsActive":true,"Clicks":0,"Impressions":0,"AdvertiserID":2,"Advertiser":{"Id":0,"Name":"","Credit":0}},{"Id":11,"Title":"test","ImagePath":"media/swoled_20240722144230_2.jpg","BidValue":12,"IsActive":true,"Clicks":0,"Impressions":0,"AdvertiserID":2,"Advertiser":{"Id":0,"Name":"","Credit":0}},{"Id":10,"Title":"first","ImagePath":"media/s.jpg","BidValue":100,"IsActive":true,"Clicks":0,"Impressions":0,"AdvertiserID":2,"Advertiser":{"Id":0,"Name":"","Credit":0}}]`)

const ADSERVER_PORT = 9090	// The port on which AdServer listens.
const FETCH_PERIOD = 60    	// How many seconds to wait between fetching
							// Ads from Panel.
const FETCH_URL = "https://panel.lontra.tech/api/v1/ads/active/" // Address from which ads are to be fetched.
const EVENT_URL = "https://eventserver.lontra.tech/"             // Address to which ads are to be sent.
const API_TEMPLATE = "/api/ads/"                                 // URL that will be routed to the getNewAd handler.
const PUBLISHER_ID_RECV_PARAM = "publisherID"                    // Name of the parameter in URL received from publisher that specifies publisher's id.

const PRINT_RESPONSE = true                                            // Whether to print allAds after it is fetched.
const USER_TOKEN_SIZE = 30                                             // User token is a random token attached to the sent click and impression link.
var JWT_ENCRYPTION_KEY = []byte("Golangers:Pooria-Mohammad-Roya-Sina") // Encryption key used to sign responses.

/* User-defined Types and Structs */

/* Holds the information contained in fetched ads from Panel. */
type FetchedAd struct {
	Id           int    `json:"Id"`
	Title        string `json:"Title"`
	ImageSource  string `json:"ImagePath"`
	Bid          int    `json:"BidValue"`
	RedirectLink string `json:"RedirectLink"`
}

/* This struct will be signed by AdServer and eventually sent to Event Server. */
type EventInfo struct {
	UserID      string
	PublisherID string
	AdID        string
	AdURL       string
	EventType   string

	jwt.StandardClaims
}

/* This information gets serialized to JSON and will be sent to Publisher. */
type ResponseInfo struct {
	Title          string `json:"Title"`
	ImagePath      string `json:"ImagePath"`
	ClickLink      string `json:"ClickLink"`
	ImpressionLink string `json:"ImpressionLink"`
}


/* Global Objects */

var allFetchedAds []FetchedAd // A slice containing all ads.

/* Functions of the Server */

/*
Issues a request to Panel and obtains all available
ads as the response. Returns the first encountered
error, if any.
*/
func fetchAdsOnce() error {
	client := http.DefaultClient
	req, err := http.NewRequest("GET", FETCH_URL, nil)
	if err != nil {
		log.Println("error in making request")
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("error in doing request")
		return err
	}
	if resp.StatusCode != http.StatusOK {
		log.Println("error in Panel:")
		return errors.New("panel sent " + resp.Status)
	}
	responseByte, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("error in reading response body:")
		return err
	}

	// You can comment the next line and uncomment its following line
	// in order to mock the response of Panel.
	err = json.Unmarshal(responseByte, &allFetchedAds)
	//err := json.Unmarshal(TEST_RAW_RESPONSE, &allFetchedAds)
	if err != nil {
		log.Println("error in parsing response:")
		return err
	}

	if PRINT_RESPONSE {
		log.Printf("Successful Ad Fetch.\nallAds: %+v\n", allFetchedAds)
	}

	return nil
}


func RemoveDisabledAds(disabledAdIds []int) {
	remainingAds := []FetchedAd{}
	for _, ad := range allFetchedAds {
		shouldRemove := false
		for _, id := range disabledAdIds {
			fmt.Println(ad.Id)
			fmt.Println(id)
			if ad.Id == id {
				shouldRemove = true
				break
			}
		}
		if !shouldRemove {
			remainingAds = append(remainingAds, ad)
			fmt.Println(remainingAds)
		}
	}
	allFetchedAds = remainingAds
}

/*
In an infinite loop, calls fetchAdsOnce and
checks if any error has occured. If so, logs the error
and waits for half of normal waiting interval. If not,
waits for `FETCH_PERIOD` seconds.
*/
func periodicallyFetchAds() {
	var err error
	for {
		err = fetchAdsOnce()
		if err == nil {
			time.Sleep(FETCH_PERIOD * time.Second)
		} else {
			log.Println("error while fetching ad:", err)
			time.Sleep(FETCH_PERIOD / 2 * time.Second)
		}
	}
}

/*
Selects best ads based on AdServer's policy.
Current policy: to select ad with highest bid.
*/
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
	return a + rand.Intn(b-a)
}

/* Generate a random token of given size.
 ASCII code of each character is between 'a' and 'z' (inclusive). */
func generateRandomToken(size int) string {
	var builder strings.Builder
	builder.Reset()
	var randomChar byte

	for i := 0; i < size; i++ {
		randomChar = byte(randomInRange('a', 'z'+1)) // Select a random alphanumeric character.
		builder.WriteByte(randomChar)
	}
	return builder.String()
}

/* Generates raw event link and signs it using the 
 private key of AdServer. */
func generateSignedEventInfo(action string, selectedAd FetchedAd, requestingPublisherId int) (string, error) {
	var eventInfo EventInfo 
	eventInfo.AdID = strconv.Itoa(selectedAd.Id)
	eventInfo.PublisherID = strconv.Itoa(requestingPublisherId)
	eventInfo.UserID = generateRandomToken(USER_TOKEN_SIZE)
	eventInfo.AdURL = selectedAd.RedirectLink
	eventInfo.EventType = action

	signedInfo, err := signEvent(&eventInfo)
	if err != nil {
		return "", err
	}
	return EVENT_URL + action + "/" + signedInfo, nil
}

/* Makes a Response instance, puts info that is to be sent 
 in it and returns it. */
func makeResopnse(selectedAd FetchedAd, requestingPublisherId int) (ResponseInfo, error) {
	var response ResponseInfo
	var err error

	response.Title					= selectedAd.Title
	response.ImagePath				= selectedAd.ImageSource
	response.ClickLink, err			= generateSignedEventInfo("click", selectedAd, requestingPublisherId)	
	if err != nil {
		return response, err
	}
	response.ImpressionLink, err	= generateSignedEventInfo("impression", selectedAd, requestingPublisherId)
	if err != nil {
		return response, err
	}
	return response, nil
}

/* Signs the information Event Server needed
 with AdServer's internal private key, so that
 it will be shown that it is really generated by AdServer. */
func signEvent(event *EventInfo) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, event)
	signedTokenString, err := token.SignedString(JWT_ENCRYPTION_KEY)
	if err != nil {
		return "", err
	}
	return signedTokenString, nil
}

/* Handels GET requests from publishers requesting
 for a new ad. */
func getNewAd(c *gin.Context) {
	selectedAd := selectAd()
	publisherId, _ := strconv.Atoi(c.Query(PUBLISHER_ID_RECV_PARAM))
	response, err := makeResopnse(selectedAd, publisherId)
	
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	c.JSON(http.StatusOK, response)
}

func main() {
	/* Configure Go's predefined logger. */
	log.SetPrefix("AdServer:")
	log.SetFlags(log.Ltime | log.Ldate)

	/* Run the two main workers: ad-fetcher
	   and query-responser. */
	go periodicallyFetchAds()
	router := gin.Default()
	router.GET(API_TEMPLATE, getNewAd)

	router.Run(":" + strconv.Itoa(ADSERVER_PORT))
}
