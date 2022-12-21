package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceGoogleBeyondcorpAppConnection_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBeyondcorpAppConnectionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleBeyondcorpAppConnection_basic(context),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceState("data.google_beyondcorp_app_connection.foo", "google_beyondcorp_app_connection.foo"),
				),
			},
		},
	})
}

func TestAccDataSourceGoogleBeyondcorpAppConnection_optionalProject(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBeyondcorpAppConnectionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleBeyondcorpAppConnection_optionalProject(context),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceState("data.google_beyondcorp_app_connection.foo", "google_beyondcorp_app_connection.foo"),
				),
			},
		},
	})
}

func TestAccDataSourceGoogleBeyondcorpAppConnection_optionalRegion(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBeyondcorpAppConnectionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleBeyondcorpAppConnection_optionalRegion(context),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceState("data.google_beyondcorp_app_connection.foo", "google_beyondcorp_app_connection.foo"),
				),
			},
		},
	})
}

func TestAccDataSourceGoogleBeyondcorpAppConnection_optionalProjectRegion(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBeyondcorpAppConnectionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleBeyondcorpAppConnection_optionalProjectRegion(context),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceState("data.google_beyondcorp_app_connection.foo", "google_beyondcorp_app_connection.foo"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleBeyondcorpAppConnection_basic(context map[string]interface{}) string {
	return Nprintf(`
resource "google_service_account" "service_account" {
	account_id   = "tf-test-my-account%{random_suffix}"
	display_name = "Test Service Account"
}

resource "google_beyondcorp_app_connector" "app_connector" {
	name = "tf-test-appconnector-%{random_suffix}"
	principal_info {
		service_account {
			email = google_service_account.service_account.email
		}
	}
}

resource "google_beyondcorp_app_connection" "foo" {
	name = "tf-test-my-app-connection-%{random_suffix}"
	type = "TCP_PROXY"
	application_endpoint {
		host = "foo-host"
		port = 8080
	}
	connectors = [google_beyondcorp_app_connector.app_connector.id]
}

data "google_beyondcorp_app_connection" "foo" {
	name    = google_beyondcorp_app_connection.foo.name
	project = google_beyondcorp_app_connection.foo.project
	region  = google_beyondcorp_app_connection.foo.region
}
`, context)
}

func testAccDataSourceGoogleBeyondcorpAppConnection_optionalProject(context map[string]interface{}) string {
	return Nprintf(`
resource "google_service_account" "service_account" {
	account_id   = "tf-test-my-account%{random_suffix}"
	display_name = "Test Service Account"
}

resource "google_beyondcorp_app_connector" "app_connector" {
	name = "tf-test-appconnector-%{random_suffix}"
	principal_info {
		service_account {
			email = google_service_account.service_account.email
		}
	}
}

resource "google_beyondcorp_app_connection" "foo" {
	name = "tf-test-my-app-connection-%{random_suffix}"
	type = "TCP_PROXY"
	application_endpoint {
		host = "foo-host"
		port = 8080
	}
	connectors = [google_beyondcorp_app_connector.app_connector.id]
}

data "google_beyondcorp_app_connection" "foo" {
	name   = google_beyondcorp_app_connection.foo.name
	region = google_beyondcorp_app_connection.foo.region
}
`, context)
}

func testAccDataSourceGoogleBeyondcorpAppConnection_optionalRegion(context map[string]interface{}) string {
	return Nprintf(`
resource "google_service_account" "service_account" {
	account_id   = "tf-test-my-account%{random_suffix}"
	display_name = "Test Service Account"
}

resource "google_beyondcorp_app_connector" "app_connector" {
	name = "tf-test-appconnector-%{random_suffix}"
	principal_info {
		service_account {
			email = google_service_account.service_account.email
		}
	}
}

resource "google_beyondcorp_app_connection" "foo" {
	name = "tf-test-my-app-connection-%{random_suffix}"
	type = "TCP_PROXY"
	application_endpoint {
		host = "foo-host"
		port = 8080
	}
	connectors = [google_beyondcorp_app_connector.app_connector.id]
}

data "google_beyondcorp_app_connection" "foo" {
	name    = google_beyondcorp_app_connection.foo.name
	project = google_beyondcorp_app_connection.foo.project
}
`, context)
}

func testAccDataSourceGoogleBeyondcorpAppConnection_optionalProjectRegion(context map[string]interface{}) string {
	return Nprintf(`
resource "google_service_account" "service_account" {
	account_id   = "tf-test-my-account%{random_suffix}"
	display_name = "Test Service Account"
}

resource "google_beyondcorp_app_connector" "app_connector" {
	name = "tf-test-appconnector-%{random_suffix}"
	principal_info {
		service_account {
			email = google_service_account.service_account.email
		}
	}
}

resource "google_beyondcorp_app_connection" "foo" {
	name = "tf-test-my-app-connection-%{random_suffix}"
	type = "TCP_PROXY"
	application_endpoint {
		host = "foo-host"
		port = 8080
	}
	connectors = [google_beyondcorp_app_connector.app_connector.id]
}

data "google_beyondcorp_app_connection" "foo" {
	name = google_beyondcorp_app_connection.foo.name
}
`, context)
}
