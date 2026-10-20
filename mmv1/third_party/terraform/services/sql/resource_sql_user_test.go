package sql_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	_ "github.com/hashicorp/terraform-provider-google/google/services/resourcemanager"
	"github.com/hashicorp/terraform-provider-google/google/services/sql"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"
)

func TestAccSqlUser_mysql(t *testing.T) {
	// Multiple fine-grained resources
	acctest.SkipIfVcr(t)
	t.Parallel()

	instance := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlUserDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlUser_mysql(instance, "password"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlUserExists(t, "google_sql_user.user1"),
					testAccCheckGoogleSqlUserExists(t, "google_sql_user.user2"),
				),
			},
			{
				// Update password
				Config: testGoogleSqlUser_mysql(instance, "new_password"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlUserExists(t, "google_sql_user.user1"),
					testAccCheckGoogleSqlUserExists(t, "google_sql_user.user2"),
					testAccCheckGoogleSqlUserExists(t, "google_sql_user.user3"),
				),
			},
			{
				ResourceName:            "google_sql_user.user2",
				ImportStateId:           fmt.Sprintf("%s/%s/gmail.com/admin", envvar.GetTestProjectFromEnv(), instance),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
			{
				ResourceName:            "google_sql_user.user3",
				ImportStateId:           fmt.Sprintf("%s/%s/10.0.0.0/24/admin", envvar.GetTestProjectFromEnv(), instance),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func TestAccSqlUser_password_wo(t *testing.T) {
	t.Parallel()

	instance := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlUserDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlUser_password_wo(instance, "password"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlUserExists(t, "google_sql_user.user1"),
				),
			},
			{
				// Update password
				Config: testGoogleSqlUser_new_password_wo(instance, "new_password"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_sql_user.user1", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlUserExists(t, "google_sql_user.user1"),
				),
			},
			{
				ResourceName:            "google_sql_user.user1",
				ImportStateId:           fmt.Sprintf("%s/%s/gmail.com/admin", envvar.GetTestProjectFromEnv(), instance),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func TestAccSqlUser_iamUser(t *testing.T) {
	// Multiple fine-grained resources
	acctest.SkipIfVcr(t)
	t.Parallel()

	instance := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlUserDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlUser_iamUser(instance),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlUserExists(t, "google_sql_user.user1"),
				),
			},
			{
				ResourceName:      "google_sql_user.user1",
				ImportStateId:     fmt.Sprintf("%s/%s/%%/%s@%s.iam.gserviceaccount.com", envvar.GetTestProjectFromEnv(), instance, instance, envvar.GetTestProjectFromEnv()),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSqlUser_iamGroupUser(t *testing.T) {
	// Multiple fine-grained resources
	acctest.SkipIfVcr(t)
	t.Parallel()

	instance := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlUserDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlUser_iamGroupUser("iam-group-auth-test-group@google.com", instance),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlUserExists(t, "google_sql_user.user"),
				),
			},
			{
				ResourceName:            "google_sql_user.user",
				ImportStateId:           fmt.Sprintf("%s/%s/iam-group-auth-test-group@google.com", envvar.GetTestProjectFromEnv(), instance),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func TestAccSqlUser_iamGroupUser_capitalizedHostName(t *testing.T) {
	// Multiple fine-grained resources
	acctest.SkipIfVcr(t)
	t.Parallel()

	instance := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlUserDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlUser_iamGroupUser("iam-group-auth-test-group@GOOGLE.com", instance),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlUserExists(t, "google_sql_user.user"),
				),
			},
			{
				ResourceName:            "google_sql_user.user",
				ImportStateId:           fmt.Sprintf("%s/%s/iam-group-auth-test-group@google.com", envvar.GetTestProjectFromEnv(), instance),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func TestAccSqlUser_postgres_iamGroupUser(t *testing.T) {
	// Multiple fine-grained resources
	acctest.SkipIfVcr(t)
	t.Parallel()

	instance := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		CheckDestroy: testAccSqlUserDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlUser_postgres_iamGroupUser("iam-group-auth-test-group@google.com", instance),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlUserExists(t, "google_sql_user.user"),
				),
			},
			{
				ResourceName:            "google_sql_user.user",
				ImportStateId:           fmt.Sprintf("%s/%s/iam-group-auth-test-group@google.com", envvar.GetTestProjectFromEnv(), instance),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func TestAccSqlUser_postgres(t *testing.T) {
	t.Parallel()

	instance := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlUserDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlUser_postgres(instance, "password"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlUserExists(t, "google_sql_user.user"),
				),
			},
			{
				// Update password
				Config: testGoogleSqlUser_postgres(instance, "new_password"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlUserExists(t, "google_sql_user.user"),
				),
			},
			{
				ResourceName:            "google_sql_user.user",
				ImportStateId:           fmt.Sprintf("%s/%s/admin", envvar.GetTestProjectFromEnv(), instance),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func TestAccSqlUser_postgresIAM(t *testing.T) {
	t.Parallel()

	instance := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	const iamUser = "admin@hashicorptest.com"
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		CheckDestroy: testAccSqlUserDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlUser_postgresIAM(instance, iamUser),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlUserExists(t, "google_sql_user.user"),
				),
			},
			{
				ResourceName:            "google_sql_user.user",
				ImportStateId:           fmt.Sprintf("%s/%s/%s", envvar.GetTestProjectFromEnv(), instance, iamUser),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func TestAccSqlUser_postgresAbandon(t *testing.T) {
	t.Parallel()

	instance := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	userName := "admin"
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlUserDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlUser_postgresAbandon(instance, userName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlUserExists(t, "google_sql_user.user"),
				),
			},
			{
				ResourceName:            "google_sql_user.user",
				ImportStateId:           fmt.Sprintf("%s/%s/admin", envvar.GetTestProjectFromEnv(), instance),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "deletion_policy"},
			},
			{
				// Abandon user
				Config: testGoogleSqlUser_postgresNoUser(instance),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlUserExistsWithName(t, instance, userName),
				),
			},
		},
	})
}

