package dataplex_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataplexEntryLink_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_number": envvar.GetTestProjectNumberFromEnv(),
		"random_suffix":  acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataplexEntryLink_dataplexEntryLinkUpdate(context),
			},
			{
				ResourceName:            "google_dataplex_entry_link.basic_entry_link",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"entry_group_id", "entry_link_id", "location"},
			},
		},
	})
}

func testAccDataplexEntryLink_dataplexEntryLinkUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dataplex_entry_group" "entry-group-basic" {
  location = "us-central1"
  entry_group_id = "tf-test-entry-group%{random_suffix}"
  project = "%{project_number}"
}
resource "google_dataplex_entry" "source" {
  location = "us-central1"
  entry_group_id = google_dataplex_entry_group.entry-group-basic.entry_group_id
  entry_id = "tf-test-source-entry%{random_suffix}"
  entry_type = google_dataplex_entry_type.entry-type-basic.name
  project = "%{project_number}"
}
resource "google_dataplex_entry_type" "entry-type-basic" {
  entry_type_id = "tf-test-entry-type%{random_suffix}"
  location = "us-central1"
  project = "%{project_number}"
}
resource "google_dataplex_entry" "target" {
  location = "us-central1"
  entry_group_id = google_dataplex_entry_group.entry-group-basic.entry_group_id
  entry_id = "tf-test-target-entry%{random_suffix}"
  entry_type = google_dataplex_entry_type.entry-type-basic.name
  project = "%{project_number}"
}
resource "google_dataplex_entry_link" "basic_entry_link" {
  project = "%{project_number}"
  location = "us-central1"
  entry_group_id = google_dataplex_entry_group.entry-group-basic.entry_group_id
  entry_link_id = "tf-test-entry-link%{random_suffix}"
  entry_link_type = "projects/655216118709/locations/global/entryLinkTypes/related"
  entry_references {
    name = google_dataplex_entry.source.name
  }
  entry_references {
    name = google_dataplex_entry.target.name
  }
}
`, context)
}
