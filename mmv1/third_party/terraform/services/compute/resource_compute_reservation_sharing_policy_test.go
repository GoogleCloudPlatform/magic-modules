package compute_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccComputeReservationSharingPolicy(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":         envvar.GetTestProjectFromEnv(),
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
		"random_suffix":   acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeReservationDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeReservationWithDefaultReservationSharingPolicy(context),
			},
			{
				Config: testAccComputeReservationWithReservationSharingPolicyAllowAll(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
					},
				},
			},
			{
				Config: testAccComputeReservationWithDefaultReservationSharingPolicy(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func testAccComputeReservationWithDefaultReservationSharingPolicy(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_reservation" "gce_reservation" {
  project = "%{project}"
  name = "tf-test-%{random_suffix}"
  zone = "us-central1-a"

  specific_reservation {
    count = 1
    instance_properties {
      machine_type     = "a2-highgpu-1g"
    }
  }
}
`, context)
}

func testAccComputeReservationWithReservationSharingPolicyAllowAll(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_reservation" "gce_reservation" {
  project = "%{project}"
  name = "tf-test-%{random_suffix}"
  zone = "us-central1-a"

  specific_reservation {
    count = 1
    instance_properties {
      machine_type     = "a2-highgpu-1g"
    }
  }

  reservation_sharing_policy {
    service_share_type = "ALLOW_ALL"
  }
}
`, context)
}
