package memorystore_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccMemorystoreAclPolicyDatasource(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMemorystoreAclPolicyDatasourceConfig(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_memorystore_acl_policy.default", "google_memorystore_acl_policy.test"),
				),
			},
		},
	})
}

func testAccMemorystoreAclPolicyDatasourceConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_memorystore_acl_policy" "test" {
  acl_policy_id = "tf-test-policy-%{random_suffix}"
  location      = "us-central1"
  rules {
    rule     = "on allkeys +get"
    username = "default"
  }
}

data "google_memorystore_acl_policy" "default" {
  acl_policy_id = google_memorystore_acl_policy.test.acl_policy_id
  location      = "us-central1"
}
`, context)
}
