package edgecontainer_test

import (
        "testing"

        "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

        "github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccEdgecontainerCluster_update(t *testing.T) {
        t.Parallel()

        context := map[string]interface{}{
                "random_suffix": acctest.RandString(t, 10),
        }

        acctest.VcrTest(t, resource.TestCase{
                PreCheck:                 func() { acctest.AccTestPreCheck(t) },
                ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
                CheckDestroy:             testAccCheckEdgecontainerClusterDestroyProducer(t),
                Steps: []resource.TestStep{
                        {
                                Config: testAccEdgecontainerCluster_full(context),
                        },
                        {
                                ResourceName:            "google_edgecontainer_cluster.default",
                                ImportState:             true,
                                ImportStateVerify:       true,
                                ImportStateVerifyIgnore: []string{"labels", "location", "name", "terraform_labels"},
                        },
                        {
                                Config: testAccEdgecontainerCluster_addedExclusionWindows(context),
                        },
                        {
                                ResourceName:            "google_edgecontainer_cluster.default",
                                ImportState:             true,
                                ImportStateVerify:       true,
                                ImportStateVerifyIgnore: []string{"labels", "location", "name", "terraform_labels"},
                        },
                },
        })
}

func testAccEdgecontainerCluster_full(context map[string]interface{}) string {
        return acctest.Nprintf(`
resource "google_edgecontainer_cluster" "default" {
  name = "tf-test-cluster-with-maintenance-exclusions%{random_suffix}"
  location = "us-central1"

  authorization {
    admin_users {
      username = "admin@hashicorptest.com"
    }
  }

  networking {
    cluster_ipv4_cidr_blocks = ["10.0.0.0/16"]
    services_ipv4_cidr_blocks = ["10.1.0.0/16"]
  }

  fleet {
    project = "projects/${data.google_project.project.number}"
  }

  maintenance_policy {
    window {
      recurring_window {
        window {
          start_time = "2023-01-01T08:00:00Z"
          end_time = "2023-01-01T17:00:00Z"
        }
        recurrence = "FREQ=WEEKLY;BYDAY=SA"
      }
    }
  }
}

data "google_project" "project" {}
`, context)
}

func testAccEdgecontainerCluster_addedExclusionWindows(context map[string]interface{}) string {
        return acctest.Nprintf(`
resource "google_edgecontainer_cluster" "default" {
  name = "tf-test-cluster-with-maintenance-exclusions%{random_suffix}"
  location = "us-central1"

  authorization {
    admin_users {
      username = "admin@hashicorptest.com"
    }
  }

  networking {
    cluster_ipv4_cidr_blocks = ["10.0.0.0/16"]
    services_ipv4_cidr_blocks = ["10.1.0.0/16"]
  }

  fleet {
    project = "projects/${data.google_project.project.number}"
  }

  maintenance_policy {
    window {
      recurring_window {
        window {
          start_time = "2023-01-01T08:00:00Z"
          end_time = "2023-01-01T17:00:00Z"
        }
        recurrence = "FREQ=WEEKLY;BYDAY=SA"
      }
    }
    maintenance_exclusions {
      window {
        start_time = "2023-01-01T08:00:00Z"
        end_time = "2023-01-01T20:00:00Z"
      }
      id = "short-exclusion"
    }
    maintenance_exclusions {
      window {
        start_time = "2023-01-02T08:00:00Z"
        end_time = "2023-01-13T08:00:00Z"
      }
      id = "long-exclusion"
    }
  }
}

data "google_project" "project" {}
`, context)
}

