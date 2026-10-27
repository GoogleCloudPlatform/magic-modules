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
	bootstrapDesignCenterLocation            = "us-central1"
	bootstrapDesignCenterSpaceID             = "tf-bootstrap-space"
	bootstrapDesignCenterApplicationTemplate = "tf-bootstrap-app-template"
)

var (
	bootstrapDesignCenterMu           sync.Mutex
	bootstrapDesignCenterProjects     = make(map[string]string)
	bootstrapDesignCenterSpaces       = make(map[string]string)
	bootstrapDesignCenterAppTemplates = make(map[string]string)

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
		ensureDesignCenterServiceIdentity(t, config, project, number)
		if strings.HasPrefix(t.Name(), "TestAccDesignCenterSpace") {
			cleanupDesignCenterSpaces(t, config, project, bootstrapDesignCenterLocation)
		}
	}

	bootstrapDesignCenterProjects[prefix] = project
	return project
}

// BootstrapDesignCenterSpace ensures a shared Space exists in the isolated
// Design Center management project and returns its space_id.
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

	bootstrapDesignCenterMu.Lock()
	defer bootstrapDesignCenterMu.Unlock()

	if space, ok := bootstrapDesignCenterSpaces[project]; ok && space != "" {
		return space
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
			Body:                 map[string]interface{}{"displayName": "Bootstrap Design Center Space"},
			Timeout:              5 * time.Minute,
			ErrorRetryPredicates: retryPredicates,
		})
		if err != nil && !transport_tpg.IsGoogleApiErrorWithCode(err, 409) {
			t.Fatalf("Bootstrapping failed. Error creating Design Center space %s: %s", createURL, err)
		}
		waitForDesignCenterOperation(t, config, project, op)
	}

	bootstrapDesignCenterSpaces[project] = bootstrapDesignCenterSpaceID
	return bootstrapDesignCenterSpaceID
}

// BootstrapDesignCenterApplicationTemplate ensures a shared ApplicationTemplate
// exists in the bootstrapped Design Center Space and returns its
// application_template_id.
func BootstrapDesignCenterApplicationTemplate(t *testing.T) string {
	t.Helper()

	project := BootstrapDesignCenterProject(t)
	space := BootstrapDesignCenterSpace(t)
	if project == "" || space == "" {
		return ""
	}
	if os.Getenv("VCR_MODE") == "REPLAYING" {
		return bootstrapDesignCenterApplicationTemplate
	}

	config := transport_tpg.BootstrapConfig(t)
	if config == nil {
		return ""
	}

	bootstrapDesignCenterMu.Lock()
	defer bootstrapDesignCenterMu.Unlock()

	cacheKey := project + "/" + space
	if appTemplate, ok := bootstrapDesignCenterAppTemplates[cacheKey]; ok && appTemplate != "" {
		return appTemplate
	}

	retryPredicates := []transport_tpg.RetryErrorPredicateFunc{isDesignCenterBootstrapRetryableError}
	templateURL := fmt.Sprintf("https://designcenter.googleapis.com/v1/projects/%s/locations/%s/spaces/%s/applicationTemplates/%s", project, bootstrapDesignCenterLocation, space, bootstrapDesignCenterApplicationTemplate)
	_, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:               config,
		Method:               "GET",
		Project:              project,
		RawURL:               templateURL,
		UserAgent:            config.UserAgent,
		Timeout:              5 * time.Minute,
		ErrorRetryPredicates: retryPredicates,
	})
	if err != nil {
		if !transport_tpg.IsGoogleApiErrorWithCode(err, 404) {
			t.Fatalf("Bootstrapping failed. Error getting Design Center application template %s: %s", templateURL, err)
		}
		createURL := fmt.Sprintf("https://designcenter.googleapis.com/v1/projects/%s/locations/%s/spaces/%s/applicationTemplates?applicationTemplateId=%s", project, bootstrapDesignCenterLocation, space, bootstrapDesignCenterApplicationTemplate)
		op, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "POST",
			Project:   project,
			RawURL:    createURL,
			UserAgent: config.UserAgent,
			Body: map[string]interface{}{
				"displayName":     "Bootstrap Design Center Application Template",
				"compositionType": "STANDARD",
			},
			Timeout:              5 * time.Minute,
			ErrorRetryPredicates: retryPredicates,
		})
		if err != nil && !transport_tpg.IsGoogleApiErrorWithCode(err, 409) {
			t.Fatalf("Bootstrapping failed. Error creating Design Center application template %s: %s", createURL, err)
		}
		waitForDesignCenterOperation(t, config, project, op)
	}

	bootstrapDesignCenterAppTemplates[cacheKey] = bootstrapDesignCenterApplicationTemplate
	return bootstrapDesignCenterApplicationTemplate
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
			t.Fatalf("Bootstrapping failed. Error waiting for Design Center operation %s: %s", name, err)
		}
		if done, _ := opRes["done"].(bool); done {
			if e, ok := opRes["error"]; ok {
				t.Fatalf("Bootstrapping failed. Design Center operation %s failed: %v", name, e)
			}
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("Bootstrapping failed. Timed out waiting for Design Center operation %s", name)
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
