<% autogen_exception -%>
package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccComputeNetworkFirewallPolicyRule_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"org_name":      fmt.Sprintf("organizations/%s", envvar.GetTestOrgFromEnv(t)),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeNetworkFirewallPolicyRule_start(context),
			},
			{
				ResourceName:      "google_compute_network_firewall_policy_rule.default",
				ImportState:       true,
				ImportStateVerify: true,
				// Referencing using ID causes import to fail
				ImportStateVerifyIgnore: []string{"firewall_policy"},
			},
			{
				Config: testAccComputeNetworkFirewallPolicyRule_update(context),
			},
			{
				ResourceName:      "google_compute_network_firewall_policy_rule.default",
				ImportState:       true,
				ImportStateVerify: true,
				// Referencing using ID causes import to fail
				ImportStateVerifyIgnore: []string{"firewall_policy", "target_resources"},
			},
			{
				Config: testAccComputeNetworkFirewallPolicyRule_removeConfigs(context),
			},
			{
				ResourceName:      "google_compute_network_firewall_policy_rule.default",
				ImportState:       true,
				ImportStateVerify: true,
				// Referencing using ID causes import to fail
				ImportStateVerifyIgnore: []string{"firewall_policy", "target_resources"},
			},
			{
				Config: testAccComputeNetworkFirewallPolicyRule_start(context),
			},
			{
				ResourceName:      "google_compute_network_firewall_policy_rule.default",
				ImportState:       true,
				ImportStateVerify: true,
				// Referencing using ID causes import to fail
				ImportStateVerifyIgnore: []string{"firewall_policy"},
			},
		},
	})
}

func TestAccComputeNetworkFirewallPolicyRule_securityProfileGroupTls_update(t *testing.T) {
  t.Parallel()

  context := map[string]interface{}{
    "random_suffix": acctest.RandString(t, 10),
    "org_name":      fmt.Sprintf("organizations/%s", envvar.GetTestOrgFromEnv(t)),
  }

  acctest.VcrTest(t, resource.TestCase{
    PreCheck:                 func() { acctest.AccTestPreCheck(t) },
    ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
    Steps: []resource.TestStep{
      {
        Config: testAccComputeNetworkFirewallPolicyRule_securityProfileGroupTls_basic(context),
      },
      {
        ResourceName:      "google_compute_network_firewall_policy_rule.default",
        ImportState:       true,
        ImportStateVerify: true,
        // Referencing using ID causes import to fail
        ImportStateVerifyIgnore: []string{"firewall_policy"},
      },
      {
        Config: testAccComputeNetworkFirewallPolicyRule_securityProfileGroupTls_update(context),
      },
      {
        ResourceName:      "google_compute_network_firewall_policy_rule.default",
        ImportState:       true,
        ImportStateVerify: true,
        // Referencing using ID causes import to fail
        ImportStateVerifyIgnore: []string{"firewall_policy", "target_resources"},
      },
    },
  })
}

func testAccComputeNetworkFirewallPolicyRule_securityProfileGroup_basic(context map[string]interface{}) string {
  return acctest.Nprintf(`
resource "google_compute_network" "network1" {
  name                    = "tf-test-%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_network_security_security_profile" "security_profile" {
    provider    = google-beta
    name        = "tf-test-my-sp%{random_suffix}"
    type        = "THREAT_PREVENTION"
    parent      = "organizations/%{org_name}"
    location    = "global"
}

resource "google_network_security_security_profile_group" "security_profile_group" {
    provider                  = google-beta
    name                      = "tf-test-my-spg%{random_suffix}"
    parent                    = "organizations/%{org_name}"
    location                  = "global"
    description               = "My security profile group."
    threat_prevention_profile = google_network_security_security_profile.security_profile.id
}

resource "google_compute_network_firewall_policy" "default" {
  parent      = google_compute_network.network1.id
  short_name  = "tf-test-policy-%{random_suffix}"
  description = "Resource created for Terraform acceptance testing"
}

resource "google_compute_network_firewall_policy_rule" "default" {
  firewall_policy        = google_compute_network_firewall_policy.default.id
  description            = "Resource created for Terraform acceptance testing"
  priority               = 9000
  enable_logging         = true
  action                 = "apply_security_profile_group"
  security_profile_group = google_network_security_security_profile_group.security_profile_group.id
  direction              = "INGRESS"
  disabled               = false
  match {
    layer4_configs {
      ip_protocol = "tcp"
      ports       = [80, 8080]
    }
    src_ip_ranges = ["11.100.0.1/32"]
  }
}
`, context)
}

