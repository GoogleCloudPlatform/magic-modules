package sql_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccSqlUserListResource_queryIdentity(t *testing.T) {
	t.Parallel()

	project := envvar.GetTestProjectFromEnv()
	region := envvar.GetTestRegionFromEnv()
	instance := fmt.Sprintf("tf-test-sql-%s", acctest.RandString(t, 10))
	name1 := fmt.Sprintf("tf_test_user_%s", acctest.RandString(t, 8))
	name2 := fmt.Sprintf("tf_test_user_%s", acctest.RandString(t, 8))

	acctest.VcrTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlUserDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSqlUserListBasic(region, instance, name1, name2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_sql_user.user1", "name", name1),
					resource.TestCheckResourceAttr("google_sql_user.user1", "instance", instance),
					resource.TestCheckResourceAttr("google_sql_user.user1", "project", project),
					resource.TestCheckResourceAttr("google_sql_user.user2", "name", name2),
					resource.TestCheckResourceAttr("google_sql_user.user2", "instance", instance),
					resource.TestCheckResourceAttr("google_sql_user.user2", "project", project),
				),
			},
			{
				Query:  true,
				Config: testAccSqlUserListQuery(project, instance),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectIdentity("google_sql_user.all_in_instance", map[string]knownvalue.Check{
						"name":     knownvalue.StringExact(name1),
						"instance": knownvalue.StringExact(instance),
						"project":  knownvalue.StringExact(project),
						"host":     knownvalue.StringExact(""),
					}),
					querycheck.ExpectIdentity("google_sql_user.all_in_instance", map[string]knownvalue.Check{
						"name":     knownvalue.StringExact(name2),
						"instance": knownvalue.StringExact(instance),
						"project":  knownvalue.StringExact(project),
						"host":     knownvalue.StringExact(""),
					}),
					querycheck.ExpectLengthAtLeast("google_sql_user.all_in_instance", 2),
				},
			},
		},
	})
}

func testAccSqlUserListBasic(region, instance, name1, name2 string) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "test" {
  name                = %q
  database_version    = "POSTGRES_15"
  region              = %q
  deletion_protection = false

  settings {
    tier = "db-f1-micro"
  }
}

resource "google_sql_user" "user1" {
	name     = %q
	instance = google_sql_database_instance.test.name
	type     = "BUILT_IN"
	password = "test-password-123"
}

resource "google_sql_user" "user2" {
	name     = %q
  instance = google_sql_database_instance.test.name
  type     = "BUILT_IN"
  password = "test-password-123"
}
`, instance, region, name1, name2)
}

func testAccSqlUserListQuery(project, instance string) string {
	return fmt.Sprintf(`
provider "google" {}

list "google_sql_user" "all_in_instance" {
  provider = google

  config {
    project  = %q
    instance = %q
  }
}
`, project, instance)
}
