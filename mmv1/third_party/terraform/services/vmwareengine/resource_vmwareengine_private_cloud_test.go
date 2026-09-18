package vmwareengine_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/services/kms"
	"github.com/hashicorp/terraform-provider-google/google/services/vmwareengine"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccVmwareenginePrivateCloud_vmwareEnginePrivateCloudUpdate(t *testing.T) {
	acctest.SkipIfVcr(t)
	t.Parallel()

	// Bootstrap KMS key in the same region as the private cloud
	kmsKey := kms.BootstrapKMSKeyInLocation(t, "me-west1")

	saSuffix := "gcp-sa-vmwareengine.iam.gserviceaccount.com"
	if strings.Contains(os.Getenv("GOOGLE_VMWAREENGINE_CUSTOM_ENDPOINT"), "autopush") {
		saSuffix = "gcp-sa-autopush-vmwareengine.iam.gserviceaccount.com"
	}

	context := map[string]interface{}{
		"region":               "me-west1", // region with allocated quota
		"random_suffix":        acctest.RandString(t, 10),
		"org_id":               envvar.GetTestOrgFromEnv(t),
		"billing_account":      envvar.GetTestBillingAccountFromEnv(t),
		"vmwareengine_project": os.Getenv("GOOGLE_VMWAREENGINE_PROJECT"),
		"kms_key_name":         kmsKey.CryptoKey.Name,
		"sa_suffix":            saSuffix,
		"use_cmek":             true,
	}

	gmekContext := make(map[string]interface{})
	for k, v := range context {
		gmekContext[k] = v
	}
	gmekContext["use_cmek"] = false

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		CheckDestroy: testAccCheckVmwareenginePrivateCloudDestroyProducer(t),
		Steps: []resource.TestStep{
			// 1. Create with CMEK
			{
				Config: testVmwareenginePrivateCloudCreateConfig(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_vmwareengine_private_cloud.ds",
						"google_vmwareengine_private_cloud.vmw-engine-pc",
						[]string{
							"deletion_delay_hours",
							"send_deletion_delay_hours_if_zero",
						}),
					testAccCheckGoogleVmwareengineNsxCredentialsMeta("data.google_vmwareengine_nsx_credentials.nsx-ds"),
					testAccCheckGoogleVmwareengineVcenterCredentialsMeta("data.google_vmwareengine_vcenter_credentials.vcenter-ds"),
					testAccCheckGoogleVmwareengineUpgradesMeta("data.google_vmwareengine_upgrades.upgrades-ds"),
					testAccCheckGoogleVmwareengineAnnouncementsMeta("data.google_vmwareengine_vcenter_credentials.announcements-ds"),
					resource.TestCheckResourceAttr("google_vmwareengine_private_cloud.vmw-engine-pc", "encryption_config.0.type", "CMEK"),
					resource.TestCheckResourceAttr("google_vmwareengine_private_cloud.vmw-engine-pc", "encryption_config.0.kms_key_name", kmsKey.CryptoKey.Name),
				),
			},
			{
				ResourceName:            "google_vmwareengine_private_cloud.vmw-engine-pc",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "update_time", "deletion_delay_hours", "send_deletion_delay_hours_if_zero"},
			},
			// 2. Perform update (change description + node count to 3, keep CMEK)
			{
				Config: testVmwareenginePrivateCloudUpdateNodeConfig(context, "Updated description, keeping CMEK"),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_vmwareengine_private_cloud.ds",
						"google_vmwareengine_private_cloud.vmw-engine-pc",
						[]string{
							"deletion_delay_hours",
							"send_deletion_delay_hours_if_zero",
						}),
					resource.TestCheckResourceAttr("google_vmwareengine_private_cloud.vmw-engine-pc", "description", "Updated description, keeping CMEK"),
					resource.TestCheckResourceAttr("google_vmwareengine_private_cloud.vmw-engine-pc", "encryption_config.0.type", "CMEK"),
					resource.TestCheckResourceAttr("google_vmwareengine_private_cloud.vmw-engine-pc", "encryption_config.0.kms_key_name", kmsKey.CryptoKey.Name),
				),
			},
			{
				ResourceName:            "google_vmwareengine_private_cloud.vmw-engine-pc",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "update_time", "deletion_delay_hours", "send_deletion_delay_hours_if_zero"},
			},

			// 3. Transition to GMEK (set type to GMEK)
			{
				Config: testVmwareenginePrivateCloudUpdateNodeConfig(gmekContext, "Reverted to GMEK"),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_vmwareengine_private_cloud.ds",
						"google_vmwareengine_private_cloud.vmw-engine-pc",
						[]string{
							"deletion_delay_hours",
							"send_deletion_delay_hours_if_zero",
						}),
					resource.TestCheckResourceAttr("google_vmwareengine_private_cloud.vmw-engine-pc", "description", "Reverted to GMEK"),
					resource.TestCheckResourceAttr("google_vmwareengine_private_cloud.vmw-engine-pc", "encryption_config.0.type", "GMEK"),
				),
			},
			{
				ResourceName:            "google_vmwareengine_private_cloud.vmw-engine-pc",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "update_time", "deletion_delay_hours", "send_deletion_delay_hours_if_zero"},
			},

			// 4. Transition back to CMEK
			{
				Config: testVmwareenginePrivateCloudUpdateNodeConfig(context, "Updated back to CMEK"),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_vmwareengine_private_cloud.ds",
						"google_vmwareengine_private_cloud.vmw-engine-pc",
						[]string{
							"deletion_delay_hours",
							"send_deletion_delay_hours_if_zero",
						}),
					resource.TestCheckResourceAttr("google_vmwareengine_private_cloud.vmw-engine-pc", "description", "Updated back to CMEK"),
					resource.TestCheckResourceAttr("google_vmwareengine_private_cloud.vmw-engine-pc", "encryption_config.0.type", "CMEK"),
					resource.TestCheckResourceAttr("google_vmwareengine_private_cloud.vmw-engine-pc", "encryption_config.0.kms_key_name", kmsKey.CryptoKey.Name),
				),
			},
			{
				ResourceName:            "google_vmwareengine_private_cloud.vmw-engine-pc",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "update_time", "deletion_delay_hours", "send_deletion_delay_hours_if_zero"},
			},

			// 5. Update Autoscale
			{
				Config: testVmwareenginePrivateCloudUpdateAutoscaleConfig(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_vmwareengine_private_cloud.ds",
						"google_vmwareengine_private_cloud.vmw-engine-pc",
						[]string{
							"deletion_delay_hours",
							"send_deletion_delay_hours_if_zero",
						}),
					resource.TestCheckResourceAttr("google_vmwareengine_private_cloud.vmw-engine-pc", "encryption_config.0.type", "CMEK"),
				),
			},
			{
				ResourceName:            "google_vmwareengine_private_cloud.vmw-engine-pc",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "update_time", "deletion_delay_hours", "send_deletion_delay_hours_if_zero"},
			},

			// 6. Delayed Delete
			{
				Config: testVmwareenginePrivateCloudDelayedDeleteConfig(context),
			},
			{
				ResourceName:            "google_vmwareengine_network.vmw-engine-nw",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name"},
			},

			// 7. Undelete
			{
				Config: testVmwareenginePrivateCloudUndeleteConfig(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_vmwareengine_private_cloud.ds",
						"google_vmwareengine_private_cloud.vmw-engine-pc",
						[]string{
							"deletion_delay_hours",
							"send_deletion_delay_hours_if_zero",
						}),
					resource.TestCheckResourceAttr("google_vmwareengine_private_cloud.vmw-engine-pc", "encryption_config.0.type", "CMEK"),
				),
			},
			{
				ResourceName:            "google_vmwareengine_private_cloud.vmw-engine-pc",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "update_time", "deletion_delay_hours", "send_deletion_delay_hours_if_zero"},
			},

			// 8. Subnet Import
			{
				Config: testVmwareengineSubnetImportConfig(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_vmwareengine_subnet.subnet-ds", "google_vmwareengine_subnet.vmw-engine-subnet"),
				),
			},
			{
				ResourceName:            "google_vmwareengine_subnet.vmw-engine-subnet",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent", "name"},
			},

			// 9. Subnet Update
			{
				Config: testVmwareengineSubnetUpdateConfig(context),
			},
			{
				ResourceName:            "google_vmwareengine_subnet.vmw-engine-subnet",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent", "name"},
			},
		},
	})
}