func testAccComputeNetworkFirewallPolicyRule_securityProfileGroup_update(context map[string]interface{}) string {
  return acctest.Nprintf(`
resource "google_compute_network" "network1" {
  name                    = "tf-test-%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_network_security_security_profile" "security_profile" {
    provider    = google-beta
    name        = "tf-test-my-security-profile%{random_suffix}"
    type        = "THREAT_PREVENTION"
    parent      = "organizations/%{org_name}"
    location    = "global"
}

resource "google_network_security_security_profile_group" "security_profile_group" {
    provider                  = google-beta
    name                      = "tf-test-my-sp-group%{random_suffix}"
    parent                    = "organizations/%{org_name}"
    location                  = "global"
    description               = "My security profile group."
    threat_prevention_profile = google_network_security_security_profile.security_profile.id
}

resource "google_network_security_security_profile_group" "security_profile_group_updated" {
    provider                  = google-beta
    name                      = "tf-test-my-spg-updated%{random_suffix}"
    parent                    = "organizations/%{org_name}"
    location                  = "global"
    description               = "My security profile group."
    threat_prevention_profile = google_network_security_security_profile.security_profile.id
}

resource "google_compute_network_firewall_policy" "default" {
  parent      = google_compute_network.network1.id
  short_name  = "tf-test-policy-%{random_suffix}"
  description = "Resource created for Terraform acceptance testing"
}

resource "google_compute_network_firewall_policy_rule" "default" {
  firewall_policy        = google_compute_network_firewall_policy.default.id
  description            = "Resource created for Terraform acceptance testing"
  priority               = 9000
  enable_logging         = true
  action                 = "apply_security_profile_group"
  security_profile_group = google_network_security_security_profile_group.security_profile_group_updated.id
  direction              = "INGRESS"
  disabled               = false
  match {
    layer4_configs {
      ip_protocol = "tcp"
      ports       = [80, 8080]
    }
    src_ip_ranges = ["11.100.0.1/32"]
  }
}
`, context)
}

func testAccComputeNetworkFirewallPolicyRule_start(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_service_account" "service_account" {
  account_id = "tf-test-sa-%{random_suffix}"
}

resource "google_service_account" "service_account2" {
  account_id = "tf-test-sa2-%{random_suffix}"
}

