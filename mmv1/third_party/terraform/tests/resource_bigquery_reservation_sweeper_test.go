package google

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// This will sweep BigqueryReservation Reservation and Assignment resources
func init() {
	resource.AddTestSweepers("BigqueryReservation", &resource.Sweeper{
		Name: "BigqueryReservation",
		F:    testSweepBigqueryReservation,
	})
}

// At the time of writing, the CI only passes us-central1 as the region
func testSweepBigqueryReservation(region string) error {
	resourceName := "BigqueryReservation"
	log.Printf("[INFO][SWEEPER_LOG] Starting sweeper for %s", resourceName)

	config, err := sharedConfigForRegion(region)
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error getting shared config for region: %s", err)
		return err
	}

	err = config.LoadAndValidate(context.Background())
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error loading: %s", err)
		return err
	}
	servicesUrl := config.BigqueryReservationBasePath + "projects/" + config.Project + "/locations/" + region + "/reservations"
	res, err := SendRequest(config, "GET", config.Project, servicesUrl, config.UserAgent, nil)
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] Error in response from request %s: %s", servicesUrl, err)
		return nil
	}

	resourceList, ok := res["reservations"]
	if !ok {
		log.Printf("[INFO][SWEEPER_LOG] Nothing found in response.")
		return nil
	}

	rl := resourceList.([]interface{})

	log.Printf("[INFO][SWEEPER_LOG] Found %d items in %s list response.", len(rl), resourceName)
	// Count items that weren't sweeped.
	nonPrefixCount := 0
	for _, ri := range rl {
		obj := ri.(map[string]interface{})
		if obj["name"] == nil {
			log.Printf("[INFO][SWEEPER_LOG] %s resource name was nil", resourceName)
			return nil
		}

		reservationName := obj["name"].(string)
		reservationNameParts := strings.Split(reservationName, "/")
		reservationShortName := reservationNameParts[len(reservationNameParts)-1]
		// Increment count and skip if resource is not sweepable.
		if !isSweepableTestResource(reservationShortName) {
			nonPrefixCount++
			continue
		}

		deleteAllAssignments(config, reservationName)

		deleteUrl := servicesUrl + "/" + reservationShortName
		// Don't wait on operations as we may have a lot to delete
		_, err = SendRequest(config, "DELETE", config.Project, deleteUrl, config.UserAgent, nil)
		if err != nil {
			log.Printf("[INFO][SWEEPER_LOG] Error deleting for url %s : %s", deleteUrl, err)
		} else {
			log.Printf("[INFO][SWEEPER_LOG] Sent delete request for %s resource: %s", resourceName, reservationShortName)
		}
	}

	if nonPrefixCount > 0 {
		log.Printf("[INFO][SWEEPER_LOG] %d items without tf-test prefix remain.", nonPrefixCount)
	}

	return nil
}

func deleteAllAssignments(config *Config, reservationName string) {
	assignmentListUrl := config.BigqueryReservationBasePath + reservationName + "/assignments"

	assignmentRes, err := SendRequest(config, "GET", config.Project, assignmentListUrl, config.UserAgent, nil)
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] Error in response from request %s: %s", assignmentListUrl, err)
		return
	}

	assignmentList, ok := assignmentRes["assignments"]
	if !ok {
		log.Printf("[INFO][SWEEPER_LOG] Nothing found in assignment response.")
		return
	}

	al := assignmentList.([]interface{})

	for _, ri := range al {
		obj := ri.(map[string]interface{})
		name := obj["name"].(string)

		deleteUrl := config.BigqueryReservationBasePath + name
		_, err = SendRequest(config, "DELETE", config.Project, deleteUrl, config.UserAgent, nil)
		if err != nil {
			log.Printf("[INFO][SWEEPER_LOG] Error deleting for url %s : %s", deleteUrl, err)
		} else {
			log.Printf("[INFO][SWEEPER_LOG] Sent delete request for bigquery reservation assignment resource: %s", name)
		}
	}
}
