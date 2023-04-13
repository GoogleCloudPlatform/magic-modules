package google_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	google "internal/terraform-provider-google"
)

func TestAccDataSourceSecretManagerSecret_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": google.RandString(t, 10),
	}

	google.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { google.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: google.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSecretManagerSecretDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSecretManagerSecret_basic(context),
				Check: resource.ComposeTestCheckFunc(
					google.CheckDataSourceStateMatchesResourceState("data.google_secret_manager_secret.foo", "google_secret_manager_secret.bar"),
				),
			},
		},
	})
}

func testAccDataSourceSecretManagerSecret_basic(context map[string]interface{}) string {
	return google.Nprintf(`
resource "google_secret_manager_secret" "bar" {
  secret_id = "tf-test-secret-%{random_suffix}"
  
  labels = {
    label = "my-label"
  }

  replication {
    user_managed {
      replicas {
        location = "us-central1"
      }
      replicas {
        location = "us-east1"
      }
    }
  }
}

data "google_secret_manager_secret" "foo" {
    secret_id = google_secret_manager_secret.bar.secret_id
}
`, context)
}
