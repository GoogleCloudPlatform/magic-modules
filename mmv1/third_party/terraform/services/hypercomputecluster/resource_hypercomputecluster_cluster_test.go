package hypercomputecluster_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccHypercomputeclusterCluster_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 8),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccHypercomputeclusterCluster_full(context),
			},
			{
				ResourceName:            "google_hypercomputecluster_cluster.cluster",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cluster_id", "labels", "location", "terraform_labels"},
			},
			{
				Config: testAccHypercomputeclusterCluster_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							"google_hypercomputecluster_cluster.cluster",
							plancheck.ResourceActionUpdate,
						),
					},
				},
			},
			{
				ResourceName:            "google_hypercomputecluster_cluster.cluster",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cluster_id", "labels", "location", "terraform_labels"},
			},
		},
	})
}

func testAccHypercomputeclusterCluster_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
}

locals {
  project_id = data.google_project.project.name
}

resource "google_hypercomputecluster_cluster" "cluster" {
  cluster_id                  = "tf%{random_suffix}"
  location                    = "us-central1"
  description                 = "Cluster Director instance created through Terraform"
  labels = {
    old-key = "old-value"
  }
  network_resources {
    id = "network-default"
    config {
      new_network {
        network = "projects/${local.project_id}/global/networks/net-%{random_suffix}"
      }
    }
  }
  compute_resources {
    id = "compute-spot"
    config {
      new_spot_instances {
        machine_type = "n2-standard-2"
        zone = "us-central1-a"
        termination_action = "STOP"
      }
    }
  }
  storage_resources {
    id = "storage-old"
    config {
      new_bucket {
        storage_class = "STANDARD"
        bucket = "bucket-old-%{random_suffix}"
      }
    }
  }
  orchestrator {
    slurm {
      login_nodes {
        machine_type = "n2-standard-2"        
        count = 1
        zone = "us-central1-a"
        boot_disk {
          size_gb = "100"
          type = "pd-balanced"
        }
        enable_os_login = "true"
        enable_public_ips = "true"
        labels = {
          old-key = "old-value"
        }
        startup_script = "#! /bin/bash"
        storage_configs {
          id = "storage-old"
          local_mount = "/home"
        }
      }
      node_sets {
        id = "nodeset"
        compute_id = "compute-spot"
        static_node_count = 1
        max_dynamic_node_count = 1
        compute_instance {
          boot_disk {
            size_gb = "100"
            type = "pd-balanced"
          }
          labels = {
            old-key = "old-value"
          }
          startup_script = "#! /bin/bash"
        }
        storage_configs {
          id = "storage-old"
          local_mount = "/home"
        }
      }
      partitions {
        id = "partition"
        node_set_ids = ["nodeset"]
      }
      default_partition = "partition"
      epilog_bash_scripts = ["#! /bin/bash"]
      prolog_bash_scripts = ["#! /bin/bash"]
    }
  }
}
`, context)
}

func testAccHypercomputeclusterCluster_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
}

locals {
  project_id = data.google_project.project.name
}

resource "google_hypercomputecluster_cluster" "cluster" {
  cluster_id                  = "tf%{random_suffix}"
  location                    = "us-central1"
  description                 = "Cluster Director instance created through Terraform (updated)"
  labels = {
    new-key = "new-value"
  }
  network_resources {
    id = "network-default"
    config {
      new_network {
        network = "projects/${local.project_id}/global/networks/net-%{random_suffix}"
      }
    }
  }
  compute_resources {
    id = "compute-spot-new"
    config {
      new_spot_instances {
        machine_type = "n2-standard-2"
        zone = "us-central1-a"
        termination_action = "DELETE"
      }
    }
  }
  storage_resources {
    id = "storage-new"
    config {
      new_bucket {
        storage_class = "STANDARD"
        bucket = "bucket-new-%{random_suffix}"
      }
    }
  }
  orchestrator {
    slurm {
      login_nodes {
        machine_type = "n2-standard-4"        
        count = 2
        zone = "us-central1-a"
        boot_disk {
          size_gb = "100"
          type = "pd-balanced"
        }
        enable_os_login = "false"
        enable_public_ips = "false"
        labels = {
          new-key = "new-value"
        }
        storage_configs {
          id = "storage-new"
          local_mount = "/home"
        }
      }
      node_sets {
        id = "nodesetnew"
        compute_id = "compute-spot-new"
        static_node_count = 2
        max_dynamic_node_count = 2
        compute_instance {
          boot_disk {
            size_gb = "100"
            type = "pd-balanced"
          }
          labels = {
            new-key = "new-value"
          }
        }
        storage_configs {
          id = "storage-new"
          local_mount = "/home"
        }
      }
      partitions {
        id = "partitionnew"
        node_set_ids = ["nodesetnew"]
      }
      default_partition = "partitionnew"
    }
  }
}
`, context)
}