func testAccCheckGoogleSqlUserExists(t *testing.T, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		name := rs.Primary.Attributes["name"]
		instance := rs.Primary.Attributes["instance"]
		host := rs.Primary.Attributes["host"]
		databaseInstance, err := sql.NewClient(config, config.UserAgent).Instances.Get(config.Project, instance).Do()
		if err != nil {
			return err
		}
		users, err := sql.NewClient(config, config.UserAgent).Users.List(config.Project,
			instance).Do()

		if err != nil {
			return err
		}

		for _, user := range users.Items {
			username := name
			if user.Type == "CLOUD_IAM_GROUP" && strings.Contains(databaseInstance.DatabaseVersion, "MYSQL") {
				splitName := strings.SplitN(name, "@", 2)
				if len(splitName) == 2 {
					groupUsername := splitName[0]
					groupHostname := splitName[1]
					username = groupUsername + "@" + strings.ToLower(groupHostname)
				}
			}
			if user.Name == username && user.Host == host {
				return nil
			}
		}

		return fmt.Errorf("Not found: %s: %s", n, err)
	}
}

func testAccCheckGoogleSqlUserExistsWithName(t *testing.T, instance, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)

		users, err := sql.NewClient(config, config.UserAgent).Users.List(config.Project,
			instance).Do()

		if err != nil {
			return err
		}

		for _, user := range users.Items {
			if user.Name == name {
				return nil
			}
		}

		return fmt.Errorf("Not found: User: %s in instance: %s: %s", name, instance, err)
	}
}

func testAccSqlUserDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			config := acctest.GoogleProviderConfig(t)
			if rs.Type != "google_sql_database" {
				continue
			}

			name := rs.Primary.Attributes["name"]
			instance := rs.Primary.Attributes["instance"]
			host := rs.Primary.Attributes["host"]
			users, err := sql.NewClient(config, config.UserAgent).Users.List(config.Project,
				instance).Do()

			if users == nil {
				return nil
			}

			for _, user := range users.Items {
				if user.Name == name && user.Host == host {
					return fmt.Errorf("User still %s exists %s", name, err)
				}
			}

			return nil
		}

		return nil
	}
}

func TestAccSqlUser_mysqlPasswordPolicy(t *testing.T) {
	// Multiple fine-grained resources
	acctest.SkipIfVcr(t)
	t.Parallel()

	instance := fmt.Sprintf("tf-test-i%d", acctest.RandInt(t))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlUserDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlUser_mysqlPasswordPolicy(instance, "password"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlUserExists(t, "google_sql_user.user1"),
					testAccCheckGoogleSqlUserExists(t, "google_sql_user.user2"),
				),
			},
			{
				// Update password
				Config: testGoogleSqlUser_mysqlPasswordPolicy(instance, "new_password"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlUserExists(t, "google_sql_user.user1"),
					testAccCheckGoogleSqlUserExists(t, "google_sql_user.user2"),
				),
			},
			{
				ResourceName:            "google_sql_user.user2",
				ImportStateId:           fmt.Sprintf("%s/%s/gmail.com/admin", envvar.GetTestProjectFromEnv(), instance),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func TestAccSqlUser_instanceWithActivationPolicy(t *testing.T) {
	// Multiple fine-grained resources
	acctest.SkipIfVcr(t)
	t.Parallel()

	instance := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlUserDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlUser_instanceWithActivationPolicy(instance, "ALWAYS"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlUserExists(t, "google_sql_user.user"),
				),
			},
			// Step 2: Update activation_policy to NEVER
			{
				Config: testGoogleSqlUser_instanceWithActivationPolicy(instance, "NEVER"),
			},
			// Step 3: Refresh with instance stopped — identity attributes must still be set
			{
				Config: testGoogleSqlUser_instanceWithActivationPolicy(instance, "NEVER"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_sql_user.user", "project"),
					resource.TestCheckResourceAttr("google_sql_user.user", "instance", instance),
					resource.TestCheckResourceAttr("google_sql_user.user", "name", "admin"),
				),
			},
			// Step 3a: ImportState with instance stopped — verifies identity survives a round-trip
			{
				ResourceName:            "google_sql_user.user",
				ImportStateId:           fmt.Sprintf("%s/%s/admin", envvar.GetTestProjectFromEnv(), instance),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "deletion_policy"},
			},
			// Step 4: Update activation_policy to ALWAYS so that post-test destroy code is able to delete the google_sql_user resource
			{
				Config: testGoogleSqlUser_instanceWithActivationPolicy(instance, "ALWAYS"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlUserExists(t, "google_sql_user.user"),
				),
			},
		},
	})
}

