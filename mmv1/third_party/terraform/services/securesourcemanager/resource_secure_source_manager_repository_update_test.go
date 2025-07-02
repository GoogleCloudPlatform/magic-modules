package securesourcemanager_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccSecureSourceManagerRepository_secureSourceManagerRepositoryBasicExample_update(t *testing.T) {
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
				Config: testAccSecureSourceManagerRepository_secureSourceManagerRepositoryBasicExample_basic(context),
			},
			{
				ResourceName:            "google_secure_source_manager_repository.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_config", "location", "repository_id"},
			},
			{
				Config: testAccSecureSourceManagerRepository_secureSourceManagerRepositoryBasicExample_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_secure_source_manager_repository.default", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_secure_source_manager_repository.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_config", "location", "repository_id"},
			},
		},
	})
}

func testAccSecureSourceManagerRepository_secureSourceManagerRepositoryBasicExample_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_secure_source_manager_instance" "instance" {
    location = "us-central1"
    instance_id = "tf-test-my-instance%{random_suffix}"

    # Prevent accidental deletions.
    lifecycle {
      prevent_destroy = "%{prevent_destroy}"
    }
}

resource "google_secure_source_manager_repository" "default" {
    location = "us-central1"
    repository_id = "tf-test-my-repository%{random_suffix}"
    instance = google_secure_source_manager_instance.instance.name

    # Prevent accidental deletions.
    lifecycle {
      prevent_destroy = "%{prevent_destroy}"
    }
}
`, context)
}

func testAccSecureSourceManagerRepository_secureSourceManagerRepositoryBasicExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_secure_source_manager_instance" "instance" {
    location = "us-central1"
    instance_id = "tf-test-my-instance%{random_suffix}"

    # Prevent accidental deletions.
    lifecycle {
      prevent_destroy = "%{prevent_destroy}"
    }
}

resource "google_secure_source_manager_repository" "default" {
    location = "us-central1"
    repository_id = "tf-test-my-repository%{random_suffix}"
    instance = google_secure_source_manager_instance.instance.name

    description = "new description"

    # Prevent accidental deletions.
    lifecycle {
      prevent_destroy = "%{prevent_destroy}"
    }
}
`, context)
}
