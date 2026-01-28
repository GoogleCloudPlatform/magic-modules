package cloudidentity_test

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"google.golang.org/api/googleapi"
)

var (
	_ = fmt.Sprintf
	_ = log.Print
	_ = strconv.Atoi
	_ = strings.Trim
	_ = time.Now
	_ = resource.TestMain
	_ = terraform.NewState
	_ = envvar.TestEnvVar
	_ = tpgresource.SetLabels
	_ = transport_tpg.Config{}
	_ = googleapi.Error{}
)

func TestAccCloudIdentityPolicy_cloudidentityPolicyBasic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderBetaFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudIdentityPolicy_cloudidentityPolicyBasic(context),
			},
			{
				Config: testAccCloudIdentityPolicy_cloudidentityPolicyBasic_update(context),
			},
			{
				ResourceName:      "google_cloud_identity_policy.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCloudIdentityPolicy_cloudidentityPolicyBasic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloud_identity_policy" "primary" {
    provider = google-beta

    customer = "customers/C01234567%{random_suffix}"

    policy_query {
        org_unit = "orgUnits/03abcxyz%{random_suffix}"
        group = "groups/0123456789%{random_suffix}"
        query = "true%{random_suffix}"
    }

    setting {
        type = "something.googleapis.com/SettingType%{random_suffix}"
	value_json = "{"enabled": true}%{random_suffix}"
    }
}
`, context)
}

func testAccCloudIdentityPolicy_cloudidentityPolicyBasic_update(context map[string]interface{}) string {
        return acctest.Nprintf(`
resource "google_cloud_identity_policy" "primary" {
    provider = google-beta

    customer = "customers/C01234567%{random_suffix}"

    policy_query {
        org_unit = "orgUnits/03abcxyz%{random_suffix}"
        group = "groups/0123456789%{random_suffix}"
        query = "true%{random_suffix}"
    }

    setting {
        type = "something.googleapis.com/SettingType%{random_suffix}"
        value_json = "{"enabled": false}"
    }
}
`, context)
}
