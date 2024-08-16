package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccComputeRegionNetworkFirewallPolicyRule_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"org_name":      fmt.Sprintf("organizations/%s", envvar.GetTestOrgFromEnv(t)),
		"region":        envvar.GetTestRegionFromEnv(),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRegionNetworkFirewallPolicyRule_start(context),
			},
			{
				ResourceName:      "google_compute_region_network_firewall_policy_rule.fw_policy_rule1",
				ImportState:       true,
				ImportStateVerify: true,
				// Referencing using ID causes import to fail
				ImportStateVerifyIgnore: []string{"firewall_policy"},
			},
			{
				Config: testAccComputeRegionNetworkFirewallPolicyRule_update(context),
			},
			{
				ResourceName:      "google_compute_region_network_firewall_policy_rule.fw_policy_rule1",
				ImportState:       true,
				ImportStateVerify: true,
				// Referencing using ID causes import to fail
				ImportStateVerifyIgnore: []string{"firewall_policy"},
			},
			{
				Config: testAccComputeRegionNetworkFirewallPolicyRule_removeConfigs(context),
			},
			{
				ResourceName:      "google_compute_region_network_firewall_policy_rule.fw_policy_rule1",
				ImportState:       true,
				ImportStateVerify: true,
				// Referencing using ID causes import to fail
				ImportStateVerifyIgnore: []string{"firewall_policy"},
			},
			{
				Config: testAccComputeRegionNetworkFirewallPolicyRule_start(context),
			},
			{
				ResourceName:      "google_compute_region_network_firewall_policy_rule.fw_policy_rule1",
				ImportState:       true,
				ImportStateVerify: true,
				// Referencing using ID causes import to fail
				ImportStateVerifyIgnore: []string{"firewall_policy"},
			},
		},
	})
}

func TestAccComputeRegionNetworkFirewallPolicyRule_multipleRules(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"org_name":      fmt.Sprintf("organizations/%s", envvar.GetTestOrgFromEnv(t)),
		"region":        envvar.GetTestRegionFromEnv(),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRegionNetworkFirewallPolicyRule_multiple(context),
			},
			{
				ResourceName:      "google_compute_region_network_firewall_policy_rule.fw_policy_rule1",
				ImportState:       true,
				ImportStateVerify: true,
				// Referencing using ID causes import to fail
				ImportStateVerifyIgnore: []string{"firewall_policy"},
			},
			{
				ResourceName:      "google_compute_region_network_firewall_policy_rule.fw_policy_rule2",
				ImportState:       true,
				ImportStateVerify: true,
				// Referencing using ID causes import to fail
				ImportStateVerifyIgnore: []string{"firewall_policy"},
			},
			{
				Config: testAccComputeRegionNetworkFirewallPolicyRule_multipleAdd(context),
			},
			{
				ResourceName:      "google_compute_region_network_firewall_policy_rule.fw_policy_rule3",
				ImportState:       true,
				ImportStateVerify: true,
				// Referencing using ID causes import to fail
				ImportStateVerifyIgnore: []string{"firewall_policy"},
			},
			{
				Config: testAccComputeRegionNetworkFirewallPolicyRule_multipleRemove(context),
			},
		},
	})
}

func testAccComputeRegionNetworkFirewallPolicyRule_start(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_service_account" "service_account1" {
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

resource "google_compute_region_network_firewall_policy" "fw_policy" {
  name        = "tf-test-policy-%{random_suffix}"
  description = "Resource created for Terraform acceptance testing"
}

resource "google_network_security_address_group" "address_group" {
  name        = "tf-test-policy%{random_suffix}"
  parent      = "%{org_name}"
  description = "Sample global networksecurity_address_group"
  location    = "%{region}"
  items       = ["208.80.154.224/32"]
  type        = "IPV4"
  capacity    = 100
}

resource "google_compute_region_network_firewall_policy_rule" "fw_policy_rule1" {
  firewall_policy = google_compute_region_network_firewall_policy.fw_policy.id
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
    dest_address_groups       = [google_network_security_address_group.address_group.id]
  }
}
`, context)
}

func testAccComputeRegionNetworkFirewallPolicyRule_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_service_account" "service_account1" {
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

resource "google_compute_region_network_firewall_policy" "fw_policy" {
  name        = "tf-test-policy-%{random_suffix}"
  description = "Resource created for Terraform acceptance testing"
}

resource "google_network_security_address_group" "address_group" {
  name        = "tf-test-policy%{random_suffix}"
  parent      = "%{org_name}"
  description = "Sample global networksecurity_address_group"
  location    = "%{region}"
  items       = ["208.80.154.224/32"]
  type        = "IPV4"
  capacity    = 100
}

resource "google_compute_region_network_firewall_policy_rule" "fw_policy_rule1" {
  firewall_policy         = google_compute_region_network_firewall_policy.fw_policy.id
  description             = "Resource created for Terraform acceptance testing"
  priority                = 9000
  enable_logging          = true
  action                  = "allow"
  direction               = "EGRESS"
  disabled                = false
  target_service_accounts = [google_service_account.service_account1.email]

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
    dest_address_groups       = [google_network_security_address_group.address_group.id]
  }
}
`, context)
}

