package iamworkforcepool

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIAMWorkforcePoolProviderSCIMTenant_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
		"org_id":       getTestOrgFromEnv(t),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIAMWorkforcePoolProviderSCIMTenantDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIAMWorkforcePoolProviderSCIMTenant_basic(context),
			},
			{
				ResourceName:            "google_iam_workforce_pool_provider_scim_tenant.tenant",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"workforce_pool_provider_id", "workforce_pool_id"},
			},
		},
	})
}

func testAccIAMWorkforcePoolProviderSCIMTenant_basic(context map[string]interface{}) string {
	return Nprintf(`
resource "google_iam_workforce_pool" "pool" {
  workforce_pool_id = "example-pool-%{random_suffix}"
  parent            = "organizations/%{org_id}"
  location         = "global"
}

resource "google_iam_workforce_pool_provider" "provider" {
  workforce_pool_id   = google_iam_workforce_pool.pool.workforce_pool_id
  location           = google_iam_workforce_pool.pool.location
  provider_id        = "example-provider-%{random_suffix}"
  attribute_mapping  = {
	"google.subject" = "assertion.sub"
  }
  saml {
	  idp_metadata_xml = file("test-fixtures/saml-metadata.xml")
  }
}

resource "google_iam_workforce_pool_provider_scim_tenant" "tenant" {
  workforce_pool_id        = google_iam_workforce_pool_provider.provider.workforce_pool_id
  provider_id              = google_iam_workforce_pool_provider.provider.provider_id
  location                 = google_iam_workforce_pool_provider.provider.location
  display_name             = "Example SCIM Tenant"
}
`, context)
}

func TestAccIAMWorkforcePoolProviderSCIMTenant_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
		"org_id":       getTestOrgFromEnv(t),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIAMWorkforcePoolProviderSCIMTenantDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIAMWorkforcePoolProviderSCIMTenant_basic(context),
			},
			{
				ResourceName:            "google_iam_workforce_pool_provider_scim_tenant.tenant",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"workforce_pool_provider_id", "workforce_pool_id"},
				Config: testAccIAMWorkforcePoolProviderSCIMTenant_update(context),
			},
		},
	})
}

func testAccIAMWorkforcePoolProviderSCIMTenant_update(context map[string]interface{}) string {
	return Nprintf(`
resource "google_iam_workforce_pool" "pool" {
  workforce_pool_id = "example-pool-%{random_suffix}"
  parent            = "organizations/%{org_id}"
  location         = "global"
}

resource "google_iam_workforce_pool_provider" "provider" {
  workforce_pool_id   = google_iam_workforce_pool.pool.workforce_pool_id
  location           = google_iam_workforce_pool.pool.location
  provider_id        = "example-provider-%{random_suffix}"
  attribute_mapping  = {
	"google.subject" = "assertion.sub"
  }
  saml {
	  idp_metadata_xml = file("test-fixtures/saml-metadata.xml")
  }
}

resource "google_iam_workforce_pool_provider_scim_tenant" "tenant" {
  workforce_pool_id        = google_iam_workforce_pool_provider.provider.workforce_pool_id
  provider_id              = google_iam_workforce_pool_provider.provider.provider_id
  location                 = google_iam_workforce_pool_provider.provider.location
  display_name             = "Updated Example SCIM Tenant"
}
`, context)
}

func TestAccIAMWorkforcePoolProviderSCIMTenant_disappears(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
		"org_id":       getTestOrgFromEnv(t),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIAMWorkforcePoolProviderSCIMTenantDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIAMWorkforcePoolProviderSCIMTenant_basic(context),
			},
			{
				ResourceName:      "google_iam_workforce_pool_provider_scim_tenant.tenant",
				ImportState:       true,
				ImportStateVerify: true,
				PreCheck: func() {
					// Delete the SCIM tenant outside of Terraform
					config := googleProviderConfig(t, testAccProvider.Meta())
					wpID := context["random_suffix"].(string)
					err := deleteIAMWorkforcePoolProviderSCIMTenant(config, "organizations/"+context["org_id"].(string), "example-pool-"+wpID, "example-provider-"+wpID)
					if err != nil {
						t.Fatalf("Error deleting SCIM tenant: %v", err)
					}
				},
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func deleteIAMWorkforcePoolProviderSCIMTenant(config *Config, parent, poolID, providerID string) error {
	// Construct the SCIM tenant ID
	tenantID := fmt.Sprintf("projects/%s/locations/%s/workforcePools/%s/providers/%s/scimConfig", config.Project, "global", poolID, providerID)

	// Attempt to delete the SCIM tenant
	_, err := config.NewIAMClient(config.userAgent).Projects.Locations.WorkforcePools.Providers.Delete(tenantID).Do()
	if err != nil {
		return fmt.Errorf("error deleting SCIM tenant %s: %w", tenantID, err)
	}

	return nil
}