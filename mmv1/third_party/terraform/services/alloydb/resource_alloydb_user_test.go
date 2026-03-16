package alloydb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccAlloydbUser_updateRoles_BuiltIn(t *testing.T) {
	t.Parallel()

	randString := acctest.RandString(t, 10)
	context := map[string]interface{}{
		"project":               envvar.GetTestProjectFromEnv(),
		"location":              "us-central1",
		"network_name":          acctest.BootstrapSharedServiceNetworkingConnection(t, "alloydb-1"),
		"random_suffix":         randString,
		"alloydb_cluster_name":  "tf-test-alloydb-cluster" + randString,
		"alloydb_cluster_pass":  "tf_test_cluster_secret" + randString,
		"alloydb_instance_name": "tf-test-alloydb-instance" + randString,
		"alloydb_user_name":     "user1" + randString,
		"alloydb_user_pass":     "tf_test_user_secret" + randString,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbUserDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbUser_alloydbUserBuiltinTestExample(context),
			},
			{
				ResourceName:            "google_alloydb_user.user1",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
			{
				Config: testAccAlloydbUser_updateRoles_BuiltIn(context),
			},
			{
				ResourceName:            "google_alloydb_user.user1",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func testAccAlloydbUser_updateRoles_BuiltIn(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_instance" "default" {
  cluster       = google_alloydb_cluster.default.name
  instance_id   = "%{alloydb_instance_name}"
  instance_type = "PRIMARY"
}

resource "google_alloydb_cluster" "default" {
  cluster_id = "%{alloydb_cluster_name}"
  location   = "us-central1"
  network_config {
    network = data.google_compute_network.default.id
  }

  initial_user {
    password = "%{alloydb_cluster_pass}"
  }

  deletion_protection = false
}

data "google_project" "project" {}

data "google_compute_network" "default" {
  name = "%{network_name}"
}

resource "google_alloydb_user" "user1" {
  cluster = google_alloydb_cluster.default.name
  user_id = "%{alloydb_user_name}"
  user_type = "ALLOYDB_BUILT_IN"

  password = "%{alloydb_user_pass}"
  database_roles = []
  depends_on = [google_alloydb_instance.default]
}`, context)
}

func TestAccAlloydbUser_updatePassword_BuiltIn(t *testing.T) {
	t.Parallel()

	randString := acctest.RandString(t, 10)
	context := map[string]interface{}{
		"project":               envvar.GetTestProjectFromEnv(),
		"location":              "us-central1",
		"network_name":          acctest.BootstrapSharedServiceNetworkingConnection(t, "alloydb-1"),
		"random_suffix":         randString,
		"alloydb_cluster_name":  "tf-test-alloydb-cluster" + randString,
		"alloydb_cluster_pass":  "tf_test_cluster_secret" + randString,
		"alloydb_instance_name": "tf-test-alloydb-instance" + randString,
		"alloydb_user_name":     "user1" + randString,
		"alloydb_user_pass":     "tf_test_user_secret" + randString,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbUserDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbUser_alloydbUserBuiltinTestExample(context),
			},
			{
				ResourceName:            "google_alloydb_user.user1",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
			{
				Config: testAccAlloydbUser_updatePass_BuiltIn(context),
			},
			{
				ResourceName:            "google_alloydb_user.user1",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func testAccAlloydbUser_updatePass_BuiltIn(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_instance" "default" {
  cluster       = google_alloydb_cluster.default.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}"
  instance_type = "PRIMARY"
}

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

data "google_project" "project" {}

data "google_compute_network" "default" {
  name = "%{network_name}"
}

resource "google_alloydb_user" "user1" {
  cluster = google_alloydb_cluster.default.name
  user_id = "%{alloydb_user_name}"
  user_type = "ALLOYDB_BUILT_IN"

  password = "%{alloydb_user_pass}-foo"
  database_roles = ["alloydbsuperuser"]
  depends_on = [google_alloydb_instance.default]
}`, context)
}

func TestAccAlloydbUser_updateRoles_IAM(t *testing.T) {
	t.Parallel()

	randString := acctest.RandString(t, 10)
	context := map[string]interface{}{
		"project":               envvar.GetTestProjectFromEnv(),
		"location":              "us-central1",
		"network_name":          acctest.BootstrapSharedServiceNetworkingConnection(t, "alloydb-1"),
		"random_suffix":         randString,
		"alloydb_cluster_name":  "tf-test-alloydb-cluster" + randString,
		"alloydb_instance_name": "tf-test-alloydb-instance" + randString,
		"alloydb_cluster_pass":  "tf_test_cluster_secret" + randString,
		"alloydb_user_name":     "user2-" + randString + "@foo.com",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbUserDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbUser_alloydbUserIamTestExample(context),
			},
			{
				ResourceName:            "google_alloydb_user.user2",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
			{
				Config: testAccAlloydbUser_updateRoles_Iam(context),
			},
			{
				ResourceName:            "google_alloydb_user.user2",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccAlloydbUser_updateRoles_Iam(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_instance" "default" {
  cluster       = google_alloydb_cluster.default.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}"
  instance_type = "PRIMARY"
}
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
data "google_project" "project" {}
data "google_compute_network" "default" {
  name = "%{network_name}"
}
resource "google_alloydb_user" "user2" {
  cluster = google_alloydb_cluster.default.name
  user_id = "%{alloydb_user_name}"
  user_type = "ALLOYDB_IAM_USER"
  database_roles = ["alloydbiamuser", "alloydbsuperuser"]
  depends_on = [google_alloydb_instance.default]
}`, context)
}

func TestAccAlloydbUser_alloydbUserBuiltinWithPasswordWo(t *testing.T) {
	t.Parallel()

	randString := acctest.RandString(t, 10)
	context := map[string]interface{}{
		"project":               envvar.GetTestProjectFromEnv(),
		"location":              "us-central1",
		"network_name":          acctest.BootstrapSharedServiceNetworkingConnection(t, "alloydb-1"),
		"random_suffix":         randString,
		"alloydb_cluster_name":  "tf-test-alloydb-cluster" + randString,
		"alloydb_cluster_pass":  "tf_test_cluster_secret" + randString,
		"alloydb_instance_name": "tf-test-alloydb-instance" + randString,
		"alloydb_user_name":     "user1" + randString,
		"alloydb_user_pass":     "tf_test_user_secret" + randString,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbUserDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbUser_alloydbUserBuiltinWithPasswordWo(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("google_alloydb_user.user1", "password_wo"),
					resource.TestCheckResourceAttr("google_alloydb_user.user1", "password_wo_version", "1"),
				),
			},
			{
				ResourceName:            "google_alloydb_user.user1",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func testAccAlloydbUser_alloydbUserBuiltinWithPasswordWo(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_instance" "default" {
  cluster       = google_alloydb_cluster.default.name
  instance_id   = "%{alloydb_instance_name}"
  instance_type = "PRIMARY"
}
resource "google_alloydb_cluster" "default" {
  cluster_id = "%{alloydb_cluster_name}"
  location   = "us-central1"
  network_config {
    network = data.google_compute_network.default.id
  }
  initial_user {
    password = "%{alloydb_cluster_pass}"
  }
  deletion_protection = false
}
data "google_project" "project" {}
data "google_compute_network" "default" {
  name = "%{network_name}"
}
resource "google_alloydb_user" "user1" {
  cluster = google_alloydb_cluster.default.name
  user_id = "%{alloydb_user_name}"
  user_type = "ALLOYDB_BUILT_IN"
  password_wo = "%{alloydb_user_pass}"
  password_wo_version = 1
  database_roles = ["alloydbsuperuser"]
  depends_on = [google_alloydb_instance.default]
}`, context)
}

func TestAccAlloydbUser_alloydbUserBuiltinWithPasswordWo_update(t *testing.T) {
	t.Parallel()

	randString := acctest.RandString(t, 10)
	context := map[string]interface{}{
		"project":               envvar.GetTestProjectFromEnv(),
		"location":              "us-central1",
		"network_name":          acctest.BootstrapSharedServiceNetworkingConnection(t, "alloydb-1"),
		"random_suffix":         randString,
		"alloydb_cluster_name":  "tf-test-alloydb-cluster" + randString,
		"alloydb_cluster_pass":  "tf_test_cluster_secret" + randString,
		"alloydb_instance_name": "tf-test-alloydb-instance" + randString,
		"alloydb_user_name":     "user1" + randString,
		"alloydb_user_pass":     "tf_test_user_secret" + randString,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbUserDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbUser_alloydbUserBuiltinWithPasswordWo(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("google_alloydb_user.user1", "password_wo"),
					resource.TestCheckResourceAttr("google_alloydb_user.user1", "password_wo_version", "1"),
				),
			},
			{
				Config: testAccAlloydbUser_alloydbUserBuiltinWithPasswordWo_update(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("google_alloydb_user.user1", "password_wo"),
					resource.TestCheckResourceAttr("google_alloydb_user.user1", "password_wo_version", "2"),
				),
			},
		},
	})
}

func testAccAlloydbUser_alloydbUserBuiltinWithPasswordWo_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_instance" "default" {
  cluster       = google_alloydb_cluster.default.name
  instance_id   = "%{alloydb_instance_name}"
  instance_type = "PRIMARY"
}
resource "google_alloydb_cluster" "default" {
  cluster_id = "%{alloydb_cluster_name}"
  location   = "us-central1"
  network_config {
    network = data.google_compute_network.default.id
  }
  initial_user {
    password = "%{alloydb_cluster_pass}"
  }
  deletion_protection = false
}
data "google_project" "project" {}
data "google_compute_network" "default" {
  name = "%{network_name}"
}
resource "google_alloydb_user" "user1" {
  cluster = google_alloydb_cluster.default.name
  user_id = "%{alloydb_user_name}"
  user_type = "ALLOYDB_BUILT_IN"
  password_wo = "tf_test_user_secret%{random_suffix}-update"
  password_wo_version = 2
  database_roles = ["alloydbsuperuser"]
  depends_on = [google_alloydb_instance.default]
}`, context)
}