func testAccComputeRegionNetworkFirewallPolicyRule_removeConfigs(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_service_account" "service_account1" {
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

resource "google_compute_region_network_firewall_policy" "fw_policy" {
  name        = "tf-test-policy-%{random_suffix}"
  description = "Resource created for Terraform acceptance testing"
}

resource "google_network_security_address_group" "address_group" {
  name        = "tf-test-policy%{random_suffix}"
  parent      = "%{org_name}"
  description = "Sample global networksecurity_address_group"
  location    = "%{region}"
  items       = ["208.80.154.224/32"]
  type        = "IPV4"
  capacity    = 100
}

resource "google_compute_region_network_firewall_policy_rule" "fw_policy_rule1" {
  firewall_policy         = google_compute_region_network_firewall_policy.fw_policy.id
  description             = "Test description"
  priority                = 9000
  enable_logging          = false
  action                  = "deny"
  direction               = "INGRESS"
  disabled                = true
  target_service_accounts = [
    google_service_account.service_account1.email,
    google_service_account.service_account2.email
  ]

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
}
`, context)
}

func testAccComputeRegionNetworkFirewallPolicyRule_multiple(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network" "network1" {
  name                    = "tf-test-%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_region_network_firewall_policy" "fw_policy" {
  name        = "tf-test-policy-%{random_suffix}"
  description = "Resource created for Terraform acceptance testing"
}

resource "google_compute_region_network_firewall_policy_association" "fw_policy_a" {
  name              = "tf-test-policy-a-%{random_suffix}"
  attachment_target = google_compute_network.network1.id
  firewall_policy   = google_compute_region_network_firewall_policy.fw_policy.id
}

resource "google_network_security_address_group" "address_group" {
  name        = "tf-test-policy%{random_suffix}"
  parent      = "%{org_name}"
  description = "Sample global networksecurity_address_group"
  location    = "%{region}"
  items       = ["208.80.154.224/32"]
  type        = "IPV4"
  capacity    = 100
}

resource "google_compute_region_network_firewall_policy_rule" "fw_policy_rule1" {
  firewall_policy = google_compute_region_network_firewall_policy.fw_policy.id
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
    dest_address_groups       = [google_network_security_address_group.address_group.id]
  }
}

resource "google_compute_region_network_firewall_policy_rule" "fw_policy_rule2" {
  firewall_policy = google_compute_region_network_firewall_policy.fw_policy.id
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
    src_address_groups       = [google_network_security_address_group.address_group.id]
  }
}
`, context)
}

func testAccComputeRegionNetworkFirewallPolicyRule_multipleAdd(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network" "network1" {
  name                    = "tf-test-%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_region_network_firewall_policy" "fw_policy" {
  name        = "tf-test-policy-%{random_suffix}"
  description = "Resource created for Terraform acceptance testing"
}

resource "google_compute_region_network_firewall_policy_association" "fw_policy_a" {
  name              = "tf-test-policy-a-%{random_suffix}"
  attachment_target = google_compute_network.network1.id
  firewall_policy   = google_compute_region_network_firewall_policy.fw_policy.id
}

resource "google_network_security_address_group" "address_group" {
  name        = "tf-test-policy%{random_suffix}"
  parent      = "%{org_name}"
  description = "Sample global networksecurity_address_group"
  location    = "%{region}"
  items       = ["208.80.154.224/32"]
  type        = "IPV4"
  capacity    = 100
}

resource "google_compute_region_network_firewall_policy_rule" "fw_policy_rule1" {
  firewall_policy = google_compute_region_network_firewall_policy.fw_policy.id
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
    dest_address_groups       = [google_network_security_address_group.address_group.id]
  }
}

resource "google_compute_region_network_firewall_policy_rule" "fw_policy_rule2" {
  firewall_policy = google_compute_region_network_firewall_policy.fw_policy.id
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
    src_address_groups       = [google_network_security_address_group.address_group.id]
  }
}

resource "google_compute_region_network_firewall_policy_rule" "fw_policy_rule3" {
  firewall_policy = google_compute_region_network_firewall_policy.fw_policy.id
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
    src_address_groups       = [google_network_security_address_group.address_group.id]
  }
}
`, context)
}

func testAccComputeRegionNetworkFirewallPolicyRule_multipleRemove(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network" "network1" {
  name                    = "tf-test-%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_region_network_firewall_policy" "fw_policy" {
  name        = "tf-test-policy-%{random_suffix}"
  description = "Resource created for Terraform acceptance testing"
}

resource "google_compute_region_network_firewall_policy_association" "fw_policy_a" {
  name              = "tf-test-policy-a-%{random_suffix}"
  attachment_target = google_compute_network.network1.id
  firewall_policy   = google_compute_region_network_firewall_policy.fw_policy.id
}

resource "google_network_security_address_group" "address_group" {
  name        = "tf-test-policy%{random_suffix}"
  parent      = "%{org_name}"
  description = "Sample global networksecurity_address_group"
  location    = "%{region}"
  items       = ["208.80.154.224/32"]
  type        = "IPV4"
  capacity    = 100
}

resource "google_compute_region_network_firewall_policy_rule" "fw_policy_rule1" {
  firewall_policy = google_compute_region_network_firewall_policy.fw_policy.id
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

resource "google_compute_region_network_firewall_policy_rule" "fw_policy_rule3" {
  firewall_policy = google_compute_region_network_firewall_policy.fw_policy.id
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
