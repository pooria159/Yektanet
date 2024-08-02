package main

import "math"

/*
This file will contain methods needed to fetch ads from panel,
together with its metadata from PostgreSQL database.

The result will be written to a `map` structure that maps a pair of (ad-id, publisher-id)
to a pair of (impression count, CTR).


IN THE FOLLOWING, A DEMO OF THE FUNCTIONALITIES
IMPLEMENTED IN THIS FILE IS DESCRIBED


* When fetching ads from panel:
	* fetch ad success statistics from metadata database
	* compute the CTR per publisher of every ad
		A SELECT statement --> map[ad-id][publisher-id] -> (impression, ctr)
		A SELECT statement  --> map[advertiser-id] -> mean CTR
	* if an ad is not displayed on this publisher yet, set CTR to the average CTR of ads of its advertiser
		for on fetched-ad, publisher: if map[ad][publisher] = null: map[ad][publisher] = (0, mean CTR of advertiser)
	* compute the revenue expected value by taking the product of CTR and bid, for each ad
		expected[ad-id][publisher-id] = bid[ad-id]  * prob[ad-id][publisher-id]
	* compute the tolerance range for each ad
		algorithm described below
	* select the ad with maximum expected value, together with ads that can 'beat' that ad as a result of tolerance.
		---
	* compute the relative weights of those selected ads
		weight[ad-id][publisher-id] = expected[ad-id][publisher-id] / NORMALIZATION_TERM
* When a publisher requests for a new ad:
	* draw a random weighted ad, encrypt and send.
		Maybe using the uniform distribution?


* Confidence interval (95%):

** METHOD ONE: ADDITION-BASED
P (d <= e) > 1 - 2exp(-2e^2N) >= 0.95
-2e^2N <= ln(0.025)
N >=  ln(40)/(2e^2) ~= 1.8444 / e^2
e^2 >= ln(40)/(2N)
e >= sqrt(ln(40) / 2)  /  sqrt(N) ~= 1.36 / sqrt(N)
upper bound: CTR + 1.36 / sqrt(N)
lower bound: CTR - 1.36 / sqrt(N)

** METHOD TWO: MULTIPLICATION-BASED
upper bound: CTR * a
lower bound: CTR / a
*/

/* Structs and Variables Relating to Ad-fetching. */

// A struct reflecting the collaboration of an ad with a publisher.
type AdPublisherCollaboration struct {
	AdID        int
	PublisherID int
}

type AdvertiserPublisherCollaboration struct {
	AdvertiserID int
	PublisherID  int
}

// Stores the statistics of a collaboration. Namely, impression count, click count and ctr.
type Statistics struct {
	Impressions int
	Clicks      int
	CTR         float64
}

// An interval in which we believe that the estimated CTR probably lies.
// This probability is a hyper-parameter that should be tuned manually.
// (A conventional value is 95%)
type ConfidenceInterval struct {
	lowerBound float64
	upperBound float64
}

// Maps collaborations to their emprical success statistics.
var evaluation map[AdPublisherCollaboration]Statistics

// Maps id of each advertiser to statistics of its ads.
var advertiserEvaluation map[AdvertiserPublisherCollaboration]Statistics

// Indicates a range containing the emprical CTR in which the real CTR will most likely lay.
var toleranceRange map[AdPublisherCollaboration]ConfidenceInterval

// What revenue is expected to gain from showing ad x to publisher y?
var expectedRevenue map[AdPublisherCollaboration]float64

// The per-publisher distribution of ads that affects the selection process.
var weight map[AdPublisherCollaboration]float64

// To some extent can the actual CTR be different, relative to the estimated CTR.
// In other words, it is assumed that actual ctr most likely lays in the interval
// [estimated_ctr / RELATIVE_TOLERACE, estimated_ctr * RELATIVE_TOLERANCE]
const RELATIVE_TOLERANCE = 2


/* Functions Used for Updating Ad Statistics */

/* Queries the metadata database and retrieves each advertiser's mean CTR per publisher. */
func fetchMeanCTRs() {
	// TODO: Update meanCtr
}

/* Queries the metadata database and computes the success statistics of each ad-publisher pair. */
func fetchAdStatistics() {
	// TODO: Update success statistics
}

