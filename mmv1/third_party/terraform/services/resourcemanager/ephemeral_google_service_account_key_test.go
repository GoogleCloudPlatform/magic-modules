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

	serviceAccount := envvar.GetTestServiceAccountFromEnv(t)
	targetServiceAccountEmail := acctest.BootstrapServiceAccount(t, "key-basic", serviceAccount)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEphemeralServiceAccountKey_create(targetServiceAccountEmail),
			},
		},
	})
}

func testAccEphemeralServiceAccountKey_create(serviceAccount string) string {
	return fmt.Sprintf(`
ephemeral "google_service_account_key" "key" {
  service_account_id = "%s"
  public_key_type = "TYPE_X509_PEM_FILE"
}
`, serviceAccount)
}

func TestAccEphemeralServiceAccountKey_upload(t *testing.T) {
	t.Parallel()

	accountID := "a" + acctest.RandString(t, 10)
	displayName := "Terraform Test"
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

func testAccEphemeralServiceAccountKey_upload_setup(accountID, name string) string {
	return fmt.Sprintf(`
resource "google_service_account" "acceptance" {
  account_id   = "%s"
  display_name = "%s"
}

resource "time_sleep" "wait_10_seconds" {
  create_duration = "10s"
  depends_on = [google_service_account.acceptance]
}
`, accountID, name)
}

func testAccEphemeralServiceAccountKey_upload(serviceAccount, name string) string {
	return fmt.Sprintf(`
resource "google_service_account" "acceptance" {
  account_id   = "%s"
  display_name = "%s"
}

resource "time_sleep" "wait_10_seconds" {
  create_duration = "10s"
  depends_on = [google_service_account.acceptance]
}

ephemeral "google_service_account_key" "key" {
  service_account_id = google_service_account.acceptance.email
  public_key_data    = filebase64("test-fixtures/public_key.pem")
}
`, serviceAccount, name)
}

func TestAccEphemeralServiceAccountKey_fetch(t *testing.T) {
	t.Parallel()

	serviceAccount := envvar.GetTestServiceAccountFromEnv(t)
	targetServiceAccountEmail := acctest.BootstrapServiceAccount(t, "key-basic", serviceAccount)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
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

resource "time_sleep" "wait_10_seconds" {
  create_duration = "10s"
  depends_on = [google_service_account_key.acceptance]
}
`, accountID)
}

func testAccEphemeralServiceAccountKey_fetch(accountID string) string {
	return fmt.Sprintf(`
resource "google_service_account_key" "acceptance" {
  service_account_id = "%s"
  public_key_type    = "TYPE_X509_PEM_FILE"
}

resource "time_sleep" "wait_10_seconds" {
  create_duration = "10s"
  depends_on = [google_service_account_key.acceptance]
}

ephemeral "google_service_account_key" "key" {
  name = google_service_account_key.acceptance.name
  fetch_key = true
}
`, accountID)
}
