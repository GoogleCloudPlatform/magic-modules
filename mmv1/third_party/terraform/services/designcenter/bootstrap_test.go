// Copyright 2026 Google Inc.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	bootstrapDesignCenterDestSpaceID         = "tf-bootstrap-dest-space"
	bootstrapDesignCenterCatalogID           = "tf-bootstrap-catalog"
	bootstrapDesignCenterCatalogTemplateID   = "tf-bootstrap-catalog-template"
	bootstrapDesignCenterAppTemplateID       = "tf-bootstrap-app-template"
	bootstrapDesignCenterComponentID         = "tf-bootstrap-src-component"
	bootstrapDesignCenterDestComponentID     = "tf-bootstrap-dest-component"
	bootstrapDesignCenterCloudRunTemplateURI = "projects/gcpdesigncenter/locations/us-central1/spaces/googlespace/sharedTemplates/1b5e09c8-780f-484e-b8ed-c7178f4e4342/revisions/r-24"
	bootstrapDesignCenterCloudSQLTemplateURI = "projects/gcpdesigncenter/locations/us-central1/spaces/googlespace/sharedTemplates/45be1bc2-89bc-477b-9b5d-64c41ff3b146/revisions/r-12"
)

var (
	bootstrapDesignCenterMu                   sync.Mutex
	bootstrapDesignCenterProjects             = make(map[string]string)
	bootstrapDesignCenterSpaces               = make(map[string]string)
	bootstrapDesignCenterDestSpaces           = make(map[string]string)
	bootstrapDesignCenterCatalogs             = make(map[string]string)
	bootstrapDesignCenterCatalogTemplates     = make(map[string]string)
	bootstrapDesignCenterAppTemplates         = make(map[string]string)
	bootstrapDesignCenterAppTemplateRevisions = make(map[string]string)
	bootstrapDesignCenterComponents           = make(map[string]string)
	bootstrapDesignCenterDestComponents       = make(map[string]string)

	bootstrapDesignCenterProjectServices = []string{
		"apphub.googleapis.com",
		"cloudresourcemanager.googleapis.com",
		"designcenter.googleapis.com",
		"iam.googleapis.com",
		"serviceusage.googleapis.com",
		"storage.googleapis.com",
	}
)

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
	case strings.HasPrefix(testName, "TestAccDesignCenterCatalogTemplate_"):
		return "tf-boot-dc-ct-"
	case strings.HasPrefix(testName, "TestAccDesignCenterCatalogTemplateRevision_"):
		return "tf-boot-dc-ctr-"
	case strings.HasPrefix(testName, "TestAccDesignCenterShare_"):
		return "tf-boot-dc-shr-"
	case strings.HasPrefix(testName, "TestAccDesignCenterApplication_"):
		return "tf-boot-dc-app-"
	case strings.HasPrefix(testName, "TestAccDesignCenterApplicationTemplate_"):
		return "tf-boot-dc-at-"
	case strings.HasPrefix(testName, "TestAccDesignCenterComponent_"):
		return "tf-boot-dc-cmp-"
	case strings.HasPrefix(testName, "TestAccDesignCenterConnection_"):
		return "tf-boot-dc-cnn-"
	case strings.HasPrefix(testName, "TestAccDesignCenterPolicy_"):
		return "tf-boot-dc-pol-"
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

	if os.Getenv("VCR_MODE") != "REPLAYING" {
		ensureDesignCenterAppHubBoundary(t, config, project, number)
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

// BootstrapDesignCenterDestinationSpace ensures a destination Space exists for
// sharing the bootstrapped Catalog and returns its space_id.
func BootstrapDesignCenterDestinationSpace(t *testing.T) string {
	t.Helper()

	return BootstrapDesignCenterSpace(t)
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

	if strings.HasPrefix(t.Name(), "TestAccDesignCenterShare_") {
		cleanupDesignCenterShares(t, config, project, bootstrapDesignCenterLocation, spaceID, bootstrapDesignCenterCatalogID)
	}

	bootstrapDesignCenterCatalogs[prefix] = bootstrapDesignCenterCatalogID
	return bootstrapDesignCenterCatalogID
}

// BootstrapDesignCenterCatalogTemplate ensures a shared CatalogTemplate exists
// in the bootstrapped Design Center Catalog and returns its catalog_template_id.
func BootstrapDesignCenterCatalogTemplate(t *testing.T) string {
	t.Helper()

	project := BootstrapDesignCenterProject(t)
	spaceID := BootstrapDesignCenterSpace(t)
	catalogID := BootstrapDesignCenterCatalog(t)
	if project == "" || spaceID == "" || catalogID == "" {
		return ""
	}
	if os.Getenv("VCR_MODE") == "REPLAYING" {
		return bootstrapDesignCenterCatalogTemplateID
	}

	config := transport_tpg.BootstrapConfig(t)
	if config == nil {
		return ""
	}

	prefix := designCenterProjectPrefixForTest(t.Name())

	bootstrapDesignCenterMu.Lock()
	defer bootstrapDesignCenterMu.Unlock()

	if templateID, ok := bootstrapDesignCenterCatalogTemplates[prefix]; ok && templateID != "" {
		return templateID
	}

	retryPredicates := []transport_tpg.RetryErrorPredicateFunc{isDesignCenterBootstrapRetryableError}
	templateURL := fmt.Sprintf("https://designcenter.googleapis.com/v1/projects/%s/locations/%s/spaces/%s/catalogs/%s/templates/%s", project, bootstrapDesignCenterLocation, spaceID, catalogID, bootstrapDesignCenterCatalogTemplateID)
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
			t.Fatalf("Bootstrapping failed. Error getting Design Center catalog template %s: %s", templateURL, err)
		}
		createURL := fmt.Sprintf("https://designcenter.googleapis.com/v1/projects/%s/locations/%s/spaces/%s/catalogs/%s/templates?catalogTemplateId=%s", project, bootstrapDesignCenterLocation, spaceID, catalogID, bootstrapDesignCenterCatalogTemplateID)
		op, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "POST",
			Project:   project,
			RawURL:    createURL,
			UserAgent: config.UserAgent,
			Body: map[string]interface{}{
				"displayName":      "Bootstrap Design Center Catalog Template",
				"description":      "Bootstrapped for Terraform acceptance tests",
				"templateCategory": "COMPONENT_TEMPLATE",
			},
			Timeout:              5 * time.Minute,
			ErrorRetryPredicates: retryPredicates,
		})
		if err != nil && !transport_tpg.IsGoogleApiErrorWithCode(err, 409) {
			t.Fatalf("Bootstrapping failed. Error creating Design Center catalog template %s: %s", templateURL, err)
		}
		waitForDesignCenterOperation(t, config, project, op)
	}

	bootstrapDesignCenterCatalogTemplates[prefix] = bootstrapDesignCenterCatalogTemplateID
	return bootstrapDesignCenterCatalogTemplateID
}

