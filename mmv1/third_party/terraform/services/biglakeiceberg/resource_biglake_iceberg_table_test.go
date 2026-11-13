package biglakeiceberg_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	_ "github.com/hashicorp/terraform-provider-google/google/services/biglakeiceberg"
	_ "github.com/hashicorp/terraform-provider-google/google/services/storage"
)

func TestAccBiglakeIcebergIcebergTable_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBiglakeIcebergIcebergTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBiglakeIcebergIcebergTable_updateInitial(context),
			},
			{
				ResourceName:            "google_biglake_iceberg_table.my_iceberg_table",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"catalog", "namespace"},
			},
			{
				Config: testAccBiglakeIcebergIcebergTable_updateUpdated(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_biglake_iceberg_table.my_iceberg_table", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_biglake_iceberg_table.my_iceberg_table",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"catalog", "namespace"},
			},
		},
	})
}

// TestAccBiglakeIcebergIcebergTable_schemaUpdate verifies that a schema change
// reaches the table: a column added in the middle of the schema, renamed,
// renamed back (a schema the table already has) and dropped. Each step is
// followed by an import, which reads the schema back from the API.
func TestAccBiglakeIcebergIcebergTable_schemaUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	added := `
    fields {
      id       = 3
      name     = "note"
      type     = "string"
      required = false
      doc      = "Added after the table was created"
    }`
	renamed := `
    fields {
      id       = 3
      name     = "comment"
      type     = "string"
      required = false
    }`

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBiglakeIcebergIcebergTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBiglakeIcebergIcebergTable_schema(context, ""),
			},
			{
				Config: testAccBiglakeIcebergIcebergTable_schema(context, added),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_biglake_iceberg_table.my_iceberg_table", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_biglake_iceberg_table.my_iceberg_table",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"catalog", "namespace"},
			},
			{
				Config: testAccBiglakeIcebergIcebergTable_schema(context, renamed),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_biglake_iceberg_table.my_iceberg_table", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_biglake_iceberg_table.my_iceberg_table",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"catalog", "namespace"},
			},
			{
				Config: testAccBiglakeIcebergIcebergTable_schema(context, added),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_biglake_iceberg_table.my_iceberg_table", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_biglake_iceberg_table.my_iceberg_table",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"catalog", "namespace"},
			},
			{
				Config: testAccBiglakeIcebergIcebergTable_schema(context, ""),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_biglake_iceberg_table.my_iceberg_table", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_biglake_iceberg_table.my_iceberg_table",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"catalog", "namespace"},
			},
		},
	})
}

// testAccBiglakeIcebergIcebergTable_schema is a two-column table with
// extraField, if any, declared between them.
func testAccBiglakeIcebergIcebergTable_schema(context map[string]interface{}, extraField string) string {
	return acctest.Nprintf(`
resource "google_storage_bucket" "bucket" {
  name          = "tf-test-my-bucket-%{random_suffix}"
  location      = "us-central1"
  force_destroy = true
  uniform_bucket_level_access = true
}

resource "google_biglake_iceberg_catalog" "catalog" {
  name = google_storage_bucket.bucket.name
  catalog_type = "CATALOG_TYPE_GCS_BUCKET"
}

resource "google_biglake_iceberg_namespace" "namespace" {
  catalog = google_biglake_iceberg_catalog.catalog.name
  namespace_id = "my_namespace_%{random_suffix}"
}

resource "google_biglake_iceberg_table" "my_iceberg_table" {
  catalog   = google_biglake_iceberg_catalog.catalog.name
  namespace = google_biglake_iceberg_namespace.namespace.namespace_id
  name      = "tf-test-my_table_%{random_suffix}"
  schema {
    type = "struct"
    fields {
      id       = 1
      name     = "id"
      type     = "long"
      required = true
    }`+extraField+`
    fields {
      id       = 2
      name     = "data"
      type     = "long"
      required = true
    }
  }
}
`, context)
}

// TestAccBiglakeIcebergIcebergTable_vendedCredentials verifies that a table can
// be created in a catalog configured with
// credential_mode = CREDENTIAL_MODE_VENDED_CREDENTIALS. In this mode every table
// operation must send the X-Iceberg-Access-Delegation: vended-credentials header;
// without it the REST catalog rejects the create call.
func TestAccBiglakeIcebergIcebergTable_vendedCredentials(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		CheckDestroy: testAccCheckBiglakeIcebergIcebergTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBiglakeIcebergIcebergTable_vendedCredentials(context),
			},
			{
				ResourceName:            "google_biglake_iceberg_table.my_iceberg_table",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"catalog", "namespace"},
			},
		},
	})
}

