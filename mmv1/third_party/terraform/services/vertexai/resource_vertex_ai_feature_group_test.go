// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package vertexai_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"testing"
)

func TestAccVertexAIFeatureGroup_vertexAiFeaturegroup_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
		"random_suffix":   acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckVertexAIFeatureGroupDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVertexAIFeatureGroup_vertexAiFeaturegroup_basic(context),
			},
			{
				ResourceName:            "google_vertex_ai_feature_group.test-tf-featuregroup",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "etag", "region", "force_destroy", "labels", "terraform_labels"},
			},
			{
				Config: testAccVertexAIFeatureGroup_vertexAiFeaturegroup_update(context),
			},
			{
				ResourceName:            "google_vertex_ai_feature_group.test-tf-featuregroup",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "etag", "region", "force_destroy", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccVertexAIFeatureGroup_vertexAiFeaturegroup_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_bigquery_dataset" "tf-test-dataset" {

  dataset_id                  = "tf_test_fg_dataset2"
  friendly_name               = "test"
  description                 = "This is a test description"
  location                    = "US"
}

resource "google_bigquery_table" "tf-test-table" {
  deletion_protection = false

  dataset_id = google_bigquery_dataset.tf-test-dataset.dataset_id
  table_id   = "tf_test_fg_table"
  schema = <<EOF
  [
  {
    "name": "entity_id",
    "mode": "NULLABLE",
    "type": "STRING",
    "description": "Test default entity_id"
  },
    {
    "name": "test_entity_column",
    "mode": "NULLABLE",
    "type": "STRING",
    "description": "test secondary entity column"
  },
  {
    "name": "feature_timestamp",
    "mode": "NULLABLE",
    "type": "TIMESTAMP",
    "description": "Default timestamp value"
  }
]
EOF
}

resource "google_vertex_ai_feature_group" "test-tf-featuregroup" {
    name     = "tf_test_test_tf_featuregroup%{random_suffix}"
	description = "test description"
    labels = {
    foo = "bar"
  }
  region   = "us-central1"
  big_query {
    big_query_source {
input_uri = "bq://${google_bigquery_table.tf-test-table.project}.${google_bigquery_table.tf-test-table.dataset_id}.${google_bigquery_table.tf-test-table.table_id}"
    }
	entity_id_columns = ["entity_id"]
  }
}
`, context)
}

func testAccVertexAIFeatureGroup_vertexAiFeaturegroup_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_bigquery_dataset" "tf-test-dataset" {

dataset_id                  = "tf_test_fg_dataset1"
friendly_name               = "test"
description                 = "This is a test description"
location                    = "US"
}

resource "google_bigquery_table" "tf-test-table" {
deletion_protection = false

dataset_id = google_bigquery_dataset.tf-test-dataset.dataset_id
table_id   = "tf_test_fg_table"
schema = <<EOF
[
{
"name": "entity_id",
"mode": "NULLABLE",
"type": "STRING",
"description": "Test default entity_id"
},
{
"name": "test_entity_column",
"mode": "NULLABLE",
"type": "STRING",
"description": "test secondary entity column"
},
{
"name": "feature_timestamp",
"mode": "NULLABLE",
"type": "TIMESTAMP",
"description": "Default timestamp value"
}
]
EOF
}

resource "google_vertex_ai_feature_group" "test-tf-featuregroup" {
name     = "tf_test_test_tf_featuregroup%{random_suffix}"
description = "updated test description"
labels = {
foo = "bar"
}
region   = "us-central1"
big_query {
big_query_source {
input_uri = "bq://${google_bigquery_table.tf-test-table.project}.${google_bigquery_table.tf-test-table.dataset_id}.${google_bigquery_table.tf-test-table.table_id}"
}
entity_id_columns = ["test_entity_column"]
}
}
`, context)
}
