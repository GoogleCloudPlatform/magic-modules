package datastream_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDatastreamPrivateConnectionDatasourceConfig(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDatastreamPrivateConnectionDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_datastream_private_connection.default", "google_datastream_private_connection.default"),
				),
			},
		},
	})
}

const testAccDatastreamPrivateConnectionDatasourceConfig = `
resource "google_datastream_private_connection" "default" {
    display_name          = "Connection profile"
    location              = "us-central1"
    private_connection_id = "my-connection"

    labels = {
        key = "value"
    }

    vpc_peering_config {
        vpc_name = "my-vpc"
        subnet = "10.0.0.0/29"
    }
}

data "google_datastream_private_connection" "default" {
	private_connection_id 	= "my-connection"
}
`
