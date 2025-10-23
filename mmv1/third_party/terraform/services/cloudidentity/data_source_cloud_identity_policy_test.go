package cloudidentity_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"google.golang.org/api/cloudidentity/v1"
)

func TestAccDataSourceGoogleCloudIdentityPolicy(t *testing.T) {
	var policyName string

	acctest.VcrTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccTestPreCheck(t)

			ctx := context.Background()
			ci, err := cloudidentity.NewService(ctx)
			if err != nil {
				t.Skipf("Cloud Identity service  not available in this env: %v", err)
			}

			lst, err := ci.Policies.List().Context(ctx).Do()
			if err != nil {
				t.Skipf("Cloud Identity Policies API not accessible in this env: %v", err)
			}
			if lst == nil || len(lst.Policies) == 0 {
				t.Skip("No Cloud Identity policies found in this customer; skipping data source test.")
			}

			policyName = lst.Policies[0].Name
			t.Logf("Discovered policy: %s", policyName)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				    data "google_cloud_identity_policy" "test" {
					    name = "%s"
				    }
				    `, policyName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_cloud_identity_policy.test", "name", policyName),
					resource.TestCheckResourceAttrSet("data.google_cloud_identity_policy.test", "customer"),
					resource.TestCheckResourceAttrSet("data.google_cloud_identity_policy.test", "type"),
				),
			},
		},
	})
}
