package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

/* Constants Configuring the Functionality of Server. */

const FETCH_PERIOD = 5  // How many minutes to wait between fetching
						// Ads from Panel.


/* User-defined Types and Structs*/

type Ad struct {
	id				string	`json:"id"`
	title			string	`json:"title"`
	redirectLink	string	`json:redirect_link`
	imageSource		string	`json:image_source`
	bid				int		`json:price`
}


/* Global Objects. */

var allAds []Ad; // A slice containing all ads.


/* Functions of the Server. */

/* In an infinite loop, waits for `FETCH_PERIOD` minutes
   and then fetches ads from Panel. */
func fetchAds() {
	/* TODO: Call the API to Panel to fetch ads. */
	/* TODO: Loop over Panel's response to update ads. */
	/* TODO: Wait for `FETCH_PERIOD` minutes. */
	var dummyAd Ad
	allAds = append(allAds, dummyAd)
}

func selectAd() Ad {
	/* TODO: Choose the Ad with highest bid. */
	var dummyAd Ad
	return dummyAd
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

	router.Run("localhost:8080")
}