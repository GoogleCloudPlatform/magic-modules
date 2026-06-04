// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0
package config_test

import (
"fmt"
"strings"
"testing"

"github.com/hashicorp/terraform-plugin-testing/helper/resource"
"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	_ "github.com/hashicorp/terraform-provider-google/google/services/config"
	_ "github.com/hashicorp/terraform-provider-google/google/services/resourcemanager"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"google.golang.org/api/googleapi"
)

func TestAccConfigDeploymentGroup_basic(t *testing.T) {
t.Parallel()

context := map[string]interface{}{
"project":       envvar.GetTestProjectFromEnv(),
"random_suffix": acctest.RandString(t, 10),
}

acctest.VcrTest(t, resource.TestCase{
PreCheck:                 func() { acctest.AccTestPreCheck(t) },
ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
CheckDestroy:             testAccCheckConfigDeploymentGroupDestroyProducer(t),
Steps: []resource.TestStep{
{
Config: testAccConfigDeploymentGroup_basic(context),
},
{
ResourceName:            "google_config_deployment_group.basic",
ImportState:             true,
ImportStateVerify:       true,
ImportStateVerifyIgnore: []string{"location", "labels", "annotations"},
},
{
Config: testAccConfigDeploymentGroup_update(context),
},
{
ResourceName:            "google_config_deployment_group.basic",
ImportState:             true,
ImportStateVerify:       true,
ImportStateVerifyIgnore: []string{"location", "labels", "annotations"},
},
},
})
}

func testAccConfigDeploymentGroup_basic(context map[string]interface{}) string {
return acctest.Nprintf(`
resource "google_config_deployment_group" "basic" {
  name     = "tf-test-dg-%{random_suffix}"
  location = "us-central1"

  labels = {
    env = "test"
  }
}
`, context)
}

func testAccConfigDeploymentGroup_update(context map[string]interface{}) string {
return acctest.Nprintf(`
resource "google_config_deployment_group" "basic" {
  name     = "tf-test-dg-%{random_suffix}"
  location = "us-central1"

  labels = {
    env = "production"
  }

  annotations = {
    purpose = "testing"
  }
}
`, context)
}

func testAccCheckConfigDeploymentGroupDestroyProducer(t *testing.T) resource.TestCheckFunc {
return func(s *terraform.State) error {
config := acctest.GoogleProviderConfig(t)

for _, rs := range s.RootModule().Resources {
if rs.Type != "google_config_deployment_group" {
continue
}

if rs.Primary.ID == "" {
return fmt.Errorf("Unable to verify delete of deployment group ID is empty")
}

project, err := acctest.GetTestProject(rs.Primary, config)
if err != nil {
return err
}

parts := strings.Split(rs.Primary.ID, "/")
deployment_group_id := parts[len(parts)-1]
location := rs.Primary.Attributes["location"]

url := fmt.Sprintf("https://config.googleapis.com/v1/projects/%s/locations/%s/deploymentGroups/%s", project, location, deployment_group_id)
_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
Config:    config,
Method:    "GET",
Project:   project,
RawURL:    url,
UserAgent: config.UserAgent,
})
if err != nil {
if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
return nil
} else if ok {
return fmt.Errorf("Error making GCP platform call: http code error : %d, http message error: %s", gerr.Code, gerr.Message)
}
return fmt.Errorf("Error making GCP platform call: %s", err.Error())
}
return fmt.Errorf("Deployment group still exists")
}

	return nil
}
}

func TestAccConfigDeploymentGroup_deploymentUnits(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckConfigDeploymentGroupDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccConfigDeploymentGroup_deploymentUnits(context),
			},
			{
				ResourceName:            "google_config_deployment_group.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "labels", "annotations", "force_destroy", "deployment_reference_policy"},
			},
			{
				Config: testAccConfigDeploymentGroup_deploymentUnitsUpdate(context),
			},
			{
				ResourceName:            "google_config_deployment_group.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "labels", "annotations", "force_destroy", "deployment_reference_policy"},
			},
		},
	})
}

func testAccConfigDeploymentGroup_deploymentUnits(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_service_account" "sa" {
  account_id   = "tf-test-dg-sa-%{random_suffix}"
  display_name = "Infra Manager Test SA for DG"
}

resource "google_project_iam_member" "binding" {
  project = "%{project}"
  role    = "roles/config.agent"
  member  = "serviceAccount:${google_service_account.sa.email}"
}

resource "google_project_iam_member" "network_admin" {
  project = "%{project}"
  role    = "roles/compute.networkAdmin"
  member  = "serviceAccount:${google_service_account.sa.email}"
}

