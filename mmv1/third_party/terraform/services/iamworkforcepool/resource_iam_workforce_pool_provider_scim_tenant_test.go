package iamworkforcepool_test

import (
	"fmt"
	"regexp"
	"testing"
	"strings"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccIAMWorkforcePoolProviderScimTenant_basic(t *testing.T) {
	t.Parallel()
	random_suffix := acctest.RandString(t, 10)
	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": random_suffix,
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIAMWorkforcePoolProviderScimTenant_basic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_iam_workforce_pool_provider_scim_tenant.default", "name"),
					resource.TestCheckResourceAttrSet("google_iam_workforce_pool_provider_scim_tenant.default", "base_uri"),
				),
			},
			{
				ResourceName:      "google_iam_workforce_pool_provider_scim_tenant.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccIAMWorkforcePoolProviderScimTenant_basic(context map[string]interface{}) string {
	return fmt.Sprintf(`resource "c" "default" {
		location                = "us-central1"
		workforce_pool_id       = "test-pool-%s"
		workforce_pool_provider_id = "test-provider-%s"
		display_name            = "Test SCIM Tenant"
		description             = "Test description"
		scim_tenant_id          = "scim-tenant-%s"
	}
	`, context["random_suffix"], context["random_suffix"], context["random_suffix"])
}

func TestAccIAMWorkforcePoolProviderScimTenant_allFields(t *testing.T) {
	t.Parallel()
	random_suffix := acctest.RandString(t, 10)
	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": random_suffix,
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIAMWorkforcePoolProviderScimTenant_allFields(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_iam_workforce_pool_provider_scim_tenant.default", "name"),
					resource.TestCheckResourceAttrSet("google_iam_workforce_pool_provider_scim_tenant.default", "base_uri"),
					resource.TestCheckResourceAttr("google_iam_workforce_pool_provider_scim_tenant.default", "display_name", "Test SCIM Tenant"),
					resource.TestCheckResourceAttr("google_iam_workforce_pool_provider_scim_tenant.default", "description", "Test description"),
					resource.TestCheckResourceAttr("google_iam_workforce_pool_provider_scim_tenant.default", "state", "ACTIVE"),
				),
			},
			{
				ResourceName:      "google_iam_workforce_pool_provider_scim_tenant.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccIAMWorkforcePoolProviderScimTenant_allFields(context map[string]interface{}) string {
	return fmt.Sprintf(`resource "google_iam_workforce_pool_provider_scim_tenant" "default" {
		location                = "us-central1"
		workforce_pool_id       = "test-pool-%s"
		workforce_pool_provider_id = "test-provider-%s"
		display_name            = "Test SCIM Tenant"
		description             = "Test description"
		state                   = "ACTIVE"
		scim_tenant_id          = "scim-tenant-%s"
	}
	`, context["random_suffix"], context["random_suffix"], context["random_suffix"])
}

func TestAccIAMWorkforcePoolProviderScimTenant_stateEnum(t *testing.T) {
	t.Parallel()
	enums := []string{"STATE_UNSPECIFIED", "ACTIVE", "DELETED"}
	for _, state := range enums {
		random_suffix := acctest.RandString(t, 10)
		context := map[string]interface{}{
			"org_id":        envvar.GetTestOrgFromEnv(t),
			"random_suffix": random_suffix,
			"state":         state,
		}
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { acctest.AccTestPreCheck(t) },
			ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
			Steps: []resource.TestStep{
				{
					Config: testAccIAMWorkforcePoolProviderScimTenant_stateEnum(context),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("google_iam_workforce_pool_provider_scim_tenant.default", "state", state),
					),
				},
			},
		})
	}
}

func testAccIAMWorkforcePoolProviderScimTenant_stateEnum(context map[string]interface{}) string {
	return fmt.Sprintf(`resource "google_iam_workforce_pool_provider_scim_tenant" "default" {
		location                = "us-central1"
		workforce_pool_id       = "test-pool-%s"
		workforce_pool_provider_id = "test-provider-%s"
		display_name            = "Test SCIM Tenant"
		description             = "Test description"
		state                   = "%s"
		scim_tenant_id          = "scim-tenant-%s"
	}
	`, context["random_suffix"], context["random_suffix"], context["state"], context["random_suffix"])
}

