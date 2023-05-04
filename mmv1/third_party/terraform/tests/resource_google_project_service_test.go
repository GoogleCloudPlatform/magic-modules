package google

import (
	"fmt"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestProjectServiceServiceValidateFunc(t *testing.T) {
	cases := map[string]struct {
		val                   interface{}
		ExpectValidationError bool
	}{
		"ignoredProjectService": {
			val:                   "dataproc-control.googleapis.com",
			ExpectValidationError: true,
		},
		"bannedProjectService": {
			val:                   "bigquery-json.googleapis.com",
			ExpectValidationError: true,
		},
		"third party API": {
			val:                   "whatever.example.com",
			ExpectValidationError: false,
		},
		"not a domain": {
			val:                   "monitoring",
			ExpectValidationError: true,
		},
		"not a string": {
			val:                   5,
			ExpectValidationError: true,
		},
	}

	for tn, tc := range cases {
		_, errs := validateProjectServiceService(tc.val, "service")
		if tc.ExpectValidationError && len(errs) == 0 {
			t.Errorf("bad: %s, %q passed validation but was expected to fail", tn, tc.val)
		}
		if !tc.ExpectValidationError && len(errs) > 0 {
			t.Errorf("bad: %s, %q failed validation but was expected to pass. errs: %q", tn, tc.val, errs)
		}
	}
}

// Test that services can be enabled and disabled on a project
func TestAccProjectService_basic(t *testing.T) {
	t.Parallel()
	// Multiple fine-grained resources
	acctest.SkipIfVcr(t)

	org := acctest.GetTestOrgFromEnv(t)
	pid := fmt.Sprintf("tf-test-%d", RandInt(t))
	services := []string{"iam.googleapis.com", "cloudresourcemanager.googleapis.com"}
	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccProjectService_basic(services, pid, org),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProjectService(t, services, pid, true),
				),
			},
			{
				ResourceName:            "google_project_service.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"disable_on_destroy"},
			},
			{
				ResourceName:            "google_project_service.test2",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"disable_on_destroy", "project"},
			},
			// Use a separate TestStep rather than a CheckDestroy because we need the project to still exist.
			{
				Config: testAccProject_create(pid, org),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProjectService(t, services, pid, false),
				),
			},
			// Create services with disabling turned off.
			{
				Config: testAccProjectService_noDisable(services, pid, org),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProjectService(t, services, pid, true),
				),
			},
			// Check that services are still enabled even after the resources are deleted.
			{
				Config: testAccProject_create(pid, org),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProjectService(t, services, pid, true),
				),
			},
		},
	})
}

func TestAccProjectService_disableDependentServices(t *testing.T) {
	// Multiple fine-grained resources
	acctest.SkipIfVcr(t)
	t.Parallel()

	org := acctest.GetTestOrgFromEnv(t)
	billingId := acctest.GetTestBillingAccountFromEnv(t)
	pid := fmt.Sprintf("tf-test-%d", RandInt(t))
	services := []string{"cloudbuild.googleapis.com", "containerregistry.googleapis.com"}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccProjectService_disableDependentServices(services, pid, org, billingId, "false"),
			},
			{
				ResourceName:            "google_project_service.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"disable_on_destroy"},
			},
			{
				Config:      testAccProjectService_dependencyRemoved(services, pid, org, billingId),
				ExpectError: regexp.MustCompile("Please specify disable_dependent_services=true if you want to proceed with disabling all services."),
			},
			{
				Config: testAccProjectService_disableDependentServices(services, pid, org, billingId, "true"),
			},
			{
				ResourceName:            "google_project_service.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"disable_on_destroy"},
			},
			{
				Config:             testAccProjectService_dependencyRemoved(services, pid, org, billingId),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccProjectService_handleNotFound(t *testing.T) {
	t.Parallel()

	org := acctest.GetTestOrgFromEnv(t)
	pid := fmt.Sprintf("tf-test-%d", RandInt(t))
	service := "iam.googleapis.com"
	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccProjectService_handleNotFound(service, pid, org),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProjectService(t, []string{service}, pid, true),
				),
			},
			// Delete the project, implicitly deletes service, expect the plan to want to create the service again
			{
				Config:             testAccProjectService_handleNotFoundNoProject(service, pid),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccProjectService_renamedService(t *testing.T) {
	t.Parallel()

	if len(renamedServices) == 0 {
		t.Skip()
	}

	var newName string
	for _, new := range renamedServices {
		newName = new
	}

	org := acctest.GetTestOrgFromEnv(t)
	pid := fmt.Sprintf("tf-test-%d", RandInt(t))
	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccProjectService_single(newName, pid, org),
			},
			{
				ResourceName:            "google_project_service.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"disable_on_destroy", "disable_dependent_services"},
			},
		},
	})
}

