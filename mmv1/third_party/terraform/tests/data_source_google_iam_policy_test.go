package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccComputeImageIamPolicy_googleIamPolicyDataSource(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
		"role":          "roles/compute.imageUser",
		"member":        "user:admin@hashicorptest.com",
		"title":         "expires_after_2019_12_31",
		"expression":    "request.time < timestamp('2020-01-01T00:00:00Z')",
	}

	// Note: `<` is changed to `\u003c` in the expression
	expectedComputedPolicyData := Nprintf(`{"bindings":[{"members":["%{member}"],"role":"%{role}"},{"condition":{"expression":"request.time \u003c timestamp('2020-01-01T00:00:00Z')","title":"%{title}"},"members":["%{member}"],"role":"%{role}"}]}`, context)

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeImageIamPolicy_googleIamPolicyDataSource(context),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_iam_policy.foo", "policy_data", expectedComputedPolicyData),
					resource.TestCheckResourceAttr("google_compute_image_iam_policy.foo", "policy_data", expectedComputedPolicyData),
				),
			},
		},
	})
}

func testAccComputeImageIamPolicy_googleIamPolicyDataSource(context map[string]interface{}) string {
	return Nprintf(`
resource "google_compute_image" "example" {
  name = "tf-test-example-image%{random_suffix}"

  raw_disk {
    source = "https://storage.googleapis.com/bosh-gce-raw-stemcells/bosh-stemcell-97.98-google-kvm-ubuntu-xenial-go_agent-raw-1557960142.tar.gz"
  }
}

data "google_iam_policy" "foo" {
  binding {
    role = "%{role}"
    members = ["%{member}"]
    condition {
        title       = "%{title}"
        expression  = "%{expression}"
    }
  }
  binding {
    role = "%{role}"
    members = ["%{member}"]
  }
}

resource "google_compute_image_iam_policy" "foo" {
  project = google_compute_image.example.project
  image = google_compute_image.example.name
  policy_data = data.google_iam_policy.foo.policy_data
}
`, context)
}
