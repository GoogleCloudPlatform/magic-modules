package designcenter_test

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/services/resourcemanager"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

const bootstrapDesignCenterProjectPrefix = "tf-boot-saas-"

var (
	bootstrapDesignCenterMu        sync.Mutex
	bootstrapDesignCenterProjectID string

	bootstrapDesignCenterProjectServices = []string{
		"apphub.googleapis.com",
		"artifactregistry.googleapis.com",
		"cloudresourcemanager.googleapis.com",
		"designcenter.googleapis.com",
		"iam.googleapis.com",
		"saasservicemgmt.googleapis.com",
		"serviceusage.googleapis.com",
		"storage.googleapis.com",
	}
)

// BootstrapDesignCenterProject returns an isolated App Hub management project
// suitable for creating Application Design Center Spaces.
//
// Application Design Center requires a Space's parent project to be an App Hub
// management project (with its App Hub boundary crmNode set to the project
// itself). Because other acceptance tests (TestAccApphubBoundary_*) modify the
// boundary on the default CI test project, Design Center tests use the
// bootstrapped isolated project.
func BootstrapDesignCenterProject(t *testing.T) string {
	t.Helper()

	config := transport_tpg.BootstrapConfig(t)
	if config == nil {
		return ""
	}

	bootstrapDesignCenterMu.Lock()
	defer bootstrapDesignCenterMu.Unlock()

	if bootstrapDesignCenterProjectID != "" {
		return bootstrapDesignCenterProjectID
	}

	p := resourcemanager.BootstrapProject(
		t,
		bootstrapDesignCenterProjectPrefix,
		envvar.GetTestBillingAccountFromEnv(t),
		bootstrapDesignCenterProjectServices,
	)
	if p == nil {
		t.Fatal("Bootstrapping failed. Could not bootstrap the isolated Design Center test project.")
	}

	project := p.ProjectId
	number := strconv.FormatInt(p.ProjectNumber, 10)
	ensureDesignCenterAppHubBoundary(t, config, project, number)

	bootstrapDesignCenterProjectID = project
	return bootstrapDesignCenterProjectID
}

func ensureDesignCenterAppHubBoundary(t *testing.T, config *transport_tpg.Config, project, number string) {
	t.Helper()

	boundaryURL := fmt.Sprintf("https://apphub.googleapis.com/v1/projects/%s/locations/global/boundary", project)
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   project,
		RawURL:    boundaryURL,
		UserAgent: config.UserAgent,
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
		Config:    config,
		Method:    "PATCH",
		Project:   project,
		RawURL:    boundaryURL + "?updateMask=crmNode",
		UserAgent: config.UserAgent,
		Body:      map[string]interface{}{"crmNode": crmNode},
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
			Config:    config,
			Method:    "GET",
			Project:   project,
			RawURL:    "https://apphub.googleapis.com/v1/" + name,
			UserAgent: config.UserAgent,
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