func TestAccSqlUser_mysql_createUserWithDatabaseRoles(t *testing.T) {
	t.Parallel()

	instance := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlUserDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlUser_mysql_userWithDatabaseRoles(instance, "testuser", "gmail.com", "BUILT_IN", "password", `["cloudsqlsuperuser"]`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlUserExists(t, "google_sql_user.testuser"),
				),
			},
			{
				ResourceName:            "google_sql_user.testuser",
				ImportStateId:           fmt.Sprintf("%s/%s/gmail.com/testuser", envvar.GetTestProjectFromEnv(), instance),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "database_roles"},
			},
			{
				Config: testGoogleSqlUser_mysql_userWithDatabaseRoles(instance, "admin@hashicorptest.com", "", "CLOUD_IAM_USER", "", `["cloudsqlsuperuser"]`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlUserExists(t, "google_sql_user.admin"),
				),
			},
			{
				ResourceName:            "google_sql_user.admin",
				ImportStateId:           fmt.Sprintf("%s/%s/%%/admin@hashicorptest.com", envvar.GetTestProjectFromEnv(), instance),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "database_roles"},
			},
			{
				Config: testGoogleSqlUser_mysql_userWithDatabaseRoles(instance, "iam-group-auth-test-group@google.com", "", "CLOUD_IAM_GROUP", "", `["cloudsqlsuperuser"]`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlUserExists(t, "google_sql_user.iam-group-auth-test-group"),
				),
			},
			{
				ResourceName:            "google_sql_user.iam-group-auth-test-group",
				ImportStateId:           fmt.Sprintf("%s/%s/%%/iam-group-auth-test-group@google.com", envvar.GetTestProjectFromEnv(), instance),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "database_roles"},
			},
		},
	})
}

func TestAccSqlUser_mysql_updateUserWithDatabaseRoles(t *testing.T) {
	t.Parallel()

	instance := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlUserDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlUser_mysql_userWithDatabaseRoles(instance, "testuser", "gmail.com", "BUILT_IN", "password", `[]`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlUserExists(t, "google_sql_user.testuser"),
				),
			},
			{
				ResourceName:            "google_sql_user.testuser",
				ImportStateId:           fmt.Sprintf("%s/%s/gmail.com/testuser", envvar.GetTestProjectFromEnv(), instance),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "database_roles"},
			},
			{
				Config: testGoogleSqlUser_mysql_userWithDatabaseRoles(instance, "testuser", "gmail.com", "BUILT_IN", "password", `["cloudsqlsuperuser"]`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlUserExists(t, "google_sql_user.testuser"),
				),
			},
			{
				ResourceName:            "google_sql_user.testuser",
				ImportStateId:           fmt.Sprintf("%s/%s/gmail.com/testuser", envvar.GetTestProjectFromEnv(), instance),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "database_roles"},
			},
			{
				Config: testGoogleSqlUser_mysql_userWithDatabaseRoles(instance, "admin@hashicorptest.com", "", "CLOUD_IAM_USER", "", `[]`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlUserExists(t, "google_sql_user.admin"),
				),
			},
			{
				ResourceName:            "google_sql_user.admin",
				ImportStateId:           fmt.Sprintf("%s/%s/%%/admin@hashicorptest.com", envvar.GetTestProjectFromEnv(), instance),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "database_roles"},
			},
			{
				Config: testGoogleSqlUser_mysql_userWithDatabaseRoles(instance, "admin@hashicorptest.com", "", "CLOUD_IAM_USER", "", `["cloudsqlsuperuser"]`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlUserExists(t, "google_sql_user.admin"),
				),
			},
			{
				ResourceName:            "google_sql_user.admin",
				ImportStateId:           fmt.Sprintf("%s/%s/%%/admin@hashicorptest.com", envvar.GetTestProjectFromEnv(), instance),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "database_roles"},
			},
			{
				Config: testGoogleSqlUser_mysql_saUserWithDatabaseRoles(instance, "sa_user", `[]`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlUserExists(t, "google_sql_user.sa_user"),
				),
			},
			{
				ResourceName:            "google_sql_user.sa_user",
				ImportStateId:           fmt.Sprintf("%s/%s/%%/%s@%s.iam.gserviceaccount.com", envvar.GetTestProjectFromEnv(), instance, instance, envvar.GetTestProjectFromEnv()),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "database_roles"},
			},
			{
				Config: testGoogleSqlUser_mysql_saUserWithDatabaseRoles(instance, "sa_user", `["cloudsqlsuperuser"]`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlUserExists(t, "google_sql_user.sa_user"),
				),
			},
			{
				ResourceName:            "google_sql_user.sa_user",
				ImportStateId:           fmt.Sprintf("%s/%s/%%/%s@%s.iam.gserviceaccount.com", envvar.GetTestProjectFromEnv(), instance, instance, envvar.GetTestProjectFromEnv()),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "database_roles"},
			},
		},
	})
}

