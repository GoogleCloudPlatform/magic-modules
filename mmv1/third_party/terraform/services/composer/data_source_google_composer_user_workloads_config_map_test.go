package composer_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceComposerUserWorkloadsConfigMap_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"env_name":        fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, acctest.RandInt(t)),
		"config_map_name": fmt.Sprintf("%s-%d", testComposerUserWorkloadsConfigMapPrefix, acctest.RandInt(t)),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceComposerUserWorkloadsConfigMap_basic(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_composer_user_workloads_config_map.test",
						"google_composer_user_workloads_config_map.test"),
				),
			},
		},
	})
}

/*
func checkConfigMapDataSourceMatchesResource() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources["data.google_composer_user_workloads_config_map.test"]
		if !ok {
			return fmt.Errorf("can't find %s in state", "data.google_composer_user_workloads_config_map.test")
		}
		rs, ok := s.RootModule().Resources["google_composer_user_workloads_config_map.test"]
		if !ok {
			return fmt.Errorf("can't find %s in state", "google_composer_user_workloads_config_map.test")
		}

		dsAttr := ds.Primary.Attributes
		rsAttr := rs.Primary.Attributes
		errMsg := ""

		for k := range rsAttr {
			if k == "%" {
				continue
			}
			// ignore diff if it's due to secrets being masked.
			if strings.HasPrefix(k, "data.") && k != "data.%" {
				if dsAttr[k] == "**********" {
					continue
				}
			}
			if dsAttr[k] != rsAttr[k] {
				errMsg += fmt.Sprintf("%s is %s; want %s\n", k, dsAttr[k], rsAttr[k])
			}
		}

		if errMsg != "" {
			return errors.New(errMsg)
		}

		return nil
	}
}
*/

func testAccDataSourceComposerUserWorkloadsConfigMap_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_composer_environment" "test" {
  name   = "%{env_name}"
  config {
    software_config {
      image_version = "composer-3-airflow-2"
    }
  }
}
resource "google_composer_user_workloads_config_map" "test" {
  environment = google_composer_environment.test.name
  name = "%{config_map_name}"
  data = {
    db_host: "dbhost:5432",
    api_host: "apihost:443",
  }
}
data "google_composer_user_workloads_config_map" "test" {
  name        = google_composer_user_workloads_config_map.test.name
  environment = google_composer_environment.test.name
}
`, context)
}