resource "google_compute_network" "network1" {
  name                    = "tf-test-%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_network" "network2" {
  name                    = "tf-test-2-%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_network_firewall_policy" "default" {
  parent      = google_compute_network.network1.id
  short_name  = "tf-test-policy-%{random_suffix}"
  description = "Resource created for Terraform acceptance testing"
}

resource "google_network_security_address_group" "basic_global_networksecurity_address_group" {
  name        = "tf-test-policy%{random_suffix}"
  parent      = "%{org_name}"
  description = "Sample global networksecurity_address_group"
  location    = "global"
  items       = ["208.80.154.224/32"]
  type        = "IPV4"
  capacity    = 100
}

resource "google_compute_network_firewall_policy_rule" "default" {
  firewall_policy = google_compute_network_firewall_policy.default.id
  description     = "Resource created for Terraform acceptance testing"
  priority        = 9000
  enable_logging  = true
  action          = "allow"
  direction       = "EGRESS"
  disabled        = false
  match {
    layer4_configs {
      ip_protocol = "tcp"
      ports       = [80, 8080]
    }
    dest_ip_ranges = ["11.100.0.1/32"]
    dest_fqdns                = []
    dest_region_codes         = []
    dest_threat_intelligences = []
    dest_address_groups       = [google_network_security_address_group.basic_global_networksecurity_address_group.id]
  }
}
`, context)
}

func testAccComputeNetworkFirewallPolicyRule_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_service_account" "service_account" {
  account_id = "tf-test-sa-%{random_suffix}"
}

resource "google_service_account" "service_account2" {
  account_id = "tf-test-sa2-%{random_suffix}"
}

resource "google_compute_network" "network1" {
  name = "tf-test-%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_network" "network2" {
  name = "tf-test-2-%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_network_firewall_policy" "default" {
  parent      = google_compute_network.network1.id
  short_name  = "tf-test-policy-%{random_suffix}"
  description = "Resource created for Terraform acceptance testing"
}

resource "google_network_security_address_group" "basic_global_networksecurity_address_group" {
  name        = "tf-test-policy%{random_suffix}"
  parent      = "%{org_name}"
  description = "Sample global networksecurity_address_group"
  location    = "global"
  items       = ["208.80.154.224/32"]
  type        = "IPV4"
  capacity    = 100
}

resource "google_compute_network_firewall_policy_rule" "default" {
  firewall_policy = google_compute_network_firewall_policy.default.id
  description     = "Resource created for Terraform acceptance testing"
  priority        = 9000
  enable_logging  = true
  action          = "allow"
  direction       = "EGRESS"
  disabled        = false
  match {
    layer4_configs {
      ip_protocol = "tcp"
      ports       = [8080]
    }
    layer4_configs {
      ip_protocol = "udp"
      ports       = [22]
    }
    dest_ip_ranges            = ["11.100.0.1/32", "10.0.0.0/24"]
    dest_fqdns                = ["google.com"]
    dest_region_codes         = ["US"]
    dest_threat_intelligences = ["iplist-known-malicious-ips"]
    src_address_groups        = []
    dest_address_groups       = [google_network_security_address_group.basic_global_networksecurity_address_group.id]
  }
  target_resources        = [google_compute_network.network1.self_link, google_compute_network.network2.self_link]
  target_service_accounts = [google_service_account.service_account.email]
}
`, context)
}

func testAccComputeNetworkFirewallPolicyRule_removeConfigs(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_service_account" "service_account" {
  account_id = "tf-test-sa-%{random_suffix}"
}

resource "google_service_account" "service_account2" {
  account_id = "tf-test-sa2-%{random_suffix}"
}