func TestAccSqlUser_postgres_updateUserWithDatabaseRoles(t *testing.T) {
	t.Parallel()

	instance := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		CheckDestroy: testAccSqlUserDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlUser_postgres_userWithDatabaseRoles(instance, "testuser", "BUILT_IN", "password", `[]`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlUserExists(t, "google_sql_user.testuser"),
				),
			},
			{
				ResourceName:            "google_sql_user.testuser",
				ImportStateId:           fmt.Sprintf("%s/%s/testuser", envvar.GetTestProjectFromEnv(), instance),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "database_roles"},
			},
			{
				Config: testGoogleSqlUser_postgres_userWithDatabaseRoles(instance, "testuser", "BUILT_IN", "password", `["cloudsqlsuperuser"]`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlUserExists(t, "google_sql_user.testuser"),
				),
			},
			{
				ResourceName:            "google_sql_user.testuser",
				ImportStateId:           fmt.Sprintf("%s/%s/testuser", envvar.GetTestProjectFromEnv(), instance),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "database_roles"},
			},
			{
				Config: testGoogleSqlUser_postgres_iamUserWithDatabaseRoles(instance, "admin@hashicorptest.com", `[]`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlUserExists(t, "google_sql_user.admin"),
				),
			},
			{
				ResourceName:            "google_sql_user.admin",
				ImportStateId:           fmt.Sprintf("%s/%s/admin@hashicorptest.com", envvar.GetTestProjectFromEnv(), instance),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "database_roles"},
			},
			{
				Config: testGoogleSqlUser_postgres_iamUserWithDatabaseRoles(instance, "admin@hashicorptest.com", `["cloudsqlsuperuser"]`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlUserExists(t, "google_sql_user.admin"),
				),
			},
			{
				ResourceName:            "google_sql_user.admin",
				ImportStateId:           fmt.Sprintf("%s/%s/admin@hashicorptest.com", envvar.GetTestProjectFromEnv(), instance),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "database_roles"},
			},
			{
				Config: testGoogleSqlUser_postgres_saUserWithDatabaseRoles(instance, "sa_user", `[]`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlUserExists(t, "google_sql_user.sa_user"),
				),
			},
			{
				ResourceName:            "google_sql_user.sa_user",
				ImportStateId:           fmt.Sprintf("%s/%s//%s@%s.iam", envvar.GetTestProjectFromEnv(), instance, instance, envvar.GetTestProjectFromEnv()),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "database_roles"},
			},
			{
				Config: testGoogleSqlUser_postgres_saUserWithDatabaseRoles(instance, "sa_user", `["cloudsqlsuperuser"]`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlUserExists(t, "google_sql_user.sa_user"),
				),
			},
			{
				ResourceName:            "google_sql_user.sa_user",
				ImportStateId:           fmt.Sprintf("%s/%s//%s@%s.iam", envvar.GetTestProjectFromEnv(), instance, instance, envvar.GetTestProjectFromEnv()),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "database_roles"},
			},
		},
	})
}

func testGoogleSqlUser_mysql_userWithDatabaseRoles(instance, username, host, usertype, password, roles string) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "instance" {
  name                = "%s"
  region              = "us-central1"
  database_version    = "MYSQL_8_0"
  deletion_protection = false
  settings {
    tier = "db-f1-micro"
    database_flags {
      name  = "cloudsql_iam_authentication"
      value = "on"
    }
  }
}

resource "google_sql_user" "%s" {
  name     = "%s"
  host     = "%s"
  password = "%s"
  type     = "%s"
  instance = google_sql_database_instance.instance.name
  database_roles = %s
}
`, instance, strings.Split(username, "@")[0], username, host, password, usertype, roles)
}

func testGoogleSqlUser_mysql_saUserWithDatabaseRoles(instance, username, roles string) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "instance" {
  name                = "%s"
  region              = "us-central1"
  database_version    = "MYSQL_8_0"
  deletion_protection = false
  settings {
    tier = "db-f1-micro"
    database_flags {
      name  = "cloudsql_iam_authentication"
      value = "on"
    }
  }
}

resource "google_service_account" "sa" {
  account_id   = "%s"
  display_name = "%s"
}

resource "google_service_account_key" "sa_key" {
  service_account_id = google_service_account.sa.email
}

resource "google_sql_user" "%s" {
  name     = google_service_account.sa.email
  instance = google_sql_database_instance.instance.name
  type     = "CLOUD_IAM_SERVICE_ACCOUNT"
  database_roles = %s
}
`, instance, instance, instance, username, roles)
}

func testGoogleSqlUser_postgres_saUserWithDatabaseRoles(instance, username, roles string) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "instance" {
  name                = "%s"
  region              = "us-central1"
  database_version    = "POSTGRES_9_6"
  deletion_protection = false
  settings {
    tier = "db-f1-micro"
    database_flags {
      name  = "cloudsql.iam_authentication"
      value = "on"
    }
  }
}

# TODO: Remove with resolution of https://github.com/hashicorp/terraform-provider-google/issues/14233
resource "time_sleep" "wait_60_seconds" {
  depends_on = [google_sql_database_instance.instance]
  create_duration = "60s"
}

resource "google_service_account" "sa" {
  account_id   = "%s"
  display_name = "%s"
}

resource "google_service_account_key" "sa_key" {
  service_account_id = google_service_account.sa.email
}

resource "google_sql_user" "%s" {
  depends_on = [time_sleep.wait_60_seconds]
  name     = trimsuffix(google_service_account.sa.email, ".gserviceaccount.com")
  instance = google_sql_database_instance.instance.name
  type     = "CLOUD_IAM_SERVICE_ACCOUNT"
  database_roles = %s
}
`, instance, instance, instance, username, roles)
}

func testGoogleSqlUser_postgres_userWithDatabaseRoles(instance, username, usertype, password, roles string) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "instance" {
  name                = "%s"
  region              = "us-central1"
  database_version    = "POSTGRES_9_6"
  deletion_protection = false
  settings {
    tier = "db-f1-micro"
    database_flags {
      name  = "cloudsql.iam_authentication"
      value = "on"
    }
  }
}

resource "google_sql_user" "%s" {
  name     = "%s"
  password = "%s"
  type     = "%s"
  instance = google_sql_database_instance.instance.name
  database_roles = %s
}
`, instance, strings.Split(username, "@")[0], username, password, usertype, roles)
}

