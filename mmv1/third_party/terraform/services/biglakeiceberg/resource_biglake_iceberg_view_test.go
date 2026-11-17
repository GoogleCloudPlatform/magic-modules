package biglakeiceberg_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	_ "github.com/hashicorp/terraform-provider-google/google/services/biglakeiceberg"
	_ "github.com/hashicorp/terraform-provider-google/google/services/storage"
)

func TestAccBiglakeIcebergIcebergView_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBiglakeIcebergIcebergViewDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBiglakeIcebergIcebergView_updateInitial(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_biglake_iceberg_view.my_iceberg_view", "view_uuid"),
					resource.TestCheckResourceAttrSet("google_biglake_iceberg_view.my_iceberg_view", "view_version.0.version_id"),
					resource.TestCheckResourceAttrSet("google_biglake_iceberg_view.my_iceberg_view", "view_version.0.timestamp_ms"),
				),
			},
			{
				ResourceName:            "google_biglake_iceberg_view.my_iceberg_view",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"catalog", "namespace"},
			},
			{
				Config: testAccBiglakeIcebergIcebergView_updateUpdated(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_biglake_iceberg_view.my_iceberg_view", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_biglake_iceberg_view.my_iceberg_view",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"catalog", "namespace"},
			},
			{
				Config: testAccBiglakeIcebergIcebergView_updateDropDialect(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_biglake_iceberg_view.my_iceberg_view", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_biglake_iceberg_view.my_iceberg_view",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"catalog", "namespace"},
			},
		},
	})
}

func testAccBiglakeIcebergIcebergView_updateInitial(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_bucket" "bucket" {
  name                        = "tf-test-my-bucket-%{random_suffix}"
  location                    = "us-central1"
  force_destroy               = true
  uniform_bucket_level_access = true
}

resource "google_biglake_iceberg_catalog" "catalog" {
  name         = google_storage_bucket.bucket.name
  catalog_type = "CATALOG_TYPE_GCS_BUCKET"
}

resource "google_biglake_iceberg_namespace" "namespace" {
  catalog      = google_biglake_iceberg_catalog.catalog.name
  namespace_id = "tf_test_my_namespace_%{random_suffix}"
}

resource "google_biglake_iceberg_view" "my_iceberg_view" {
  catalog   = google_biglake_iceberg_catalog.catalog.name
  namespace = google_biglake_iceberg_namespace.namespace.namespace_id
  name      = "tf_test_my_view_%{random_suffix}"

  schema {
    type = "struct"
    fields {
      id       = 1
      name     = "id"
      type     = "long"
      required = true
      doc      = "Initial doc"
    }
    fields {
      id   = 2
      name = "tags"
      type = jsonencode({
        type               = "list"
        "element-id"       = 3
        element            = "string"
        "element-required" = false
      })
      required = false
    }
  }

  view_version {
    default_catalog   = google_biglake_iceberg_catalog.catalog.name
    default_namespace = [google_biglake_iceberg_namespace.namespace.namespace_id]
    summary = {
      "engine-name" = "spark"
    }
    representations {
      type    = "sql"
      sql     = "SELECT id, tags FROM base_table"
      dialect = "spark"
    }
  }

  properties = {
    key       = "initial"
    removable = "temp"
  }
}
`, context)
}

func testAccBiglakeIcebergIcebergView_updateUpdated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_bucket" "bucket" {
  name                        = "tf-test-my-bucket-%{random_suffix}"
  location                    = "us-central1"
  force_destroy               = true
  uniform_bucket_level_access = true
}

resource "google_biglake_iceberg_catalog" "catalog" {
  name         = google_storage_bucket.bucket.name
  catalog_type = "CATALOG_TYPE_GCS_BUCKET"
}

resource "google_biglake_iceberg_namespace" "namespace" {
  catalog      = google_biglake_iceberg_catalog.catalog.name
  namespace_id = "tf_test_my_namespace_%{random_suffix}"
}

resource "google_biglake_iceberg_view" "my_iceberg_view" {
  catalog   = google_biglake_iceberg_catalog.catalog.name
  namespace = google_biglake_iceberg_namespace.namespace.namespace_id
  name      = "tf_test_my_view_%{random_suffix}"

  schema {
    type = "struct"
    fields {
      id       = 1
      name     = "id"
      type     = "long"
      required = true
    }
    fields {
      id   = 2
      name = "tags"
      type = jsonencode({
        type               = "list"
        "element-id"       = 3
        element            = "string"
        "element-required" = false
      })
      required = false
    }
    fields {
      id       = 4
      name     = "label"
      type     = "string"
      required = false
    }
  }

  view_version {
    default_namespace = [google_biglake_iceberg_namespace.namespace.namespace_id]
    representations {
      type    = "sql"
      sql     = "SELECT id, tags, label FROM base_table"
      dialect = "spark"
    }
    representations {
      type    = "sql"
      sql     = "SELECT id, tags, label FROM base_table"
      dialect = "trino"
    }
  }

  properties = {
    key                                = "updated"
    "write.metadata.compression-codec" = "none"
  }
}
`, context)
}

func testAccBiglakeIcebergIcebergView_updateDropDialect(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_bucket" "bucket" {
  name                        = "tf-test-my-bucket-%{random_suffix}"
  location                    = "us-central1"
  force_destroy               = true
  uniform_bucket_level_access = true
}

resource "google_biglake_iceberg_catalog" "catalog" {
  name         = google_storage_bucket.bucket.name
  catalog_type = "CATALOG_TYPE_GCS_BUCKET"
}

resource "google_biglake_iceberg_namespace" "namespace" {
  catalog      = google_biglake_iceberg_catalog.catalog.name
  namespace_id = "tf_test_my_namespace_%{random_suffix}"
}

resource "google_biglake_iceberg_view" "my_iceberg_view" {
  catalog   = google_biglake_iceberg_catalog.catalog.name
  namespace = google_biglake_iceberg_namespace.namespace.namespace_id
  name      = "tf_test_my_view_%{random_suffix}"

  schema {
    type = "struct"
    fields {
      id       = 1
      name     = "id"
      type     = "long"
      required = true
    }
    fields {
      id   = 2
      name = "tags"
      type = jsonencode({
        type               = "list"
        "element-id"       = 3
        element            = "string"
        "element-required" = false
      })
      required = false
    }
    fields {
      id       = 4
      name     = "label"
      type     = "string"
      required = false
    }
  }

  view_version {
    default_namespace = [google_biglake_iceberg_namespace.namespace.namespace_id]
    representations {
      type    = "sql"
      sql     = "SELECT id, tags, label FROM base_table"
      dialect = "spark"
    }
  }

  properties = {
    key                                = "updated"
    "write.metadata.compression-codec" = "none"
    "replace.drop-dialect.allowed"     = "true"
  }
}
`, context)
}
