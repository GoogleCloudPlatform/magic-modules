package vertexai_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccVertexAINotebookRuntimeTemplate_vertexAiNotebookRuntimeTemplateFullExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"kms_key_name":         acctest.BootstrapKMSKeyInLocation(t, "europe-west4").CryptoKey.Name,
		"kms_key_name_updated": acctest.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", "europe-west4", "tf-bootstrap-colab-key1").CryptoKey.Name,
		"random_suffix":        acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVertexAINotebookRuntimeTemplate_vertexAiNotebookRuntimeTemplateFullExample_full(context),
			},
			{
				ResourceName:            "google_vertex_ai_notebook_runtime_template.template",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "name", "region", "terraform_labels"},
			},
			{
				Config: testAccVertexAINotebookRuntimeTemplate_vertexAiNotebookRuntimeTemplateFullExample_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_vertex_ai_notebook_runtime_template.template", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_vertex_ai_notebook_runtime_template.template",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "name", "region", "terraform_labels"},
			},
		},
	})
}

func testAccVertexAINotebookRuntimeTemplate_vertexAiNotebookRuntimeTemplateFullExample_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network" "network" {
  name                    = "tf-test-my-network%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "subnetwork" {
  name          = "tf-test-my-subnetwork%{random_suffix}"
  ip_cidr_range = "10.3.0.0/16"
  region        = "europe-west4"
  private_ip_google_access = true
  network       = google_compute_network.network.id
}

resource "google_vertex_ai_notebook_runtime_template" "template" {
  name = "tf-test-my-template%{random_suffix}"
  display_name = "Sample template"
  description = "Sample template description"

  machine_spec {
    machine_type = "e2-standard-4"
  }

  data_persistent_disk_spec {
    disk_size_gb = 100
    disk_type    = "pd-standard"
  }

  network_spec {
    enable_internet_access = false
    network = google_compute_network.network.id
    subnetwork = google_compute_subnetwork.subnetwork.id
  }

  idle_shutdown_config {
    idle_timeout = "10800s"
    idle_shutdown_disabled = false
  }

  notebook_runtime_type = "USER_DEFINED"

  shielded_vm_config {
    enable_secure_boot = true
  }

  network_tags = ["colab-notebook"]

  encryption_spec {
    kms_key_name = "%{kms_key_name}"
  }
}
`, context)
}

func testAccVertexAINotebookRuntimeTemplate_vertexAiNotebookRuntimeTemplateFullExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network" "network" {
  name                    = "tf-test-my-network%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "subnetwork" {
  name          = "tf-test-my-subnetwork%{random_suffix}"
  ip_cidr_range = "10.3.0.0/16"
  region        = "europe-west4"
  private_ip_google_access = true
  network       = google_compute_network.network.id
}

resource "google_vertex_ai_notebook_runtime_template" "template" {
  name = "tf-test-my-template%{random_suffix}"
  display_name = "Sample template"
  description = "Sample template description"

  machine_spec {
    machine_type = "e2-standard-4"
  }

  data_persistent_disk_spec {
    disk_size_gb = 100
    disk_type    = "pd-standard"
  }

  network_spec {
    enable_internet_access = false
    network = google_compute_network.network.id
    subnetwork = google_compute_subnetwork.subnetwork.id
  }

  idle_shutdown_config {
    idle_timeout = "10800s"
    idle_shutdown_disabled = false
  }

  notebook_runtime_type = "USER_DEFINED"

  shielded_vm_config {
    enable_secure_boot = true
  }

  network_tags = ["colab-notebook"]

  encryption_spec {
    kms_key_name = "%{kms_key_name_updated}"
  }
}
`, context)
}
