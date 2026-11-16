package memorystore_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	tpgmemorystore "github.com/hashicorp/terraform-provider-google/google/services/memorystore"
)

func TestUnitMemorystoreTokenAuthUser_userIdValidation(t *testing.T) {
	t.Parallel()
	schema := tpgmemorystore.ResourceMemorystoreTokenAuthUser().Schema["user_id"]
	validateFn := schema.ValidateFunc

	for _, invalid := range []string{"default"} {
		_, errs := validateFn(invalid, "user_id")
		if len(errs) == 0 {
			t.Errorf("expected validation error for user_id=%q, got none", invalid)
		}
	}

	for _, valid := range []string{"valid-user-1", "user-prod", "alice"} {
		warns, errs := validateFn(valid, "user_id")
		if len(errs) != 0 || len(warns) != 0 {
			t.Errorf("expected no validation errors for valid user_id=%q, got errs=%v, warns=%v", valid, errs, warns)
		}
	}
}

func TestAccMemorystoreTokenAuthUser_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMemorystoreTokenAuthUserDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccMemorystoreTokenAuthUser_invalidUserId(context, "default"),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected user_id to not match any of \[default\]`),
			},
			{
				Config: testAccMemorystoreTokenAuthUser_basic(context),
			},
			{
				ResourceName:      "google_memorystore_token_auth_user.user-basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccMemorystoreTokenAuthUser_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_memorystore_instance" "instance" {
  instance_id                 = "tf-test-ms-instance-%{random_suffix}"
  shard_count                 = 1
  node_type                   = "SHARED_CORE_NANO"
  location                    = "europe-west4"
  authorization_mode          = "TOKEN_AUTH"
  deletion_protection_enabled = false
}

resource "google_memorystore_token_auth_user" "user-basic" {
  instance                    = google_memorystore_instance.instance.name
  location                    = "europe-west4"
  user_id                     = "tf-test-user-%{random_suffix}"
}
`, context)
}

func testAccMemorystoreTokenAuthUser_invalidUserId(context map[string]interface{}, userId string) string {
	context["user_id"] = userId
	return acctest.Nprintf(`
resource "google_memorystore_instance" "instance" {
  instance_id                 = "tf-test-ms-instance-%{random_suffix}"
  shard_count                 = 1
  node_type                   = "SHARED_CORE_NANO"
  location                    = "europe-west4"
  authorization_mode          = "TOKEN_AUTH"
  deletion_protection_enabled = false
}

resource "google_memorystore_token_auth_user" "user-invalid" {
  instance                    = google_memorystore_instance.instance.name
  location                    = "europe-west4"
  user_id                     = "%{user_id}"
}
`, context)
}