resource "google_compute_network" "network1" {
  name                    = "tf-test-%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_network" "network2" {
  name                    = "tf-test-2-%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_network_firewall_policy" "default" {
  parent      = google_compute_network.network1.id
  short_name  = "tf-test-policy-%{random_suffix}"
  description = "Resource created for Terraform acceptance testing"
}

resource "google_network_security_address_group" "basic_global_networksecurity_address_group" {
  name        = "tf-test-policy%{random_suffix}"
  parent      = "%{org_name}"
  description = "Sample global networksecurity_address_group"
  location    = "global"
  items       = ["208.80.154.224/32"]
  type        = "IPV4"
  capacity    = 100
}

resource "google_compute_network_firewall_policy_rule" "default" {
  firewall_policy = google_compute_network_firewall_policy.default.id
  description     = "Test description"
  priority        = 9000
  enable_logging  = false
  action          = "deny"
  direction       = "INGRESS"
  disabled        = true
  match {
    layer4_configs {
      ip_protocol = "udp"
      ports       = [22]
    }
    src_ip_ranges            = ["11.100.0.1/32", "10.0.0.0/24"]
    src_fqdns                = ["google.com"]
    src_region_codes         = ["US"]
    src_threat_intelligences = ["iplist-known-malicious-ips"]
  }
  target_resources        = [google_compute_network.network1.self_link]
  target_service_accounts = [google_service_account.service_account.email, google_service_account.service_account2.email]
}
`, context)
}

func TestAccComputeNetworkFirewallPolicyRule_multipleRules(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"org_name":      fmt.Sprintf("organizations/%s", envvar.GetTestOrgFromEnv(t)),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeNetworkFirewallPolicyRule_multiple(context),
			},
			{
				ResourceName:      "google_compute_network_firewall_policy_rule.rule1",
				ImportState:       true,
				ImportStateVerify: true,
				// Referencing using ID causes import to fail
				ImportStateVerifyIgnore: []string{"firewall_policy"},
			},
			{
				ResourceName:      "google_compute_network_firewall_policy_rule.rule2",
				ImportState:       true,
				ImportStateVerify: true,
				// Referencing using ID causes import to fail
				ImportStateVerifyIgnore: []string{"firewall_policy"},
			},
			{
				Config: testAccComputeNetworkFirewallPolicyRule_multipleAdd(context),
			},
			{
				ResourceName:      "google_compute_network_firewall_policy_rule.rule3",
				ImportState:       true,
				ImportStateVerify: true,
				// Referencing using ID causes import to fail
				ImportStateVerifyIgnore: []string{"firewall_policy"},
			},
			{
				Config: testAccComputeNetworkFirewallPolicyRule_multipleRemove(context),
			},
		},
	})
}

func testAccComputeNetworkFirewallPolicyRule_multiple(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network_firewall_policy" "default" {
  parent      = google_compute_network.network1.id
  short_name  = "tf-test-policy-%{random_suffix}"
  description = "Resource created for Terraform acceptance testing"
}

resource "google_network_security_address_group" "basic_global_networksecurity_address_group" {
  name        = "tf-test-policy%{random_suffix}"
  parent      = "%{org_name}"
  description = "Sample global networksecurity_address_group"
  location    = "global"
  items       = ["208.80.154.224/32"]
  type        = "IPV4"
  capacity    = 100
}

resource "google_compute_network_firewall_policy_rule" "rule1" {
  firewall_policy = google_compute_network_firewall_policy.default.id
  description     = "Resource created for Terraform acceptance testing"
  priority        = 9000
  enable_logging  = true
  action          = "allow"
  direction       = "EGRESS"
  disabled        = false
  match {
    layer4_configs {
      ip_protocol = "tcp"
      ports       = [80, 8080]
    }
    dest_ip_ranges            = ["11.100.0.1/32"]
    dest_fqdns                = ["google.com"]
    dest_region_codes         = ["US"]
    dest_threat_intelligences = ["iplist-known-malicious-ips"]
    dest_address_groups       = [google_network_security_address_group.basic_global_networksecurity_address_group.id]
  }
}

resource "google_compute_network_firewall_policy_rule" "rule2" {
  firewall_policy = google_compute_network_firewall_policy.default.id
  description     = "Resource created for Terraform acceptance testing"
  priority        = 9001
  enable_logging  = false
  action          = "deny"
  direction       = "INGRESS"
  disabled        = false
  match {
    layer4_configs {
      ip_protocol = "tcp"
      ports       = [80, 8080]
    }
    layer4_configs {
      ip_protocol = "all"
    }
    src_ip_ranges            = ["11.100.0.1/32"]
    src_fqdns                = ["google.com"]
    src_region_codes         = ["US"]
    src_threat_intelligences = ["iplist-known-malicious-ips"]
    src_address_groups       = [google_network_security_address_group.basic_global_networksecurity_address_group.id]
  }
}
`, context)
}

func testAccComputeNetworkFirewallPolicyRule_multipleAdd(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network_firewall_policy" "default" {
  parent      = google_compute_network.network1.id
  short_name  = "tf-test-policy-%{random_suffix}"
  description = "Description Update"
}

resource "google_network_security_address_group" "basic_global_networksecurity_address_group" {
  name        = "tf-test-policy%{random_suffix}"
  parent      = "%{org_name}"
  description = "Sample global networksecurity_address_group"
  location    = "global"
  items       = ["208.80.154.224/32"]
  type        = "IPV4"
  capacity    = 100
}

