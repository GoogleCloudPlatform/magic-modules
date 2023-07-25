// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccAlloydbUser_alloydbUserExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbUser_alloydbUserExample(context),
			},
			{
				ResourceName:            "google_alloydb_user.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "user_id", "cluster"},
			},
			{
				Config: testAccAlloydbUser_alloydbUserExample(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlloydbUserDestroyProducer(t),
				),
			},
		},
	})
}

func testAccAlloydbUser_alloydbUserExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_user" "default" {
  user_id     = "me%{random_suffix}"
  password    = "changeme%{random_suffix}"

  database_roles = [
    "postgres"
  ]

  cluster = google_alloydb_cluster.default.id

  depends_on = [google_alloydb_cluster.default]
}

resource "google_compute_network" "default" {
  name = "alloydb-cluster"
}

resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network    = google_compute_network.default.id
}

resource "google_alloydb_instance" "default" {
  cluster       = google_alloydb_cluster.default.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}"
  instance_type = "PRIMARY"

  depends_on = [google_service_networking_connection.vpc_connection]
}

resource "google_compute_global_address" "private_ip_alloc" {
  name          =  "tf-test-alloydb-cluster%{random_suffix}"
  address_type  = "INTERNAL"
  purpose       = "VPC_PEERING"
  prefix_length = 16
  network       = google_compute_network.default.id
}

resource "google_service_networking_connection" "vpc_connection" {
  network                 = google_compute_network.default.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_ip_alloc.name]
}
`, context)
}

func TestAccAlloydbUser_alloydbUserIamUserExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbUser_alloydbUserIamUserExample(context),
			},
			{
				ResourceName:            "google_alloydb_user.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "user_id", "cluster"},
			},
			{
				Config: testAccAlloydbUser_alloydbUserIamUserExample(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlloydbUserDestroyProducer(t),
				),
			},
		},
	})
}

func testAccAlloydbUser_alloydbUserIamUserExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_user" "default" {
  user_id      = "me@example.com%{random_suffix}"
  user_type    = "ALLOYDB_IAM_USER"

  database_roles = [
    "alloydbiamuser"
  ]

  cluster = google_alloydb_cluster.default.id

  depends_on = [google_alloydb_instance.default]
}

resource "google_compute_network" "default" {
  name = "alloydb-cluster"
}

resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster-iam%{random_suffix}"
  location   = "us-central1"
  network    = google_compute_network.default.id
}

resource "google_alloydb_instance" "default" {
  cluster       = google_alloydb_cluster.default.name
  instance_id   = "tf-test-alloydb-instance-iam%{random_suffix}"
  instance_type = "PRIMARY"

  depends_on = [google_service_networking_connection.vpc_connection]
}

resource "google_compute_global_address" "private_ip_alloc" {
  name          =  "tf-test-alloydb-cluster-iam%{random_suffix}"
  address_type  = "INTERNAL"
  purpose       = "VPC_PEERING"
  prefix_length = 16
  network       = google_compute_network.default.id
}

resource "google_service_networking_connection" "vpc_connection" {
  network                 = google_compute_network.default.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_ip_alloc.name]
}
`, context)
}

func testAccCheckAlloydbUserDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_alloydb_cluster" {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{AlloydbBasePath}}{{cluster}}/users/{{user_id}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("AlloydbUser still exists at %s", url)
			}
		}

		return nil
	}
}