func testVmwareenginePrivateCloudCreateConfig(context map[string]interface{}) string {
	return testVmwareenginePrivateCloudConfig(context, "sample description", "TIME_LIMITED", 1, 0) +
		testVmwareengineVcenterNSXCredentialsConfig(context) +
		testVmwareengineUpgradesConfig(context) +
		testVmwareengineAnnouncementsConfig(context)
}

func testVmwareenginePrivateCloudUpdateNodeConfig(context map[string]interface{}, description string) string {
	return testVmwareenginePrivateCloudConfig(context, description, "STANDARD", 3, 8) + testVmwareengineVcenterNSXCredentialsConfig(context)
}

func testVmwareenginePrivateCloudUpdateAutoscaleConfig(context map[string]interface{}) string {
	return testVmwareenginePrivateCloudAutoscaleConfig(context, "sample updated description", "", 3, 8) + testVmwareengineVcenterNSXCredentialsConfig(context)
}

func testVmwareenginePrivateCloudDelayedDeleteConfig(context map[string]interface{}) string {
	return testVmwareenginePrivateCloudDeletedConfig(context)
}

func testVmwareenginePrivateCloudUndeleteConfig(context map[string]interface{}) string {
	return testVmwareenginePrivateCloudAutoscaleConfig(context, "sample updated description", "STANDARD", 3, 0) + testVmwareengineVcenterNSXCredentialsConfig(context)
}

