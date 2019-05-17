package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccContainerClusterDatasource_zonal(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerClusterDatasource_zonal(),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceStateWithIgnores("data.google_container_cluster.kubes", "google_container_cluster.kubes",
						// Remove once https://github.com/hashicorp/terraform/issues/21347 is fixed.
						map[string]struct{}{
							"enable_tpu":                   {},
							"enable_binary_authorization":  {},
							"pod_security_policy_config.#": {},
						},
					),
				),
			},
		},
	})
}

func TestAccContainerClusterDatasource_regional(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerClusterDatasource_regional(),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceStateWithIgnores("data.google_container_cluster.kubes", "google_container_cluster.kubes",
						// Remove once https://github.com/hashicorp/terraform/issues/21347 is fixed.
						map[string]struct{}{
							"enable_tpu":                   {},
							"enable_binary_authorization":  {},
							"pod_security_policy_config.#": {},
						},
					),
				),
			},
		},
	})
}

func testAccDataSourceGoogleContainerClusterCheck(dataSourceName string, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources[dataSourceName]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", dataSourceName)
		}

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("can't find %s in state", resourceName)
		}

		dsAttr := ds.Primary.Attributes
		rsAttr := rs.Primary.Attributes

		// Remove once https://github.com/hashicorp/terraform/issues/21347 is fixed.
		ignoreFields := map[string]struct{}{
			"enable_tpu":                   {},
			"enable_binary_authorization":  {},
			"pod_security_policy_config.#": {},
		}

		errMsg := ""
		for k, attr := range rsAttr {
			if _, ok := ignoreFields[k]; ok {
				continue
			}
			if dsAttr[k] != attr {
				errMsg += fmt.Sprintf("%s is %s; want %s\n", k, dsAttr[k], attr)
			}
		}

		if errMsg != "" {
			return fmt.Errorf(errMsg)
		}

		return nil
	}
}

func testAccContainerClusterDatasource_zonal() string {
	return fmt.Sprintf(`
resource "google_container_cluster" "kubes" {
	name               = "cluster-test-%s"
	location           = "us-central1-a"
	initial_node_count = 1

	master_auth {
		username = "mr.yoda"
		password = "adoy.rm.123456789"
	}
}

data "google_container_cluster" "kubes" {
	name     = "${google_container_cluster.kubes.name}"
	location = "${google_container_cluster.kubes.zone}"
}
`, acctest.RandString(10))
}

func testAccContainerClusterDatasource_regional() string {
	return fmt.Sprintf(`
resource "google_container_cluster" "kubes" {
	name               = "cluster-test-%s"
	location           = "us-central1"
	initial_node_count = 1
}

data "google_container_cluster" "kubes" {
	name     = "${google_container_cluster.kubes.name}"
	location = "${google_container_cluster.kubes.region}"
}
`, acctest.RandString(10))
}
