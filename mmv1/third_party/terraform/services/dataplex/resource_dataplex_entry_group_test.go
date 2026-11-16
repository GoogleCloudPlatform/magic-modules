package dataplex_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	_ "github.com/hashicorp/terraform-provider-google/google/services/dataplex"
)

func TestAccDataplexEntryGroup_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataplexEntryGroupDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataplexEntryGroup_full(context),
			},
			{
				ResourceName:            "google_dataplex_entry_group.test_entry_group",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "entry_group_id", "labels", "terraform_labels"},
			},
			{
				Config: testAccDataplexEntryGroup_update(context),
			},
			{
				ResourceName:            "google_dataplex_entry_group.test_entry_group",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "entry_group_id", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccDataplexEntryGroup_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dataplex_entry_group" "test_entry_group" {
  entry_group_id = "tf-test-entry-group%{random_suffix}"
  project = "%{project_name}"
  location = "us-central1"
}
`, context)
}

func testAccDataplexEntryGroup_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dataplex_entry_group" "test_entry_group" {
  entry_group_id = "tf-test-entry-group%{random_suffix}"
  project = "%{project_name}"
  location = "us-central1"
	
  labels = { "tag": "test-tf" }
  display_name = "terraform entry group"
  description = "entry group created by Terraform"
}
`, context)
}

func TestAccDataplexEntryGroup_type(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataplexEntryGroupDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataplexEntryGroup_typeContextSet(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_dataplex_entry_group.test_entry_group", "type", "context-set"),
				),
			},
			{
				ResourceName:            "google_dataplex_entry_group.test_entry_group",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "entry_group_id", "labels", "terraform_labels"},
			},
			{
				Config:   testAccDataplexEntryGroup_typeContextSet(context),
				PlanOnly: true,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				Config: testAccDataplexEntryGroup_typeContextSetUpdated(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_dataplex_entry_group.test_entry_group", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_dataplex_entry_group.test_entry_group", "type", "context-set"),
					resource.TestCheckResourceAttr("google_dataplex_entry_group.test_entry_group", "description", "updated context set entry group"),
				),
			},
			{
				Config: testAccDataplexEntryGroup_typeDefault(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_dataplex_entry_group.test_entry_group", plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_dataplex_entry_group.test_entry_group", "type", ""),
				),
			},
			{
				ResourceName:            "google_dataplex_entry_group.test_entry_group",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "entry_group_id", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccDataplexEntryGroup_typeContextSet(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dataplex_entry_group" "test_entry_group" {
  entry_group_id = "tf-test-entry-group%{random_suffix}"
  project        = "%{project_name}"
  location       = "us-central1"
  display_name   = "terraform context set"
  description    = "context set entry group"
  type           = "context-set"
}
`, context)
}

func testAccDataplexEntryGroup_typeContextSetUpdated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dataplex_entry_group" "test_entry_group" {
  entry_group_id = "tf-test-entry-group%{random_suffix}"
  project        = "%{project_name}"
  location       = "us-central1"
  display_name   = "terraform context set"
  description    = "updated context set entry group"
  type           = "context-set"
}
`, context)
}

func testAccDataplexEntryGroup_typeDefault(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dataplex_entry_group" "test_entry_group" {
  entry_group_id = "tf-test-entry-group%{random_suffix}"
  project        = "%{project_name}"
  location       = "us-central1"
  display_name   = "terraform context set"
  description    = "updated context set entry group"
}
`, context)
}
