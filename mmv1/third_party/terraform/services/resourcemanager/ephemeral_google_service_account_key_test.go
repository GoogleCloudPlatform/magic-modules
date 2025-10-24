package resourcemanager_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccEphemeralServiceAccountKey_basic(t *testing.T) {
	t.Parallel()

	serviceAccount := envvar.GetTestServiceAccountFromEnv(t)
	targetServiceAccountEmail := acctest.BootstrapServiceAccount(t, "key-basic", serviceAccount)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEphemeralServiceAccountKey_setup(targetServiceAccountEmail),
			},
			{
				Config: testAccEphemeralServiceAccountKey_basic(targetServiceAccountEmail),
			},
		},
	})
}

func TestAccEphemeralServiceAccountKey_publicKey(t *testing.T) {
	t.Parallel()

	echoResourceName := "echo.test_public_key"
	accountID := "a" + acctest.RandString(t, 10)
	displayName := "Terraform Test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEphemeralServiceAccountKey_publicKey(accountID, displayName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(echoResourceName, "data"),
				),
			},
		},
	})
}

func testAccEphemeralServiceAccountKey_publicKey(account, name string) string {
	return fmt.Sprintf(`
resource "google_service_account" "acceptance" {
  account_id   = "%s"
  display_name = "%s"
}

ephemeral "google_service_account_key" "key" {
  service_account_id            = google_service_account.acceptance.email
  public_key    = filebase64("test-fixtures/public_key.pem")
}
`, account, name)
}

func TestAccEphemeralServiceAccountKey_privateKey(t *testing.T) {
	t.Parallel()
	echoResourceName := "echo.test_private_key"
	accountID := "a" + acctest.RandString(t, 10)
	displayName := "Terraform Test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEphemeralServiceAccountKey_privateKey(accountID, displayName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(echoResourceName, "data"),
				),
			},
		},
	})
}

func testAccEphemeralServiceAccountKey_privateKey(account, name string) string {
	return fmt.Sprintf(`
resource "google_service_account" "acceptance" {
  account_id   = "%s"
  display_name = "%s"
}

ephemeral "google_service_account_key" "key" {
  service_account_id            = google_service_account.acceptance.email
  private_key_type = "TYPE_GOOGLE_CREDENTIALS_FILE"
  key_algorithm = "KEY_ALG_RSA_2048"
}
`, account, name)
}

func testAccEphemeralServiceAccountKey_setup(serviceAccount string) string {
	return fmt.Sprintf(`
resource "google_service_account_key" "key" {
  service_account_id = "%s"
  public_key_type = "TYPE_X509_PEM_FILE"
}
`, serviceAccount)
}

func testAccEphemeralServiceAccountKey_basic(serviceAccount string) string {
	return fmt.Sprintf(`
resource "google_service_account_key" "key" {
  service_account_id = "%s"
  public_key_type = "TYPE_X509_PEM_FILE"
}

ephemeral "google_service_account_key" "key" {
  name            = google_service_account_key.key.name
  public_key_type = "TYPE_X509_PEM_FILE"
}
`, serviceAccount)
}
