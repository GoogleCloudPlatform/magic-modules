package designcenter_test

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/services/resourcemanager"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

const (
	bootstrapDesignCenterLocation  = "us-central1"
	bootstrapDesignCenterSpaceID   = "tf-bootstrap-space"
	bootstrapDesignCenterCatalogID = "tf-bootstrap-catalog"
)

var (
	bootstrapDesignCenterMu       sync.Mutex
	bootstrapDesignCenterProjects = make(map[string]string)
	bootstrapDesignCenterSpaces   = make(map[string]string)
	bootstrapDesignCenterCatalogs = make(map[string]string)

	bootstrapDesignCenterProjectServices = []string{
		"apphub.googleapis.com",
		"cloudresourcemanager.googleapis.com",
		"designcenter.googleapis.com",
		"iam.googleapis.com",
		"serviceusage.googleapis.com",
		"storage.googleapis.com",
	}
)

// designCenterProjectPrefixForTest returns a dedicated project prefix (<= 18
// chars) per Design Center acceptance test. Application Design Center allows
// only one Space per region within a management project and only one Catalog
// per Space.
func designCenterProjectPrefixForTest(testName string) string {
	switch {
	case strings.HasPrefix(testName, "TestAccDesignCenterSpace_"):
		return "tf-boot-dc-spc-"
	case strings.HasPrefix(testName, "TestAccDesignCenterSpaceIamBinding"):
		return "tf-boot-dc-ib-"
	case strings.HasPrefix(testName, "TestAccDesignCenterSpaceIamMember"):
		return "tf-boot-dc-im-"
	case strings.HasPrefix(testName, "TestAccDesignCenterSpaceIamPolicy"):
		return "tf-boot-dc-ip-"
	case strings.HasPrefix(testName, "TestAccDesignCenterCatalog_"):
		return "tf-boot-dc-cat-"
	default:
		return "tf-boot-dc-"
	}
}

// BootstrapDesignCenterProject returns an isolated App Hub management project
// suitable for creating Application Design Center Spaces.
//
// Application Design Center requires a Space's parent project to be an App Hub
// management project (with its App Hub boundary crmNode set to the project
// itself) and allows only one Space per region within that management project.
func BootstrapDesignCenterProject(t *testing.T) string {
	t.Helper()

	config := transport_tpg.BootstrapConfig(t)
	if config == nil {
		return ""
	}

	prefix := designCenterProjectPrefixForTest(t.Name())

	bootstrapDesignCenterMu.Lock()
	defer bootstrapDesignCenterMu.Unlock()

	if project, ok := bootstrapDesignCenterProjects[prefix]; ok && project != "" {
		return project
	}

	p := resourcemanager.BootstrapProject(
		t,
		prefix,
		envvar.GetTestBillingAccountFromEnv(t),
		bootstrapDesignCenterProjectServices,
	)
	if p == nil {
		t.Fatal("Bootstrapping failed. Could not bootstrap the isolated Design Center test project.")
	}

	project := p.ProjectId
	number := strconv.FormatInt(p.ProjectNumber, 10)
	ensureDesignCenterAppHubBoundary(t, config, project, number)

	if os.Getenv("VCR_MODE") != "REPLAYING" {
		ensureDesignCenterServiceIdentity(t, config, project, number)
		if strings.HasPrefix(t.Name(), "TestAccDesignCenterSpace") {
			cleanupDesignCenterSpaces(t, config, project, bootstrapDesignCenterLocation)
		}
	}

	bootstrapDesignCenterProjects[prefix] = project
	return project
}