func testVmwareengineSubnetImportConfig(context map[string]interface{}) string {
	return testVmwareenginePrivateCloudAutoscaleConfig(context, "sample updated description", "STANDARD", 3, 0) + testVmwareengineSubnetConfig(context, "192.168.1.0/26")
}

func testVmwareengineSubnetUpdateConfig(context map[string]interface{}) string {
	return testVmwareenginePrivateCloudAutoscaleConfig(context, "sample updated description", "STANDARD", 3, 0) + testVmwareengineSubnetConfig(context, "192.168.2.0/26")
}

func testVmwareenginePrivateCloudCmekSetupConfig(context map[string]interface{}) string {
	if useCmek, ok := context["use_cmek"]; ok && useCmek.(bool) {
		return acctest.Nprintf(`
data "google_project" "project" {
  project_id = "%{vmwareengine_project}"
}

resource "google_kms_crypto_key_iam_member" "vmwareengine-key" {
  crypto_key_id = "%{kms_key_name}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:service-${data.google_project.project.number}@%{sa_suffix}"
}
`, context)
	}
	return ""
}

func testVmwareenginePrivateCloudConfig(context map[string]interface{}, description, pcType string, nodeCount, delayHours int) string {
	context["node_count"] = nodeCount
	context["delay_hrs"] = delayHours
	context["description"] = description
	context["type"] = pcType

	encryptionConfigBlock := ""
	dependsOnLine := ""
	if useCmekVal, ok := context["use_cmek"]; ok {
		if useCmek, ok := useCmekVal.(bool); ok {
			if useCmek {
				encryptionConfigBlock = acctest.Nprintf(`
  encryption_config {
    type         = "CMEK"
    kms_key_name = "%{kms_key_name}"
  }
`, context)
				dependsOnLine = "depends_on = [google_kms_crypto_key_iam_member.vmwareengine-key]"
			} else {
				encryptionConfigBlock = `
  encryption_config {
    type = "GMEK"
  }
`
			}
		}
	}
	context["encryption_config_block"] = encryptionConfigBlock
	context["depends_on_line"] = dependsOnLine

	return testVmwareenginePrivateCloudCmekSetupConfig(context) + acctest.Nprintf(`
resource "google_vmwareengine_network" "vmw-engine-nw" {
  project = "%{vmwareengine_project}"
  name              = "tf-test-pc-nw-%{random_suffix}"
  location          = "global"
  type              = "STANDARD"
  description = "PC network description."
}

resource "google_vmwareengine_private_cloud" "vmw-engine-pc" {
  project = "%{vmwareengine_project}"
  location = "%{region}-b"
  name = "tf-test-sample-pc%{random_suffix}"
  description = "%{description}"
  type = "%{type}"
  deletion_delay_hours = "%{delay_hrs}"
  send_deletion_delay_hours_if_zero = true
  network_config {
    management_cidr = "192.168.0.0/24"
    vmware_engine_network = google_vmwareengine_network.vmw-engine-nw.id
  }
  management_cluster {
    cluster_id = "tf-test-sample-mgmt-cluster-custom-core-count%{random_suffix}"
    node_type_configs {
      node_type_id = "standard-72"
      node_count = "%{node_count}"
      custom_core_count = 32
    }
  }

  %{encryption_config_block}
  %{depends_on_line}
}

data "google_vmwareengine_private_cloud" "ds" {
    project = "%{vmwareengine_project}"
	location = "%{region}-b"
	name = "tf-test-sample-pc%{random_suffix}"
	depends_on = [
   	google_vmwareengine_private_cloud.vmw-engine-pc,
  ]
}
`, context)
}

