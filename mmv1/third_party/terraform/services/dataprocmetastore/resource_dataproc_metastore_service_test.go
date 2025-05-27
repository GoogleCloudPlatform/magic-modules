package dataprocmetastore_test

import (
	"fmt"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataprocMetastoreService_updateAndImport(t *testing.T) {
	t.Parallel()

	name := "tf-test-metastore-" + acctest.RandString(t, 10)
	tier := [2]string{"DEVELOPER", "ENTERPRISE"}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocMetastoreService_updateAndImport(name, tier[0]),
			},
			{
				ResourceName:      "google_dataproc_metastore_service.my_metastore",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataprocMetastoreService_updateAndImport(name, tier[1]),
			},
			{
				ResourceName:      "google_dataproc_metastore_service.my_metastore",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDataprocMetastoreService_updateAndImport(name, tier string) string {
	return fmt.Sprintf(`
resource "google_dataproc_metastore_service" "my_metastore" {
	service_id = "%s"
	location   = "us-central1"
	tier       = "%s"

	hive_metastore_config {
		version = "2.3.6"
	}
}
`, name, tier)
}

func TestAccDataprocMetastoreService_dataprocMetastoreServiceScheduledBackupExampleUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataprocMetastoreServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocMetastoreService_dataprocMetastoreServiceScheduledBackupExample(context),
			},
			{
				ResourceName:            "google_dataproc_metastore_service.backup",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"service_id", "location", "labels", "terraform_labels"},
			},
			{
				Config: testAccDataprocMetastoreService_dataprocMetastoreServiceScheduledBackupExampleUpdate(context),
			},
		},
	})
}

func TestAccDataprocMetastoreService_PrivateServiceConnect(t *testing.T) {
	t.Skip("Skipping due to https://github.com/hashicorp/terraform-provider-google/issues/13710")
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataprocMetastoreServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocMetastoreService_PrivateServiceConnect(context),
			},
			{
				ResourceName:            "google_dataproc_metastore_service.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"service_id", "location"},
			},
		},
	})
}

func testAccDataprocMetastoreService_PrivateServiceConnect(context map[string]interface{}) string {
	return acctest.Nprintf(`
// Use data source instead of creating a subnetwork due to a bug on API side.
// With the bug, the new created subnetwork cannot be deleted when deleting the dataproc metastore service.
data "google_compute_subnetwork" "subnet" {
  name   = "default"
  region = "us-central1"
}

resource "google_dataproc_metastore_service" "default" {
  service_id = "tf-test-metastore-srv%{random_suffix}"
  location   = "us-central1"

  hive_metastore_config {
    version = "3.1.2"
  }

  network_config {
    consumers {
      subnetwork = data.google_compute_subnetwork.subnet.id
    }
  }
}
`, context)
}

func testAccDataprocMetastoreService_dataprocMetastoreServiceScheduledBackupExampleUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dataproc_metastore_service" "backup" {
  service_id = "tf-test-backup%{random_suffix}"
  location   = "us-central1"
  port       = 9080
  tier       = "DEVELOPER"

  maintenance_window {
    hour_of_day = 2
    day_of_week = "SUNDAY"
  }

  hive_metastore_config {
    version = "2.3.6"
  }

  scheduled_backup {
    enabled         = true
    cron_schedule   = "0 0 * * 0"
    time_zone       = "America/Los_Angeles"
    backup_location = "gs://${google_storage_bucket.bucket.name}"
  }

  labels = {
    env = "test"
  }
}