// BootstrapDesignCenterSpace ensures a shared Space exists in the isolated
// Design Center test project and returns its space_id.
func BootstrapDesignCenterSpace(t *testing.T) string {
	t.Helper()

	project := BootstrapDesignCenterProject(t)
	if project == "" {
		return ""
	}
	if os.Getenv("VCR_MODE") == "REPLAYING" {
		return bootstrapDesignCenterSpaceID
	}

	config := transport_tpg.BootstrapConfig(t)
	if config == nil {
		return ""
	}

	prefix := designCenterProjectPrefixForTest(t.Name())

	bootstrapDesignCenterMu.Lock()
	defer bootstrapDesignCenterMu.Unlock()

	if spaceID, ok := bootstrapDesignCenterSpaces[prefix]; ok && spaceID != "" {
		return spaceID
	}

	retryPredicates := []transport_tpg.RetryErrorPredicateFunc{isDesignCenterBootstrapRetryableError}
	spaceURL := fmt.Sprintf("https://designcenter.googleapis.com/v1/projects/%s/locations/%s/spaces/%s", project, bootstrapDesignCenterLocation, bootstrapDesignCenterSpaceID)
	_, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:               config,
		Method:               "GET",
		Project:              project,
		RawURL:               spaceURL,
		UserAgent:            config.UserAgent,
		Timeout:              5 * time.Minute,
		ErrorRetryPredicates: retryPredicates,
	})
	if err != nil {
		if !transport_tpg.IsGoogleApiErrorWithCode(err, 404) {
			t.Fatalf("Bootstrapping failed. Error getting Design Center space %s: %s", spaceURL, err)
		}
		cleanupDesignCenterSpaces(t, config, project, bootstrapDesignCenterLocation)
		createURL := fmt.Sprintf("https://designcenter.googleapis.com/v1/projects/%s/locations/%s/spaces?spaceId=%s", project, bootstrapDesignCenterLocation, bootstrapDesignCenterSpaceID)
		op, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:               config,
			Method:               "POST",
			Project:              project,
			RawURL:               createURL,
			UserAgent:            config.UserAgent,
			Body:                 map[string]interface{}{"displayName": "Bootstrap Design Center Space", "enableGcpSharedTemplates": true},
			Timeout:              5 * time.Minute,
			ErrorRetryPredicates: retryPredicates,
		})
		if err != nil && !transport_tpg.IsGoogleApiErrorWithCode(err, 409) {
			t.Fatalf("Bootstrapping failed. Error creating Design Center space %s: %s", spaceURL, err)
		}
		waitForDesignCenterOperation(t, config, project, op)
	}

	if strings.HasPrefix(t.Name(), "TestAccDesignCenterCatalog_") {
		cleanupDesignCenterCatalogs(t, config, project, bootstrapDesignCenterLocation, bootstrapDesignCenterSpaceID)
	}

	bootstrapDesignCenterSpaces[prefix] = bootstrapDesignCenterSpaceID
	return bootstrapDesignCenterSpaceID
}