func testVmwareenginePrivateCloudAutoscaleConfig(context map[string]interface{}, description, pcType string, nodeCount, delayHours int) string {
	context["node_count"] = nodeCount
	context["delay_hrs"] = delayHours
	context["description"] = description
	context["type"] = pcType

	encryptionConfigBlock := ""
	dependsOnLine := ""
	if useCmekVal, ok := context["use_cmek"]; ok {
		if useCmek, ok := useCmekVal.(bool); ok {
			if useCmek {
				encryptionConfigBlock = acctest.Nprintf(`
  encryption_config {
    type         = "CMEK"
    kms_key_name = "%{kms_key_name}"
  }
`, context)
				dependsOnLine = "depends_on = [google_kms_crypto_key_iam_member.vmwareengine-key]"
			} else {
				encryptionConfigBlock = `
  encryption_config {
    type = "GMEK"
  }
`
			}
		}
	}
	context["encryption_config_block"] = encryptionConfigBlock
	context["depends_on_line"] = dependsOnLine

	return testVmwareenginePrivateCloudCmekSetupConfig(context) + acctest.Nprintf(`
resource "google_vmwareengine_network" "vmw-engine-nw" {
  project = "%{vmwareengine_project}"
  name              = "tf-test-pc-nw-%{random_suffix}"
  location          = "global"
  type              = "STANDARD"
  description = "PC network description."
}

resource "google_vmwareengine_private_cloud" "vmw-engine-pc" {
  project = "%{vmwareengine_project}"
  location = "%{region}-b"
  name = "tf-test-sample-pc%{random_suffix}"
  description = "%{description}"
  type = "%{type}"
  deletion_delay_hours = "%{delay_hrs}"
  send_deletion_delay_hours_if_zero = true
  network_config {
    management_cidr = "192.168.0.0/24"
    vmware_engine_network = google_vmwareengine_network.vmw-engine-nw.id
  }
  management_cluster {
    cluster_id = "tf-test-sample-mgmt-cluster-custom-core-count%{random_suffix}"
    node_type_configs {
      node_type_id = "standard-72"
      node_count = "%{node_count}"
      custom_core_count = 32
    }
    autoscaling_settings {
      autoscaling_policies {
        autoscale_policy_id = "autoscaling-policy"
        node_type_id = "standard-72"
        scale_out_size = 1
        cpu_thresholds {
          scale_out = 80
          scale_in  = 15
        }
        consumed_memory_thresholds {
          scale_out = 75
          scale_in  = 20
        }
        storage_thresholds {
          scale_out = 80
          scale_in  = 20
        }
      }
      min_cluster_node_count = 3
      max_cluster_node_count = 8
      cool_down_period = "1800s"
    }
  }

  %{encryption_config_block}
  %{depends_on_line}
}

data "google_vmwareengine_private_cloud" "ds" {
    project = "%{vmwareengine_project}"
	location = "%{region}-b"
	name = "tf-test-sample-pc%{random_suffix}"
	depends_on = [
   	google_vmwareengine_private_cloud.vmw-engine-pc,
  ]
}
`, context)
}

func testVmwareenginePrivateCloudDeletedConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vmwareengine_network" "vmw-engine-nw" {
  project = "%{vmwareengine_project}"
  name              = "tf-test-pc-nw-%{random_suffix}"
  location          = "global"
  type              = "STANDARD"
  description = "PC network description."
}
`, context)
}

func testVmwareengineVcenterNSXCredentialsConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_vmwareengine_nsx_credentials" "nsx-ds" {
	parent =  google_vmwareengine_private_cloud.vmw-engine-pc.id
}

data "google_vmwareengine_vcenter_credentials" "vcenter-ds" {
	parent =  google_vmwareengine_private_cloud.vmw-engine-pc.id
}
`, context)
}

func testVmwareengineUpgradesConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_vmwareengine_upgrades" "upgrades-ds" {
	parent =  google_vmwareengine_private_cloud.vmw-engine-pc.id
}
`, context)
}

func testVmwareengineAnnouncementsConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_vmwareengine_announcements" "announcements-ds" {
	parent =  "projects/%{vmwareengine_project}/locations/%{region}-b"
}
`, context)
}

func testVmwareengineSubnetConfig(context map[string]interface{}, ipCidrRange string) string {
	context["ip_cidr_range"] = ipCidrRange
	return acctest.Nprintf(`
resource "google_vmwareengine_subnet" "vmw-engine-subnet" {
  name = "service-2"
  parent =  google_vmwareengine_private_cloud.vmw-engine-pc.id
  ip_cidr_range = "%{ip_cidr_range}"
}

data "google_vmwareengine_subnet" "subnet-ds" {
  name = "service-2"
  parent = google_vmwareengine_private_cloud.vmw-engine-pc.id
  depends_on = [
    google_vmwareengine_subnet.vmw-engine-subnet,
  ]
}
`, context)
}

func testAccCheckGoogleVmwareengineNsxCredentialsMeta(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find nsx credentials data source: %s", n)
		}
		_, ok = rs.Primary.Attributes["username"]
		if !ok {
			return fmt.Errorf("can't find 'username' attribute in data source: %s", n)
		}
		_, ok = rs.Primary.Attributes["password"]
		if !ok {
			return fmt.Errorf("can't find 'password' attribute in data source: %s", n)
		}
		return nil
	}
}

func testAccCheckGoogleVmwareengineVcenterCredentialsMeta(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find vcenter credentials data source: %s", n)
		}
		_, ok = rs.Primary.Attributes["username"]
		if !ok {
			return fmt.Errorf("can't find 'username' attribute in data source: %s", n)
		}
		_, ok = rs.Primary.Attributes["password"]
		if !ok {
			return fmt.Errorf("can't find 'password' attribute in data source: %s", n)
		}
		return nil
	}
}

func testAccCheckGoogleVmwareengineUpgradesMeta(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find upgrades data source: %s", n)
		}
		if _, ok := rs.Primary.Attributes["parent"]; !ok {
			return fmt.Errorf("can't find 'parent' attribute in data source: %s", n)
		}
		if _, ok := rs.Primary.Attributes["upgrades.#"]; !ok {
			return fmt.Errorf("can't find 'upgrades' attribute in data source: %s", n)
		}
		return nil
	}
}

func testAccCheckGoogleVmwareengineAnnouncementsMeta(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find announcements data source: %s", n)
		}
		if _, ok := rs.Primary.Attributes["parent"]; !ok {
			return fmt.Errorf("can't find 'parent' attribute in data source: %s", n)
		}
		if _, ok := rs.Primary.Attributes["announcements.#"]; !ok {
			return fmt.Errorf("can't find 'announcements' attribute in data source: %s", n)
		}
		return nil
	}
}

func testAccCheckVmwareenginePrivateCloudDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_vmwareengine_private_cloud" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}
			config := acctest.GoogleProviderConfig(t)
			url, err := tpgresource.ReplaceVarsForTest(config, rs, transport_tpg.BaseUrl(vmwareengine.Product, config)+"projects/{{project}}/locations/{{location}}/privateClouds/{{name}}")
			if err != nil {
				return err
			}
			billingProject := ""
			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}
			res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				pcState, ok := res["state"]
				if !ok {
					return fmt.Errorf("Unable to fetch state for existing VmwareenginePrivateCloud %s", url)
				}
				if pcState.(string) != "DELETED" {
					return fmt.Errorf("VmwareenginePrivateCloud still exists at %s", url)
				}
			}
		}
		return nil
	}
}
