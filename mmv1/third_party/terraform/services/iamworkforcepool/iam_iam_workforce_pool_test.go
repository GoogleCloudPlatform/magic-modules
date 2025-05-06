package iamworkforcepool_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccIAMWorkforcePoolIamMember(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"role":          "roles/iam.workforcePoolViewer",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIAMWorkforcePoolIamMember_basic(context),
			},
			{
				ResourceName: "google_iam_workforce_pool_iam_member.foo",
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					id := state.RootModule().Resources["google_iam_workforce_pool.my_pool"].Primary.Attributes["id"]
					return fmt.Sprintf("%s %s user:admin@hashicorptest.com",
						id,
						context["role"],
					), nil
				},
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccIAMWorkforcePoolIamBinding(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"role":          "roles/iam.workforcePoolViewer",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIAMWorkforcePoolIamBinding_basic(context),
			},
			{
				ResourceName: "google_iam_workforce_pool_iam_binding.foo",
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					id := state.RootModule().Resources["google_iam_workforce_pool.my_pool"].Primary.Attributes["id"]
					return fmt.Sprintf("%s %s",
						id,
						context["role"],
					), nil
				},
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test Iam Binding update
				Config: testAccIAMWorkforcePoolIamBinding_update(context),
			},
			{
				ResourceName: "google_iam_workforce_pool_iam_binding.foo",
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					id := state.RootModule().Resources["google_iam_workforce_pool.my_pool"].Primary.Attributes["id"]
					return fmt.Sprintf("%s %s",
						id,
						context["role"],
					), nil
				},
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccIAMWorkforcePoolIamPolicy(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"role":          "roles/iam.workforcePoolViewer",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIAMWorkforcePoolIamPolicy_basic(context),
			},
			{
				ResourceName:      "google_iam_workforce_pool_iam_policy.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccIAMWorkforcePoolIamPolicy_emptyBinding(context),
			},
			{
				ResourceName:      "google_iam_workforce_pool_iam_policy.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccIAMWorkforcePoolIamMember_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_iam_workforce_pool" "my_pool" {
  workforce_pool_id = "my-pool-%{random_suffix}"
  parent            = "organizations/%{org_id}"
  location          = "global"
}

resource "google_iam_workforce_pool_iam_member" "foo" {
  location          = google_iam_workforce_pool.my_pool.location
  workforce_pool_id = google_iam_workforce_pool.my_pool.workforce_pool_id
  role              = "%{role}"
  member            = "user:admin@hashicorptest.com"
}
`, context)
}

func testAccIAMWorkforcePoolIamBinding_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_iam_workforce_pool" "my_pool" {
  workforce_pool_id = "my-pool-%{random_suffix}"
  parent            = "organizations/%{org_id}"
  location          = "global"
}

resource "google_iam_workforce_pool_iam_binding" "foo" {
  location          = google_iam_workforce_pool.my_pool.location
  workforce_pool_id = google_iam_workforce_pool.my_pool.workforce_pool_id
  role              = "%{role}"
  members           = ["user:admin@hashicorptest.com"]
}
`, context)
}

func testAccIAMWorkforcePoolIamBinding_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_iam_workforce_pool" "my_pool" {
  workforce_pool_id = "my-pool-%{random_suffix}"
  parent            = "organizations/%{org_id}"
  location          = "global"
}

resource "google_iam_workforce_pool_iam_binding" "foo" {
  location          = google_iam_workforce_pool.my_pool.location
  workforce_pool_id = google_iam_workforce_pool.my_pool.workforce_pool_id
  role              = "%{role}"
  members           = ["user:admin@hashicorptest.com", "user:gterraformtest1@gmail.com"]
}
`, context)
}

func testAccIAMWorkforcePoolIamPolicy_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_iam_workforce_pool" "my_pool" {
  workforce_pool_id = "my-pool-%{random_suffix}"
  parent            = "organizations/%{org_id}"
  location          = "global"
}

data "google_iam_policy" "foo" {
  binding {
    role    = "%{role}"
    members = ["user:admin@hashicorptest.com"]
  }
}

resource "google_iam_workforce_pool_iam_policy" "foo" {
  location          = google_iam_workforce_pool.my_pool.location
  workforce_pool_id = google_iam_workforce_pool.my_pool.workforce_pool_id
  policy_data       = data.google_iam_policy.foo.policy_data
}
`, context)
}

func testAccIAMWorkforcePoolIamPolicy_emptyBinding(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_iam_workforce_pool" "my_pool" {
  workforce_pool_id = "my-pool-%{random_suffix}"
  parent            = "organizations/%{org_id}"
  location          = "global"
}

data "google_iam_policy" "foo" {
}

resource "google_iam_workforce_pool_iam_policy" "foo" {
  location          = google_iam_workforce_pool.my_pool.location
  workforce_pool_id = google_iam_workforce_pool.my_pool.workforce_pool_id
  policy_data       = data.google_iam_policy.foo.policy_data
}
`, context)
}