func testAccBiglakeIcebergIcebergTable_vendedCredentials(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_bucket" "bucket" {
  name          = "tf-test-my-bucket-%{random_suffix}"
  location      = "us-central1"
  force_destroy = true
  uniform_bucket_level_access = true
}

resource "google_biglake_iceberg_catalog" "catalog" {
  name = google_storage_bucket.bucket.name
  catalog_type = "CATALOG_TYPE_GCS_BUCKET"
  credential_mode = "CREDENTIAL_MODE_VENDED_CREDENTIALS"
}

# The catalog's generated service account is not immediately usable in an IAM
# binding right after the catalog is created, so wait for it to propagate before
# referencing it.
resource "time_sleep" "wait_for_biglake_service_account" {
  depends_on      = [google_biglake_iceberg_catalog.catalog]
  create_duration = "60s"
}

# Credential vending downscopes a generated service account to access the GCS
# bucket, so it must be granted access before any table operation runs.
resource "google_storage_bucket_iam_member" "cv_sa_storage_admin" {
  bucket = google_storage_bucket.bucket.name
  role   = "roles/storage.admin"
  member = "serviceAccount:${google_biglake_iceberg_catalog.catalog.biglake_service_account}"

  depends_on = [time_sleep.wait_for_biglake_service_account]
}

resource "google_biglake_iceberg_namespace" "namespace" {
  catalog = google_biglake_iceberg_catalog.catalog.name
  namespace_id = "my_namespace_%{random_suffix}"
}

resource "google_biglake_iceberg_table" "my_iceberg_table" {
  catalog   = google_biglake_iceberg_catalog.catalog.name
  namespace = google_biglake_iceberg_namespace.namespace.namespace_id
  name      = "tf-test-my_table_%{random_suffix}"
  schema {
    type = "struct"
    fields {
      id       = 1
      name     = "id"
      type     = "long"
      required = true
    }
  }

  properties = {
    key = "initial"
  }

  depends_on = [google_storage_bucket_iam_member.cv_sa_storage_admin]
}
`, context)
}

func testAccBiglakeIcebergIcebergTable_updateInitial(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_bucket" "bucket" {
  name          = "my-bucket-%{random_suffix}"
  location      = "us-central1"
  force_destroy = true
  uniform_bucket_level_access = true
}

resource "google_biglake_iceberg_catalog" "catalog" {
  name = google_storage_bucket.bucket.name
  catalog_type = "CATALOG_TYPE_GCS_BUCKET"
}

resource "google_biglake_iceberg_namespace" "namespace" {
  catalog = google_biglake_iceberg_catalog.catalog.name
  namespace_id = "my_namespace_%{random_suffix}"
}

resource "google_biglake_iceberg_table" "my_iceberg_table" {
  catalog   = google_biglake_iceberg_catalog.catalog.name
  namespace = google_biglake_iceberg_namespace.namespace.namespace_id
  name      = "my_table_%{random_suffix}"
  schema {
    type = "struct"
    fields {
      id       = 1
      name     = "id"
      type     = "long"
      required = true
    }
    fields {
      id       = 2
      name     = "data"
      type     = "long"
      required = true
    }
  }
  sort_order {
    fields {
      source_id  = 1
      transform  = "identity"
      direction  = "asc"
      null_order = "nulls-first"
    }
  }

  properties = {
    key = "initial"
  }
}
`, context)
}

func testAccBiglakeIcebergIcebergTable_updateUpdated(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_bucket" "bucket" {
  name          = "my-bucket-%{random_suffix}"
  location      = "us-central1"
  force_destroy = true
  uniform_bucket_level_access = true
}

resource "google_biglake_iceberg_catalog" "catalog" {
  name = google_storage_bucket.bucket.name
  catalog_type = "CATALOG_TYPE_GCS_BUCKET"
}

resource "google_biglake_iceberg_namespace" "namespace" {
  catalog = google_biglake_iceberg_catalog.catalog.name
  namespace_id = "my_namespace_%{random_suffix}"
}

resource "google_biglake_iceberg_table" "my_iceberg_table" {
  catalog   = google_biglake_iceberg_catalog.catalog.name
  namespace = google_biglake_iceberg_namespace.namespace.namespace_id
  name      = "my_table_%{random_suffix}"
  schema {
    type = "struct"
    fields {
      id       = 1
      name     = "id"
      type     = "long"
      required = true
    }
    fields {
      id       = 2
      name     = "data"
      type     = "long"
      required = true
    }
  }
  sort_order {
    fields {
      source_id  = 2
      transform  = "bucket[4]"
      direction  = "desc"
      null_order = "nulls-last"
    }
  }

  properties = {
    key = "updated"
  }
}
`, context)
}
