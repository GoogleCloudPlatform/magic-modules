package bigtable_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccBigtableInstanceIamBinding(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	acctest.SkipIfVcr(t)
	t.Parallel()

	randomString := acctest.RandString(t, 10)

	context := map[string]interface{}{
		"instance": "tf-bigtable-iam-" + randomString,
		"cluster":  "c-" + randomString,
		"account":  "tf-bigtable-iam-" + randomString,
		"role":     "roles/bigtable.user",
	}

	importId := fmt.Sprintf("projects/%s/instances/%s %s",
		envvar.GetTestProjectFromEnv(), context["instance"].(string), context["role"].(string))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Test IAM Binding creation
				Config: testAccBigtableInstanceIamBinding_basic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_bigtable_instance_iam_binding.binding", "role", context["role"].(string)),
				),
			},
			{
				ResourceName:      "google_bigtable_instance_iam_binding.binding",
				ImportStateId:     importId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test IAM Binding update
				Config: testAccBigtableInstanceIamBinding_update(context),
			},
			{
				ResourceName:      "google_bigtable_instance_iam_binding.binding",
				ImportStateId:     importId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBigtableInstanceIamMember(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	acctest.SkipIfVcr(t)
	t.Parallel()

	randomString := acctest.RandString(t, 10)

	context := map[string]interface{}{
		"instance": "tf-bigtable-iam-" + randomString,
		"cluster":  "c-" + randomString,
		"account":  "tf-bigtable-iam-" + randomString,
		"role":     "roles/bigtable.user",
	}

	importId := fmt.Sprintf("projects/%s/instances/%s %s serviceAccount:%s",
		envvar.GetTestProjectFromEnv(),
		context["instance"].(string),
		context["role"].(string),
		envvar.ServiceAccountCanonicalEmail(context["account"].(string)))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Test IAM Binding creation
				Config: testAccBigtableInstanceIamMember(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_bigtable_instance_iam_member.member", "role", context["role"].(string)),
					resource.TestCheckResourceAttr(
						"google_bigtable_instance_iam_member.member", "member", "serviceAccount:"+envvar.ServiceAccountCanonicalEmail(context["account"].(string))),
				),
			},
			{
				ResourceName:      "google_bigtable_instance_iam_member.member",
				ImportStateId:     importId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBigtableInstanceIamPolicy(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	acctest.SkipIfVcr(t)
	t.Parallel()

	randomString := acctest.RandString(t, 10)

	context := map[string]interface{}{
		"instance": "tf-bigtable-iam-" + randomString,
		"cluster":  "c-" + randomString,
		"account":  "tf-bigtable-iam-" + randomString,
		"role":     "roles/bigtable.user",
	}

	importId := fmt.Sprintf("projects/%s/instances/%s",
		envvar.GetTestProjectFromEnv(), context["instance"].(string))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Test IAM Binding creation
				Config: testAccBigtableInstanceIamPolicy(context),
				Check:  resource.TestCheckResourceAttrSet("data.google_bigtable_instance_iam_policy.policy", "policy_data"),
			},
			{
				ResourceName:      "google_bigtable_instance_iam_policy.policy",
				ImportStateId:     importId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccBigtableInstanceIamBinding_basic(context map[string]interface{}) string {
	return testBigtableInstanceIam(context) + acctest.Nprintf(`
resource "google_service_account" "test-account1" {
  account_id   = "%{account}-1"
  display_name = "Bigtable Instance IAM Testing Account"
}

resource "google_service_account" "test-account2" {
  account_id   = "%{account}-2"
  display_name = "Bigtable instance Iam Testing Account"
}

resource "google_bigtable_instance_iam_binding" "binding" {
  instance = google_bigtable_instance.instance.name
  role     = "%{role}"
  members = [
    "serviceAccount:${google_service_account.test-account1.email}",
  ]
}
`, context)
}

func testAccBigtableInstanceIamBinding_update(context map[string]interface{}) string {
	return testBigtableInstanceIam(context) + acctest.Nprintf(`
resource "google_service_account" "test-account1" {
  account_id   = "%{account}-1"
  display_name = "Bigtable Instance IAM Testing Account"
}

resource "google_service_account" "test-account2" {
  account_id   = "%{account}-2"
  display_name = "Bigtable Instance IAM Testing Account"
}

resource "google_bigtable_instance_iam_binding" "binding" {
  instance = google_bigtable_instance.instance.name
  role     = "%{role}"
  members = [
    "serviceAccount:${google_service_account.test-account1.email}",
    "serviceAccount:${google_service_account.test-account2.email}",
  ]
}
`, context)
}

func testAccBigtableInstanceIamMember(context map[string]interface{}) string {
	return testBigtableInstanceIam(context) + acctest.Nprintf(`
resource "google_service_account" "test-account" {
  account_id   = "%{account}"
  display_name = "Bigtable Instance IAM Testing Account"
}

resource "google_bigtable_instance_iam_member" "member" {
  instance = google_bigtable_instance.instance.name
  role     = "%{role}"
  member   = "serviceAccount:${google_service_account.test-account.email}"
}
`, context)
}

func testAccBigtableInstanceIamPolicy(context map[string]interface{}) string {
	return testBigtableInstanceIam(context) + acctest.Nprintf(`
resource "google_service_account" "test-account" {
  account_id   = "%{account}"
  display_name = "Bigtable Instance IAM Testing Account"
}

data "google_iam_policy" "policy" {
  binding {
    role    = "%{role}"
    members = ["serviceAccount:${google_service_account.test-account.email}"]
  }
}

resource "google_bigtable_instance_iam_policy" "policy" {
  instance    = google_bigtable_instance.instance.name
  policy_data = data.google_iam_policy.policy.policy_data
}

data "google_bigtable_instance_iam_policy" "policy" {
  instance    = google_bigtable_instance.instance.name
}
`, context)
}

// Smallest instance possible for testing
func testBigtableInstanceIam(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_bigtable_instance" "instance" {
	name                  = "%{instance}"
    instance_type = "DEVELOPMENT"

    cluster {
      cluster_id   = "%{cluster}"
      zone         = "us-central1-b"
      storage_type = "HDD"
    }

    deletion_protection = false
}
`, context)
}
