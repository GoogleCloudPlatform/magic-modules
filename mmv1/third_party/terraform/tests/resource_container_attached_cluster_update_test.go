package google

import (
  "testing"

  "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccContainerAttachedCluster_update(t *testing.T) {
  t.Parallel()

  context := map[string]interface{}{
    "random_suffix": randString(t, 10),
  }

  vcrTest(t, resource.TestCase{
    PreCheck:     func() { testAccPreCheck(t) },
    Providers:    testAccProviders,
    CheckDestroy: testAccCheckContainerAttachedClusterDestroyProducer(t),
    Steps: []resource.TestStep{
      {
        Config: testAccContainerAttachedCluster_containerAttachedCluster_full(context),
      },
      {
        ResourceName:            "google_container_attached_cluster.primary",
        ImportState:             true,
        ImportStateVerify:       true,
        ImportStateVerifyIgnore: []string{"location"},
      },
      {
        Config: testAccContainerAttachedCluster_containerAttachedCluster_update(context),
      },
      {
        ResourceName:            "google_container_attached_cluster.primary",
        ImportState:             true,
        ImportStateVerify:       true,
        ImportStateVerifyIgnore: []string{"location"},
      },
    },
  })
}

func testAccContainerAttachedCluster_containerAttachedCluster_full(context map[string]interface{}) string {
  return Nprintf(`
data "google_project" "project" {
}

resource "google_container_attached_cluster" "primary" {
  name     = "update%{random_suffix}"
  description = "Test cluster"
  distribution = "aks"
  annotations = {
    label-one = "value-one"
  }
  authorization {
    admin_users = [ "user1@example.com", "user2@example.com"]
  }
  oidc_config {
      issuer_url = "https://oidc.issuer.url"
      jwks = base64encode("{\"keys\":[{\"use\":\"sig\",\"kty\":\"RSA\",\"kid\":\"testid\",\"alg\":\"RS256\",\"n\":\"somedata\",\"e\":\"AQAB\"}]}")
  }
  platform_version = "1.24.0-gke.1"
  fleet {
      project = data.google_project.project.number
  }
  logging_config {
    component_config {
      enable_components = ["SYSTEM_COMPONENTS", "WORKLOADS"]
    }
  }
  monitoring_config {
    managed_prometheus_config {
      enabled = true
    }
  }
  project = data.google_project.project.project_id
  location = "us-west1"
}
`, context)
}

func testAccContainerAttachedCluster_containerAttachedCluster_update(context map[string]interface{}) string {
  return Nprintf(`
data "google_project" "project" {
}

resource "google_container_attached_cluster" "primary" {
  name     = "update%{random_suffix}"
  description = "Test cluster updated"
  distribution = "aks"
  annotations = {
    label-one = "value-one"
  label-two = "value-two"
  }
  authorization {
    admin_users = [ "user2@example.com", "user3@example.com"]
  }
  oidc_config {
      issuer_url = "https://oidc.issuer.url"
      jwks = base64encode("{\"keys\":[{\"use\":\"sig\",\"kty\":\"RSA\",\"kid\":\"testid\",\"alg\":\"RS256\",\"n\":\"somedata\",\"e\":\"AQAB\"}]}")
  }
  platform_version = "1.24.0-gke.1"
  fleet {
      project = data.google_project.project.number
  }
  monitoring_config {
    managed_prometheus_config {}
  }
  project = data.google_project.project.project_id
  location = "us-west1"
}
`, context)
}
