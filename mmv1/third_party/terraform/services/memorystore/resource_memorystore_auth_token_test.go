package memorystore_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	_ "github.com/hashicorp/terraform-provider-google/google/services/memorystore"
)

func TestAccMemorystoreAuthToken_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMemorystoreInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMemorystoreAuthToken_basic(context),
			},
			{
				ResourceName:      "google_memorystore_auth_token.token-basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccMemorystoreAuthToken_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_memorystore_instance" "instance-basic" {
  instance_id                 = "tf-test-memorystore-instance%{random_suffix}"
  shard_count                 = 1
  location                    = "us-central1"
  authorization_mode          = "TOKEN_AUTH"
  transit_encryption_mode     = "SERVER_AUTHENTICATION"
  deletion_protection_enabled = false
}

resource "google_memorystore_token_auth_user" "user-basic" {
  instance = google_memorystore_instance.instance-basic.name
  user_id  = "tf-test-user-%{random_suffix}"
}

resource "google_memorystore_auth_token" "token-basic" {
  token_auth_user = google_memorystore_token_auth_user.user-basic.id
  token_id        = "tf-test-token-%{random_suffix}"
}
`, context)
}