func TestAccIAMWorkforcePoolProviderScimTenant_invalidScimTenantId(t *testing.T) {
	t.Parallel()
	invalidIds := []string{"abc", strings.Repeat("a", 33), "invalid!id", ""}
	for _, id := range invalidIds {
		random_suffix := acctest.RandString(t, 10)
		context := map[string]interface{}{
			"org_id":        envvar.GetTestOrgFromEnv(t),
			"random_suffix": random_suffix,
			"scim_tenant_id": id,
		}
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { acctest.AccTestPreCheck(t) },
			ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
			Steps: []resource.TestStep{
				{
					Config: testAccIAMWorkforcePoolProviderScimTenant_invalidScimTenantId(context),
					ExpectError: regexp.MustCompile("(?i)scim_tenant_id"),
				},
			},
		})
	}
}

func testAccIAMWorkforcePoolProviderScimTenant_invalidScimTenantId(context map[string]interface{}) string {
	return fmt.Sprintf(`resource "google_iam_workforce_pool_provider_scim_tenant" "default" {
		location                = "us-central1"
		workforce_pool_id       = "test-pool-%s"
		workforce_pool_provider_id = "test-provider-%s"
		display_name            = "Test SCIM Tenant"
		description             = "Test description"
		state                   = "ACTIVE"
		scim_tenant_id          = "%s"
	}
	`, context["random_suffix"], context["random_suffix"], context["scim_tenant_id"])
}

func TestAccIAMWorkforcePoolProviderScimTenant_updateFields(t *testing.T) {
	t.Parallel()
	random_suffix := acctest.RandString(t, 10)
	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": random_suffix,
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIAMWorkforcePoolProviderScimTenant_allFields(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_iam_workforce_pool_provider_scim_tenant.default", "display_name", "Test SCIM Tenant"),
					resource.TestCheckResourceAttr("google_iam_workforce_pool_provider_scim_tenant.default", "description", "Test description"),
				),
			},
			{
				Config: testAccIAMWorkforcePoolProviderScimTenant_updateFields(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_iam_workforce_pool_provider_scim_tenant.default", "display_name", "Updated SCIM Tenant"),
					resource.TestCheckResourceAttr("google_iam_workforce_pool_provider_scim_tenant.default", "description", "Updated description"),
					resource.TestCheckResourceAttr("google_iam_workforce_pool_provider_scim_tenant.default", "state", "DELETED"),
				),
			},
		},
	})
}

func testAccIAMWorkforcePoolProviderScimTenant_updateFields(context map[string]interface{}) string {
	return fmt.Sprintf(`resource "google_iam_workforce_pool_provider_scim_tenant" "default" {
		location                = "us-central1"
		workforce_pool_id       = "test-pool-%s"
		workforce_pool_provider_id = "test-provider-%s"
		display_name            = "Updated SCIM Tenant"
		description             = "Updated description"
		state                   = "DELETED"
		scim_tenant_id          = "scim-tenant-%s"
	}
	`, context["random_suffix"], context["random_suffix"], context["random_suffix"])
}

func TestAccIAMWorkforcePoolProviderScimTenant_immutableFields(t *testing.T) {
	t.Parallel()
	random_suffix := acctest.RandString(t, 10)
	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": random_suffix,
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIAMWorkforcePoolProviderScimTenant_allFields(context),
			},
			{
				Config: testAccIAMWorkforcePoolProviderScimTenant_immutableFields(context),
				ExpectError: regexp.MustCompile("(?i)cannot update immutable field"),
			},
		},
	})
}

func testAccIAMWorkforcePoolProviderScimTenant_immutableFields(context map[string]interface{}) string {
	return fmt.Sprintf(`resource "google_iam_workforce_pool_provider_scim_tenant" "default" {
		location                = "europe-west1" // attempt to change immutable field
		workforce_pool_id       = "test-pool-%s"
		workforce_pool_provider_id = "test-provider-%s"
		display_name            = "Test SCIM Tenant"
		description             = "Test description"
		state                   = "ACTIVE"
		scim_tenant_id          = "scim-tenant-%s"
	}
	`, context["random_suffix"], context["random_suffix"], context["random_suffix"])
}
