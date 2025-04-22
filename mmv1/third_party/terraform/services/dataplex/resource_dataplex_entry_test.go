package dataplex_test

import (
	//"fmt"
	//"strings"
	"testing"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	//"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	//"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	//transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccDataplexEntry_dataplexEntryUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_number": envvar.GetTestProjectNumberFromEnv(),
		"random_suffix":  acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataplexEntryDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataplexEntry_dataplexEntryFullUpdatePreapre(context),
			},
			{
				ResourceName:            "google_dataplex_entry.test_entry_full",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"aspects", "entry_group_id", "entry_id", "location"},
			},
			{
				Config: testAccDataplexEntry_dataplexEntryUpdate(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_dataplex_entry.test_entry_full", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_dataplex_entry.test_entry_full",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"aspects", "entry_group_id", "entry_id", "location"},
			},
		},
	})
}

func testAccDataplexEntry_dataplexEntryFullUpdatePreapre(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dataplex_aspect_type" "tf-test-aspect-type-full%{random_suffix}-one" {
  aspect_type_id         = "tf-test-aspect-type-full%{random_suffix}-one"
  location     = "us-central1"
  project      = "%{project_number}"

  metadata_template = <<EOF
{
  "name": "tf-test-template",
  "type": "record",
  "recordFields": [
    {
      "name": "type",
      "type": "enum",
      "annotations": {
        "displayName": "Type",
        "description": "Specifies the type of view represented by the entry."
      },
      "index": 1,
      "constraints": {
        "required": true
      },
      "enumValues": [
        {
          "name": "VIEW",
          "index": 1
        }
      ]
    }
  ]
}
EOF
}

resource "google_dataplex_aspect_type" "tf-test-aspect-type-full%{random_suffix}-two" {
  aspect_type_id         = "tf-test-aspect-type-full%{random_suffix}-two"
  location     = "us-central1"
  project      = "%{project_number}"

  metadata_template = <<EOF
{
  "name": "tf-test-template",
  "type": "record",
  "recordFields": [
    {
      "name": "story",
      "type": "enum",
      "annotations": {
        "displayName": "Story",
        "description": "Specifies the story of an entry."
      },
      "index": 1,
      "constraints": {
        "required": true
      },
      "enumValues": [
        {
          "name": "SEQUENCE",
          "index": 1
        },
        {
          "name": "DESERT_ISLAND",
          "index": 2
        }
      ]
    }
  ]
}
EOF
}

resource "google_dataplex_entry_group" "tf-test-entry-group-full%{random_suffix}" {
  entry_group_id = "tf-test-entry-group-full%{random_suffix}"
  project = "%{project_number}"
  location = "us-central1"
}

resource "google_dataplex_entry_type" "tf-test-entry-type-full%{random_suffix}" {
  entry_type_id = "tf-test-entry-type-full%{random_suffix}"
  project = "%{project_number}"
  location = "us-central1"

  labels = { "tag": "test-tf" }
  display_name = "terraform entry type"
  description = "entry type created by Terraform"

  type_aliases = ["TABLE", "DATABASE"]
  platform = "GCS"
  system = "CloudSQL"

  required_aspects {
    type = google_dataplex_aspect_type.tf-test-aspect-type-full%{random_suffix}-one.name
  }
}

