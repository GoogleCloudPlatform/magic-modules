package google

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccBigtableAppProfile_basic(t *testing.T) {
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigtableAppProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBigtableAppProfile_multiClusterRouting(instanceName),
			},
			{
				ResourceName:      "google_bigtable_app_profile.ap",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBigtableAppProfile_singleClusterRouting(t *testing.T) {
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigtableAppProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBigtableAppProfile_singleClusterRouting(instanceName),
			},
			{
				ResourceName:      "google_bigtable_app_profile.ap",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckBigtableAppProfileDestroy(s *terraform.State) error {
	var ctx = context.Background()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_bigtable_app_profile" {
			continue
		}

		config := testAccProvider.Meta().(*Config)
		c, err := config.bigtableClientFactory.NewInstanceAdminClient(config.Project)
		if err != nil {
			return fmt.Errorf("Error starting instance admin client. %s", err)
		}

		defer c.Close()

		_, err = c.GetAppProfile(ctx, rs.Primary.Attributes["instance"], rs.Primary.Attributes["name"])
		if err == nil {
			return fmt.Errorf("Instance %s and app profile %s still exists.", rs.Primary.Attributes["instance"], rs.Primary.Attributes["name"])
		}
	}

	return nil
}

func testAccBigtableAppProfile_multiClusterRouting(instanceName string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
	name = "%s"
	cluster {
		cluster_id   = "%s"
		zone         = "us-central1-b"
		num_nodes    = 3
		storage_type = "HDD"
	}
}

resource "google_bigtable_app_profile" "ap" {
	instance = google_bigtable_instance.instance.name
	app_profile_id = "test"

	multi_cluster_routing_use_any = true
}
`, instanceName, instanceName)
}

func testAccBigtableAppProfile_singleClusterRouting(instanceName string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
	name = "%s"
	cluster {
		cluster_id   = "%s"
		zone         = "us-central1-b"
		num_nodes    = 3
		storage_type = "HDD"
	}
}

resource "google_bigtable_app_profile" "ap" {
	instance = google_bigtable_instance.instance.name
	app_profile_id = "test"

	single_cluster_routing {
		cluster_id = "%s"
		allow_transactional_writes = true
	}
}
`, instanceName, instanceName, instanceName)
}
