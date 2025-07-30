resource "google_network_security_server_tls_policy" "lb_mtls_policy" {
  allow_open  = "false"
  description = "my description"

  labels = {
    foo = "bar"
  }

  location = "global"

  mtls_policy {
    client_validation_mode         = "REJECT_INVALID"
    client_validation_trust_config = "projects/307841421122/locations/global/trustConfigs/id-4adf7779-1e9f-4124-9438-652c80886074"
  }

  name = "lb_mtls_policy"
}

resource "google_network_security_server_tls_policy" "td_mtls_policy" {
  allow_open  = "false"
  description = "my description"

  labels = {
    foo = "bar"
  }

  location = "global"

  mtls_policy {
    client_validation_ca {
      certificate_provider_instance {
        plugin_instance = "google_cloud_private_spiffe"
      }
    }
  }

  name = "td_mtls_policy"

  server_certificate {
    certificate_provider_instance {
      plugin_instance = "google_cloud_private_spiffe"
    }
  }
}

resource "google_network_security_server_tls_policy" "td_with_server_cert_policy" {
  allow_open  = "false"
  description = "my description"
  location    = "global"
  name        = "td_with_server_cert_policy"

  server_certificate {
    grpc_endpoint {
      target_uri = "unix:mypath"
    }
  }
}

resource "google_network_security_server_tls_policy" "empty_description_policy" {
  allow_open = "false"

  labels = {
    foo = "bar"
  }

  location = "global"

  mtls_policy {
    client_validation_mode         = "REJECT_INVALID"
    client_validation_trust_config = "projects/307841421122/locations/global/trustConfigs/id-4adf7779-1e9f-4124-9438-652c80886074"
  }

  name = "empty_description_policy"
}

resource "google_network_security_server_tls_policy" "empty_labels_policy" {
  allow_open  = "false"
  description = "my description"
  location    = "global"

  mtls_policy {
    client_validation_mode         = "REJECT_INVALID"
    client_validation_trust_config = "projects/307841421122/locations/global/trustConfigs/id-4adf7779-1e9f-4124-9438-652c80886074"
  }

  name = "empty_labels_policy"
}

resource "google_network_security_server_tls_policy" "regional_location_policy" {
  allow_open  = "false"
  description = "my description"

  labels = {
    foo = "bar"
  }

  location = "us-central1"

  mtls_policy {
    client_validation_mode         = "REJECT_INVALID"
    client_validation_trust_config = "projects/307841421122/locations/us-central1/trustConfigs/tsmx-20250609-tc1"
  }

  name = "regional_location_policy"
}

resource "google_network_security_server_tls_policy" "lb_mtls_allow_invalid_cert_policy" {
  allow_open  = "false"
  description = "my description"

  labels = {
    foo = "bar"
  }

  location = "global"

  mtls_policy {
    client_validation_mode = "ALLOW_INVALID_OR_MISSING_CLIENT_CERT"
  }

  name = "lb_mtls_allow_invalid_cert_policy"
}

resource "google_network_security_server_tls_policy" "td_allow_open_policy" {
  allow_open  = "true"
  description = "my description"
  location    = "global"

  mtls_policy {
    client_validation_ca {
      certificate_provider_instance {
        plugin_instance = "google_cloud_private_spiffe"
      }
    }
  }

  name = "td_allow_open_policy"

  server_certificate {
    grpc_endpoint {
      target_uri = "unix:mypath"
    }
  }
}

resource "google_network_security_server_tls_policy" "td_with_cert_provider_policy" {
  allow_open  = "false"
  description = "my description"
  location    = "global"
  name        = "td_with_cert_provider_policy"

  server_certificate {
    certificate_provider_instance {
      plugin_instance = "google_cloud_private_spiffe"
    }
  }
}

resource "google_network_security_server_tls_policy" "td_mtls_client_validation_grpc_policy" {
  allow_open  = "false"
  description = "my description"

  labels = {
    foo = "bar"
  }

  location = "global"

  mtls_policy {
    client_validation_ca {
      grpc_endpoint {
        target_uri = "unix:mypath"
      }
    }
  }

  name = "td_mtls_client_validation_grpc_policy"

  server_certificate {
    certificate_provider_instance {
      plugin_instance = "google_cloud_private_spiffe"
    }
  }
}
