package vmwareengine_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	_ "github.com/hashicorp/terraform-provider-google/google/services/resourcemanager"
	_ "github.com/hashicorp/terraform-provider-google/google/services/vmwareengine"
)

func TestAccDataSourceVmwareEngineNetwork_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":        acctest.RandString(t, 10),
		"org_id":               envvar.GetTestOrgFromEnv(t),
		"billing_account":      envvar.GetTestBillingAccountFromEnv(t),
		"vmwareengine_project": os.Getenv("GOOGLE_VMWAREENGINE_PROJECT"),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		CheckDestroy: testAccCheckVmwareengineNetworkDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVmwareEngineNetworkConfig(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_vmwareengine_network.ds", "google_vmwareengine_network.nw"),
				),
			},
		},
	})
}

func testAccDataSourceVmwareEngineNetworkConfig(context map[string]interface{}) string {
	projectSetup := `
resource "google_project" "project" {
  project_id      = "tf-test%{random_suffix}"
  name            = "tf-test%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
  deletion_policy = "DELETE"
}

resource "google_project_service" "vmwareengine" {
  project = google_project.project.project_id
  service = "vmwareengine.googleapis.com"
}

resource "time_sleep" "sleep" {
  create_duration = "1m"

  depends_on = [
    google_project_service.vmwareengine,
  ]
}
`
	projectVar := "google_project.project.project_id"
	dependsOnLine := `
  depends_on = [
    time_sleep.sleep # Sleep allows permissions in the new project to propagate
  ]
`

	if isProjectCreationDisabled() {
		projectSetup = ""
		projectVar = fmt.Sprintf(`"%s"`, context["vmwareengine_project"].(string))
		dependsOnLine = ""
	} else {
		projectSetup = acctest.Nprintf(projectSetup, context)
	}

	context["project_setup"] = projectSetup
	context["project_var"] = projectVar
	context["depends_on_line"] = dependsOnLine

	return acctest.Nprintf(`
%{project_setup}

resource "google_vmwareengine_network" "nw" {
  project           = %{project_var}
  name              = "tf-test-sample-network%{random_suffix}"
  location          = "global" # Standard network needs to be global
  type              = "STANDARD"
  description       = "VMwareEngine standard network sample"
  %{depends_on_line}
}

data "google_vmwareengine_network" "ds" {
  name     = google_vmwareengine_network.nw.name
  project  = %{project_var}
  location = "global"
}
`, context)
}
