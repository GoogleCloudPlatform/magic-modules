package resourcemanager_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

// Test that a service account resource can be created, updated, and destroyed
func TestAccServiceAccount_basic(t *testing.T) {
	t.Parallel()

	accountId := "a" + acctest.RandString(t, 10)
	uniqueId := ""
	displayName := "Terraform Test"
	displayName2 := "Terraform Test Update"
	desc := "test description"
	desc2 := ""
	project := envvar.GetTestProjectFromEnv()
	expectedEmail := fmt.Sprintf("%s@%s.iam.gserviceaccount.com", accountId, project)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			// The first step creates a basic service account
			{
				Config: testAccServiceAccountBasic(accountId, displayName, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_service_account.acceptance", "project", project),
					resource.TestCheckResourceAttr(
						"google_service_account.acceptance", "member", "serviceAccount:"+expectedEmail),
				),
			},
			{
				ResourceName:      "google_service_account.acceptance",
				ImportStateId:     fmt.Sprintf("projects/%s/serviceAccounts/%s", project, expectedEmail),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_service_account.acceptance",
				ImportStateId:     fmt.Sprintf("%s/%s", project, expectedEmail),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_service_account.acceptance",
				ImportStateId:     expectedEmail,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// The second step updates the service account
			{
				Config: testAccServiceAccountBasic(accountId, displayName2, desc2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_service_account.acceptance", "project", project),
					testAccStoreServiceAccountUniqueId(&uniqueId),
				),
			},
			{
				ResourceName:      "google_service_account.acceptance",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// The third step explicitly adds the same default project to the service account configuration
			// and ensure the service account is not recreated by comparing the value of its unique_id with the one from the previous step
			{
				Config: testAccServiceAccountWithProject(project, accountId, displayName2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_service_account.acceptance", "project", project),
					resource.TestCheckResourceAttrPtr(
						"google_service_account.acceptance", "unique_id", &uniqueId),
				),
			},
			{
				ResourceName:      "google_service_account.acceptance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test the option to ignore ALREADY_EXISTS error from creating a service account.
func TestAccServiceAccount_createIgnoreAlreadyExists(t *testing.T) {
	t.Parallel()

	accountId := "a" + acctest.RandString(t, 10)
	displayName := "Terraform Test"
	desc := "test description"
	project := envvar.GetTestProjectFromEnv()
	expectedEmail := fmt.Sprintf("%s@%s.iam.gserviceaccount.com", accountId, project)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			// The first step creates a basic service account
			{
				Config: testAccServiceAccountBasic(accountId, displayName, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_service_account.acceptance", "project", project),
					resource.TestCheckResourceAttr(
						"google_service_account.acceptance", "member", "serviceAccount:"+expectedEmail),
				),
			},
			{
				ResourceName:      "google_service_account.acceptance",
				ImportStateId:     fmt.Sprintf("projects/%s/serviceAccounts/%s", project, expectedEmail),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// The second step creates a new resource that duplicates with the existing service account.
			{
				Config: testAccServiceAccountCreateIgnoreAlreadyExists(accountId, displayName, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_service_account.duplicate", "member", "serviceAccount:"+expectedEmail),
				),
			},
		},
	})
}

// Test setting create_ignore_already_exists on an existing resource
func TestAccServiceAccount_existingResourceCreateIgnoreAlreadyExists(t *testing.T) {
	t.Parallel()

	project := envvar.GetTestProjectFromEnv()

	serviceAccount := map[string]interface{}{
		"account_id":   "a" + acctest.RandString(t, 10),
		"display_name": "Terraform Test",
		"description":  "test description",
	}

	serviceAccountConfig := newServiceAccountResource(serviceAccount)

	serviceAccount["create_ignore_already_exists"] = true
	serviceAccountConfigIgnoreAlreadyExists := newServiceAccountResource(serviceAccount)

	expectedEmail := fmt.Sprintf("%s@%s.iam.gserviceaccount.com", serviceAccount["account_id"], project)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			// The first step creates a new resource with create_ignore_already_exists omitted
			{
				Config: serviceAccountConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_service_account.acceptance", "project", project),
					resource.TestCheckResourceAttr(
						"google_service_account.acceptance", "member", "serviceAccount:"+expectedEmail),
				),
			},
			{
				ResourceName:      "google_service_account.acceptance",
				ImportStateId:     fmt.Sprintf("projects/%s/serviceAccounts/%s", project, expectedEmail),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// The second step updates the resource to have create_ignore_already_exists=true
			{
				Config: serviceAccountConfigIgnoreAlreadyExists,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(
					"google_service_account.acceptance", "project", project),
					resource.TestCheckResourceAttr(
						"google_service_account.acceptance", "member", "serviceAccount:"+expectedEmail),
				),
			},
		},
	})
}