resource "google_config_deployment" "dep1" {
  name            = "tf-test-dg-dep1-%{random_suffix}"
  location        = "us-central1"
  service_account = "projects/%{project}/serviceAccounts/${google_service_account.sa.email}"
  force_destroy   = true

  terraform_blueprint {
    git_source {
      repo      = "https://github.com/terraform-google-modules/terraform-google-network"
      directory = "modules/vpc"
      ref       = "main"
    }
    
    input_values {
      variable_name = "project_id"
      input_value   = jsonencode("%{project}")
    }
    input_values {
      variable_name = "network_name"
      input_value   = jsonencode("tf-test-dg-net1-%{random_suffix}")
    }
  }

  depends_on = [
    google_project_iam_member.binding,
    google_project_iam_member.network_admin
  ]
}

resource "google_config_deployment" "dep2" {
  name            = "tf-test-dg-dep2-%{random_suffix}"
  location        = "us-central1"
  service_account = "projects/%{project}/serviceAccounts/${google_service_account.sa.email}"
  force_destroy   = true

  terraform_blueprint {
    git_source {
      repo      = "https://github.com/terraform-google-modules/terraform-google-network"
      directory = "modules/vpc"
      ref       = "main"
    }
    
    input_values {
      variable_name = "project_id"
      input_value   = jsonencode("%{project}")
    }
    input_values {
      variable_name = "network_name"
      input_value   = jsonencode("tf-test-dg-net2-%{random_suffix}")
    }
  }

  depends_on = [
    google_project_iam_member.binding,
    google_project_iam_member.network_admin
  ]
}

resource "google_config_deployment_group" "basic" {
  name     = "tf-test-dg-%{random_suffix}"
  location = "us-central1"

  deployment_units {
    id           = "unit-1"
    deployment   = google_config_deployment.dep1.id
    dependencies = []
  }

  deployment_units {
    id           = "unit-2"
    deployment   = google_config_deployment.dep2.id
    dependencies = ["unit-1"]
  }

  labels = {
    env = "test"
  }
}
`, context)
}

func testAccConfigDeploymentGroup_deploymentUnitsUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_service_account" "sa" {
  account_id   = "tf-test-dg-sa-%{random_suffix}"
  display_name = "Infra Manager Test SA for DG"
}

resource "google_project_iam_member" "binding" {
  project = "%{project}"
  role    = "roles/config.agent"
  member  = "serviceAccount:${google_service_account.sa.email}"
}

resource "google_project_iam_member" "network_admin" {
  project = "%{project}"
  role    = "roles/compute.networkAdmin"
  member  = "serviceAccount:${google_service_account.sa.email}"
}

resource "google_config_deployment" "dep1" {
  name            = "tf-test-dg-dep1-%{random_suffix}"
  location        = "us-central1"
  service_account = "projects/%{project}/serviceAccounts/${google_service_account.sa.email}"
  force_destroy   = true

  terraform_blueprint {
    git_source {
      repo      = "https://github.com/terraform-google-modules/terraform-google-network"
      directory = "modules/vpc"
      ref       = "main"
    }
    
    input_values {
      variable_name = "project_id"
      input_value   = jsonencode("%{project}")
    }
    input_values {
      variable_name = "network_name"
      input_value   = jsonencode("tf-test-dg-net1-%{random_suffix}")
    }
  }

  depends_on = [
    google_project_iam_member.binding,
    google_project_iam_member.network_admin
  ]
}

resource "google_config_deployment" "dep2" {
  name            = "tf-test-dg-dep2-%{random_suffix}"
  location        = "us-central1"
  service_account = "projects/%{project}/serviceAccounts/${google_service_account.sa.email}"
  force_destroy   = true

  terraform_blueprint {
    git_source {
      repo      = "https://github.com/terraform-google-modules/terraform-google-network"
      directory = "modules/vpc"
      ref       = "main"
    }
    
    input_values {
      variable_name = "project_id"
      input_value   = jsonencode("%{project}")
    }
    input_values {
      variable_name = "network_name"
      input_value   = jsonencode("tf-test-dg-net2-%{random_suffix}")
    }
  }

  depends_on = [
    google_project_iam_member.binding,
    google_project_iam_member.network_admin
  ]
}

resource "google_config_deployment_group" "basic" {
  name     = "tf-test-dg-%{random_suffix}"
  location = "us-central1"

  deployment_units {
    id           = "unit-1"
    deployment   = google_config_deployment.dep1.id
    dependencies = []
  }

  labels = {
    env     = "test"
    updated = "true"
  }

  force_destroy = true
}
`, context)
}