func testAccCheckProjectService(t *testing.T, services []string, pid string, expectEnabled bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := GoogleProviderConfig(t)
		currentlyEnabled, err := ListCurrentlyEnabledServices(pid, "", config.UserAgent, config, time.Minute*10)
		if err != nil {
			return fmt.Errorf("Error listing services for project %q: %v", pid, err)
		}

		for _, expected := range services {
			exists := false
			for actual := range currentlyEnabled {
				if expected == actual {
					exists = true
				}
			}
			if expectEnabled && !exists {
				return fmt.Errorf("Expected service %s is not enabled server-side", expected)
			}
			if !expectEnabled && exists {
				return fmt.Errorf("Expected disabled service %s is enabled server-side", expected)
			}
		}

		return nil
	}
}

func testAccProjectService_basic(services []string, pid, org string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  project_id = "%s"
  name       = "%s"
  org_id     = "%s"
}

resource "google_project_service" "test" {
  project = google_project.acceptance.project_id
  service = "%s"
}

resource "google_project_service" "test2" {
  project = google_project.acceptance.id
  service = "%s"
}
`, pid, pid, org, services[0], services[1])
}

func testAccProjectService_disableDependentServices(services []string, pid, org, billing, disableDependentServices string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  project_id      = "%s"
  name            = "%s"
  org_id          = "%s"
  billing_account = "%s"
}

resource "google_project_service" "test" {
  project = google_project.acceptance.project_id
  service = "%s"
}

resource "google_project_service" "test2" {
  project                    = google_project.acceptance.project_id
  service                    = "%s"
  disable_dependent_services = %s
}
`, pid, pid, org, billing, services[0], services[1], disableDependentServices)
}

func testAccProjectService_dependencyRemoved(services []string, pid, org, billing string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  project_id      = "%s"
  name            = "%s"
  org_id          = "%s"
  billing_account = "%s"
}

resource "google_project_service" "test" {
  project = google_project.acceptance.project_id
  service = "%s"
}
`, pid, pid, org, billing, services[0])
}

func testAccProjectService_noDisable(services []string, pid, org string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  project_id = "%s"
  name       = "%s"
  org_id     = "%s"
}

resource "google_project_service" "test" {
  project            = google_project.acceptance.project_id
  service            = "%s"
  disable_on_destroy = false
}

resource "google_project_service" "test2" {
  project            = google_project.acceptance.project_id
  service            = "%s"
  disable_on_destroy = false
}
`, pid, pid, org, services[0], services[1])
}

func testAccProjectService_handleNotFound(service, pid, org string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  project_id = "%s"
  name       = "%s"
  org_id     = "%s"
}

// by passing through locals, we break the dependency chain
// see terraform-provider-google#1292
locals {
  project_id = google_project.acceptance.project_id
}

resource "google_project_service" "test" {
  project = local.project_id
  service = "%s"
}
`, pid, pid, org, service)
}

func testAccProjectService_handleNotFoundNoProject(service, pid string) string {
	return fmt.Sprintf(`
resource "google_project_service" "test" {
  project = "%s"
  service = "%s"
}
`, pid, service)
}

func testAccProjectService_single(service string, pid, org string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  project_id = "%s"
  name       = "%s"
  org_id     = "%s"
}

resource "google_project_service" "test" {
  project = google_project.acceptance.project_id
  service = "%s"

  disable_dependent_services = true
}
`, pid, pid, org, service)
}
