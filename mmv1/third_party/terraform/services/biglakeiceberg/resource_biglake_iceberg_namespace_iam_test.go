// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package biglakeiceberg_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccBiglakeIcebergIcebergNamespaceIamBinding(t *testing.T) {
	t.Parallel()
	acctest.SkipIfVcr(t)

	suffix := acctest.RandString(t, 10)
	catalogId := fmt.Sprintf("tf-test-catalog-%s", suffix)
	namespaceId := "tf_test_namespace"
	role := "roles/viewer"
	project := envvar.GetTestProjectFromEnv()

	context := map[string]interface{}{
		"catalog_id":    catalogId,
		"namespace_id":  namespaceId,
		"role":          role,
		"random_suffix": suffix,
	}

	importId := fmt.Sprintf("projects/%s/catalogs/%s/namespaces/%s", project, catalogId, namespaceId)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBiglakeIcebergIcebergNamespace_catalogOnly(context),
			},
			{
				PreConfig: func() {
					acctest.BootstrapBigLakeIcebergNamespace(t, catalogId, namespaceId)
				},
				Config: testAccBiglakeIcebergIcebergNamespaceIamBinding_basic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_biglake_iceberg_namespace_iam_binding.binding", "role", role),
				),
			},
			{
				ResourceName:      "google_biglake_iceberg_namespace_iam_binding.binding",
				ImportStateId:     fmt.Sprintf("%s %s", importId, role),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test Iam Binding update
				Config: testAccBiglakeIcebergIcebergNamespaceIamBinding_update(context),
			},
			{
				ResourceName:      "google_biglake_iceberg_namespace_iam_binding.binding",
				ImportStateId:     fmt.Sprintf("%s %s", importId, role),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBiglakeIcebergIcebergNamespace_catalogOnly(context),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						acctest.CleanupBigLakeIcebergNamespace(t, catalogId, namespaceId)
						return nil
					},
				),
			},
		},
	})

	t.Cleanup(func() {
		acctest.CleanupBigLakeIcebergNamespace(t, catalogId, namespaceId)
	})
}

func TestAccBiglakeIcebergIcebergNamespaceIamMember(t *testing.T) {
	t.Parallel()
	acctest.SkipIfVcr(t)

	suffix := acctest.RandString(t, 10)
	catalogId := fmt.Sprintf("tf-test-catalog-%s", suffix)
	namespaceId := "tf_test_namespace"
	role := "roles/viewer"
	project := envvar.GetTestProjectFromEnv()

	context := map[string]interface{}{
		"catalog_id":    catalogId,
		"namespace_id":  namespaceId,
		"role":          role,
		"random_suffix": suffix,
	}

	importId := fmt.Sprintf("projects/%s/catalogs/%s/namespaces/%s", project, catalogId, namespaceId)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBiglakeIcebergIcebergNamespace_catalogOnly(context),
			},
			{
				PreConfig: func() {
					acctest.BootstrapBigLakeIcebergNamespace(t, catalogId, namespaceId)
				},
				Config: testAccBiglakeIcebergIcebergNamespaceIamMember_basic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_biglake_iceberg_namespace_iam_member.member", "role", role),
					resource.TestCheckResourceAttr("google_biglake_iceberg_namespace_iam_member.member", "member", "user:admin@hashicorptest.com"),
				),
			},
			{
				ResourceName:      "google_biglake_iceberg_namespace_iam_member.member",
				ImportStateId:     fmt.Sprintf("%s %s user:admin@hashicorptest.com", importId, role),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBiglakeIcebergIcebergNamespace_catalogOnly(context),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						acctest.CleanupBigLakeIcebergNamespace(t, catalogId, namespaceId)
						return nil
					},
				),
			},
		},
	})

	t.Cleanup(func() {
		acctest.CleanupBigLakeIcebergNamespace(t, catalogId, namespaceId)
	})
}

