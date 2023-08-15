package gkehub2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccGKEHub2ScopeRBACRoleBinding_gkehubScopeRbacRoleBindingBasicExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGKEHub2ScopeRBACRoleBinding_gkehubScopeRbacRoleBindingBasicExample_basic(context),
			},
			{
				ResourceName:            "google_gke_hub_scope_rbac_role_binding.scoperbacrolebinding",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"scope_rbac_role_binding_id", "scope_id"},
			},
			{
				Config: testAccGKEHub2ScopeRBACRoleBinding_gkehubScopeRbacRoleBindingBasicExample_update(context),
			},
			{
				ResourceName:            "google_gke_hub_scope_rbac_role_binding.scoperbacrolebinding",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"scope_rbac_role_binding_id", "scope_id"},
			},
		},
	})
}

func testAccGKEHub2ScopeRBACRoleBinding_gkehubScopeRbacRoleBindingBasicExample_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_gke_hub_scope" "scoperbacrolebinding" {
  scope_id = "tf-test-scope%{random_suffix}"
  all_memberships = false
}

resource "google_gke_hub_scope_rbac_role_binding" "scoperbacrolebinding" {
  scope_rbac_role_binding_id = "tf-test-scope-rbac-role-binding%{random_suffix}"
  scope_id = "tf-test-scope%{random_suffix}"
  user = "test-email@gmail.com"
  role {
    predefined_role = "ADMIN"
  }
  depends_on = [google_gke_hub_scope.scoperbacrolebinding]
}
`, context)
}

func testAccGKEHub2ScopeRBACRoleBinding_gkehubScopeRbacRoleBindingBasicExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_gke_hub_scope" "scoperbacrolebinding" {
  scope_id = "tf-test-scope%{random_suffix}"
  all_memberships = false
}

resource "google_gke_hub_scope_rbac_role_binding" "scoperbacrolebinding" {
  scope_rbac_role_binding_id = "tf-test-scope-rbac-role-binding%{random_suffix}"
  scope_id = "tf-test-scope%{random_suffix}"
  user = "test-email@gmail.com"
  role {
    predefined_role = "VIEW"
  }
  depends_on = [google_gke_hub_scope.scoperbacrolebinding]
}
`, context)
}
