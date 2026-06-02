package compute_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	_ "github.com/hashicorp/terraform-provider-google/google/services/compute"
	_ "github.com/hashicorp/terraform-provider-google/google/services/resourcemanager"
)

func TestAccComputeNetworkAttachment_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"random_suffix":   acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeNetworkAttachmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeNetworkAttachment_full(context),
			},
			{
				ResourceName:            "google_compute_network_attachment.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"region"},
			},
			{
				Config: testAccComputeNetworkAttachment_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_compute_network_attachment.default", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_compute_network_attachment.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"region"},
			},
		},
	})
}

func testAccComputeNetworkAttachment_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network_attachment" "default" {
    name = "tf-test-basic-network-attachment%{random_suffix}"
    region = "us-central1"
    description = "basic network attachment description"
    connection_preference = "ACCEPT_MANUAL"

    subnetworks = [
        google_compute_subnetwork.net1.self_link
    ]

    producer_accept_lists = [
        google_project.accepted_producer_project1.project_id
    ]

    producer_reject_lists = [
        google_project.rejected_producer_project1.project_id
    ]
}

resource "google_compute_network" "default" {
    name = "tf-test-basic-network%{random_suffix}"
    auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "net1" {
    name = "tf-test-basic-subnetwork1-%{random_suffix}"
    region = "us-central1"

    network = google_compute_network.default.id
    ip_cidr_range = "10.0.0.0/16"
}

resource "google_compute_subnetwork" "net2" {
    name = "tf-test-basic-subnetwork2-%{random_suffix}"
    region = "us-central1"

    network = google_compute_network.default.id
    ip_cidr_range = "10.1.0.0/16"
}

resource "google_project" "rejected_producer_project1" {
    project_id      = "tf-test-prj-reject1-%{random_suffix}"
    name            = "tf-test-prj-reject1-%{random_suffix}"
    org_id          = "%{org_id}"
    billing_account = "%{billing_account}"
    deletion_policy = "DELETE"
}

resource "google_project" "rejected_producer_project2" {
    project_id      = "tf-test-prj-reject2-%{random_suffix}"
    name            = "tf-test-prj-reject2-%{random_suffix}"
    org_id          = "%{org_id}"
    billing_account = "%{billing_account}"
    deletion_policy = "DELETE"
}

resource "google_project" "accepted_producer_project1" {
    project_id      = "tf-test-prj-accept1-%{random_suffix}"
    name            = "tf-test-prj-accept1-%{random_suffix}"
    org_id          = "%{org_id}"
    billing_account = "%{billing_account}"
    deletion_policy = "DELETE"
}

resource "google_project" "accepted_producer_project2" {
    project_id      = "tf-test-prj-accept2-%{random_suffix}"
    name            = "tf-test-prj-accept2-%{random_suffix}"
    org_id          = "%{org_id}"
    billing_account = "%{billing_account}"
    deletion_policy = "DELETE"
}
`, context)
}

func testAccComputeNetworkAttachment_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network_attachment" "default" {
    name = "tf-test-basic-network-attachment%{random_suffix}"
    region = "us-central1"
    description = "basic network attachment description"
    connection_preference = "ACCEPT_MANUAL"

    subnetworks = [
        google_compute_subnetwork.net2.self_link
    ]

    producer_accept_lists = [
        google_project.accepted_producer_project2.project_id
    ]

    producer_reject_lists = [
        google_project.rejected_producer_project2.project_id
    ]
}

resource "google_compute_network" "default" {
    name = "tf-test-basic-network%{random_suffix}"
    auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "net1" {
    name = "tf-test-basic-subnetwork1-%{random_suffix}"
    region = "us-central1"

    network = google_compute_network.default.id
    ip_cidr_range = "10.0.0.0/16"
}

resource "google_compute_subnetwork" "net2" {
    name = "tf-test-basic-subnetwork2-%{random_suffix}"
    region = "us-central1"

    network = google_compute_network.default.id
    ip_cidr_range = "10.1.0.0/16"
}

resource "google_project" "rejected_producer_project1" {
    project_id      = "tf-test-prj-reject1-%{random_suffix}"
    name            = "tf-test-prj-reject1-%{random_suffix}"
    org_id          = "%{org_id}"
    billing_account = "%{billing_account}"
    deletion_policy = "DELETE"
}

resource "google_project" "rejected_producer_project2" {
    project_id      = "tf-test-prj-reject2-%{random_suffix}"
    name            = "tf-test-prj-reject2-%{random_suffix}"
    org_id          = "%{org_id}"
    billing_account = "%{billing_account}"
    deletion_policy = "DELETE"
}

resource "google_project" "accepted_producer_project1" {
    project_id      = "tf-test-prj-accept1-%{random_suffix}"
    name            = "tf-test-prj-accept1-%{random_suffix}"
    org_id          = "%{org_id}"
    billing_account = "%{billing_account}"
    deletion_policy = "DELETE"
}

resource "google_project" "accepted_producer_project2" {
    project_id      = "tf-test-prj-accept2-%{random_suffix}"
    name            = "tf-test-prj-accept2-%{random_suffix}"
    org_id          = "%{org_id}"
    billing_account = "%{billing_account}"
    deletion_policy = "DELETE"
}
`, context)
}