func testGoogleSqlUser_postgres_iamUserWithDatabaseRoles(instance, username, roles string) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "instance" {
  name                = "%s"
  region              = "us-central1"
  database_version    = "POSTGRES_9_6"
  deletion_protection = false
  settings {
    tier = "db-f1-micro"
    database_flags {
      name  = "cloudsql.iam_authentication"
      value = "on"
    }
  }
}

# TODO: Remove with resolution of https://github.com/hashicorp/terraform-provider-google/issues/14233
resource "time_sleep" "wait_60_seconds" {
  depends_on = [google_sql_database_instance.instance]
  create_duration = "60s"
}

resource "google_sql_user" "%s" {
  depends_on = [time_sleep.wait_60_seconds]
  name     = "%s"
  type     = "CLOUD_IAM_USER"
  instance = google_sql_database_instance.instance.name
  database_roles = %s
}
`, instance, strings.Split(username, "@")[0], username, roles)
}

func testGoogleSqlUser_instanceWithActivationPolicy(instance, activationPolicy string) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "instance" {
  name             = "%s"
  database_version = "MYSQL_5_7"
  region          = "us-central1"
  deletion_protection = false
  settings {
    tier = "db-f1-micro"
    availability_type = "ZONAL"
    activation_policy = "%s"
  }
}

resource "google_sql_user" "user" {
	name     = "admin"
	instance = google_sql_database_instance.instance.name
	password = "password"
  }
`, instance, activationPolicy)
}

func testGoogleSqlUser_password_wo(instance, password string) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "instance" {
  name                = "%s"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false
  settings {
    tier = "db-f1-micro"
  }
}

resource "google_sql_user" "user1" {
  name     = "admin"
  instance = google_sql_database_instance.instance.name
  host     = "gmail.com"
  password_wo = "%s"
  password_wo_version = 1
}
`, instance, password)
}

func testGoogleSqlUser_new_password_wo(instance, password string) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "instance" {
  name                = "%s"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false
  settings {
    tier = "db-f1-micro"
  }
}

resource "google_sql_user" "user1" {
  name     = "admin"
  instance = google_sql_database_instance.instance.name
  host     = "gmail.com"
  password_wo = "%s"
  password_wo_version = 2
}
`, instance, password)
}

func testGoogleSqlUser_mysql(instance, password string) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "instance" {
  name                = "%s"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false
  settings {
    tier = "db-f1-micro"
  }
}

resource "google_sql_user" "user1" {
  name     = "admin"
  instance = google_sql_database_instance.instance.name
  host     = "google.com"
  password = "%s"
}

resource "google_sql_user" "user2" {
  name     = "admin"
  instance = google_sql_database_instance.instance.name
  host     = "gmail.com"
  password = "hunter2"
}

resource "google_sql_user" "user3" {
  name     = "admin"
  instance = google_sql_database_instance.instance.name
  host     = "10.0.0.0/24"
  password = "hunter3"
}
`, instance, password)
}

func testGoogleSqlUser_mysqlPasswordPolicy(instance, password string) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "instance" {
  name                = "%s"
  region              = "us-central1"
  database_version    = "MYSQL_8_0"
  deletion_protection = false
  settings {
    tier = "db-f1-micro"
  }
}

resource "google_sql_user" "user1" {
  name     = "admin"
  instance = google_sql_database_instance.instance.name
  host     = "google.com"
  password = "%s"

  password_policy {
    allowed_failed_attempts  = 6
    password_expiration_duration  =  "2592000s"
    enable_failed_attempts_check = true
    enable_password_verification = true
  }
}

resource "google_sql_user" "user2" {
  name     = "admin"
  instance = google_sql_database_instance.instance.name
  host     = "gmail.com"
  password = "hunter2"
  password_policy {
    allowed_failed_attempts  = 6
    enable_failed_attempts_check = true
  }
}
`, instance, password)
}

func testGoogleSqlUser_postgres(instance, password string) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "instance" {
  name             = "%s"
  region           = "us-central1"
  database_version = "POSTGRES_9_6"
  deletion_protection = false

  settings {
    tier = "db-f1-micro"
  }
}

resource "google_sql_user" "user" {
  name     = "admin"
  instance = google_sql_database_instance.instance.name
  password = "%s"
}
`, instance, password)
}

func testGoogleSqlUser_postgresIAM(instance, iamUser string) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "instance" {
  name             = "%s"
  region           = "us-central1"
  database_version = "POSTGRES_9_6"
  deletion_protection = false

  settings {
    tier = "db-f1-micro"
    database_flags {
      name  = "cloudsql.iam_authentication"
      value = "on"
    }
  }
}

# TODO: Remove with resolution of https://github.com/hashicorp/terraform-provider-google/issues/14233
resource "time_sleep" "wait_60_seconds" {
  depends_on = [google_sql_database_instance.instance]

  create_duration = "60s"
}

resource "google_sql_user" "user" {
  depends_on = [time_sleep.wait_60_seconds]
  name     = "%s"
  instance = google_sql_database_instance.instance.name
  type     = "CLOUD_IAM_USER"
}
`, instance, iamUser)
}

