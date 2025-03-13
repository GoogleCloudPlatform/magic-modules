package securesourcemanager_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccSecureSourceManagerHook_secureSourceManagerHookWithFieldsExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"prevent_destroy": false,
		"random_suffix":   acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecureSourceManagerHook_secureSourceManagerHookWithFieldsExample_full(context),
			},
			{
				ResourceName:            "google_secure_source_manager_hook.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"hook_id", "location", "repository_id"},
			},
			{
				Config: testAccSecureSourceManagerHook_secureSourceManagerHookWithFieldsExample_update(context),
			},
			{
				ResourceName:            "google_secure_source_manager_hook.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"hook_id", "location", "repository_id"},
			},
		},
	})
}

func testAccSecureSourceManagerHook_secureSourceManagerHookWithFieldsExample_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_secure_source_manager_instance" "instance" {
    location = "us-central1"
    instance_id = "tf-test-my-initial-instance%{random_suffix}"
    # Prevent accidental deletions.
    lifecycle {
        prevent_destroy = "%{prevent_destroy}"
    }
}

resource "google_secure_source_manager_repository" "repository" {
    repository_id = "tf-test-my-initial-repository%{random_suffix}"
    instance = google_secure_source_manager_instance.instance.name
    location = google_secure_source_manager_instance.instance.location
    # Prevent accidental deletions.
    lifecycle {
        prevent_destroy = "%{prevent_destroy}"
    }
}

resource "google_secure_source_manager_hook" "default" {
    hook_id = "tf-test-my-initial-hook%{random_suffix}"
    location = google_secure_source_manager_repository.repository.location
    repository_id = google_secure_source_manager_repository.repository.repository_id
    events = ["PUSH", "PULL_REQUEST"]
    push_option {
        branch_filter = "main"
    }
    target_uri = "https://www.example.com"
    disabled = false
}
`, context)
}

func testAccSecureSourceManagerHook_secureSourceManagerHookWithFieldsExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_secure_source_manager_instance" "instance" {
    location = "us-central1"
    instance_id = "tf-test-my-initial-instance%{random_suffix}"
    # Prevent accidental deletions.
    lifecycle {
        prevent_destroy = "%{prevent_destroy}"
    }
}

resource "google_secure_source_manager_repository" "repository" {
    repository_id = "tf-test-my-initial-repository%{random_suffix}"
    instance = google_secure_source_manager_instance.instance.name
    location = google_secure_source_manager_instance.instance.location
    # Prevent accidental deletions.
    lifecycle {
        prevent_destroy = "%{prevent_destroy}"
    }
}

resource "google_secure_source_manager_hook" "default" {
    hook_id = "tf-test-my-initial-hook%{random_suffix}"
    location = google_secure_source_manager_repository.repository.location
    repository_id = google_secure_source_manager_repository.repository.repository_id
    events = ["PUSH", "PULL_REQUEST"]
    push_option {
        branch_filter = "main"
    }
    target_uri = "https://www.update.com"
    disabled = true
}
`, context)
}