func TestAccServiceAccount_Disabled(t *testing.T) {
	t.Parallel()

	accountId := "a" + acctest.RandString(t, 10)
	uniqueId := ""
	displayName := "Terraform Test"
	desc := "test description"
	project := envvar.GetTestProjectFromEnv()
	expectedEmail := fmt.Sprintf("%s@%s.iam.gserviceaccount.com", accountId, project)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			// The first step creates a basic service account
			{
				Config: testAccServiceAccountBasic(accountId, displayName, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_service_account.acceptance", "project", project),
					resource.TestCheckResourceAttr(
						"google_service_account.acceptance", "member", "serviceAccount:"+expectedEmail),
				),
			},
			{
				ResourceName:      "google_service_account.acceptance",
				ImportStateId:     fmt.Sprintf("projects/%s/serviceAccounts/%s", project, expectedEmail),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// The second step disables the service account
			{
				Config: testAccServiceAccountDisabled(accountId, displayName, desc, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_service_account.acceptance", "project", project),
					testAccStoreServiceAccountUniqueId(&uniqueId),
				),
			},
			{
				ResourceName:      "google_service_account.acceptance",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// The third step enables the disabled service account
			{
				Config: testAccServiceAccountDisabled(accountId, displayName, desc, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_service_account.acceptance", "project", project),
					testAccStoreServiceAccountUniqueId(&uniqueId),
				),
			},
			{
				ResourceName:      "google_service_account.acceptance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccStoreServiceAccountUniqueId(uniqueId *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		*uniqueId = s.RootModule().Resources["google_service_account.acceptance"].Primary.Attributes["unique_id"]
		return nil
	}
}

func newServiceAccountResource(options map[string]interface{}) string {
	var sb strings.Builder
	sb.WriteString(`resource "google_service_account" "acceptance" {`)
	sb.WriteString("\n")
	for key, value := range options {
		sb.WriteString(fmt.Sprintf("\t%v\t= \"%v\"\n", key, value))
	}
	sb.WriteString(`}`)
	return sb.String()
}

func testAccServiceAccountBasic(account, name, desc string) string {
	return fmt.Sprintf(`
resource "google_service_account" "acceptance" {
  account_id   = "%v"
  display_name = "%v"
  description  = "%v"
}
`, account, name, desc)
}

func testAccServiceAccountCreateIgnoreAlreadyExists(account, name, desc string) string {
	return fmt.Sprintf(`
resource "google_service_account" "acceptance" {
  account_id   = "%v"
  display_name = "%v"
  description  = "%v"
}
resource "google_service_account" "duplicate" {
  account_id   = "%v"
  display_name = "%v"
  description  = "%v"
  create_ignore_already_exists = true
}
`, account, name, desc, account, name, desc)
}

func testAccServiceAccountWithProject(project, account, name string) string {
	return fmt.Sprintf(`
resource "google_service_account" "acceptance" {
  project      = "%v"
  account_id   = "%v"
  display_name = "%v"
  description  = "foo"
}
`, project, account, name)
}

func testAccServiceAccountDisabled(account, name, desc string, disabled bool) string {
	return fmt.Sprintf(`
resource "google_service_account" "acceptance" {
  account_id   = "%v"
  display_name = "%v"
  description  = "%v"
  disabled      = "%t"
}
`, account, name, desc, disabled)
}