// BootstrapDesignCenterApplicationTemplate ensures an ApplicationTemplate
// exists inside the bootstrapped Space and returns its application_template_id.
func BootstrapDesignCenterApplicationTemplate(t *testing.T) string {
	t.Helper()

	project := BootstrapDesignCenterProject(t)
	if project == "" {
		return ""
	}
	spaceID := BootstrapDesignCenterSpace(t)
	if spaceID == "" {
		return ""
	}

	config := transport_tpg.BootstrapConfig(t)
	if config == nil || os.Getenv("VCR_MODE") == "REPLAYING" {
		return bootstrapDesignCenterAppTemplateID
	}

	bootstrapDesignCenterMu.Lock()
	defer bootstrapDesignCenterMu.Unlock()

	cacheKey := project + "/" + spaceID
	if appTemplateID, ok := bootstrapDesignCenterAppTemplates[cacheKey]; ok && appTemplateID != "" {
		return appTemplateID
	}

	retryPredicates := []transport_tpg.RetryErrorPredicateFunc{isDesignCenterBootstrapRetryableError}
	appTemplateURL := fmt.Sprintf("https://designcenter.googleapis.com/v1/projects/%s/locations/%s/spaces/%s/applicationTemplates/%s", project, bootstrapDesignCenterLocation, spaceID, bootstrapDesignCenterAppTemplateID)
	_, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:               config,
		Method:               "GET",
		Project:              project,
		RawURL:               appTemplateURL,
		UserAgent:            config.UserAgent,
		Timeout:              5 * time.Minute,
		ErrorRetryPredicates: retryPredicates,
	})
	if err != nil {
		if !transport_tpg.IsGoogleApiErrorWithCode(err, 404) {
			t.Fatalf("Bootstrapping failed. Error getting Design Center application template %s: %s", appTemplateURL, err)
		}

		createURL := fmt.Sprintf("https://designcenter.googleapis.com/v1/projects/%s/locations/%s/spaces/%s/applicationTemplates?applicationTemplateId=%s", project, bootstrapDesignCenterLocation, spaceID, bootstrapDesignCenterAppTemplateID)
		op, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "POST",
			Project:   project,
			RawURL:    createURL,
			UserAgent: config.UserAgent,
			Body: map[string]interface{}{
				"displayName": "Bootstrap Application Template",
				"description": "Bootstrapped for Terraform acceptance tests",
			},
			Timeout:              5 * time.Minute,
			ErrorRetryPredicates: retryPredicates,
		})
		if err != nil {
			t.Fatalf("Bootstrapping failed. Error creating Design Center application template %s: %s", createURL, err)
		}
		waitForDesignCenterOperation(t, config, project, op)
	}

	bootstrapDesignCenterAppTemplates[cacheKey] = bootstrapDesignCenterAppTemplateID
	return bootstrapDesignCenterAppTemplateID
}

