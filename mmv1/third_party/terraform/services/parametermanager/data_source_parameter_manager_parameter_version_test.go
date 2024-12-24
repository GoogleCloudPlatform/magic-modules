package parametermanager_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceParameterManagerParameterVersion_basicWithResourceReference(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckParameterManagerParameterVersionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAppParameterManagerParameterVersion_basicWithResourceReference(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckParameterManagerParameterDataDataSourceMatchesResource("data.google_parameter_manager_parameter_version.parameter-version-basic", "google_parameter_manager_parameter_version.parameter-version-basic"),
				),
			},
		},
	})

}

func testAccAppParameterManagerParameterVersion_basicWithResourceReference(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_parameter_manager_parameter" "parameter-basic" {
  parameter_id = "tf_test_parameter%{random_suffix}"
}

resource "google_parameter_manager_parameter_version" "parameter-version-basic" {
  parameter = google_parameter_manager_parameter.parameter-basic.id
  parameter_version_id = "tf_test_parameter_version%{random_suffix}"
  parameter_data = "test-parameter-data-with-resource-reference"
}

data "google_parameter_manager_parameter_version" "parameter-version-basic" {
  parameter = google_parameter_manager_parameter_version.parameter-version-basic.parameter
  parameter_version_id = google_parameter_manager_parameter_version.parameter-version-basic.parameter_version_id
}
`, context)
}

func TestAccDataSourceAppParameterManagerParameterVersion_basicWithParameterName(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckParameterManagerParameterVersionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAppParameterManagerParameterVersion_basicWithParameterName(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckParameterManagerParameterDataDataSourceMatchesResource("data.google_parameter_manager_parameter_version.parameter-version-basic", "google_parameter_manager_parameter_version.parameter-version-basic"),
				),
			},
		},
	})

}

func testAccAppParameterManagerParameterVersion_basicWithParameterName(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_parameter_manager_parameter" "parameter-basic" {
  parameter_id = "tf_test_parameter%{random_suffix}"
}

resource "google_parameter_manager_parameter_version" "parameter-version-basic" {
  parameter = google_parameter_manager_parameter.parameter-basic.id
  parameter_version_id = "tf_test_parameter_version%{random_suffix}"
  parameter_data = "test-parameter-data-with-parameter-name"
}

data "google_parameter_manager_parameter_version" "parameter-version-basic" {
  parameter = google_parameter_manager_parameter.parameter-basic.parameter_id
  parameter_version_id = google_parameter_manager_parameter_version.parameter-version-basic.parameter_version_id
}
`, context)
}

func TestAccDataSourceAppParameterManagerParameterVersion_withJsonData(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"data":          "./test-fixtures/parameter_data_json_format.json",
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckParameterManagerParameterVersionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAppParameterManagerParameterVersion_withJsonData(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckParameterManagerParameterDataDataSourceMatchesResource("data.google_parameter_manager_parameter_version.parameter-version-with-json-data", "google_parameter_manager_parameter_version.parameter-version-with-json-data"),
				),
			},
		},
	})

}

func testAccAppParameterManagerParameterVersion_withJsonData(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_parameter_manager_parameter" "parameter-basic" {
  parameter_id = "tf_test_parameter%{random_suffix}"
  format = "JSON"
}

resource "google_parameter_manager_parameter_version" "parameter-version-with-json-data" {
  parameter = google_parameter_manager_parameter.parameter-basic.id
  parameter_version_id = "tf_test_parameter_version%{random_suffix}"
  parameter_data = file("%{data}")
}

data "google_parameter_manager_parameter_version" "parameter-version-with-json-data" {
  parameter = google_parameter_manager_parameter_version.parameter-version-with-json-data.parameter
  parameter_version_id = google_parameter_manager_parameter_version.parameter-version-with-json-data.parameter_version_id
}
`, context)
}

func TestAccDataSourceAppParameterManagerParameterVersion_withYamlData(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"data":          "./test-fixtures/parameter_data_yaml_format.yaml",
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckParameterManagerParameterVersionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAppParameterManagerParameterVersion_withYamlData(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckParameterManagerParameterDataDataSourceMatchesResource("data.google_parameter_manager_parameter_version.parameter-version-with-yaml-data", "google_parameter_manager_parameter_version.parameter-version-with-yaml-data"),
				),
			},
		},
	})

}

func testAccAppParameterManagerParameterVersion_withYamlData(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_parameter_manager_parameter" "parameter-basic" {
  parameter_id = "tf_test_parameter%{random_suffix}"
  format = "YAML"
}

resource "google_parameter_manager_parameter_version" "parameter-version-with-yaml-data" {
  parameter = google_parameter_manager_parameter.parameter-basic.id
  parameter_version_id = "tf_test_parameter_version%{random_suffix}"
  parameter_data = file("%{data}")
}

data "google_parameter_manager_parameter_version" "parameter-version-with-yaml-data" {
  parameter = google_parameter_manager_parameter_version.parameter-version-with-yaml-data.parameter
  parameter_version_id = google_parameter_manager_parameter_version.parameter-version-with-yaml-data.parameter_version_id
}
`, context)
}

func testAccCheckParameterManagerParameterDataDataSourceMatchesResource(dataSource, resource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("can't find Parameter Version resource: %s", resource)
		}

		ds, ok := s.RootModule().Resources[dataSource]
		if !ok {
			return fmt.Errorf("can't find Parameter Version data source: %s", dataSource)
		}

		if rs.Primary.ID == "" {
			return errors.New("resource ID not set.")
		}

		if ds.Primary.ID == "" {
			return errors.New("data source ID not set.")
		}

		resourceParameterData, ok := rs.Primary.Attributes["parameter_data"]
		if !ok {
			return errors.New("can't find 'parameter_data' attribute in Parameter Version resource")
		}

		dataSourceParameterData, ok := ds.Primary.Attributes["parameter_data"]
		if !ok {
			return errors.New("can't find 'parameter_data' attribute in Parameter Version data source")
		}

		if resourceParameterData != dataSourceParameterData {
			return fmt.Errorf("expected %s, got %s, parameter_data doesn't match", resourceParameterData, dataSourceParameterData)
		}
		return nil
	}
}
