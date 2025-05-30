package gkehub2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

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
				ImportStateVerifyIgnore: []string{"scope_rbac_role_binding_id", "scope_id", "labels", "terraform_labels"},
			},
			{
				Config: testAccGKEHub2ScopeRBACRoleBinding_gkehubScopeRbacRoleBindingBasicExample_update(context),
			},
			{
				ResourceName:            "google_gke_hub_scope_rbac_role_binding.scoperbacrolebinding",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"scope_rbac_role_binding_id", "scope_id", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccGKEHub2ScopeRBACRoleBinding_gkehubScopeRbacRoleBindingBasicExample_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_gke_hub_scope" "scoperbacrolebinding" {
  scope_id = "tf-test-scope%{random_suffix}"
}

resource "google_gke_hub_scope_rbac_role_binding" "scoperbacrolebinding" {
  scope_rbac_role_binding_id = "tf-test-scope-rbac-role-binding%{random_suffix}"
  scope_id = "tf-test-scope%{random_suffix}"
  user = "test-email@gmail.com"
  role {
    predefined_role = "ADMIN"
  }
  labels = {
      key = "value" 
  }
  depends_on = [google_gke_hub_scope.scoperbacrolebinding]
}
`, context)
}

func testAccGKEHub2ScopeRBACRoleBinding_gkehubScopeRbacRoleBindingBasicExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_gke_hub_scope" "scoperbacrolebinding" {
  scope_id = "tf-test-scope%{random_suffix}"
}

resource "google_gke_hub_scope_rbac_role_binding" "scoperbacrolebinding" {
  scope_rbac_role_binding_id = "tf-test-scope-rbac-role-binding%{random_suffix}"
  scope_id = "tf-test-scope%{random_suffix}"
  group = "test-email2@gmail.com"
  role {
    predefined_role = "VIEW"
  }
  labels = {
      key = "updated_value" 
  }
  depends_on = [google_gke_hub_scope.scoperbacrolebinding]
}
`, context)
}

func TestAccGKEHub2ScopeRBACRoleBinding_gkehubScopeRbacCustomRoleBindingBasicExample_update(t *testing.T) {
        t.Parallel()

        context := map[string]interface{}{
                "project":       envvar.GetTestProjectFromEnv(),
                "random_suffix": acctest.RandString(t, 10),
        }

        acctest.VcrTest(t, resource.TestCase{
                PreCheck:                 func() { acctest.AccTestPreCheck(t) },
                ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
                CheckDestroy:             testAccCheckGKEHub2ScopeRBACRoleBindingDestroyProducer(t),
                Steps: []resource.TestStep{
                        {
                                Config: testAccGKEHub2ScopeRBACRoleBinding_gkehubScopeRbacCustomRoleBindingBasicExample_basic(context),
                        },
                        {
                                ResourceName:            "google_gke_hub_scope_rbac_role_binding.scope_rbac_custom_role_binding",
                                ImportState:             true,
                                ImportStateVerify:       true,
                                ImportStateVerifyIgnore: []string{"labels", "scope_id", "scope_rbac_role_binding_id", "terraform_labels"},
                        },
			{
				Config: testAccGKEHub2ScopeRBACRoleBinding_gkehubScopeRbacCustomRoleBindingBasicExample_update(context),
			},
			{
				ResourceName:            "google_gke_hub_scope_rbac_role_binding.scope_rbac_custom_role_binding",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"scope_rbac_role_binding_id", "scope_id", "labels", "terraform_labels"},
			},
                },
        })
}

func testAccGKEHub2ScopeRBACRoleBinding_gkehubScopeRbacCustomRoleBindingBasicExample_basic(context map[string]interface{}) string {
        return acctest.Nprintf(`
resource "google_gke_hub_scope" "scope" {
  scope_id = "tf-test-scope%{random_suffix}"
}

resource "google_gke_hub_feature" "rbacrolebindingactuation" {
  name = "rbacrolebindingactuation"
  location = "global"
  spec {
    rbacrolebindingactuation {
      allowed_custom_roles = ["my-custom-role", "my-custom-role2"]
    }
  }
}

resource "google_gke_hub_scope_rbac_role_binding" "scope_rbac_custom_role_binding" {
  scope_rbac_role_binding_id = "tf-test-scope-rbac-role-binding%{random_suffix}"
  scope_id = google_gke_hub_scope.scope.scope_id
  user = "test-email@gmail.com"
  role {
    custom_role = "my-custom-role"
  }
  labels = {
      key = "value" 
  }
  depends_on = [google_gke_hub_scope.scope, google_gke_hub_feature.rbacrolebindingactuation]
}
`, context)
}

func testAccGKEHub2ScopeRBACRoleBinding_gkehubScopeRbacCustomRoleBindingBasicExample_update(context map[string]interface{}) string {
        return acctest.Nprintf(`
resource "google_gke_hub_scope" "scope" {
  scope_id = "tf-test-scope%{random_suffix}"
}

resource "google_gke_hub_feature" "rbacrolebindingactuation" {
  name = "rbacrolebindingactuation"
  location = "global"
  spec {
    rbacrolebindingactuation {
      allowed_custom_roles = ["my-custom-role", "my-custom-role2"]
    }
  }
}

resource "google_gke_hub_scope_rbac_role_binding" "scope_rbac_custom_role_binding" {
  scope_rbac_role_binding_id = "tf-test-scope-rbac-role-binding%{random_suffix}"
  scope_id = google_gke_hub_scope.scope.scope_id
  user = "test-email@gmail.com"
  role {
    custom_role = "my-custom-role-2"
  }
  labels = {
      key = "value" 
  }
  depends_on = [google_gke_hub_scope.scope, google_gke_hub_feature.rbacrolebindingactuation]
}
`, context)
}
