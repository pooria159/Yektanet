package main

import (
	"dummies/http"
	"fmt"
	"testing"
	"time"
)

const DUMMY_PANEL_PORT = 8080
const DUMMY_PANEL_FETCH_URL = "/api/v1/ads/active/"

/* Makes AdServer recieve a non-OK status code after fetching ads. */
func TestNonOKFetchStatus(t *testing.T) {
	http.SetDoReturn("[]", 400, nil)
	err := fetchAdsOnce()
	fmt.Println("returned error: ", err)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

/* Makes AdServer recieve an empty response after fetching ads. */
func TestEpsilonResponse(t *testing.T) {
	http.SetDoReturn("", 200, nil)
	err := fetchAdsOnce()
	fmt.Print("returned error: ", err)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestEmptySliceResponse(t *testing.T) {
	http.SetDoReturn("[]", 200, nil)
	err := fetchAdsOnce()
	if err != nil {
		t.Errorf("Unexpected error: " + err.Error())
	}
	if len(allFetchedAds) != 0 {
		t.Log(allFetchedAds)
		t.Errorf("Expected allAds to be empty")
	}
}

/* Checks if a request is really made after FETCH_PERIOD seconds. */
func TestIfFetchedInTime(t *testing.T) {
	go periodicallyFetchAds()
	http.SetDoReturn(`[{"Id":333,"Title":"strangeTitle","ImagePath":"eeps-eeps.jpg","BidValue":312,"IsActive":true,"Clicks":4,"Impressions":0,"AdvertiserID":2,"Advertiser":{"Id":0,"Name":"","Credit":0}}]`, 200, nil)
	time.Sleep(time.Second * FETCH_PERIOD * 2)
	/* Now it is expected for allFetchedAds to be updates. */
	if len(allFetchedAds) != 1 {
		t.Errorf("allAds is not updated properly.")
	}
	var theAd = allFetchedAds[0]
	if 	theAd.Id 			!= 333 ||
		theAd.Bid			!= 312 ||
		theAd.Title			!= "strangeTitle" ||
		theAd.ImageSource	!= "eeps-eeps.jpg" {
		t.Logf("Ad info: %+v", theAd)
		t.Errorf("Ad structure was not as expected")
	}
}

/* Changes the best ads and checks if the selected one really changes. */
func TestAlterBestAd(t *testing.T) {
	var ad1 = `{"Id":123,"Title":"strangeTitle","ImagePath":"eeps-eeps.jpg","BidValue":123,"IsActive":true,"Clicks":4,"Impressions":0,"AdvertiserID":2,"Advertiser":{"Id":0,"Name":"","Credit":0}}`
	var ad2 = `{"Id":321,"Title":"strangeTitle","ImagePath":"oops-oops.jpg","BidValue":321,"IsActive":true,"Clicks":4,"Impressions":0,"AdvertiserID":2,"Advertiser":{"Id":0,"Name":"","Credit":0}}`
	http.SetDoReturn("[" + ad1 + "]", 200, nil)
	err := fetchAdsOnce()
	if (err != nil) {
		t.Errorf("Unexpected error: %v", err)
	}
	if (len(allFetchedAds) != 1) {
		t.Errorf("Expected allAds to have one element, %d found.", len(allFetchedAds))
	}
	http.SetDoReturn("[" + ad1 + "," + ad2 + "]", 200, nil)
	err = fetchAdsOnce()
	if (err != nil) {
		t.Errorf("Unexpected error: %v", err)
	}
	if (len(allFetchedAds) != 2) {
		t.Errorf("Expected allAds to have two elements, %d found.", len(allFetchedAds))
	}
	selectedAD := selectAd()
	if selectedAD.Id != 321 {
		t.Errorf("Expected highest bid to be equal to 321, found %d", selectedAD.Bid)
	}
}

/* Checks if generated random token is of proper size and format.*/
func TestTokenGeneration(t *testing.T) {
	var randomToken string

	for i := 200; i < 300; i++ {
		randomToken = generateRandomToken(i)
		if (len(randomToken) != i) {
			t.Errorf("Expected random token to be of length %d, actual length was %d.", i, len(randomToken))
		}
		for j := 0; j < i; j++ {
			if randomToken[j] == '/' {
				t.Errorf("Unexpected '/' in randomToken")
			}
		}
	}
}

/* Checks if the generated random number is outside of the given range. */
func TestRandomInt(t *testing.T) {
	var randomInt int
	for i := 0; i < 200; i++ {
		randomInt = randomInRange(i, 2*i + 1)
		if randomInt < i || randomInt >= 2*i + 1 {
			t.Errorf("Randomly generated number %d is not in the range [%d, %d)", randomInt, i, 2*i + 1)
		}
	}
}
