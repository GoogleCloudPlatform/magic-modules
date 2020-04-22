package google

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"google.golang.org/api/logging/v2"
)

func TestAccLoggingFolderBucketConfig_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(10),
		"folder_name":   "tf-test-" + acctest.RandString(10),
		"org_id":        getTestOrgFromEnv(t),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingFolderBucketConfig_basic(context, 30),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoggingBucketConfig(getFolderBucket, "google_folder.default", "id", "google_logging_folder_bucket_config.basic", 30),
				),
			},
			{
				Config: testAccLoggingFolderBucketConfig_basic(context, 40),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoggingBucketConfig(getFolderBucket, "google_folder.default", "id", "google_logging_folder_bucket_config.basic", 40),
				),
			},
		},
	})
}

func TestAccLoggingProjectBucketConfig_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(10),
		"project_name":  "tf-test-" + acctest.RandString(10),
		"org_id":        getTestOrgFromEnv(t),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingProjectBucketConfig_basic(context, 30),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoggingBucketConfig(getProjectBucket, "google_project.default", "id", "google_logging_project_bucket_config.basic", 30),
				),
			},
			{
				Config: testAccLoggingProjectBucketConfig_basic(context, 40),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoggingBucketConfig(getProjectBucket, "google_project.default", "id", "google_logging_project_bucket_config.basic", 40),
				),
			},
		},
	})
}

func TestAccLoggingBillingAccountBucketConfig_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":        acctest.RandString(10),
		"billing_account_name": "billingAccounts/" + getTestBillingAccountFromEnv(t),
		"org_id":               getTestOrgFromEnv(t),
	}

	fmt.Println(testAccLoggingBillingAccountBucketConfig_basic(context, 30))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingBillingAccountBucketConfig_basic(context, 30),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoggingBucketConfig(getBillingAccountBucket, "data.google_billing_account.default", "name", "google_logging_billing_account_bucket_config.basic", 30),
				),
			},
			{
				Config: testAccLoggingBillingAccountBucketConfig_basic(context, 40),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoggingBucketConfig(getBillingAccountBucket, "data.google_billing_account.default", "name", "google_logging_billing_account_bucket_config.basic", 40),
				),
			},
		},
	})
}

func TestAccLoggingOrganizationBucketConfig_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(10),
		"org_id":        getTestOrgFromEnv(t),
	}

	fmt.Println(testAccLoggingOrganizationBucketConfig_basic(context, 30))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingOrganizationBucketConfig_basic(context, 30),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoggingBucketConfig(getOrganizationBucket, "data.google_organization.default", "name", "google_logging_organization_bucket_config.basic", 30),
				),
			},
			{
				Config: testAccLoggingOrganizationBucketConfig_basic(context, 40),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoggingBucketConfig(getOrganizationBucket, "data.google_organization.default", "name", "google_logging_organization_bucket_config.basic", 40),
				),
			},
		},
	})
}

func testAccLoggingFolderBucketConfig_basic(context map[string]interface{}, retention int) string {
	return fmt.Sprintf(Nprintf(`
resource "google_folder" "default" {
	display_name = "%{folder_name}"
	parent       = "organizations/%{org_id}"
}

resource "google_logging_folder_bucket_config" "basic" {
	folder    = google_folder.default.name
	location  = "global"
	retention_days = %d
	description = "retention test %d days"
	bucket_id = "_Default"
}
`, context), retention, retention)
}

func testAccLoggingProjectBucketConfig_basic(context map[string]interface{}, retention int) string {
	return fmt.Sprintf(Nprintf(`
resource "google_project" "default" {
	project_id = "%{project_name}"
	name       = "%{project_name}"
	org_id     = "%{org_id}"
}

resource "google_logging_project_bucket_config" "basic" {
	project    = google_project.default.name
	location  = "global"
	retention_days = %d
	description = "retention test %d days"
	bucket_id = "_Default"
}
`, context), retention, retention)
}

func testAccLoggingBillingAccountBucketConfig_basic(context map[string]interface{}, retention int) string {
	return fmt.Sprintf(Nprintf(`

data "google_billing_account" "default" {
	billing_account = "%{billing_account_name}"
}

resource "google_logging_billing_account_bucket_config" "basic" {
	billing_account    = data.google_billing_account.default.billing_account
	location  = "global"
	retention_days = %d
	description = "retention test %d days"
	bucket_id = "_Default"
}
`, context), retention, retention)
}

func testAccLoggingOrganizationBucketConfig_basic(context map[string]interface{}, retention int) string {
	return fmt.Sprintf(Nprintf(`
data "google_organization" "default" {
	organization = "%{org_id}"
}

resource "google_logging_organization_bucket_config" "basic" {
	organization    = data.google_organization.default.organization
	location  = "global"
	retention_days = %d
	description = "retention test %d days"
	bucket_id = "_Default"
}
`, context), retention, retention)
}

type bucketConfigFetcher func(string) (*logging.LogBucket, error)

func getFolderBucket(name string) (*logging.LogBucket, error) {
	config := testAccProvider.Meta().(*Config)
	return config.clientLogging.Folders.Locations.Buckets.Get(name).Do()
}

func getProjectBucket(name string) (*logging.LogBucket, error) {
	config := testAccProvider.Meta().(*Config)
	return config.clientLogging.Projects.Locations.Buckets.Get(name).Do()
}

func getBillingAccountBucket(name string) (*logging.LogBucket, error) {
	config := testAccProvider.Meta().(*Config)
	return config.clientLogging.BillingAccounts.Buckets.Get(name).Do()
}

func getOrganizationBucket(name string) (*logging.LogBucket, error) {
	config := testAccProvider.Meta().(*Config)
	return config.clientLogging.Organizations.Locations.Buckets.Get(name).Do()
}

// testAccCheckLoggingBucketConfig is a generic function that will fetch a retention bucket from the SDK and compare it
// with a retention bucket in terraform state. We can do this because each of these parent objects return the same proto
// for the bucket config and the only difference is the url to fetch it at.
// Parents can be folders, organizations, projects or billing accounts.
func testAccCheckLoggingBucketConfig(f bucketConfigFetcher, parent, parentID, bucket string, retention int64) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		pr, ok := s.RootModule().Resources[parent]
		if !ok {
			return fmt.Errorf("Unable to fetch resource %s from state", parent)
		}
		pa := pr.Primary.Attributes
		br, ok := s.RootModule().Resources[bucket]
		if !ok {
			return fmt.Errorf("Unable to fetch resource %s from state", bucket)
		}
		ba := br.Primary.Attributes

		bucket, err := f(pa[parentID] + "/locations/global/buckets/_Default")
		if err != nil {
			return err
		}

		retentionInt, err := strconv.ParseInt(ba["retention_days"], 10, 64)
		if err != nil {
			return err
		}

		if retentionInt != bucket.RetentionDays {
			return fmt.Errorf("retention days in resource didn't match API. resource: %d API: %d", retentionInt, bucket.RetentionDays)
		}

		return nil
	}
}
