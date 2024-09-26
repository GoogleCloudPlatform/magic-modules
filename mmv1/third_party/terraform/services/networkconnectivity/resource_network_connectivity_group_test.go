// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package networkconnectivity_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccNetworkConnectivityGroup_BasicGroup(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkConnectivityGroupDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkConnectivityGroup_BasicGroup(context),
			},
			{
				ResourceName:            "google_network_connectivity_group.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccNetworkConnectivityGroup_BasicGroupUpdate0(context),
			},
			{
				ResourceName:            "google_network_connectivity_group.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func testAccNetworkConnectivityGroup_BasicGroup(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_network_connectivity_hub" "basic_hub" {
  name        = "tf-test-hub%{random_suffix}"
  description = "A sample hub"
  project     = "%{project_name}"

  labels = {
    label-one = "value-one"
  }
}

resource "google_network_connectivity_group" "primary" {
  hub = google_network_connectivity_hub.basic_hub.id
  name = "default"
  auto_accept {
    auto_accept_projects = ["tf-test-name%{random_suffix}"]
  }
}


`, context)
}

func testAccNetworkConnectivityGroup_BasicGroupUpdate0(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_network_connectivity_hub" "basic_hub" {
  name        = "tf-test-hub%{random_suffix}"
  description = "A sample hub"
  project     = "%{project_name}"

  labels = {
    label-one = "value-one"
  }
}

resource "google_network_connectivity_group" "primary" {
  hub = google_network_connectivity_hub.basic_hub.id
  name = "default"
  auto_accept {
    auto_accept_projects = ["tf-test-name%{random_suffix}", "tf-test-name%{random_suffix}"]
  }
}


`, context)
}
