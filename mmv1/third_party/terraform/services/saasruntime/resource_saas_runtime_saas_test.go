package saasruntime_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccSaasRuntimeSaas_saasRuntimeSaasBasicExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderBetaFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSaasRuntimeSaas_saasRuntimeSaasBasicExample_basic(context),
			},
			{
				ResourceName:            "google_saas_runtime_saas.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "labels", "location", "saas_id", "terraform_labels"},
			},
			{
				Config: testAccSaasRuntimeSaas_saasRuntimeSaasBasicExample_update(context),
			},
			{
				ResourceName:            "google_saas_runtime_saas.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "labels", "location", "saas_id", "terraform_labels"},
			},
		},
	})
}

func testAccSaasRuntimeSaas_saasRuntimeSaasBasicExample_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_saas_runtime_saas" "example" {
  provider = google-beta
  saas_id  = "tf-test-test-saas%{random_suffix}"
  location = "global"

  locations {
    name = "us-central1"
  }
  locations {
    name = "europe-west1"
  }
}
`, context)
}

func testAccSaasRuntimeSaas_saasRuntimeSaasBasicExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_saas_runtime_saas" "example" {
  provider = google-beta
  saas_id  = "tf-test-test-saas%{random_suffix}"
  location = "global"

  locations {
    name = "us-central1"
  }
  locations {
    name = "europe-west1"
  }
  locations {
    name = "us-east1"
  }
  labels = {
    "label-one": "value-one"
  }
  annotations = {
    "annotation-one": "value-one"
  }
}
`, context)
}