resource "google_storage_bucket" "bucket" {
  name     = "tf-test-backup%{random_suffix}"
  location = "us-central1"
}
`, context)
}

func TestAccMetastoreService_tags(t *testing.T) {
	t.Parallel()

	// Bootstrap the new Tag Key and two distinct Tag Values
	tagKeyURI := acctest.BootstrapSharedTestTagKey(t, "metastore-org-policy-tagkey")
	allowedTagValueURI := acctest.BootstrapSharedTestTagValue(t, "metastore-org-policy-allowed-value", tagKeyURI)
	disallowedTagValueURI := acctest.BootstrapSharedTestTagValue(t, "metastore-org-policy-disallowed-value", tagKeyURI)

	allowedTagValueCanonicalName := strings.TrimPrefix(allowedTagValueURI, "//cloudresourcemanager.googleapis.com/")
	tagKeyShortName := strings.Split(tagKeyURI, "/")[len(strings.Split(tagKeyURI, "/"))-1]

	orgID := envvar.GetTestOrgFromEnv(t)
	if orgID == "" {
		orgID = "735183260412" // Default value if env var is not set
		t.Logf("GOOGLE_ORGANIZATION environment variable not set, using default org id: %s", orgID)
	}
	contextData := map[string]interface{}{
		"random_suffix":                acctest.RandString(t, 10),
		"org_id":                       orgID,
		"project":                      "tags-blr-test-ap",
		"location":                     "us-central1",
		"tagKey":                       tagKeyShortName,
		"allowedTagValue":              allowedTagValueURI,
		"disallowedTagValue":           disallowedTagValueURI,
		"allowedTagValueCanonicalName": allowedTagValueCanonicalName,
		"tagKeyURI":                    tagKeyURI,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckDataprocMetastoreServiceDestroyProducer(t),
		),
		Steps: []resource.TestStep{
			// Step 0: Define the custom constraint
			{
				Config: testAccCustomConstraint(contextData),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_org_policy_custom_constraint.constraint",
						"display_name",
						"Metastore Service Custom Constraint",
					),
					resource.TestCheckResourceAttr(
						"google_org_policy_custom_constraint.constraint",
						"condition",
						"resource.management.autoUpgrade == false",
					),
				),
			},
			// Step 1: Define the Organization Policy with rules inside spec and a condition
			{
				Config: testAccOrgPolicy(contextData),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_org_policy_policy.primary",
						"constraint",
						fmt.Sprintf("custom.metastoreService%s", contextData["random_suffix"]),
					),
					resource.TestCheckResourceAttr(
						"google_org_policy_policy.primary",
						"spec.0.rules.0.enforce",
						"false",
					),
				),
			},
			// Step 2: Create a Metastore Service with the ALLOWED tag (should succeed)
			{
				Config: testAccMetastoreServiceTags(contextData),
				Check: resource.ComposeTestCheckFunc(
					// Verify the tag is present in the Terraform state
					resource.TestCheckResourceAttr(
						"google_dataproc_metastore_service.default",
						"tags.%",
						"1",
					),
					resource.TestCheckResourceAttr(
						"google_dataproc_metastore_service.default",
						"tags."+contextData["org"].(string)+"/"+contextData["tagKey"].(string),
						contextData["allowedTagValue"].(string),
					),
				),
			},
			// Step 3: Attempt to create a Metastore Service with the DISALLOWED tag (should fail)
			{
				Config: testAccMetastoreServiceTagsDisallowed(contextData),
				// Expect an error indicating policy violation
				ExpectError: regexp.MustCompile(`(?i)(CONSTRAINT_VIOLATION|policy violation|denied by a policy)`),
			},
			// Step 4: Delete the Organization Policy (added step)
			{
				Config: testAccOrgPolicyDelete(contextData),
			},
		},
	})
}
func testAccMetastoreServiceTags(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dataproc_metastore_service" "default" {
  service_id   = "tf-test-my-service-%{random_suffix}"
  location   = "us-central1"
  port       = 9080
  tier       = "DEVELOPER"

  maintenance_window {
    hour_of_day = 2
    day_of_week = "SUNDAY"
   }

  hive_metastore_config {
    version = "2.3.6"
  }

  labels = {
    env = "test"
  }
  tags = {
    "%{org_id}/%{tagKey}" = "%{tagValue}"
  }
}
`, context)
}

// Metastore service with a disallowed tag value (for negative test)
func testAccMetastoreServiceTagsDisallowed(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dataproc_metastore_service" "disallowed_service" {
  service_id   = "tf-test-disallowed-%{random_suffix}" # Unique ID for this
  location   = "%{location}"
  port       = 9080
  tier       = "DEVELOPER"

  maintenance_window {
    hour_of_day = 2
    day_of_week = "SUNDAY"
   }

  hive_metastore_config {
    version = "2.3.6"
  }

  labels = {
    env = "test"
  }
  tags = {
    "%{org_id}/%{tagKey}" = "%{disallowedTagValue}" # Using the disallowed tag value
  }
}
`, context)
}

func testAccOrgPolicy(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_org_policy_policy" "primary" {
  name = "organizations/%{org_id}/policies/${google_org_policy_custom_constraint.constraint.name}"
  parent = "organizations/%{org_id}"
  spec {
      rules {
          enforce = "FALSE"
      }
  }
}
`, context)
}

func testAccCustomConstraint(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_org_policy_custom_constraint" "constraint" {
    name         = "custom.metastoreService%{random_suffix}"
    parent       = "organizations/%{org_id}"
    display_name = "Metastore Service Custom Constraint"
    description  = "Only allow Metastore resource to be created or updated if AutoUpgrade is not enabled where this custom constraint is enforced."
  
    action_type    = "ALLOW"
    condition      = "resource.management.autoUpgrade == false"
    method_types   = ["CREATE"]
    resource_types = ["metastore.googleapis.com/Service"]
  }
`, context)
}

// Configuration to delete the Organization Policy
func testAccOrgPolicyDelete(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_org_policy_policy" "primary" {
  org_id = "%{org_id}"
  constraint = "serviceuser.services"
  policy_type = "unset"
}
`, context)
}

