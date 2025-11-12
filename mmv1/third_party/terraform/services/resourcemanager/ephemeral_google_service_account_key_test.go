package resourcemanager_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccEphemeralServiceAccountKey_create(t *testing.T) {
	t.Parallel()

	accountID := "a" + acctest.RandString(t, 10)
	displayName := "Terraform Test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEphemeralServiceAccountKey_create_setup(accountID, displayName),
			},
			{
				Config: testAccEphemeralServiceAccountKey_create(accountID, displayName),
			},
		},
	})
}

func testAccEphemeralServiceAccountKey_create_setup(accountID, displayName string) string {
	return fmt.Sprintf(`
resource "google_service_account" "service_account" {
	account_id   = "%s"
	display_name = "%s"
}

`, accountID, displayName)
}

func testAccEphemeralServiceAccountKey_create(accountID, displayName string) string {
	return fmt.Sprintf(`
resource "google_service_account" "service_account" {
	account_id   = "%s"
	display_name = "%s"
}

ephemeral "google_service_account_key" "key" {
  service_account_id = google_service_account.service_account.email
  public_key_type = "TYPE_X509_PEM_FILE"
}
`, accountID, displayName)
}

func TestAccEphemeralServiceAccountKey_upload(t *testing.T) {
	t.Parallel()

	accountID := "b" + acctest.RandString(t, 10)
	displayName := "Terraform Test Two"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccEphemeralServiceAccountKey_upload_setup(accountID, displayName),
			},
			{
				Config: testAccEphemeralServiceAccountKey_upload(accountID, displayName),
			},
		},
	})
}

func testAccEphemeralServiceAccountKey_upload_setup(accountID, displayName string) string {
	return fmt.Sprintf(`
resource "google_service_account" "service_account" {
	account_id   = "%s"
	display_name = "%s"
}
resource "time_sleep" "wait_30_seconds" {
  create_duration = "30s"
  depends_on = [google_service_account.service_account]
}
`, accountID, displayName)
}

func testAccEphemeralServiceAccountKey_upload(accountID, displayName string) string {
	return fmt.Sprintf(`
resource "google_service_account" "service_account" {
	account_id   = "%s"
	display_name = "%s"
}
resource "time_sleep" "wait_30_seconds" {
  create_duration = "30s"
  depends_on = [google_service_account.service_account]
}

ephemeral "google_service_account_key" "key" {
  service_account_id = google_service_account.service_account.email
  public_key_data    = filebase64("test-fixtures/public_key.pem")
}
`, accountID, displayName)
}

func TestAccEphemeralServiceAccountKey_fetch(t *testing.T) {
	t.Parallel()

	serviceAccount := envvar.GetTestServiceAccountFromEnv(t)
	targetServiceAccountEmail := acctest.BootstrapServiceAccount(t, "key-basic", serviceAccount)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEphemeralServiceAccountKey_fetch_setup(targetServiceAccountEmail),
			},
			{
				Config: testAccEphemeralServiceAccountKey_fetch(targetServiceAccountEmail),
			},
		},
	})
}

func testAccEphemeralServiceAccountKey_fetch_setup(accountID string) string {
	return fmt.Sprintf(`
resource "google_service_account_key" "acceptance" {
  service_account_id = "%s"
  public_key_type    = "TYPE_X509_PEM_FILE"
}
`, accountID)
}

func testAccEphemeralServiceAccountKey_fetch(accountID string) string {
	return fmt.Sprintf(`
resource "google_service_account_key" "acceptance" {
  service_account_id = "%s"
  public_key_type    = "TYPE_X509_PEM_FILE"
}

ephemeral "google_service_account_key" "key" {
  name = google_service_account_key.acceptance.name
  fetch_key = true
}
`, accountID)
}