func TestAccBiglakeIcebergIcebergNamespaceIamPolicy(t *testing.T) {
	t.Parallel()
	acctest.SkipIfVcr(t)

	suffix := acctest.RandString(t, 10)
	catalogId := fmt.Sprintf("tf-test-catalog-%s", suffix)
	namespaceId := "tf_test_namespace"
	role := "roles/viewer"
	project := envvar.GetTestProjectFromEnv()

	context := map[string]interface{}{
		"catalog_id":    catalogId,
		"namespace_id":  namespaceId,
		"role":          role,
		"random_suffix": suffix,
	}

	importId := fmt.Sprintf("projects/%s/catalogs/%s/namespaces/%s", project, catalogId, namespaceId)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBiglakeIcebergIcebergNamespace_catalogOnly(context),
			},
			{
				PreConfig: func() {
					acctest.BootstrapBigLakeIcebergNamespace(t, catalogId, namespaceId)
				},
				Config: testAccBiglakeIcebergIcebergNamespaceIamPolicy_basic(context),
				Check:  resource.TestCheckResourceAttrSet("google_biglake_iceberg_namespace_iam_policy.policy", "policy_data"),
			},
			{
				ResourceName:      "google_biglake_iceberg_namespace_iam_policy.policy",
				ImportStateId:     importId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBiglakeIcebergIcebergNamespaceIamPolicy_emptyBinding(context),
			},
			{
				ResourceName:      "google_biglake_iceberg_namespace_iam_policy.policy",
				ImportStateId:     importId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBiglakeIcebergIcebergNamespace_catalogOnly(context),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						acctest.CleanupBigLakeIcebergNamespace(t, catalogId, namespaceId)
						return nil
					},
				),
			},
		},
	})

	t.Cleanup(func() {
		acctest.CleanupBigLakeIcebergNamespace(t, catalogId, namespaceId)
	})
}

func testAccBiglakeIcebergIcebergNamespace_catalogOnly(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_bucket" "bucket" {
  name          = "%{catalog_id}"
  location      = "us-central1"
  force_destroy = true
  uniform_bucket_level_access = true
}

resource "google_biglake_iceberg_catalog" "catalog" {
    name = "%{catalog_id}"
    catalog_type = "CATALOG_TYPE_GCS_BUCKET"
    depends_on = [
      google_storage_bucket.bucket
    ]
}
`, context)
}

func testAccBiglakeIcebergIcebergNamespaceIamMember_basic(context map[string]interface{}) string {
	return testAccBiglakeIcebergIcebergNamespace_catalogOnly(context) + acctest.Nprintf(`
resource "google_biglake_iceberg_namespace_iam_member" "member" {
  catalog = google_biglake_iceberg_catalog.catalog.name
  name    = "%{namespace_id}"
  role    = "%{role}"
  member  = "user:admin@hashicorptest.com"
}
`, context)
}

func testAccBiglakeIcebergIcebergNamespaceIamPolicy_basic(context map[string]interface{}) string {
	return testAccBiglakeIcebergIcebergNamespace_catalogOnly(context) + acctest.Nprintf(`
data "google_iam_policy" "policy" {
  binding {
    role = "%{role}"
    members = ["user:admin@hashicorptest.com"]
  }
}

resource "google_biglake_iceberg_namespace_iam_policy" "policy" {
  catalog     = google_biglake_iceberg_catalog.catalog.name
  name        = "%{namespace_id}"
  policy_data = data.google_iam_policy.policy.policy_data
}
`, context)
}

func testAccBiglakeIcebergIcebergNamespaceIamPolicy_emptyBinding(context map[string]interface{}) string {
	return testAccBiglakeIcebergIcebergNamespace_catalogOnly(context) + acctest.Nprintf(`
data "google_iam_policy" "policy" {
}

resource "google_biglake_iceberg_namespace_iam_policy" "policy" {
  catalog     = google_biglake_iceberg_catalog.catalog.name
  name        = "%{namespace_id}"
  policy_data = data.google_iam_policy.policy.policy_data
}
`, context)
}

func testAccBiglakeIcebergIcebergNamespaceIamBinding_basic(context map[string]interface{}) string {
	return testAccBiglakeIcebergIcebergNamespace_catalogOnly(context) + acctest.Nprintf(`
resource "google_biglake_iceberg_namespace_iam_binding" "binding" {
  catalog = google_biglake_iceberg_catalog.catalog.name
  name    = "%{namespace_id}"
  role    = "%{role}"
  members = ["user:admin@hashicorptest.com"]
}
`, context)
}

func testAccBiglakeIcebergIcebergNamespaceIamBinding_update(context map[string]interface{}) string {
	return testAccBiglakeIcebergIcebergNamespace_catalogOnly(context) + acctest.Nprintf(`
resource "google_biglake_iceberg_namespace_iam_binding" "binding" {
  catalog = google_biglake_iceberg_catalog.catalog.name
  name    = "%{namespace_id}"
  role    = "%{role}"
  members = ["user:admin@hashicorptest.com", "user:gterraformtest1@gmail.com"]
}
`, context)
}