resource "google_compute_network_firewall_policy_rule" "rule1" {
  firewall_policy = google_compute_network_firewall_policy.default.id
  description     = "Resource created for Terraform acceptance testing"
  priority        = 9000
  enable_logging  = true
  action          = "allow"
  direction       = "EGRESS"
  disabled        = false
  match {
    layer4_configs {
      ip_protocol = "tcp"
    }
    dest_ip_ranges            = ["11.100.0.1/32"]
    dest_fqdns                = ["google.com"]
    dest_region_codes         = ["US"]
    dest_threat_intelligences = ["iplist-known-malicious-ips"]
    dest_address_groups       = [google_network_security_address_group.basic_global_networksecurity_address_group.id]
  }
}

resource "google_compute_network_firewall_policy_rule" "rule2" {
  firewall_policy = google_compute_network_firewall_policy.default.id
  description     = "Resource created for Terraform acceptance testing"
  priority        = 9001
  enable_logging  = false
  action          = "deny"
  direction       = "INGRESS"
  disabled        = false
  match {
    layer4_configs {
      ip_protocol = "tcp"
      ports       = [80, 8080]
    }
    layer4_configs {
      ip_protocol = "all"
    }
    src_ip_ranges            = ["11.100.0.1/32"]
    src_fqdns                = ["google.com"]
    src_region_codes         = ["US"]
    src_threat_intelligences = ["iplist-known-malicious-ips"]
    src_address_groups       = [google_network_security_address_group.basic_global_networksecurity_address_group.id]
  }
}

resource "google_compute_network_firewall_policy_rule" "rule3" {
  firewall_policy = google_compute_network_firewall_policy.default.id
  description     = "Resource created for Terraform acceptance testing"
  priority        = 40
  enable_logging  = true
  action          = "allow"
  direction       = "INGRESS"
  disabled        = true
  match {
    layer4_configs {
      ip_protocol = "udp"
      ports       = [8000]
    }
    src_ip_ranges            = ["11.100.0.1/32", "10.0.0.0/24"]
    src_fqdns                = ["google.com"]
    src_region_codes         = ["US"]
    src_threat_intelligences = ["iplist-known-malicious-ips"]
    src_address_groups       = [google_network_security_address_group.basic_global_networksecurity_address_group.id]
  }
}
`, context)
}

func testAccComputeNetworkFirewallPolicyRule_multipleRemove(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network_firewall_policy" "default" {
  parent      = google_compute_network.network1.id
  short_name  = "tf-test-policy-%{random_suffix}"
  description = "Resource created for Terraform acceptance testing"
}

resource "google_network_security_address_group" "basic_global_networksecurity_address_group" {
  name        = "tf-test-policy%{random_suffix}"
  parent      = "%{org_name}"
  description = "Sample global networksecurity_address_group"
  location    = "global"
  items       = ["208.80.154.224/32"]
  type        = "IPV4"
  capacity    = 100
}

resource "google_compute_network_firewall_policy_rule" "rule1" {
  firewall_policy = google_compute_network_firewall_policy.default.id
  description     = "Resource created for Terraform acceptance testing"
  priority        = 9000
  enable_logging  = true
  action          = "allow"
  direction       = "EGRESS"
  disabled        = false
  match {
    layer4_configs {
      ip_protocol = "tcp"
      ports       = [80, 8080]
    }
    dest_ip_ranges            = ["11.100.0.1/32"]
    dest_fqdns                = ["google.com"]
    dest_region_codes         = ["US"]
    dest_threat_intelligences = ["iplist-known-malicious-ips"]
  }
}

resource "google_compute_network_firewall_policy_rule" "rule3" {
  firewall_policy = google_compute_network_firewall_policy.default.id
  description     = "Resource created for Terraform acceptance testing"
  priority        = 40
  enable_logging  = true
  action          = "allow"
  direction       = "INGRESS"
  disabled        = true
  match {
    layer4_configs {
      ip_protocol = "udp"
      ports       = [8000]
    }
    src_ip_ranges            = ["11.100.0.1/32", "10.0.0.0/24"]
    src_fqdns                = ["google.com"]
    src_region_codes         = ["US"]
    src_threat_intelligences = ["iplist-known-malicious-ips"]
  }
}
`, context)
}
