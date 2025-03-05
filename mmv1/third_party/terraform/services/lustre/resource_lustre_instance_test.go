package lustre_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccLustreInstance_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLustreInstance_full(context),
			},
			{
				ResourceName:            "google_lustre_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance_id", "labels", "location", "terraform_labels"},
			},
			{
				Config: testAccLustreInstance_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							"google_lustre_instance.description",
							plancheck.ResourceActionUpdate,
						),
					},
				},
			},
			{
				ResourceName:            "google_lustre_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance_id", "labels", "location", "terraform_labels"},
			},
		},
	})
}

func testAccLustreInstance_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_lustre_instance" "instance" {
  instance_id  = "tf-test-my-instance%{random_suffix}"
  location     = "us-central1-a"
  filesystem   = "testfs"
  capacity_gib = 18000
  network      = "${data.google_project.project.id}/global/networks/default"
	timeouts {
		create = "180m"
	}
}

# Create an IP address
resource "google_compute_global_address" "private_ip_alloc" {
  name          = "tf-test-my-ip-range%{random_suffix}"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 24
  network       = "${data.google_project.project.id}/global/networks/default"
}

# Create a private connection
resource "google_service_networking_connection" "default" {
  network                 = "${data.google_project.project.id}/global/networks/default"
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_ip_alloc.name]
	update_on_creation_fail = true
}

data "google_project" "project" {
}
`, context)
}

func testAccLustreInstance_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_lustre_instance" "instance" {
  instance_id  = "tf-test-my-instance%{random_suffix}"
  location     = "us-central1-a"
  description  = "description updated"
  filesystem   = "testfs"
  capacity_gib = 18000
  network      = "${data.google_project.project.id}/global/networks/default"
  labels       = {
    test = "newLabel"
  }
	timeouts {
		create = "180m"
  }
}

# Create an IP address
resource "google_compute_global_address" "private_ip_alloc" {
  name          = "tf-test-my-ip-range%{random_suffix}"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 24
  network       = "${data.google_project.project.id}/global/networks/default"
}

# Create a private connection
resource "google_service_networking_connection" "default" {
  network                 = "${data.google_project.project.id}/global/networks/default"
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_ip_alloc.name]
	update_on_creation_fail = true
}

data "google_project" "project" {
}
`, context)
}

func TestAccLustreInstanceBeta_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLustreInstance_full(context),
			},
			{
				ResourceName:            "google_lustre_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance_id", "labels", "location", "terraform_labels"},
			},
			{
				Config: testAccLustreInstanceBeta_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							"google_lustre_instance.description",
							plancheck.ResourceActionUpdate,
						),
					},
				},
			},
			{
				ResourceName:            "google_lustre_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance_id", "labels", "location", "terraform_labels"},
			},
		},
	})
}

func testAccLustreInstanceBeta_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_lustre_instance" "instance" {
  provider     = google-beta
  instance_id  = "tf-test-my-instance%{random_suffix}"
  location     = "us-central1-a"
  filesystem   = "testfs"
  capacity_gib = 18000
  network      = "${data.google_project.project.id}/global/networks/default"
	timeouts {
		create = "180m"
	}
}

# Create an IP address
resource "google_compute_global_address" "private_ip_alloc" {
  provider      = google-beta
  name          = "tf-test-my-ip-range%{random_suffix}"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 24
  network       = "${data.google_project.project.id}/global/networks/default"
}

# Create a private connection
resource "google_service_networking_connection" "default" {
  provider                = google-beta
  network                 = "${data.google_project.project.id}/global/networks/default"
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_ip_alloc.name]
	update_on_creation_fail = true
}

data "google_project" "project" {
}
`, context)
}

func testAccLustreInstanceBeta_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_lustre_instance" "instance" {
  provider     = google-beta
  instance_id  = "tf-test-my-instance%{random_suffix}"
  location     = "us-central1-a"
  description  = "description updated"
  filesystem   = "testfs"
  capacity_gib = 18000
  network      = "${data.google_project.project.id}/global/networks/default"
  labels       = {
    test = "newLabel"
  }
	timeouts {
		create = "180m"
  }
}

# Create an IP address
resource "google_compute_global_address" "private_ip_alloc" {
  provider      = google-beta
  name          = "tf-test-my-ip-range%{random_suffix}"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 24
  network       = "${data.google_project.project.id}/global/networks/default"
}

# Create a private connection
resource "google_service_networking_connection" "default" {
  provider                = google-beta
  network                 = "${data.google_project.project.id}/global/networks/default"
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_ip_alloc.name]
	update_on_creation_fail = true
}

data "google_project" "project" {
}
`, context)
}
