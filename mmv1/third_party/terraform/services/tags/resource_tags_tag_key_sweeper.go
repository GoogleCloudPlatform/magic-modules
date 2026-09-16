package tags

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/sweeper"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func init() {
	sweeper.AddTestSweepersLegacy("TagKey", testSweepTagKey)
}

const tagsBaseUrl = "https://cloudresourcemanager.googleapis.com/v3/"

// This sweeper is handwritten rather than generated because the Resource Manager
// v3 tag APIs are listed by query parameter: tagKeys.list and tagValues.list both
// require ?parent=..., and the generated sweeper truncates the list URL at the
// first "?" (see mmv1/templates/terraform/sweeper_file.go.tmpl). A generated
// sweeper for these resources would issue an unscoped GET and fail. TagKey.yaml,
// TagValue.yaml and TagBinding.yaml therefore keep exclude_sweeper: true.
//
// One sweeper covers both TagKeys and TagValues because a TagKey cannot be
// deleted while it still has TagValues, and the sweeper framework's dependency
// ordering only applies to generated sweepers.
func testSweepTagKey(region string) error {
	resourceName := "TagKey"
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

	// Tags are not project-scoped in the usual sense: a TagKey hangs off either an
	// organization or a project, and the list call has to name that parent
	// explicitly. Sweep both of the parents our tests actually create keys under.
	var parents []string
	if org := envvar.UnsafeGetTestOrgFromEnv(); org != "" {
		parents = append(parents, "organizations/"+org)
	} else {
		log.Printf("[INFO][SWEEPER_LOG] no org set in %v, skipping organization-parented tag keys", envvar.OrgEnvVars)
	}
	if project := envvar.GetTestProjectFromEnv(); project != "" {
		parents = append(parents, "projects/"+project)
	}
	if len(parents) == 0 {
		log.Printf("[INFO][SWEEPER_LOG] no tag parents to sweep")
		return nil
	}

	nonPrefixCount := 0
	for _, parent := range parents {
		tagKeys, err := listTagResources(config, "tagKeys", parent)
		if err != nil {
			log.Printf("[INFO][SWEEPER_LOG] error listing tag keys under %s: %s", parent, err)
			continue
		}
		log.Printf("[INFO][SWEEPER_LOG] Found %d tag keys under %s", len(tagKeys), parent)

		for _, tagKey := range tagKeys {
			keyName, _ := tagKey["name"].(string)
			keyShortName, _ := tagKey["shortName"].(string)
			if keyName == "" {
				log.Printf("[INFO][SWEEPER_LOG] tag key with no name field, skipping")
				continue
			}

			// Reclaim sweepable TagValues even when the key itself must be kept --
			// tests create tf-test values under the long-lived bootstrapped key.
			keySweepable := sweeper.IsSweepableTestResource(keyShortName)
			deletedValues, remainingValues := sweepTagValues(config, keyName)
			if !keySweepable {
				nonPrefixCount++
				continue
			}

			if remainingValues > 0 {
				// TagValue deletes are asynchronous and we deliberately don't wait on
				// them, so the key may still look non-empty. It will be picked up by
				// the next sweeper run once its values are really gone.
				log.Printf("[INFO][SWEEPER_LOG] Deleted %d tag values under %s, %d not deleted; deferring key delete", deletedValues, keyName, remainingValues)
				continue
			}

			if err := deleteTagResource(config, keyName); err != nil {
				log.Printf("[INFO][SWEEPER_LOG] Error deleting tag key %s (%s): %s", keyShortName, keyName, err)
			} else {
				log.Printf("[INFO][SWEEPER_LOG] Sent delete request for tag key %s (%s)", keyShortName, keyName)
			}
		}
	}

	if nonPrefixCount > 0 {
		log.Printf("[INFO][SWEEPER_LOG] %d tag keys were non-sweepable and skipped.", nonPrefixCount)
	}

	return nil
}

// sweepTagValues deletes every sweepable TagValue under a TagKey. It returns the
// number deleted and the number left behind, so the caller can decide whether the
// parent key is safe to delete.
func sweepTagValues(config *transport_tpg.Config, tagKeyName string) (deleted, remaining int) {
	tagValues, err := listTagResources(config, "tagValues", tagKeyName)
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error listing tag values under %s: %s", tagKeyName, err)
		// Unknown -- report a non-zero remainder so we don't try to delete the key.
		return 0, 1
	}

	for _, tagValue := range tagValues {
		valueName, _ := tagValue["name"].(string)
		valueShortName, _ := tagValue["shortName"].(string)
		if valueName == "" {
			continue
		}
		if !sweeper.IsSweepableTestResource(valueShortName) {
			remaining++
			continue
		}
		if err := deleteTagResource(config, valueName); err != nil {
			// A TagValue that is still bound to a resource cannot be deleted. The
			// binding goes away with the resource it is attached to, so this
			// resolves itself on a later run.
			log.Printf("[INFO][SWEEPER_LOG] Error deleting tag value %s (%s): %s", valueShortName, valueName, err)
			remaining++
			continue
		}
		log.Printf("[INFO][SWEEPER_LOG] Sent delete request for tag value %s (%s)", valueShortName, valueName)
		deleted++
	}
	return deleted, remaining
}

// listTagResources pages through tagKeys.list or tagValues.list for a parent.
func listTagResources(config *transport_tpg.Config, collection, parent string) ([]map[string]interface{}, error) {
	var out []map[string]interface{}
	pageToken := ""

	for {
		listUrl := fmt.Sprintf("%s%s?parent=%s&pageSize=300", tagsBaseUrl, collection, url.QueryEscape(parent))
		if pageToken != "" {
			listUrl += "&pageToken=" + url.QueryEscape(pageToken)
		}

		res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			Project:   config.Project,
			RawURL:    listUrl,
			UserAgent: config.UserAgent,
		})
		if err != nil {
			return nil, err
		}

		if items, ok := res[collection].([]interface{}); ok {
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
		// Guard against a server that keeps handing back the same token: without
		// this the sweep would spin here forever, burning CI time and producing no
		// signal. Stopping early is the better failure.
		if next == pageToken {
			log.Printf("[INFO][SWEEPER_LOG] %s list pageToken did not advance for parent %s; stopping", collection, parent)
			return out, nil
		}
		pageToken = next
	}
}

// deleteTagResource deletes by relative resource name, e.g. "tagKeys/123".
func deleteTagResource(config *transport_tpg.Config, name string) error {
	// Deletes are asynchronous; we don't wait on the operation because a sweep may
	// have a lot to delete.
	_, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "DELETE",
		Project:   config.Project,
		RawURL:    tagsBaseUrl + name,
		UserAgent: config.UserAgent,
	})
	return err
}
