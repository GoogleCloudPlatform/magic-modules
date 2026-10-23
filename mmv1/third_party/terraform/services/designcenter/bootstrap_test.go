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

var (
	bootstrapDesignCenterMu       sync.Mutex
	bootstrapDesignCenterProjects = make(map[string]string)

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
// chars) per Space acceptance test. Application Design Center allows only one
// Space per region within a management project, and the resource test plus
// three generated IAM tests all run in parallel in us-central1.
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
		cleanupDesignCenterSpaces(t, config, project, "us-central1")
	}

	bootstrapDesignCenterProjects[prefix] = project
	return project
}

func isDesignCenterBootstrapRetryableError(err error) (bool, string) {
	if transport_tpg.IsGoogleApiErrorWithCode(err, 403) {
		return true, "Waiting for project IAM permissions and service activation to propagate"
	}
	return false, ""
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
		_, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
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
	}
}