func testGoogleSqlUser_postgresAbandon(instance, name string) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "instance" {
  name             = "%s"
  region           = "us-central1"
  database_version = "POSTGRES_9_6"
  deletion_protection = false

  settings {
    tier = "db-f1-micro"
  }
}

resource "google_sql_user" "user" {
  name     = "%s"
  instance = google_sql_database_instance.instance.name
  password = "password"
  deletion_policy = "ABANDON"
}
`, instance, name)
}

func testGoogleSqlUser_postgresNoUser(instance string) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "instance" {
  name             = "%s"
  region           = "us-central1"
  database_version = "POSTGRES_9_6"
  deletion_protection = false

  settings {
    tier = "db-f1-micro"
  }
}
`, instance)
}

func testGoogleSqlUser_iamUser(instance string) string {
	return fmt.Sprintf(`
data "google_project" "project" {}

resource "google_sql_database_instance" "instance" {
  database_version = "MYSQL_8_0"
  name             = "%s"
  region           = "us-central1"

  settings {
    tier              = "db-f1-micro"
    availability_type = "REGIONAL"

    backup_configuration {
      enabled            = true
      binary_log_enabled = true
    }

    database_flags {
      name  = "cloudsql_iam_authentication"
      value = "on"
    }
  }

  deletion_protection = false
}

resource "google_sql_database" "db" {
  name     = "%s"
  instance = google_sql_database_instance.instance.name
}

resource "google_service_account" "sa" {
  account_id   = "%s"
  display_name = "%s"
}

resource "google_service_account_key" "sa_key" {
  service_account_id = google_service_account.sa.email
}

resource "google_sql_user" "user1" {
  name     = google_service_account.sa.email
  instance = google_sql_database_instance.instance.name
  type     = "CLOUD_IAM_SERVICE_ACCOUNT"
}

resource "google_project_iam_member" "instance_user" {
  project = data.google_project.project.project_id
  role    = "roles/cloudsql.instanceUser"
  member  = "serviceAccount:${google_service_account.sa.email}"
}

resource "google_project_iam_member" "sa_user" {
  project = data.google_project.project.project_id
  role    = "roles/iam.serviceAccountUser"
  member  = "serviceAccount:${google_service_account.sa.email}"
}
`, instance, instance, instance, instance)
}

func testGoogleSqlUser_iamGroupUser(username, instance string) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "instance" {
  name                = "%s"
  region              = "us-central1"
  database_version    = "MYSQL_8_0"
  deletion_protection = false
  settings {
    tier = "db-f1-micro"
    database_flags {
      name  = "cloudsql_iam_authentication"
      value = "on"
    }
  }
}

resource "google_sql_user" "user" {
  name     = "%s"
  instance = google_sql_database_instance.instance.name
  type     = "CLOUD_IAM_GROUP"
}
`, instance, username)
}

func testGoogleSqlUser_postgres_iamGroupUser(username, instance string) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "instance" {
  name                = "%s"
  region              = "us-central1"
  database_version    = "POSTGRES_9_6"
  deletion_protection = false
  settings {
    tier = "db-f1-micro"
    database_flags {
      name  = "cloudsql.iam_authentication"
      value = "on"
    }
  }
}

# TODO: Remove with resolution of https://github.com/hashicorp/terraform-provider-google/issues/14233
resource "time_sleep" "wait_60_seconds" {
  depends_on = [google_sql_database_instance.instance]

  create_duration = "60s"
}

resource "google_sql_user" "user" {
  depends_on = [time_sleep.wait_60_seconds]
  name     = "%s"
  instance = google_sql_database_instance.instance.name
  type     = "CLOUD_IAM_GROUP"
}
`, instance, username)
}

func TestAccSqlUser_importBlockWithResourceIdentity(t *testing.T) {
	t.Parallel()
	instance := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	acctest.VcrTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_12_0),
		},
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlUserDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlUser_postgres(instance, "password"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlUserExists(t, "google_sql_user.user"),
				),
			},
			{
				Config:          testGoogleSqlUser_postgresNoPassword(instance),
				ResourceName:    "google_sql_user.user",
				ImportState:     true,
				ImportStateKind: resource.ImportBlockWithResourceIdentity,
			},
		},
	})
}

