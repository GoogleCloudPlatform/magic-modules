package alloydb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	_ "github.com/hashicorp/terraform-provider-google/google/services/alloydb"
	_ "github.com/hashicorp/terraform-provider-google/google/services/compute"
	"github.com/hashicorp/terraform-provider-google/google/services/servicenetworking"
)

func TestAccAlloydbUserDatasourceConfig(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"network_name":  servicenetworking.BootstrapSharedServiceNetworkingConnection(t, "alloydb-1"),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbUserDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbUserDatasourceConfig(context),
				Check: resource.ComposeTestCheckFunc(
					// password is a write-only (ignore_read) client-side field not returned by the AlloyDB API.
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_alloydb_user.default",
						"google_alloydb_user.default",
						[]string{"password"},
					),
					// password is a write-only (ignore_read) client-side field not returned by the AlloyDB API.
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_alloydb_user.without_optional",
						"google_alloydb_user.default",
						[]string{"password"},
					),
				),
			},
		},
	})
}

func testAccAlloydbUserDatasourceConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network_config {
    network = data.google_compute_network.default.id
  }
  initial_user {
    password = "tf_test_cluster_secret%{random_suffix}"
  }

  deletion_protection = false
}

resource "google_alloydb_instance" "default" {
  cluster       = google_alloydb_cluster.default.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}"
  instance_type = "PRIMARY"
}

resource "google_alloydb_user" "default" {
  cluster        = google_alloydb_cluster.default.name
  user_id        = "user1%{random_suffix}"
  user_type      = "ALLOYDB_BUILT_IN"
  password       = "tf_test_user_secret%{random_suffix}"
  database_roles = ["alloydbsuperuser"]
  depends_on     = [google_alloydb_instance.default]
}

data "google_compute_network" "default" {
  name = "%{network_name}"
}

data "google_alloydb_user" "default" {
  cluster_id = google_alloydb_cluster.default.cluster_id
  user_id    = google_alloydb_user.default.user_id
  location   = google_alloydb_cluster.default.location
  project    = google_alloydb_cluster.default.project
}

data "google_alloydb_user" "without_optional" {
  cluster_id = google_alloydb_cluster.default.cluster_id
  user_id    = google_alloydb_user.default.user_id
}
`, context)
}
