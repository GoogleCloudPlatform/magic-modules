package privilegedaccessmanager_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceGooglePrivilegedAccessManagerEntitlement_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPrivilegedAccessManagerEntitlemetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGooglePrivilegedAccessManagerEntitlement_basic(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_privileged_access_manager_entitlement.tfentitlement", "google_privileged_access_manager_entitlement.tfentitlement"),
				),
			},
		},
	})
}

func TestAccDataSourceGooglePrivilegedAccessManagerEntitlement_optionalProject(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPrivilegedAccessManagerEntitlemetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGooglePrivilegedAccessManagerEntitlement_optionalProject(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_privileged_access_manager_entitlement.tfentitlement", "google_privileged_access_manager_entitlement.tfentitlement"),
				),
			},
		},
	})
}

func testAccDataSourceGooglePrivilegedAccessManagerEntitlement_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_privileged_access_manager_entitlement" "tfentitlement" {
	entitlement_id = "tf-test-example-entitlement%{random_suffix}"
	location = "global"
	max_request_duration = "43200s"
	parent = "projects/%{project_name}"
	requester_justification_config { 
		unstructured{}
	}
	eligible_users {
		principals = [
			"group:test@google.com"
		]
	}
	privileged_access{
		gcp_iam_access{
			role_bindings{
				role = "roles/storage.admin"
				condition_expression = "request.time < timestamp(\"2024-04-23T18:30:00.000Z\")"
			}
			resource = "//cloudresourcemanager.googleapis.com/projects/%{project_name}"
			resource_type = "cloudresourcemanager.googleapis.com/Project"
		}
	}
}

data "google_privileged_access_manager_entitlement" "tfentitlement" {
  entitlement_id     = google_privileged_access_manager_entitlement.tfentitlement.entitlement_id
  parent  = google_privileged_access_manager_entitlement.tfentitlement.parent
}
`, context)
}

func testAccDataSourceGooglePrivilegedAccessManagerEntitlement_optionalProject(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_privileged_access_manager_entitlement" "tfentitlement" {
	entitlement_id = "tf-test-example-entitlement%{random_suffix}"
	location = "global"
	max_request_duration = "43200s"
	parent = "projects/%{project_name}"
	requester_justification_config { 
		unstructured{}
	}
	eligible_users {
		principals = [
			"group:test@google.com"
		]
	}
	privileged_access{
		gcp_iam_access{
			role_bindings{
				role = "roles/storage.admin"
			}
			resource = "//cloudresourcemanager.googleapis.com/projects/%{project_name}"
			resource_type = "cloudresourcemanager.googleapis.com/Project"
		}
	}
}

data "google_privileged_access_manager_entitlement" "tfentitlement" {
  entitlement_id     = google_privileged_access_manager_entitlement.tfentitlement.entitlement_id
  parent  = google_privileged_access_manager_entitlement.tfentitlement.parent
}
`, context)
}