// BootstrapDesignCenterCatalog ensures a shared Catalog exists in the
// bootstrapped Design Center Space and returns its catalog_id.
func BootstrapDesignCenterCatalog(t *testing.T) string {
	t.Helper()

	project := BootstrapDesignCenterProject(t)
	spaceID := BootstrapDesignCenterSpace(t)
	if project == "" || spaceID == "" {
		return ""
	}
	if os.Getenv("VCR_MODE") == "REPLAYING" {
		return bootstrapDesignCenterCatalogID
	}

	config := transport_tpg.BootstrapConfig(t)
	if config == nil {
		return ""
	}

	prefix := designCenterProjectPrefixForTest(t.Name())

	bootstrapDesignCenterMu.Lock()
	defer bootstrapDesignCenterMu.Unlock()

	if catalogID, ok := bootstrapDesignCenterCatalogs[prefix]; ok && catalogID != "" {
		return catalogID
	}

	retryPredicates := []transport_tpg.RetryErrorPredicateFunc{isDesignCenterBootstrapRetryableError}
	catalogURL := fmt.Sprintf("https://designcenter.googleapis.com/v1/projects/%s/locations/%s/spaces/%s/catalogs/%s", project, bootstrapDesignCenterLocation, spaceID, bootstrapDesignCenterCatalogID)
	_, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:               config,
		Method:               "GET",
		Project:              project,
		RawURL:               catalogURL,
		UserAgent:            config.UserAgent,
		Timeout:              5 * time.Minute,
		ErrorRetryPredicates: retryPredicates,
	})
	if err != nil {
		if !transport_tpg.IsGoogleApiErrorWithCode(err, 404) {
			t.Fatalf("Bootstrapping failed. Error getting Design Center catalog %s: %s", catalogURL, err)
		}
		cleanupDesignCenterCatalogs(t, config, project, bootstrapDesignCenterLocation, spaceID)
		createURL := fmt.Sprintf("https://designcenter.googleapis.com/v1/projects/%s/locations/%s/spaces/%s/catalogs?catalogId=%s", project, bootstrapDesignCenterLocation, spaceID, bootstrapDesignCenterCatalogID)
		op, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:               config,
			Method:               "POST",
			Project:              project,
			RawURL:               createURL,
			UserAgent:            config.UserAgent,
			Body:                 map[string]interface{}{"displayName": "Bootstrap Design Center Catalog"},
			Timeout:              5 * time.Minute,
			ErrorRetryPredicates: retryPredicates,
		})
		if err != nil && !transport_tpg.IsGoogleApiErrorWithCode(err, 409) {
			t.Fatalf("Bootstrapping failed. Error creating Design Center catalog %s: %s", catalogURL, err)
		}
		waitForDesignCenterOperation(t, config, project, op)
	}

	bootstrapDesignCenterCatalogs[prefix] = bootstrapDesignCenterCatalogID
	return bootstrapDesignCenterCatalogID
}

func isDesignCenterBootstrapRetryableError(err error) (bool, string) {
	if err == nil {
		return false, ""
	}
	if transport_tpg.IsGoogleApiErrorWithCode(err, 403) {
		return true, "Waiting for project IAM permissions and service activation to propagate"
	}
	msg := err.Error()
	if strings.Contains(msg, "gcp-sa-designcenter.iam.gserviceaccount.com") &&
		strings.Contains(msg, "ServiceAccountInfoDataType") {
		return true, "Waiting for Design Center service agent to propagate"
	}
	return false, ""
}

func waitForDesignCenterOperation(t *testing.T, config *transport_tpg.Config, project string, op map[string]interface{}) {
	t.Helper()
	if op == nil {
		return
	}
	name, _ := op["name"].(string)
	if !strings.Contains(name, "/operations/") {
		return
	}

	retryPredicates := []transport_tpg.RetryErrorPredicateFunc{isDesignCenterBootstrapRetryableError}
	deadline := time.Now().Add(5 * time.Minute)
	for {
		opRes, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:               config,
			Method:               "GET",
			Project:              project,
			RawURL:               "https://designcenter.googleapis.com/v1/" + name,
			UserAgent:            config.UserAgent,
			Timeout:              5 * time.Minute,
			ErrorRetryPredicates: retryPredicates,
		})
		if err != nil {
			t.Fatalf("Bootstrapping failed. Error waiting for operation %s: %s", name, err)
		}
		if done, _ := opRes["done"].(bool); done {
			if e, ok := opRes["error"]; ok {
				t.Fatalf("Bootstrapping failed. Operation %s failed: %v", name, e)
			}
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("Bootstrapping failed. Timed out waiting for operation %s", name)
		}
		time.Sleep(5 * time.Second)
	}
}

