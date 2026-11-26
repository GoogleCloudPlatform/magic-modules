package storage

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-provider-google/google/envvar"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

const SharedTestAnywhereCacheBucketPrefix = "tf-bootstrap-ac-"

type BootstrappedAnywhereCache struct {
	Bucket          string
	AnywhereCacheID string
	Zone            string
}

func sharedTestAnywhereCacheBucketName(projectID, testId string) string {
	sanitizedProject := strings.ReplaceAll(strings.ToLower(projectID), "google", "gcp")
	bucketName := fmt.Sprintf("%s%s-%s", SharedTestAnywhereCacheBucketPrefix, sanitizedProject, strings.ToLower(testId))
	if len(bucketName) > 63 {
		bucketName = strings.TrimRight(bucketName[:63], "-")
	}
	return bucketName
}

// BootstrapSharedTestAnywhereCache returns a persistent GCS bucket (with uniform
// bucket-level access enabled) and AnywhereCache in the specified location and
// zone for tests that reference an AnywhereCache without mutating or deleting it
// (such as google_storage_managed_folder rapid_cache_config tests).
//
// Provisioning a fresh AnywhereCache LRO can take 1.5–4+ hours when zonal SSD
// capacity is constrained. Because bucket metadata registers the AnywhereCache
// instance immediately once creation is initiated (while state is "creating" or
// "running"), this helper initiates cache creation if absent without blocking on
// the multi-hour LRO, and reuses the persistent tf-bootstrap-ac-* bucket and
// cache across subsequent test runs.
func BootstrapSharedTestAnywhereCache(t *testing.T, testId, location, zone string) BootstrappedAnywhereCache {
	t.Helper()

	projectID := envvar.GetTestProjectFromEnv()
	bucketName := sharedTestAnywhereCacheBucketName(projectID, testId)
	if os.Getenv("VCR_PATH") != "" && os.Getenv("VCR_MODE") == "REPLAYING" {
		return BootstrappedAnywhereCache{
			Bucket:          bucketName,
			AnywhereCacheID: zone,
			Zone:            zone,
		}
	}

	config := transport_tpg.BootstrapConfig(t)
	if config == nil {
		t.Fatal("Could not bootstrap config.")
	}

	baseURL := transport_tpg.BaseUrl(Product, config)
	bucketURL := fmt.Sprintf("%sb/%s", baseURL, bucketName)

	log.Printf("[DEBUG] Getting shared test bucket %q", bucketName)
	_, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   projectID,
		RawURL:    bucketURL,
		UserAgent: config.UserAgent,
	})
	if err != nil {
		if !transport_tpg.IsGoogleApiErrorWithCode(err, 404) {
			t.Fatalf("Error getting shared test bucket %q: %s", bucketName, err)
		}

		log.Printf("[DEBUG] Shared test bucket %q not found, bootstrapping in %q", bucketName, location)
		createBucketURL, err := transport_tpg.AddQueryParams(baseURL+"b", map[string]string{
			"project": projectID,
		})
		if err != nil {
			t.Fatalf("Error building create URL for shared test bucket %q: %s", bucketName, err)
		}

		bucketObj := map[string]interface{}{
			"name":     bucketName,
			"location": location,
			"iamConfiguration": map[string]interface{}{
				"uniformBucketLevelAccess": map[string]interface{}{
					"enabled": true,
				},
			},
		}

		_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "POST",
			Project:   projectID,
			RawURL:    createBucketURL,
			UserAgent: config.UserAgent,
			Body:      bucketObj,
			Timeout:   4 * time.Minute,
		})
		if err != nil && !transport_tpg.IsGoogleApiErrorWithCode(err, 409) {
			t.Fatalf("Error bootstrapping shared test bucket %q: %s", bucketName, err)
		}
	}

	cacheURL := fmt.Sprintf("%sb/%s/anywhereCaches/%s", baseURL, bucketName, zone)
	log.Printf("[DEBUG] Getting shared test AnywhereCache %q in bucket %q", zone, bucketName)
	cacheRes, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   projectID,
		RawURL:    cacheURL,
		UserAgent: config.UserAgent,
	})
	if err != nil {
		if !transport_tpg.IsGoogleApiErrorWithCode(err, 404) {
			t.Fatalf("Error getting shared test AnywhereCache %q in bucket %q: %s", zone, bucketName, err)
		}

		log.Printf("[DEBUG] Shared test AnywhereCache %q not found in bucket %q, bootstrapping", zone, bucketName)
		createCacheURL := fmt.Sprintf("%sb/%s/anywhereCaches/", baseURL, bucketName)
		cacheObj := map[string]interface{}{
			"zone":          zone,
			"ttl":           "3601s",
			"ingestOnWrite": false,
		}

		_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "POST",
			Project:   projectID,
			RawURL:    createCacheURL,
			UserAgent: config.UserAgent,
			Body:      cacheObj,
			Timeout:   4 * time.Minute,
		})
		if err != nil && !transport_tpg.IsGoogleApiErrorWithCode(err, 409) {
			t.Fatalf("Error bootstrapping shared test AnywhereCache %q in bucket %q: %s", zone, bucketName, err)
		}

		cacheRes, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			Project:   projectID,
			RawURL:    cacheURL,
			UserAgent: config.UserAgent,
		})
		if err != nil {
			t.Fatalf("Error verifying bootstrapped AnywhereCache %q in bucket %q: %s", zone, bucketName, err)
		}
	} else {
		switch state := strings.ToLower(fmt.Sprint(cacheRes["state"])); state {
		case "paused", "disabled":
			resumeURL := fmt.Sprintf("%sb/%s/anywhereCaches/%s/resume", baseURL, bucketName, zone)
			log.Printf("[DEBUG] Resuming shared test AnywhereCache %q in bucket %q (state=%q)", zone, bucketName, state)
			if _, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "POST",
				Project:   projectID,
				RawURL:    resumeURL,
				UserAgent: config.UserAgent,
				Timeout:   4 * time.Minute,
			}); err != nil {
				t.Fatalf("Error resuming shared test AnywhereCache %q in bucket %q: %s", zone, bucketName, err)
			}
		}
	}

	anywhereCacheID := zone
	if id, ok := cacheRes["anywhereCacheId"].(string); ok && id != "" {
		anywhereCacheID = id
	}

	return BootstrappedAnywhereCache{
		Bucket:          bucketName,
		AnywhereCacheID: anywhereCacheID,
		Zone:            zone,
	}
}