/* For each publisher, sets CTR of its new ads to the mean CTR of its advertiser. */
func usePriorsForNewAds() {
	var adPubCollab AdPublisherCollaboration
	var statistics Statistics
	var exists bool
	var advertiserPubCollab AdvertiserPublisherCollaboration

	for _, publisherID := range allPublisherIDs {
		adPubCollab.PublisherID = publisherID
		for _, ad := range allFetchedAds {
			adPubCollab.AdID = ad.Id
			statistics, exists = evaluation[adPubCollab]
			if !exists || statistics.Impressions == 0 {
				statistics.Impressions = 0
				statistics.Clicks = 0
				advertiserPubCollab.AdvertiserID = ad.AdvertiserID
				advertiserPubCollab.PublisherID = publisherID
				statistics.CTR = advertiserEvaluation[advertiserPubCollab].CTR
				evaluation[adPubCollab] = statistics
			}
		}
	}
}

/* Updates the expected gain revenue for each ad-publisher pair by
 taking the product of the ad's bid value and the ad-publisher pair's
 estimated CTR. */
func updateExpectedRevenues() {
	var adPubCollab AdPublisherCollaboration

	for _, publisherID := range allPublisherIDs {
		adPubCollab.PublisherID = publisherID
		for _, ad := range allFetchedAds {
			adPubCollab.AdID = ad.Id
			expectedRevenue[adPubCollab] = float64(ad.Bid) * evaluation[adPubCollab].CTR
		}
	}
}

/* Calculates a confidence interval for each ad-publisher estimated CTR.
 The method relies in part on the Hoeffding’s inequality. */
func calculateToleranceRanges() {
	var adPubCollab AdPublisherCollaboration
	var absoluteConfInterval ConfidenceInterval
	var relativeConfInterval ConfidenceInterval
	var finalConfInterval ConfidenceInterval
	var ctr float64
	var N int

	for _, publisherID := range allPublisherIDs {
		adPubCollab.PublisherID = publisherID
		for _, ad := range allFetchedAds {
			adPubCollab.AdID = ad.Id
			ctr = evaluation[adPubCollab].CTR
			N = evaluation[adPubCollab].Impressions
			if N > 0 {
				absoluteConfInterval.upperBound = ctr + 1.36 / math.Sqrt(float64(N))
				absoluteConfInterval.lowerBound = ctr - 1.36 / math.Sqrt(float64(N))
			} else {
				absoluteConfInterval.upperBound = 1	
				absoluteConfInterval.lowerBound = 0
			}

			relativeConfInterval.upperBound = ctr * RELATIVE_TOLERANCE
			relativeConfInterval.lowerBound = ctr / RELATIVE_TOLERANCE
			finalConfInterval.upperBound = min(absoluteConfInterval.upperBound, relativeConfInterval.upperBound)
			finalConfInterval.lowerBound = max(absoluteConfInterval.lowerBound, relativeConfInterval.lowerBound)
			toleranceRange[adPubCollab] = finalConfInterval
		}
	}
}


/* Re-calculates per-publisher distributions on ads based on
the computed confidence intervals for ad-publisher pairs. */
func updatePerPublisherDistributions() {
	var adPubCollab AdPublisherCollaboration
	var maxLowerBound float64
	var winnerAdsRevenueSum float64

	for _, publisherID := range allPublisherIDs {
		adPubCollab.PublisherID = publisherID
		maxLowerBound = 0
		for _, ad := range allFetchedAds {
			adPubCollab.AdID = ad.Id
			if toleranceRange[adPubCollab].lowerBound > maxLowerBound {
				maxLowerBound =  toleranceRange[adPubCollab].lowerBound
			}
		}
		winnerAdsRevenueSum = 0
		for _, ad := range allFetchedAds {
			adPubCollab.AdID = ad.Id
			if toleranceRange[adPubCollab].upperBound >= maxLowerBound {
				winnerAdsRevenueSum += expectedRevenue[adPubCollab]
			}
		}
		for _, ad := range allFetchedAds {
			adPubCollab.AdID = ad.Id
			if toleranceRange[adPubCollab].upperBound >= maxLowerBound {
				weight[adPubCollab] = expectedRevenue[adPubCollab] / winnerAdsRevenueSum
			}
		}
	}
}