func ensureDesignCenterServiceIdentity(t *testing.T, config *transport_tpg.Config, project, number string) {
	t.Helper()

	retryPredicates := []transport_tpg.RetryErrorPredicateFunc{isDesignCenterBootstrapRetryableError}
	genURL := fmt.Sprintf("https://serviceusage.googleapis.com/v1beta1/projects/%s/services/designcenter.googleapis.com:generateServiceIdentity", number)
	op, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:               config,
		Method:               "POST",
		Project:              project,
		RawURL:               genURL,
		UserAgent:            config.UserAgent,
		Timeout:              5 * time.Minute,
		ErrorRetryPredicates: retryPredicates,
	})
	if err != nil {
		t.Fatalf("Bootstrapping failed. Error generating Design Center service identity for %s: %s", project, err)
	}

	name, _ := op["name"].(string)
	if !strings.Contains(name, "/operations/") || name == "operations/noop.DONE_OPERATION" {
		return
	}

	deadline := time.Now().Add(5 * time.Minute)
	for {
		opRes, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:               config,
			Method:               "GET",
			Project:              project,
			RawURL:               "https://serviceusage.googleapis.com/v1beta1/" + name,
			UserAgent:            config.UserAgent,
			Timeout:              5 * time.Minute,
			ErrorRetryPredicates: retryPredicates,
		})
		if err != nil {
			t.Fatalf("Bootstrapping failed. Error waiting for operation %s: %s", name, err)
		}
		if done, _ := opRes["done"].(bool); done {
			if e, ok := opRes["error"]; ok {
				t.Fatalf("Bootstrapping failed. Operation %s failed: %v", name, e)
			}
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("Bootstrapping failed. Timed out waiting for operation %s", name)
		}
		time.Sleep(5 * time.Second)
	}
}

func ensureDesignCenterAppHubBoundary(t *testing.T, config *transport_tpg.Config, project, number string) {
	t.Helper()

	retryPredicates := []transport_tpg.RetryErrorPredicateFunc{isDesignCenterBootstrapRetryableError}
	boundaryURL := fmt.Sprintf("https://apphub.googleapis.com/v1/projects/%s/locations/global/boundary", project)
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:               config,
		Method:               "GET",
		Project:              project,
		RawURL:               boundaryURL,
		UserAgent:            config.UserAgent,
		Timeout:              5 * time.Minute,
		ErrorRetryPredicates: retryPredicates,
	})
	if err != nil {
		t.Fatalf("Bootstrapping failed. Error getting App Hub boundary %s: %s", boundaryURL, err)
	}

	crmNode := "projects/" + number
	if res["crmNode"] == crmNode {
		return
	}

	log.Printf("[DEBUG] Setting App Hub boundary of %s to %s", project, crmNode)
	op, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:               config,
		Method:               "PATCH",
		Project:              project,
		RawURL:               boundaryURL + "?updateMask=crmNode",
		UserAgent:            config.UserAgent,
		Body:                 map[string]interface{}{"crmNode": crmNode},
		Timeout:              5 * time.Minute,
		ErrorRetryPredicates: retryPredicates,
	})
	if err != nil {
		t.Fatalf("Bootstrapping failed. Error patching App Hub boundary %s: %s", boundaryURL, err)
	}

	name, _ := op["name"].(string)
	if !strings.Contains(name, "/operations/") {
		return
	}

	deadline := time.Now().Add(5 * time.Minute)
	for {
		opRes, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:               config,
			Method:               "GET",
			Project:              project,
			RawURL:               "https://apphub.googleapis.com/v1/" + name,
			UserAgent:            config.UserAgent,
			Timeout:              5 * time.Minute,
			ErrorRetryPredicates: retryPredicates,
		})
		if err != nil {
			t.Fatalf("Bootstrapping failed. Error waiting for operation %s: %s", name, err)
		}
		if done, _ := opRes["done"].(bool); done {
			if e, ok := opRes["error"]; ok {
				t.Fatalf("Bootstrapping failed. Operation %s failed: %v", name, e)
			}
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("Bootstrapping failed. Timed out waiting for operation %s", name)
		}
		time.Sleep(5 * time.Second)
	}
}

