package integrationconnectors_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccIntegrationConnectorsConnection_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckIntegrationConnectorsConnectionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIntegrationConnectorsConnection_full(context),
			},
			{
				ResourceName:            "google_integration_connectors_connection.zendeskconnection",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "labels", "terraform_labels"},
			},
			{
				Config: testAccIntegrationConnectorsConnection_update(context),
			},
			{
				ResourceName:            "google_integration_connectors_connection.zendeskconnection",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccIntegrationConnectorsConnection_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "test_project" {
}

resource "google_secret_manager_secret" "secret-basic" {
  secret_id = "tf-test-test-secret%{random_suffix}"
  replication {
    user_managed {
      replicas {
        location = "us-central1"
      }
    }
  }
}


resource "google_secret_manager_secret_version" "secret-version-basic" {
  secret = google_secret_manager_secret.secret-basic.id
  secret_data = "dummypassword"
}

resource "google_secret_manager_secret_iam_member" "secret_iam" {
  secret_id  = google_secret_manager_secret.secret-basic.id
  role       = "roles/secretmanager.admin"
  member     = "serviceAccount:${data.google_project.test_project.number}-compute@developer.gserviceaccount.com"
  depends_on = [google_secret_manager_secret_version.secret-version-basic]
}

resource "google_integration_connectors_connection" "zendeskconnection" {
  name     = "tf-test-test-zendesk%{random_suffix}"
  description = "tf updated description"
  location = "us-central1"
  service_account = "${data.google_project.test_project.number}-compute@developer.gserviceaccount.com"
  connector_version = "projects/connectors-example/locations/global/providers/zendesk/connectors/zendesk/versions/1"
  config_variable {
      key = "proxy_enabled"
      boolean_value = false
  }
  config_variable {
    key = "sample_integer_value"
    integer_value = 1
  }

  config_variable {
      key = "sample_encryption_key_value"
      encryption_key_value {
        type = "GOOGLE_MANAGED"
        kms_key_name = "sampleKMSKkey"
      }
  }

  config_variable {
    key = "sample_secret_value"
    secret_value {
      secret_version = google_secret_manager_secret_version.secret-version-basic.name
    }
  }

  suspended = false
  auth_config {
    auth_type = "USER_PASSWORD"
    user_password {
      username = "user@xyz.com"
      password {
        secret_version = google_secret_manager_secret_version.secret-version-basic.name
      }
    }
  }

  destination_config {
    key = "url"
    destination {
        host = "https://test.zendesk.com"
    }
  }
  lock_config {
    locked = false
    reason = "Its not locked"
  }
  node_config {
    min_node_count = 2
    max_node_count = 50
  }
  labels = {
    foo = "bar"
  }
  ssl_config {
    additional_variable {
      key = "sample_string"
      string_value = "sampleString"
    }
    additional_variable {
      key = "sample_boolean"
      boolean_value = false
    }
    additional_variable {
      key = "sample_integer"
      integer_value = 1
    }
    additional_variable {
      key = "sample_secret_value"
      secret_value {
        secret_version = google_secret_manager_secret_version.secret-version-basic.name
      }
    }
    additional_variable {
      key = "sample_encryption_key_value"
      encryption_key_value {
        type = "GOOGLE_MANAGED"
        kms_key_name = "sampleKMSKkey"
      }
    }
    client_cert_type = "PEM"
    client_certificate {
      secret_version = google_secret_manager_secret_version.secret-version-basic.name
    }
    client_private_key {
      secret_version = google_secret_manager_secret_version.secret-version-basic.name
    }
    client_private_key_pass {
      secret_version = google_secret_manager_secret_version.secret-version-basic.name
    }
    private_server_certificate {
      secret_version = google_secret_manager_secret_version.secret-version-basic.name
    }
    server_cert_type = "PEM"
    trust_model      = "PRIVATE"
    type             = "TLS"
    use_ssl          = true
  }

  eventing_enablement_type = "EVENTING_AND_CONNECTION"
  eventing_config {
    additional_variable {
      key = "sample_string"
      string_value = "sampleString"
    }
    additional_variable {
      key = "sample_boolean"
      boolean_value = false
    }
    additional_variable {
      key = "sample_integer"
      integer_value = 1
    }
    additional_variable {
      key = "sample_secret_value"
      secret_value {
        secret_version = google_secret_manager_secret_version.secret-version-basic.name
      }
    }
    additional_variable {
      key = "sample_encryption_key_value"
      encryption_key_value {
        type = "GOOGLE_MANAGED"
        kms_key_name = "sampleKMSKkey"
      }
    }
    registration_destination_config {
      key = "registration_destination_config"
      destinations {
          host = "https://test.zendesk.com"
        }
    }
    auth_config {
      auth_type = "USER_PASSWORD"
      user_password {
        username = "user@xyz.com"
        password {
          secret_version = google_secret_manager_secret_version.secret-version-basic.name
        }
      }
    }
    enrichment_enabled = true
  }
}
`, context)
}

func testAccIntegrationConnectorsConnection_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "test_project" {
}

resource "google_secret_manager_secret" "secret-basic" {
  secret_id = "tf-test-test-secret%{random_suffix}"
  replication {
    user_managed {
      replicas {
        location = "us-central1"
      }
    }
  }
}


resource "google_secret_manager_secret_version" "secret-version-basic" {
  secret = google_secret_manager_secret.secret-basic.id
  secret_data = "dummypassword"
}

resource "google_secret_manager_secret_iam_member" "secret_iam" {
  secret_id  = google_secret_manager_secret.secret-basic.id
  role       = "roles/secretmanager.admin"
  member     = "serviceAccount:${data.google_project.test_project.number}-compute@developer.gserviceaccount.com"
  depends_on = [google_secret_manager_secret_version.secret-version-basic]
}

resource "google_integration_connectors_connection" "zendeskconnection" {
  name     = "tf-test-test-zendesk%{random_suffix}"
  description = "tf updated description"
  location = "us-central1"
  service_account = "${data.google_project.test_project.number}-compute@developer.gserviceaccount.com"
  connector_version = "projects/connectors-example/locations/global/providers/zendesk/connectors/zendesk/versions/1"
  config_variable {
      key = "proxy_enabled"
      boolean_value = true
  }

  auth_config {
    auth_type = "USER_PASSWORD"
    user_password {
      username = "user1@xyz.com"
      password {
        secret_version = google_secret_manager_secret_version.secret-version-basic.name
      }
    }
  }

  destination_config {
    key = "url"
    destination {
        host = "https://test1.zendesk.com"
    }
  }
  lock_config {
    locked = false
    reason = "Its not locked"
  }
  log_config {
    enabled = true
  }
  node_config {
    min_node_count = 3
    max_node_count = 49
  }
  labels = {
    foo = "bar1"
  }
  ssl_config {

    client_cert_type = "PEM"
    client_certificate {
      secret_version = google_secret_manager_secret_version.secret-version-basic.name
    }
    client_private_key {
      secret_version = google_secret_manager_secret_version.secret-version-basic.name
    }
    client_private_key_pass {
      secret_version = google_secret_manager_secret_version.secret-version-basic.name
    }
    private_server_certificate {
      secret_version = google_secret_manager_secret_version.secret-version-basic.name
    }
    server_cert_type = "PEM"
    trust_model      = "PRIVATE"
    type             = "TLS"
    use_ssl          = false
  }

  eventing_enablement_type = "EVENTING_AND_CONNECTION"
  eventing_config {
    registration_destination_config {
      key = "registration_destination_config"
      destinations {
          host = "https://test1.zendesk.com"
        }
    }
    auth_config {
      auth_type = "USER_PASSWORD"
      user_password {
        username = "user1@xyz.com"
        password {
          secret_version = google_secret_manager_secret_version.secret-version-basic.name
        }
      }
    }
    enrichment_enabled = true
  }
}
`, context)
}
