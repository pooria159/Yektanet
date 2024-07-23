package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

/* Constants Configuring Functionality of the Server. */

const FETCH_PERIOD = 5 // How many minutes to wait between fetching
// Ads from Panel.
const FETCH_URL = "http://localhost:8080/api/v1/ads/active"

const PRINT_RESPONSE = true // Whether to print allAds after it is fetched.

/* User-defined Types and Structs*/

type Ad struct {
	Id           int    `json:"Id"`
	Title        string `json:"Title"`
	ImageSource  string `json:"ImagePath"`
	Bid          int    `json:"BidValue"`
	IsActive     bool   `json:"IsActive"`
	RedirectLink string `json:"RedirectLink"`
}

/* Global Objects. */

var allAds []Ad            // A slice containing all ads.
var allAdsMutex sync.Mutex // Mutex object to syncronize working with allAds.

/* Functions of the Server. */

/*
In an infinite loop, waits for `FETCH_PERIOD` minutes

	and then fetches ads from Panel.
*/
func fetchAds() error {
	for {

		/* Sleep for FETCH_PERIOD minutes. */
		time.Sleep(1 * time.Second) // For demonstration purposes, we just wait
		// a single second instead of FETCH_PERIOD minutes.
		// TODO: Revert the waiting interval to the original FETCH_PERIOD.

		client := http.DefaultClient
		req, err := http.NewRequest("GET", FETCH_URL, nil)
		if err != nil {
			log.Print("error in making request:", err)
			continue
		}

		resp, err := client.Do(req)

		if err != nil {
			log.Print("error in doing request:", err)
			continue
		}

		responseByte, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Print("error in reading response body:", err)
			continue
		}

		allAdsMutex.Lock()
		json.Unmarshal(responseByte, &allAds)
		allAdsMutex.Unlock()

		if PRINT_RESPONSE {
			log.Printf("Successful Ad Fetch.\nallAds: %+v\n", allAds)
		}
	}
}

func selectAd() Ad {
	var bestAd Ad
	var maxBid int = 0

	allAdsMutex.Lock()
	for _, ad := range allAds {
		if ad.Bid > maxBid {
			maxBid = ad.Bid
			bestAd = ad
		}
	}
	allAdsMutex.Unlock()
	
	return bestAd
}

func getNewAd(c *gin.Context) {
	selectedAd := selectAd()
	c.IndentedJSON(http.StatusOK, selectedAd)
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