// cleanupDesignCenterSpaces removes any leftover Spaces in the target region
// from previous failed or interrupted test runs so that creating a new Space
// does not fail with 409 ("only one space is allowed per region within a
// management project").
func cleanupDesignCenterSpaces(t *testing.T, config *transport_tpg.Config, project, location string) {
	t.Helper()

	retryPredicates := []transport_tpg.RetryErrorPredicateFunc{isDesignCenterBootstrapRetryableError}
	listURL := fmt.Sprintf("https://designcenter.googleapis.com/v1/projects/%s/locations/%s/spaces", project, location)
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:               config,
		Method:               "GET",
		Project:              project,
		RawURL:               listURL,
		UserAgent:            config.UserAgent,
		Timeout:              5 * time.Minute,
		ErrorRetryPredicates: retryPredicates,
	})
	if err != nil {
		if transport_tpg.IsGoogleApiErrorWithCode(err, 404) {
			return
		}
		t.Fatalf("Bootstrapping failed. Error listing Design Center spaces in %s/%s: %s", project, location, err)
	}

	spaces, _ := res["spaces"].([]interface{})
	for _, s := range spaces {
		spaceMap, _ := s.(map[string]interface{})
		name, _ := spaceMap["name"].(string)
		if name == "" {
			continue
		}
		log.Printf("[DEBUG] Deleting leftover Design Center space %s before test", name)
		deleteURL := fmt.Sprintf("https://designcenter.googleapis.com/v1/%s?force=true", name)
		op, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:               config,
			Method:               "DELETE",
			Project:              project,
			RawURL:               deleteURL,
			UserAgent:            config.UserAgent,
			Timeout:              5 * time.Minute,
			ErrorRetryPredicates: retryPredicates,
		})
		if err != nil && !transport_tpg.IsGoogleApiErrorWithCode(err, 404) {
			t.Fatalf("Bootstrapping failed. Error deleting leftover Design Center space %s: %s", name, err)
		}
		waitForDesignCenterOperation(t, config, project, op)
	}
}

// cleanupDesignCenterCatalogs removes any leftover Catalogs in the target Space
// from previous failed or interrupted test runs so that creating a new Catalog
// does not fail with FAILED_PRECONDITION ("Only one catalog is allowed per
// space").
func cleanupDesignCenterCatalogs(t *testing.T, config *transport_tpg.Config, project, location, spaceID string) {
	t.Helper()

	retryPredicates := []transport_tpg.RetryErrorPredicateFunc{isDesignCenterBootstrapRetryableError}
	listURL := fmt.Sprintf("https://designcenter.googleapis.com/v1/projects/%s/locations/%s/spaces/%s/catalogs", project, location, spaceID)
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:               config,
		Method:               "GET",
		Project:              project,
		RawURL:               listURL,
		UserAgent:            config.UserAgent,
		Timeout:              5 * time.Minute,
		ErrorRetryPredicates: retryPredicates,
	})
	if err != nil {
		if transport_tpg.IsGoogleApiErrorWithCode(err, 404) {
			return
		}
		t.Fatalf("Bootstrapping failed. Error listing Design Center catalogs in %s/%s/%s: %s", project, location, spaceID, err)
	}

	catalogs, _ := res["catalogs"].([]interface{})
	for _, c := range catalogs {
		catalogMap, _ := c.(map[string]interface{})
		name, _ := catalogMap["name"].(string)
		if name == "" {
			continue
		}
		log.Printf("[DEBUG] Deleting leftover Design Center catalog %s before test", name)
		deleteURL := fmt.Sprintf("https://designcenter.googleapis.com/v1/%s?force=true", name)
		op, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:               config,
			Method:               "DELETE",
			Project:              project,
			RawURL:               deleteURL,
			UserAgent:            config.UserAgent,
			Timeout:              5 * time.Minute,
			ErrorRetryPredicates: retryPredicates,
		})
		if err != nil && !transport_tpg.IsGoogleApiErrorWithCode(err, 404) {
			t.Fatalf("Bootstrapping failed. Error deleting leftover Design Center catalog %s: %s", name, err)
		}
		waitForDesignCenterOperation(t, config, project, op)
	}
}