resource "google_dataplex_entry" "test_entry_full" {
  entry_group_id = google_dataplex_entry_group.tf-test-entry-group-full%{random_suffix}.entry_group_id
  project = "%{project_number}"
  location = "us-central1"
  entry_id = "tf-test-entry-full%{random_suffix}"
  entry_type = google_dataplex_entry_type.tf-test-entry-type-full%{random_suffix}.name
  fully_qualified_name = "bigquery:%{project_number}.test-dataset"
  parent_entry = "projects/%{project_number}/locations/us-central1/entryGroups/tf-test-entry-group-full%{random_suffix}/entries/some-other-entry"
  entry_source {
    resource = "bigquery:%{project_number}.test-dataset"
    system = "System III"
    platform = "BigQuery"
    display_name = "Human readable name"
    description = "Description from source system"
    labels = {
      "old-label": "old-value"
      "some-label": "some-value"
    }

    ancestors {
      name = "ancestor-one"
      type = "type-one"
    }

    ancestors {
      name = "ancestor-two"
      type = "type-two"
    }

    create_time = "2023-08-03T19:19:00.094Z"
    update_time = "2023-08-03T20:19:00.094Z"
  }

  aspects {
    aspect_key = "%{project_number}.us-central1.tf-test-aspect-type-full%{random_suffix}-one"
    aspect_value {
      data = <<EOF
          {"type": "VIEW"    }
        EOF
    }
  }

  aspects {
    aspect_key = "%{project_number}.us-central1.tf-test-aspect-type-full%{random_suffix}-two"
    aspect_value {
      data = <<EOF
          {"story": "SEQUENCE"    }
        EOF
    }
  }
}
`, context)
}

func testAccDataplexEntry_dataplexEntryUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dataplex_aspect_type" "tf-test-aspect-type-full%{random_suffix}-one" {
  aspect_type_id         = "tf-test-aspect-type-full%{random_suffix}-one"
  location     = "us-central1"
  project      = "%{project_number}"

  metadata_template = <<EOF
{
  "name": "tf-test-template",
  "type": "record",
  "recordFields": [
    {
      "name": "type",
      "type": "enum",
      "annotations": {
        "displayName": "Type",
        "description": "Specifies the type of view represented by the entry."
      },
      "index": 1,
      "constraints": {
        "required": true
      },
      "enumValues": [
        {
          "name": "VIEW",
          "index": 1
        }
      ]
    }
  ]
}
EOF
}

resource "google_dataplex_aspect_type" "tf-test-aspect-type-full%{random_suffix}-two" {
  aspect_type_id         = "tf-test-aspect-type-full%{random_suffix}-two"
  location     = "us-central1"
  project      = "%{project_number}"

  metadata_template = <<EOF
{
  "name": "tf-test-template",
  "type": "record",
  "recordFields": [
    {
      "name": "story",
      "type": "enum",
      "annotations": {
        "displayName": "Story",
        "description": "Specifies the story of an entry."
      },
      "index": 1,
      "constraints": {
        "required": true
      },
      "enumValues": [
        {
          "name": "SEQUENCE",
          "index": 1
        },
        {
          "name": "DESERT_ISLAND",
          "index": 2
        }
      ]
    }
  ]
}
EOF
}

resource "google_dataplex_entry_group" "tf-test-entry-group-full%{random_suffix}" {
  entry_group_id = "tf-test-entry-group-full%{random_suffix}"
  project = "%{project_number}"
  location = "us-central1"
}

resource "google_dataplex_entry_type" "tf-test-entry-type-full%{random_suffix}" {
  entry_type_id = "tf-test-entry-type-full%{random_suffix}"
  project = "%{project_number}"
  location = "us-central1"

  labels = { "tag": "test-tf" }
  display_name = "terraform entry type"
  description = "entry type created by Terraform"

  type_aliases = ["TABLE", "DATABASE"]
  platform = "GCS"
  system = "CloudSQL"

  required_aspects {
    type = google_dataplex_aspect_type.tf-test-aspect-type-full%{random_suffix}-one.name
  }
}

resource "google_dataplex_entry" "test_entry_full" {
  entry_group_id = google_dataplex_entry_group.tf-test-entry-group-full%{random_suffix}.entry_group_id
  project = "%{project_number}"
  location = "us-central1"
  entry_id = "tf-test-entry-full%{random_suffix}"
  entry_type = google_dataplex_entry_type.tf-test-entry-type-full%{random_suffix}.name
  fully_qualified_name = "bigquery:%{project_number}.test-dataset-modified"
  parent_entry = "projects/%{project_number}/locations/us-central1/entryGroups/tf-test-entry-group-full%{random_suffix}/entries/some-other-entry"
  entry_source {
    resource = "bigquery:%{project_number}.test-dataset-modified"
    system = "System III - modified"
    platform = "BigQuery-modified"
    display_name = "Human readable name-modified"
    description = "Description from source system-modified"
    labels = {
      "some-label": "some-value-modified"
      "new-label": "new-value"
    }

    ancestors {
      name = "ancestor-one"
      type = "type-one"
    }

    ancestors {
      name = "ancestor-two"
      type = "type-two"
    }

    create_time = "2024-08-03T19:19:00.094Z"
    update_time = "2024-08-03T20:19:00.094Z"
  }

  aspects {
    aspect_key = "%{project_number}.us-central1.tf-test-aspect-type-full%{random_suffix}-two"
    aspect_value {
      data = <<EOF
          {"story": "DESERT_ISLAND"    }
        EOF
    }
  }
}
`, context)
}