func testGoogleSqlUser_postgresNoPassword(instance string) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "instance" {
  name             = "%s"
  region           = "us-central1"
  database_version = "POSTGRES_9_6"
  deletion_protection = false
  settings {
    tier = "db-f1-micro"
  }
}
resource "google_sql_user" "user" {
  name     = "admin"
  instance = google_sql_database_instance.instance.name
}
`, instance)
}

func TestSqlUser_findSqlUser(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name            string
		databaseVersion string
		userName        string
		host            string
		users           []*sqladmin.User
		wantName        string
		wantType        string
	}{
		{
			name:            "SQL Server ENTRAID_USER preserves full email UPN",
			databaseVersion: "SQLSERVER_2022_STANDARD",
			userName:        "user@tenant.onmicrosoft.com",
			host:            "",
			users: []*sqladmin.User{
				{Name: "sqlserver", Type: "BUILT_IN"},
				{Name: "user@tenant.onmicrosoft.com", Type: "ENTRAID_USER"},
			},
			wantName: "user@tenant.onmicrosoft.com",
			wantType: "ENTRAID_USER",
		},
		{
			name:            "SQL Server ENTRAID_USER does not collide with same-prefix BUILT_IN login",
			databaseVersion: "SQLSERVER_2022_STANDARD",
			userName:        "sqlserver@tenant.onmicrosoft.com",
			host:            "",
			users: []*sqladmin.User{
				{Name: "sqlserver", Type: "BUILT_IN"},
				{Name: "sqlserver@tenant.onmicrosoft.com", Type: "ENTRAID_USER"},
			},
			wantName: "sqlserver@tenant.onmicrosoft.com",
			wantType: "ENTRAID_USER",
		},
		{
			name:            "SQL Server ENTRAID_USER missing from list returns nil",
			databaseVersion: "SQLSERVER_2022_STANDARD",
			userName:        "missing@tenant.onmicrosoft.com",
			host:            "",
			users: []*sqladmin.User{
				{Name: "sqlserver", Type: "BUILT_IN"},
			},
			wantName: "",
		},
		{
			name:            "SQL Server BUILT_IN user matches exact username",
			databaseVersion: "SQLSERVER_2022_STANDARD",
			userName:        "appuser",
			host:            "",
			users: []*sqladmin.User{
				{Name: "appuser", Type: "BUILT_IN"},
				{Name: "sqlserver", Type: "BUILT_IN"},
			},
			wantName: "appuser",
			wantType: "BUILT_IN",
		},
		{
			name:            "MySQL BUILT_IN user matches username and host",
			databaseVersion: "MYSQL_8_0",
			userName:        "admin",
			host:            "gmail.com",
			users: []*sqladmin.User{
				{Name: "admin", Host: "google.com", Type: "BUILT_IN"},
				{Name: "admin", Host: "gmail.com", Type: "BUILT_IN"},
			},
			wantName: "admin",
			wantType: "BUILT_IN",
		},
		{
			name:            "MySQL CLOUD_IAM_USER truncates domain at @",
			databaseVersion: "MYSQL_8_0",
			userName:        "admin@hashicorptest.com",
			host:            "",
			users: []*sqladmin.User{
				{Name: "admin", Host: "%", Type: "CLOUD_IAM_USER"},
			},
			wantName: "admin",
			wantType: "CLOUD_IAM_USER",
		},
		{
			name:            "MySQL CLOUD_IAM_SERVICE_ACCOUNT truncates domain at @",
			databaseVersion: "MYSQL_8_0",
			userName:        "tf-test-sa@tf-test-project.iam.gserviceaccount.com",
			host:            "",
			users: []*sqladmin.User{
				{Name: "tf-test-sa", Host: "%", Type: "CLOUD_IAM_SERVICE_ACCOUNT"},
			},
			wantName: "tf-test-sa",
			wantType: "CLOUD_IAM_SERVICE_ACCOUNT",
		},
		{
			name:            "MySQL CLOUD_IAM_GROUP preserves @ and lowercases domain",
			databaseVersion: "MYSQL_8_0",
			userName:        "iam-group-auth-test-group@GOOGLE.com",
			host:            "",
			users: []*sqladmin.User{
				{Name: "iam-group-auth-test-group@google.com", Host: "%", Type: "CLOUD_IAM_GROUP"},
			},
			wantName: "iam-group-auth-test-group@google.com",
			wantType: "CLOUD_IAM_GROUP",
		},
		{
			name:            "PostgreSQL CLOUD_IAM_USER preserves full email",
			databaseVersion: "POSTGRES_15",
			userName:        "admin@hashicorptest.com",
			host:            "",
			users: []*sqladmin.User{
				{Name: "postgres", Type: "BUILT_IN"},
				{Name: "admin@hashicorptest.com", Type: "CLOUD_IAM_USER"},
			},
			wantName: "admin@hashicorptest.com",
			wantType: "CLOUD_IAM_USER",
		},
		{
			name:            "PostgreSQL CLOUD_IAM_GROUP preserves exact name",
			databaseVersion: "POSTGRES_15",
			userName:        "iam-group-auth-test-group@GOOGLE.com",
			host:            "",
			users: []*sqladmin.User{
				{Name: "iam-group-auth-test-group@GOOGLE.com", Type: "CLOUD_IAM_GROUP"},
			},
			wantName: "iam-group-auth-test-group@GOOGLE.com",
			wantType: "CLOUD_IAM_GROUP",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := sql.FindSqlUser(tc.databaseVersion, tc.userName, tc.host, tc.users)
			if tc.wantName == "" {
				if got != nil {
					t.Fatalf("FindSqlUser() = %+v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatalf("FindSqlUser() = nil, want user %q (%s)", tc.wantName, tc.wantType)
			}
			if got.Name != tc.wantName || got.Type != tc.wantType {
				t.Fatalf("FindSqlUser() = {Name: %q, Type: %q}, want {Name: %q, Type: %q}", got.Name, got.Type, tc.wantName, tc.wantType)
			}
		})
	}
}

func TestSqlUser_entraidUserLifecycle(t *testing.T) {
	t.Parallel()

	const (
		project      = "tf-test-project"
		instanceName = "tf-test-sqlserver"
		entraUser    = "user@tenant.onmicrosoft.com"
	)

	testCases := []struct {
		name  string
		users []*sqladmin.User
	}{
		{
			name: "default built-in login plus ENTRAID_USER",
			users: []*sqladmin.User{
				{
					Name:     "sqlserver",
					Instance: instanceName,
					Type:     "BUILT_IN",
					SqlserverUserDetails: &sqladmin.SqlServerUserDetails{
						Disabled:    false,
						ServerRoles: []string{"CustomerDbRootRole"},
					},
				},
				{
					Name:     entraUser,
					Instance: instanceName,
					Type:     "ENTRAID_USER",
					SqlserverUserDetails: &sqladmin.SqlServerUserDetails{
						Disabled: false,
					},
				},
			},
		},
		{
			name: "same-prefix built-in login plus ENTRAID_USER",
			users: []*sqladmin.User{
				{
					Name:     "sqlserver",
					Instance: instanceName,
					Type:     "BUILT_IN",
					SqlserverUserDetails: &sqladmin.SqlServerUserDetails{
						Disabled:    false,
						ServerRoles: []string{"CustomerDbRootRole"},
					},
				},
				{
					Name:     "user",
					Instance: instanceName,
					Type:     "BUILT_IN",
					SqlserverUserDetails: &sqladmin.SqlServerUserDetails{
						Disabled:    false,
						ServerRoles: []string{"CustomerDbRootRole"},
					},
				},
				{
					Name:     entraUser,
					Instance: instanceName,
					Type:     "ENTRAID_USER",
					SqlserverUserDetails: &sqladmin.SqlServerUserDetails{
						Disabled: false,
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var deletedName string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				switch {
				case r.Method == http.MethodPost && r.URL.Path == fmt.Sprintf("/sql/v1beta4/projects/%s/instances/%s/users", project, instanceName):
					_ = json.NewEncoder(w).Encode(&sqladmin.Operation{
						Name:   "op-insert",
						Status: "DONE",
					})
				case r.Method == http.MethodGet && r.URL.Path == fmt.Sprintf("/sql/v1beta4/projects/%s/operations/op-insert", project):
					_ = json.NewEncoder(w).Encode(&sqladmin.Operation{
						Name:   "op-insert",
						Status: "DONE",
					})
				case r.Method == http.MethodGet && r.URL.Path == fmt.Sprintf("/sql/v1beta4/projects/%s/operations/op-delete", project):
					_ = json.NewEncoder(w).Encode(&sqladmin.Operation{
						Name:   "op-delete",
						Status: "DONE",
					})
				case r.Method == http.MethodGet && r.URL.Path == fmt.Sprintf("/sql/v1beta4/projects/%s/instances/%s", project, instanceName):
					_ = json.NewEncoder(w).Encode(&sqladmin.DatabaseInstance{
						Name:            instanceName,
						DatabaseVersion: "SQLSERVER_2022_STANDARD",
						Settings: &sqladmin.Settings{
							ActivationPolicy: "ALWAYS",
						},
					})
				case r.Method == http.MethodGet && r.URL.Path == fmt.Sprintf("/sql/v1beta4/projects/%s/instances/%s/users", project, instanceName):
					_ = json.NewEncoder(w).Encode(&sqladmin.UsersListResponse{
						Items: tc.users,
					})
				case r.Method == http.MethodDelete && r.URL.Path == fmt.Sprintf("/sql/v1beta4/projects/%s/instances/%s/users", project, instanceName):
					deletedName = r.URL.Query().Get("name")
					_ = json.NewEncoder(w).Encode(&sqladmin.Operation{
						Name:   "op-delete",
						Status: "DONE",
					})
				default:
					http.NotFound(w, r)
				}
			}))
			defer server.Close()

			cfg := &transport_tpg.Config{
				Context:      context.Background(),
				Client:       server.Client(),
				Project:      project,
				PollInterval: time.Millisecond,
				CustomEndpoints: map[string]string{
					"sql_custom_endpoint": server.URL + "/sql/v1beta4/",
				},
			}

			res := sql.ResourceSqlUser()
			d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
				"instance": instanceName,
				"name":     entraUser,
				"type":     "ENTRAID_USER",
			})

			if err := res.Create(d, cfg); err != nil {
				t.Fatalf("Create failed: %v", err)
			}

			wantID := fmt.Sprintf("%s//%s", entraUser, instanceName)
			if d.Id() != wantID {
				t.Fatalf("Create/Read produced ID %q, want %q", d.Id(), wantID)
			}
			if got := d.Get("name").(string); got != entraUser {
				t.Fatalf("Read produced name %q, want %q", got, entraUser)
			}
			if got := d.Get("type").(string); got != "ENTRAID_USER" {
				t.Fatalf("Read produced type %q, want %q", got, "ENTRAID_USER")
			}

			// Verify import by project/instance/name preserves full email UPN.
			importData := res.Data(nil)
			importData.SetId(fmt.Sprintf("%s/%s/%s", project, instanceName, entraUser))
			imported, err := res.Importer.State(importData, cfg)
			if err != nil {
				t.Fatalf("Importer failed: %v", err)
			}
			if err := res.Read(imported[0], cfg); err != nil {
				t.Fatalf("Post-import Read failed: %v", err)
			}
			if got := imported[0].Get("type").(string); got != "ENTRAID_USER" {
				t.Fatalf("Post-import Read produced type %q, want %q", got, "ENTRAID_USER")
			}

			if err := res.Delete(d, cfg); err != nil {
				t.Fatalf("Delete failed: %v", err)
			}
			if deletedName != entraUser {
				t.Fatalf("Delete sent query param name=%q, want %q", deletedName, entraUser)
			}
		})
	}
}
