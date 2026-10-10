package bigquery

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/hashicorp/terraform-provider-google/google/sweeper"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func init() {
	sweeper.AddTestSweepersLegacy("BigqueryDataset", testSweepBigqueryDataset)
}

// maxResults for datasets.list. The API is free to return a smaller page than
// asked for, so the caller pages until nextPageToken is empty regardless.
const bigqueryDatasetListPageSize = 1000

// This sweeper is handwritten rather than generated because a generated sweeper
// for Dataset cannot work, for two independent reasons:
//
//  1. datasets.list items have no top-level "name" field; the id lives at
//     datasetReference.datasetId. The generated sweeper falls back to obj["id"]
//     (of the form "my-project:tf_test_foo") and runs it through
//     GetResourceNameFromSelfLink, which only splits on "/" -- so the prefix
//     check would be false for every dataset. identifier_field is a flat map
//     lookup and cannot reach a nested field.
//  2. Dataset's delete_url carries a query string
//     ("...datasets/{{dataset_id}}?deleteContents={{delete_contents_on_destroy}}").
//     The generated sweeper renders the delete URL and then appends the resource
//     name to the end of it, which would produce ".../datasets/?deleteContents=<id>".
//
// Dataset.yaml therefore keeps exclude_sweeper: true.
//
// At the time of writing, the CI only passes us-central1 as the region. Datasets
// are listed per-project rather than per-region, so this sweeper deliberately
// sweeps datasets in every location; filtering on the swept region would leave
// datasets outside us-central1 behind forever.
func testSweepBigqueryDataset(region string) error {
	resourceName := "BigqueryDataset"
	log.Printf("[INFO][SWEEPER_LOG] Starting sweeper for %s", resourceName)

	config, err := sweeper.SharedConfigForRegion(region)
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error getting shared config for region: %s", err)
		return err
	}

	if err := config.LoadAndValidate(context.Background()); err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error loading: %s", err)
		return err
	}

	// SharedConfigForRegion populates config.Project from
	// envvar.GetTestProjectFromEnv() and fails if it is unset, so there is no
	// need to read the environment again here.
	project := config.Project

	datasets, err := listBigqueryDatasets(config, project)
	if err != nil {
		// Sweep whatever was collected before the failure rather than dropping
		// the whole run; a later page failing shouldn't strand earlier pages.
		log.Printf("[INFO][SWEEPER_LOG] Error listing datasets in %s: %s", project, err)
		if len(datasets) == 0 {
			return nil
		}
	}

	log.Printf("[INFO][SWEEPER_LOG] Found %d items in %s list response.", len(datasets), resourceName)
	// Count items that weren't sweeped.
	nonPrefixCount := 0
	for _, obj := range datasets {
		ref, ok := obj["datasetReference"].(map[string]interface{})
		if !ok {
			log.Printf("[INFO][SWEEPER_LOG] %s resource had no datasetReference, skipping", resourceName)
			continue
		}
		datasetId, _ := ref["datasetId"].(string)
		if datasetId == "" {
			log.Printf("[INFO][SWEEPER_LOG] %s resource datasetId was empty, skipping", resourceName)
			continue
		}

		// Skip resources that shouldn't be sweeped. Dataset ids may only contain
		// letters, numbers and underscores, so in practice this matches the
		// "tf_test" prefix, but IsSweepableTestResource covers every prefix the
		// test suite is allowed to use.
		if !sweeper.IsSweepableTestResource(datasetId) {
			nonPrefixCount++
			continue
		}

		datasetProject, _ := ref["projectId"].(string)
		if datasetProject == "" {
			datasetProject = project
		}

		// deleteContents=true is required: a dataset that still holds tables,
		// views or routines returns a 400 otherwise, and test datasets nearly
		// always hold something. Everything reached here already passed the test
		// prefix check, so its contents are test contents too.
		deleteUrl := fmt.Sprintf("%sprojects/%s/datasets/%s?deleteContents=true", transport_tpg.BaseUrl(Product, config), datasetProject, datasetId)
		_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "DELETE",
			Project:   datasetProject,
			RawURL:    deleteUrl,
			UserAgent: config.UserAgent,
		})
		if err != nil {
			log.Printf("[INFO][SWEEPER_LOG] Error deleting for url %s : %s", deleteUrl, err)
		} else {
			log.Printf("[INFO][SWEEPER_LOG] Deleted a %s resource: %s", resourceName, datasetId)
		}
	}

	if nonPrefixCount > 0 {
		log.Printf("[INFO][SWEEPER_LOG] %d items were non-sweepable and skipped.", nonPrefixCount)
	}

	return nil
}

// listBigqueryDatasets pages through datasets.list for a project. Any datasets
// collected before an error are returned alongside it.
func listBigqueryDatasets(config *transport_tpg.Config, project string) ([]map[string]interface{}, error) {
	var out []map[string]interface{}
	pageToken := ""

	for {
		// all=true is deliberately not set. It only adds hidden datasets, i.e.
		// the "_"-prefixed ones BigQuery creates for itself (_scripts, session
		// and anonymous query result datasets). Those are never sweepable by
		// name and are managed by the service.
		listUrl := fmt.Sprintf("%sprojects/%s/datasets?maxResults=%d", transport_tpg.BaseUrl(Product, config), project, bigqueryDatasetListPageSize)
		if pageToken != "" {
			listUrl += "&pageToken=" + url.QueryEscape(pageToken)
		}

		res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			Project:   project,
			RawURL:    listUrl,
			UserAgent: config.UserAgent,
		})
		if err != nil {
			return out, err
		}

		if items, ok := res["datasets"].([]interface{}); ok {
			for _, item := range items {
				if obj, ok := item.(map[string]interface{}); ok {
					out = append(out, obj)
				}
			}
		}

		next, _ := res["nextPageToken"].(string)
		if next == "" {
			return out, nil
		}
		// Guard, not redundancy: if the API ever hands back the token we just
		// sent, this loop would spin forever in an unattended nightly job and
		// produce no signal. Stopping early only costs us part of one sweep.
		if next == pageToken {
			log.Printf("[INFO][SWEEPER_LOG] pageToken did not advance; stopping dataset listing to avoid a loop")
			return out, nil
		}
		pageToken = next
	}
}
