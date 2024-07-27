package main

import (
	"bytes"
	"encoding/json"
	"errors"
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

const ADSERVER_PORT = 9090 // The port on which AdServer listens.
const FETCH_PERIOD = 60    // How many seconds to wait between fetching
// Ads from Panel.
const FETCH_URL = "http://panel.lontra.tech/api/v1/ads/active/" // Address from which ads are to be fetched.
const EVENT_URL = "http://eventserver.lontra.tech/"             // Address to which ads are to be sent.
const API_TEMPLATE = "/api/ads/"                                // URL that will be routed to the getNewAd handler.
const PUBLISHER_ID_RECV_PARAM = "publisherID"                   // Name of the parameter in URL received from publisher that specifies publisher's id.

/* Parameter names of the URL sent to publisher event server. */
const PUBLISHER_ID_SEND_PARAM = "publisher_id"
const USER_ID_SEND_PARAM = "user_id"
const AD_ID_SEND_PARAM = "ad_id"
const AD_URL_SEND_PARAM = "ad_url"

const PRINT_RESPONSE = true // Whether to print allAds after it is fetched.
const USER_TOKEN_SIZE = 30  // User token is a random token attached to the sent click and impression link.

/* User-defined Types and Structs */

type FetchedAd struct {
	Id           int    `json:"Id"`
	Title        string `json:"Title"`
	ImageSource  string `json:"ImagePath"`
	Bid          int    `json:"BidValue"`
	RedirectLink string `json:"RedirectLink"`
}

type ResponseInfo struct {
	Title          string `json:"Title"`
	ImagePath      string `json:"ImagePath"`
	ClickLink      string `json:"ClickLink"`
	ImpressionLink string `json:"ImpressionLink"`
	RedirectLink   string `json:"RedirectLink"`
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

/*
Returns a uniformly random int

	in the interval [a, b).
*/
func randomInRange(a, b int) int {
	return a + rand.Intn(b-a)
}

/*
Generate a random token of given size.

	ASCII code of each character is between 'a' and 'z' (inclusive).
*/
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

/*
Generates link to be sent in the response to publisher which

	in turn will be requested form event server.
	`action` determines the meaning of this link, by specifying
	the situation in which this link is requested from event server.
	Current values are `click` or `impression` for now.
*/
func generateEventServerLink(action string, selectedAd FetchedAd, requestingPublisherId int) string {
	var builder strings.Builder
	builder.Reset()
	builder.WriteString(EVENT_URL)
	builder.WriteString(action)
	builder.WriteString("?" + USER_ID_SEND_PARAM + "=")
	userToken := generateRandomToken(USER_TOKEN_SIZE)
	builder.WriteString(userToken)
	builder.WriteString("&" + PUBLISHER_ID_SEND_PARAM + "=")
	builder.WriteString(strconv.Itoa(requestingPublisherId))
	builder.WriteString("&" + AD_ID_SEND_PARAM + "=")
	builder.WriteString(strconv.Itoa(selectedAd.Id))
	builder.WriteString("&" + AD_URL_SEND_PARAM + "=")
	builder.WriteString(selectedAd.RedirectLink)
	return builder.String()
}

/*
Makes a Response instance, puts info that is to be sent

	in it and returns it.
*/
func generateResponse(selectedAd FetchedAd, requestingPublisherId int) ResponseInfo {
	var response ResponseInfo
	response.Title = selectedAd.Title
	response.ImagePath = selectedAd.ImageSource
	response.RedirectLink = selectedAd.RedirectLink
	response.ClickLink = generateEventServerLink("click", selectedAd, requestingPublisherId)
	response.ImpressionLink = generateEventServerLink("impression", selectedAd, requestingPublisherId)
	return response
}

/*
Handels GET requests from publishers requesting

	for a new ad.
*/
func getNewAd(c *gin.Context) {
	selectedAd := selectAd()
	publisherId, _ := strconv.Atoi(c.Query(PUBLISHER_ID_RECV_PARAM))
	response := generateResponse(selectedAd, publisherId)

	/* Gin's default JSON serializer escapes the '&' character.
	   Hence, we use a costum serializer to generate the response
	   string from the instantiated `response` object. */
	var buf bytes.Buffer
	buf.Reset()
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)

	err := encoder.Encode(response)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.String(http.StatusOK, buf.String())

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
