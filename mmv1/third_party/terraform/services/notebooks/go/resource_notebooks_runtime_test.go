package notebooks_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNotebooksRuntime_update(t *testing.T) {
	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNotebooksRuntimeDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNotebooksRuntime_basic(context),
			},
			{
				ResourceName:            "google_notebooks_runtime.runtime",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccNotebooksRuntime_update(context),
			},
			{
				ResourceName:            "google_notebooks_runtime.runtime",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccNotebooksRuntime_basic(context),
			},
			{
				ResourceName:            "google_notebooks_runtime.runtime",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}


func testAccNotebooksRuntime_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_notebooks_runtime" "runtime" {
  name = "tf-test-notebooks-runtime%{random_suffix}"
  location = "us-central1"
  access_config {
    access_type = "SINGLE_USER"
    runtime_owner = "admin@hashicorptest.com"
  }
  software_config {}
  virtual_machine {
    virtual_machine_config {
     machine_type = "n1-standard-4"
      data_disk {
        initialize_params {
          disk_size_gb = "100"
          disk_type = "PD_STANDARD"
        }
      }
      reserved_ip_range = "192.168.255.0/24"
    }
  }
}
`, context)
}

func testAccNotebooksRuntime_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_notebooks_runtime" "runtime" {
  name = "tf-test-notebooks-runtime%{random_suffix}"
  location = "us-central1"
  access_config {
    access_type = "SINGLE_USER"
    runtime_owner = "admin@hashicorptest.com"
  }
  software_config {
    idle_shutdown_timeout = "80"
  }
  virtual_machine {
    virtual_machine_config {
     machine_type = "n1-standard-8"
      data_disk {
        initialize_params {
          disk_size_gb = "100"
          disk_type = "PD_STANDARD"
        }
      }
      accelerator_config {
        core_count = "1"
        type = "NVIDIA_TESLA_V100"
      }
      reserved_ip_range = "192.168.255.0/24"
    }
  }
  labels = {
    k = "val"
  }
}
`, context)
}
