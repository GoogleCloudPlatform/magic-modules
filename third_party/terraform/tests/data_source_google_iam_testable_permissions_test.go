package google

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDataSourceGoogleIamTestablePermissions_basic(t *testing.T) {
	t.Parallel()

	project := getTestProjectFromEnv()
	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
			 data "google_iam_testable_permissions" "perms" {
				full_resource_name = "//cloudresourcemanager.googleapis.com/projects/%s"
			}
		`, project),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleIamTestablePermissionsMeta(
						project,
						"data.google_iam_testable_permissions.perms",
						"GA",
						"",
					),
				),
			},
			{
				Config: fmt.Sprintf(`
			 data "google_iam_testable_permissions" "perms" {
				full_resource_name   = "//cloudresourcemanager.googleapis.com/projects/%s"
				custom_support_level = "NOT_SUPPORTED"
				stage                = "BETA"
			}
		`, project),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleIamTestablePermissionsMeta(
						project,
						"data.google_iam_testable_permissions.perms",
						"BETA",
						"NOT_SUPPORTED",
					),
				),
			},
			{
				Config: fmt.Sprintf(`
			 data "google_iam_testable_permissions" "perms" {
				full_resource_name   = "//cloudresourcemanager.googleapis.com/projects/%s"
				custom_support_level = "not_supported"
				stage                = "beta"
			}
		`, project),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleIamTestablePermissionsMeta(
						project,
						"data.google_iam_testable_permissions.perms",
						"BETA",
						"NOT_SUPPORTED",
					),
				),
			},
		},
	})
}

func testAccCheckGoogleIamTestablePermissionsMeta(project string, n string, expectedStage string, expectedSupportLevel string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find perms data source: %s", n)
		}
		expectedId := fmt.Sprintf("//cloudresourcemanager.googleapis.com/projects/%s", project)
		if rs.Primary.ID != expectedId {
			return fmt.Errorf("perms data source ID not set.")
		}
		attrs := rs.Primary.Attributes
		count, ok := attrs["permissions.#"]
		if !ok {
			return fmt.Errorf("can't find 'permsissions' attribute")
		}
		permCount, err := strconv.Atoi(count)
		if err != nil {
			return err
		}
		if permCount < 2 {
			return fmt.Errorf("count should be greater than 2")
		}
		foundStage := false
		foundSupport := false

		for i := 0; i < permCount; i++ {
			stageKey := "permissions." + strconv.Itoa(i) + ".stage"
			supportKey := "permissions." + strconv.Itoa(i) + ".custom_support_level"
			if attrs[stageKey] == expectedStage {
				foundStage = true
			}
			if attrs[supportKey] == expectedSupportLevel {
				foundSupport = true
			}
			if foundSupport && foundStage {
				return nil
			}
		}

		if foundSupport {
			return fmt.Errorf("Could not find stage %s in output", expectedStage)
		}
		if foundStage {
			return fmt.Errorf("Could not find custom_support_level %s in output", expectedSupportLevel)
		}
		return fmt.Errorf("Unable to find customeSupportLevel or stage")
	}
}