func TestAccComputeNetworkAttachment_stableIpProjectIdAuth(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"random_suffix":   acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders:        map[string]resource.ExternalProvider{"time": {}},
		CheckDestroy:             testAccCheckComputeNetworkAttachmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeNetworkAttachment_stableIpProjectIdAuth(context),
			},
			{
				ResourceName:            "google_compute_network_attachment.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"region"},
			},
		},
	})
}

func testAccComputeNetworkAttachment_stableIpProjectIdAuth(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "accepted_producer_project1" {
  project_id      = "tf-test-prj-accept1-%{random_suffix}"
  name            = "tf-test-prj-accept1-%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
  deletion_policy = "DELETE"
}

resource "google_project_service" "compute_service" {
  project = google_project.accepted_producer_project1.project_id
  service = "compute.googleapis.com"
}

resource "time_sleep" "wait_300_seconds" {
  create_duration = "300s"
  depends_on = [google_project_service.compute_service]
}

resource "google_compute_network" "default" {
  name = "tf-test-basic-network%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "net1" {
  name = "tf-test-basic-subnetwork1-%{random_suffix}"
  region = "us-central1"

  network = google_compute_network.default.id
  ip_cidr_range = "10.0.0.0/16"
}

resource "google_compute_network_attachment" "default" {
  name = "tf-test-basic-network-attachment%{random_suffix}"
  region = "us-central1"
  description = "basic network attachment description"
  connection_preference = "ACCEPT_MANUAL"

  subnetworks = [
    google_compute_subnetwork.net1.self_link
  ]

  producer_accept_lists = [
    google_project.accepted_producer_project1.project_id
  ]

	depends_on = [time_sleep.wait_300_seconds]
}

resource "google_compute_address" "producer_address" {
  project = google_project.accepted_producer_project1.project_id
  name = "tf-test-basic-address%{random_suffix}"
	network_attachment = google_compute_network_attachment.default.id
}

resource "google_compute_instance" "producer_instance" {
  project = google_project.accepted_producer_project1.project_id
  name         = "tf-test-basic-instance%{random_suffix}"
  machine_type = "e2-medium"
  zone         = "us-central1-a"
  tags         = ["foo", "bar"]

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = google_compute_network.default.id
    subnetwork = google_compute_subnetwork.net1.id
  }

  network_interface {
    network_attachment = google_compute_network_attachment.default.self_link
		network_ip = google_compute_address.producer_address.self_link
  }
  
  metadata = {
    foo = "bar"
  }
}

data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}
`, context)
}
