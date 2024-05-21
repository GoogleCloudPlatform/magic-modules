package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccComputeStoragePoolIam(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"role":          "roles/viewer",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Test Iam Policy
				Config: testAccComputeStoragePoolIamPolicy_basic(context),
				Check:  resource.TestCheckResourceAttrSet("data.google_compute_storage_pool_iam_policy.foo", "policy_data"),
			},
			{
				ResourceName:      "google_compute_storage_pool_iam_policy.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/zones/%s/storagePools/%s", envvar.GetTestProjectFromEnv(), envvar.GetTestZoneFromEnv(), fmt.Sprintf("tf-test-test-storage-pool%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeStoragePoolIamPolicy_emptyBinding(context),
			},
			{
				ResourceName:      "google_compute_storage_pool_iam_policy.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/zones/%s/storagePools/%s", envvar.GetTestProjectFromEnv(), envvar.GetTestZoneFromEnv(), fmt.Sprintf("tf-test-test-storage-pool%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test Iam Binding
				Config: testAccComputeStoragePoolIamBinding_basic(context),
			},
			{
				ResourceName:      "google_compute_storage_pool_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/zones/%s/storagePools/%s roles/viewer", envvar.GetTestProjectFromEnv(), envvar.GetTestZoneFromEnv(), fmt.Sprintf("tf-test-test-storage-pool%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test Iam Binding update
				Config: testAccComputeStoragePoolIamBinding_update(context),
			},
			{
				ResourceName:      "google_compute_storage_pool_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/zones/%s/storagePools/%s roles/viewer", envvar.GetTestProjectFromEnv(), envvar.GetTestZoneFromEnv(), fmt.Sprintf("tf-test-test-storage-pool%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test Iam Member creation (no update for member, no need to test)
				Config: testAccComputeStoragePoolIamMember_basic(context),
			},
			{
				ResourceName:      "google_compute_storage_pool_iam_member.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/zones/%s/storagePools/%s roles/viewer user:admin@hashicorptest.com", envvar.GetTestProjectFromEnv(), envvar.GetTestZoneFromEnv(), fmt.Sprintf("tf-test-test-storage-pool%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeStoragePoolIamMember_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_storage_pool" "default" {
  name  = "tf-test-test-storage-pool%{random_suffix}"
  zone  = "us-central1-a"
  storage_pool_type = "hyperdisk-throughput"
  pool_provisioned_capacity_gb = 10240
  pool_provisioned_throughput = 100
}

resource "google_compute_storage_pool_iam_member" "foo" {
  project = google_compute_storage_pool.default.project
  zone = google_compute_storage_pool.default.zone
  name = google_compute_storage_pool.default.name
  role = "%{role}"
  member = "user:admin@hashicorptest.com"
}
`, context)
}

func testAccComputeStoragePoolIamPolicy_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_storage_pool" "default" {
  name  = "tf-test-test-storage-pool%{random_suffix}"
  zone  = "us-central1-a"
  storage_pool_type = "hyperdisk-throughput"
  pool_provisioned_capacity_gb = 10240
  pool_provisioned_throughput = 100
}

data "google_iam_policy" "foo" {
  binding {
    role = "%{role}"
    members = ["user:admin@hashicorptest.com"]
  }
}

resource "google_compute_storage_pool_iam_policy" "foo" {
  project = google_compute_storage_pool.default.project
  zone = google_compute_storage_pool.default.zone
  name = google_compute_storage_pool.default.name
  policy_data = data.google_iam_policy.foo.policy_data
}

data "google_compute_storage_pool_iam_policy" "foo" {
  project = google_compute_storage_pool.default.project
  zone = google_compute_storage_pool.default.zone
  name = google_compute_storage_pool.default.name
  depends_on = [
    google_compute_storage_pool_iam_policy.foo
  ]
}
`, context)
}

func testAccComputeStoragePoolIamPolicy_emptyBinding(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_storage_pool" "default" {
  name  = "tf-test-test-storage-pool%{random_suffix}"
  zone  = "us-central1-a"
  storage_pool_type = "hyperdisk-throughput"
  pool_provisioned_capacity_gb = 10240
  pool_provisioned_throughput = 100
}

data "google_iam_policy" "foo" {
}

resource "google_compute_storage_pool_iam_policy" "foo" {
  project = google_compute_storage_pool.default.project
  zone = google_compute_storage_pool.default.zone
  name = google_compute_storage_pool.default.name
  policy_data = data.google_iam_policy.foo.policy_data
}
`, context)
}

func testAccComputeStoragePoolIamBinding_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_storage_pool" "default" {
  name  = "tf-test-test-storage-pool%{random_suffix}"
  zone  = "us-central1-a"
  storage_pool_type = "hyperdisk-throughput"
  pool_provisioned_capacity_gb = 10240
  pool_provisioned_throughput = 100
}

resource "google_compute_storage_pool_iam_binding" "foo" {
  project = google_compute_storage_pool.default.project
  zone = google_compute_storage_pool.default.zone
  name = google_compute_storage_pool.default.name
  role = "%{role}"
  members = ["user:admin@hashicorptest.com"]
}
`, context)
}

func testAccComputeStoragePoolIamBinding_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_storage_pool" "default" {
  name  = "tf-test-test-storage-pool%{random_suffix}"
  zone  = "us-central1-a"
  storage_pool_type = "hyperdisk-throughput"
  pool_provisioned_capacity_gb = 10240
  pool_provisioned_throughput = 100
}

resource "google_compute_storage_pool_iam_binding" "foo" {
  project = google_compute_storage_pool.default.project
  zone = google_compute_storage_pool.default.zone
  name = google_compute_storage_pool.default.name
  role = "%{role}"
  members = ["user:admin@hashicorptest.com", "user:gterraformtest1@gmail.com"]
}
`, context)
}
