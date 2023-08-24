package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-google-beta/google-beta/acctest"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/envvar"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestProvider_versionfive_upgrade(t *testing.T) {
	acctest.SkipIfVcr(t)
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	billingId := envvar.GetTestBillingAccountFromEnv(t)
	project := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	name1 := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	name2 := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	name3 := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	name4 := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck: func() { acctest.AccTestPreCheck(t) },
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"google": {
						VersionConstraint: "4.58.0",
						Source:            "hashicorp/google-beta",
					},
				},
				Config: testProvider_versionfive_upgrades(project, org, billingId, name1, name2, name3, name4),
			},
			{
				ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
				Config:                   testProvider_versionfive_upgrades(project, org, billingId, name1, name2, name3, name4),
			},
			{
				ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
				ResourceName:             "google_data_fusion_instance.unset",
				ImportState:              true,
				ImportStateVerify:        true,
			},
			{
				ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
				ResourceName:             "google_data_fusion_instance.set",
				ImportState:              true,
				ImportStateVerify:        true,
			},
			{
				ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
				ResourceName:             "google_data_fusion_instance.reference",
				ImportState:              true,
				ImportStateVerify:        true,
			},
		},
	})
}

func TestProvider_versionfive_upgrades_ignorereads(t *testing.T) {
	acctest.SkipIfVcr(t)
	t.Parallel()

	var itName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	var tpName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	var igmName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	var autoscalerName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	diskName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	policyName := fmt.Sprintf("tf-test-policy-%s", acctest.RandString(t, 10))

	endpointContext := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	var itNameRegion = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	var tpNameRegion = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	var igmNameRegion = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	var autoscalerNameRegion = fmt.Sprintf("tf-test-region-autoscaler-%s", acctest.RandString(t, 10))

	policyContext := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	attachmentContext := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck: func() { acctest.AccTestPreCheck(t) },
		CheckDestroy: resource.ComposeTestCheckFunc(testAccCheckComputeResourcePolicyDestroyProducer(t),
			testAccCheckComputeRegionAutoscalerDestroyProducer(t),
			testAccCheckComputeNetworkEndpointGroupDestroyProducer(t),
			testAccCheckComputeAutoscalerDestroyProducer(t),
		),
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"google": {
						VersionConstraint: "4.58.0",
						Source:            "hashicorp/google-beta",
					},
				},
				Config: testProvider_versionfive_upgrades_ignorereads(itName, tpName, igmName, autoscalerName, diskName, policyName,
					itNameRegion, tpNameRegion, igmNameRegion, autoscalerNameRegion, endpointContext,
					policyContext, attachmentContext),
			},
			{
				ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
				Config: testProvider_versionfive_upgrades_ignorereads(itName, tpName, igmName, autoscalerName, diskName, policyName,
					itNameRegion, tpNameRegion, igmNameRegion, autoscalerNameRegion, endpointContext,
					policyContext, attachmentContext),
			},
		},
	})
}

func testProvider_versionfive_upgrades(project, org, billing, name1, name2, name3, name4 string) string {
	return fmt.Sprintf(`
resource "google_project" "host" {
  project_id      = "%s"
  name            = "%s"
  org_id          = "%s"
  billing_account = "%s"
}

resource "google_project_service" "dfapi" {
  project = google_project.host.project_id
  service = "datafusion.googleapis.com"

  disable_dependent_services = false
}

resource "google_data_fusion_instance" "unset" {
  name   = "%s"
  type   = "BASIC"
  options = {
  	prober_test_run = "true"
  }
}

resource "google_data_fusion_instance" "set" {
  name   = "%s"
  region = "us-central1"
  type   = "BASIC"
  options = {
  	prober_test_run = "true"
  }
}

resource "google_data_fusion_instance" "reference" {
  project = google_project.host.project_id
  name   = "%s"
  type   = "DEVELOPER"
  options = {
  	prober_test_run = "true"
  }
  zone   = "us-central1-a"
  depends_on = [
    google_project_service.dfapi
  ]
}

resource "google_redis_instance" "overridewithnonstandardlogic" {
  name           = "%s"
  memory_size_gb = 1
  location_id    = "us-south1-a"
}


`, project, project, org, billing, name1, name2, name3, name4)
}

func testProvider_versionfive_upgrades_ignorereads(itName, tpName, igmName, autoscalerName, diskName, policyName, itNameRegion, tpNameRegion, igmNameRegion, autoscalerNameRegion string, endpointContext, policyContext, attachmentContext map[string]interface{}) string {
	return testAccComputeAutoscaler_basic(itName, tpName, igmName, autoscalerName) +
		testAccComputeDiskResourcePolicyAttachment_basic(diskName, policyName) +
		testAccComputeNetworkEndpointGroup_networkEndpointGroup(endpointContext) +
		testAccComputeRegionAutoscaler_basic(itNameRegion, tpNameRegion, igmNameRegion, autoscalerNameRegion) +
		testAccComputeResourcePolicy_resourcePolicyBasicExample(policyContext) +
		testAccComputeServiceAttachment_serviceAttachmentBasicExample(attachmentContext)
}