// BootstrapDesignCenterApplicationTemplateRevisionURI ensures a committed
// revision exists for the bootstrapped ApplicationTemplate and returns its
// full revision URI.
func BootstrapDesignCenterApplicationTemplateRevisionURI(t *testing.T) string {
	t.Helper()

	project := BootstrapDesignCenterProject(t)
	if project == "" {
		return ""
	}
	spaceID := BootstrapDesignCenterSpace(t)
	if spaceID == "" {
		return ""
	}
	appTemplateID := BootstrapDesignCenterApplicationTemplate(t)
	if appTemplateID == "" {
		return ""
	}

	defaultRevisionURI := fmt.Sprintf("projects/%s/locations/%s/spaces/%s/applicationTemplates/%s/revisions/r-1", project, bootstrapDesignCenterLocation, spaceID, appTemplateID)
	config := transport_tpg.BootstrapConfig(t)
	if config == nil || os.Getenv("VCR_MODE") == "REPLAYING" {
		return defaultRevisionURI
	}

	bootstrapDesignCenterMu.Lock()
	defer bootstrapDesignCenterMu.Unlock()

	cacheKey := project + "/" + spaceID + "/" + appTemplateID
	if revURI, ok := bootstrapDesignCenterAppTemplateRevisions[cacheKey]; ok && revURI != "" {
		return revURI
	}

	ensureDesignCenterComponent(t, config, project, spaceID, appTemplateID, bootstrapDesignCenterComponentID, "Bootstrap Source Component", bootstrapDesignCenterCloudRunTemplateURI)

	retryPredicates := []transport_tpg.RetryErrorPredicateFunc{isDesignCenterBootstrapRetryableError}
	revURL := fmt.Sprintf("https://designcenter.googleapis.com/v1/%s", defaultRevisionURI)
	_, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:               config,
		Method:               "GET",
		Project:              project,
		RawURL:               revURL,
		UserAgent:            config.UserAgent,
		Timeout:              5 * time.Minute,
		ErrorRetryPredicates: retryPredicates,
	})
	if err != nil {
		if !transport_tpg.IsGoogleApiErrorWithCode(err, 404) {
			t.Fatalf("Bootstrapping failed. Error getting Design Center application template revision %s: %s", revURL, err)
		}

		commitURL := fmt.Sprintf("https://designcenter.googleapis.com/v1/projects/%s/locations/%s/spaces/%s/applicationTemplates/%s:commit", project, bootstrapDesignCenterLocation, spaceID, appTemplateID)
		res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:               config,
			Method:               "POST",
			Project:              project,
			RawURL:               commitURL,
			UserAgent:            config.UserAgent,
			Body:                 map[string]interface{}{},
			Timeout:              5 * time.Minute,
			ErrorRetryPredicates: retryPredicates,
		})
		if err != nil {
			t.Fatalf("Bootstrapping failed. Error committing Design Center application template %s: %s", commitURL, err)
		}
		waitForDesignCenterOperation(t, config, project, res)
		if name, ok := res["name"].(string); ok && strings.Contains(name, "/revisions/") {
			defaultRevisionURI = name
		}
	}

	bootstrapDesignCenterAppTemplateRevisions[cacheKey] = defaultRevisionURI
	return defaultRevisionURI
}

