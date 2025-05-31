// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package datastream_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDatastreamPrivateConnectionDatasourceConfig(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"network_name":  acctest.BootstrapSharedTestNetwork(t, "datastream-network"),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDatastreamPrivateConnectionDatasourceConfig(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_datastream_private_connection.default", "google_datastream_private_connection.default"),
				),
			},
		},
	})
}

func testAccDatastreamPrivateConnectionDatasourceConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_datastream_private_connection" "default" {
    display_name          = "Connection profile"
    location              = "us-central1"
    private_connection_id = "tf-test-my-connection%{random_suffix}"

    labels = {
        key = "value"
    }

    vpc_peering_config {
        vpc 	= data.google_compute_network.default.id
        subnet	= "10.0.0.0/20"
    }
}

data "google_compute_network" "default" {
  name = "%{network_name}"
}

data "google_datastream_private_connection" "default" {
	location              	= "us-central1"
	private_connection_id 	= "tf-test-my-connection%{random_suffix}"
}
`, context)
}
