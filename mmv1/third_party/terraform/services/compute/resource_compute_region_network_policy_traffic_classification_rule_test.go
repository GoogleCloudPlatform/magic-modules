package compute_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccComputeRegionNetworkPolicyTrafficClassificationRule_regionNetworkPolicyTrafficClassificationRuleBasicUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"region":        envvar.GetTestRegionFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeRegionNetworkPolicyTrafficClassificationRuleDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRegionNetworkPolicyTrafficClassificationRule_regionNetworkPolicyTrafficClassificationRuleBasicCreate(context),
			},
			{
				ResourceName:            "google_compute_region_network_policy_traffic_classification_rule.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"network_policy", "region"},
			},
			{
				Config: testAccComputeRegionNetworkPolicyTrafficClassificationRule_regionNetworkPolicyTrafficClassificationRuleBasicUpdate(context),
			},
			{
				ResourceName:            "google_compute_region_network_policy_traffic_classification_rule.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"network_policy", "region"},
			},
		},
	})
}

func testAccComputeRegionNetworkPolicyTrafficClassificationRule_regionNetworkPolicyTrafficClassificationRuleBasicCreate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_region_network_policy" "basic_regional_network_policy" {
  name        = "tf-test-nw-policy%{random_suffix}"
  description = "Sample regional network firewall policy"
  project     = "%{project_name}"
  region      = "%{region}"
}

resource "google_compute_region_network_policy_traffic_classification_rule" "primary" {
  rule_name               = "test-rule"
  description             = "This is a simple rule description"
  disabled                = false
  network_policy          = google_compute_region_network_policy.basic_regional_network_policy.name
  priority                = 1000
  region                  = "%{region}"

  action {
    traffic_class = "TC1"
    dscp_mode = "AUTO"
  }
  match {
    src_ip_ranges            = ["10.100.0.1/32"]
    dest_ip_ranges            = ["11.100.0.1/32"]
    layer4_configs {
      ip_protocol = "all"
    }
  }
}
`, context)
}

func testAccComputeRegionNetworkPolicyTrafficClassificationRule_regionNetworkPolicyTrafficClassificationRuleBasicUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_region_network_policy" "basic_regional_network_policy" {
  name        = "tf-test-nw-policy%{random_suffix}"
  description = "Sample regional network firewall policy"
  project     = "%{project_name}"
  region      = "%{region}"
}

resource "google_compute_region_network_policy_traffic_classification_rule" "primary" {
  rule_name               = "test-rule"
  description             = "This is a simple rule description"
  disabled                = false
  network_policy          = google_compute_region_network_policy.basic_regional_network_policy.name
  priority                = 1000
  region                  = "%{region}"

  action {
    dscp_mode = "CUSTOM"
    dscp_value = 47
    traffic_class = "TC5"
  }
  match {
    src_ip_ranges            = ["10.101.0.1/32"]
    dest_ip_ranges            = ["11.102.0.1/32"]
    layer4_configs {
      ip_protocol = "tcp"
      ports = ["80","443"]
    }
  }
}
`, context)
}