// BootstrapDesignCenterComponent ensures a source Component exists inside the
// bootstrapped ApplicationTemplate and returns its component_id.
func BootstrapDesignCenterComponent(t *testing.T) string {
	t.Helper()

	project := BootstrapDesignCenterProject(t)
	if project == "" {
		return ""
	}
	spaceID := BootstrapDesignCenterSpace(t)
	if spaceID == "" {
		return ""
	}
	appTemplateID := BootstrapDesignCenterApplicationTemplate(t)
	if appTemplateID == "" {
		return ""
	}

	config := transport_tpg.BootstrapConfig(t)
	if config == nil || os.Getenv("VCR_MODE") == "REPLAYING" {
		return bootstrapDesignCenterComponentID
	}

	bootstrapDesignCenterMu.Lock()
	defer bootstrapDesignCenterMu.Unlock()

	cacheKey := project + "/" + spaceID + "/" + appTemplateID
	if componentID, ok := bootstrapDesignCenterComponents[cacheKey]; ok && componentID != "" {
		return componentID
	}

	ensureDesignCenterComponent(t, config, project, spaceID, appTemplateID, bootstrapDesignCenterComponentID, "Bootstrap Source Component", bootstrapDesignCenterCloudRunTemplateURI)
	bootstrapDesignCenterComponents[cacheKey] = bootstrapDesignCenterComponentID
	return bootstrapDesignCenterComponentID
}

// BootstrapDesignCenterDestinationComponentURI ensures a destination Component
// exists inside the bootstrapped ApplicationTemplate and returns its full
// component URI.
func BootstrapDesignCenterDestinationComponentURI(t *testing.T) string {
	t.Helper()

	project := BootstrapDesignCenterProject(t)
	if project == "" {
		return ""
	}
	spaceID := BootstrapDesignCenterSpace(t)
	if spaceID == "" {
		return ""
	}
	appTemplateID := BootstrapDesignCenterApplicationTemplate(t)
	if appTemplateID == "" {
		return ""
	}

	destURI := fmt.Sprintf("projects/%s/locations/%s/spaces/%s/applicationTemplates/%s/components/%s", project, bootstrapDesignCenterLocation, spaceID, appTemplateID, bootstrapDesignCenterDestComponentID)
	config := transport_tpg.BootstrapConfig(t)
	if config == nil || os.Getenv("VCR_MODE") == "REPLAYING" {
		return destURI
	}

	bootstrapDesignCenterMu.Lock()
	defer bootstrapDesignCenterMu.Unlock()

	cacheKey := project + "/" + spaceID + "/" + appTemplateID
	if uri, ok := bootstrapDesignCenterDestComponents[cacheKey]; ok && uri != "" {
		return uri
	}

	ensureDesignCenterComponent(t, config, project, spaceID, appTemplateID, bootstrapDesignCenterDestComponentID, "Bootstrap Destination Component", bootstrapDesignCenterCloudSQLTemplateURI)
	bootstrapDesignCenterDestComponents[cacheKey] = destURI
	return destURI
}

func ensureDesignCenterComponent(t *testing.T, config *transport_tpg.Config, project, spaceID, appTemplateID, componentID, displayName, sharedTemplateRevisionURI string) {
	t.Helper()

	retryPredicates := []transport_tpg.RetryErrorPredicateFunc{isDesignCenterBootstrapRetryableError}
	componentURL := fmt.Sprintf("https://designcenter.googleapis.com/v1/projects/%s/locations/%s/spaces/%s/applicationTemplates/%s/components/%s", project, bootstrapDesignCenterLocation, spaceID, appTemplateID, componentID)
	_, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:               config,
		Method:               "GET",
		Project:              project,
		RawURL:               componentURL,
		UserAgent:            config.UserAgent,
		Timeout:              5 * time.Minute,
		ErrorRetryPredicates: retryPredicates,
	})
	if err == nil {
		return
	}
	if !transport_tpg.IsGoogleApiErrorWithCode(err, 404) {
		t.Fatalf("Bootstrapping failed. Error getting Design Center component %s: %s", componentURL, err)
	}

	createURL := fmt.Sprintf("https://designcenter.googleapis.com/v1/projects/%s/locations/%s/spaces/%s/applicationTemplates/%s/components?componentId=%s", project, bootstrapDesignCenterLocation, spaceID, appTemplateID, componentID)
	op, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "POST",
		Project:   project,
		RawURL:    createURL,
		UserAgent: config.UserAgent,
		Body: map[string]interface{}{
			"displayName":               displayName,
			"sharedTemplateRevisionUri": sharedTemplateRevisionURI,
		},
		Timeout:              5 * time.Minute,
		ErrorRetryPredicates: retryPredicates,
	})
	if err != nil {
		t.Fatalf("Bootstrapping failed. Error creating Design Center component %s: %s", createURL, err)
	}
	waitForDesignCenterOperation(t, config, project, op)
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

// cleanupDesignCenterShares removes any leftover Shares in the target Catalog
// from previous failed or interrupted test runs so that creating a new Share
// does not fail with FAILED_PRECONDITION ("CATALOG_ALREADY_SHARED").
func cleanupDesignCenterShares(t *testing.T, config *transport_tpg.Config, project, location, spaceID, catalogID string) {
	t.Helper()

	retryPredicates := []transport_tpg.RetryErrorPredicateFunc{isDesignCenterBootstrapRetryableError}
	listURL := fmt.Sprintf("https://designcenter.googleapis.com/v1/projects/%s/locations/%s/spaces/%s/catalogs/%s/shares", project, location, spaceID, catalogID)
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
		t.Fatalf("Bootstrapping failed. Error listing Design Center shares in %s/%s/%s/%s: %s", project, location, spaceID, catalogID, err)
	}

	shares, _ := res["shares"].([]interface{})
	for _, s := range shares {
		shareMap, _ := s.(map[string]interface{})
		name, _ := shareMap["name"].(string)
		if name == "" {
			continue
		}
		log.Printf("[DEBUG] Deleting leftover Design Center share %s before test", name)
		deleteURL := fmt.Sprintf("https://designcenter.googleapis.com/v1/%s", name)
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
			t.Fatalf("Bootstrapping failed. Error deleting leftover Design Center share %s: %s", name, err)
		}
		waitForDesignCenterOperation(t, config, project, op)
	}
}
