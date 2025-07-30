resource "google_network_security_server_tls_policy" "lb_mtls_policy" {
  name                   = "lb_mtls_policy"
  description            = "my description"
  location               = "global"
  labels                 = {
    foo = "bar"
  }
  mtls_policy {
    client_validation_mode         = "REJECT_INVALID"
    client_validation_trust_config = "projects/307841421122/locations/global/trustConfigs/id-4adf7779-1e9f-4124-9438-652c80886074"
  }
}

resource "google_network_security_server_tls_policy" "td_mtls_policy" {
  name                   = "td_mtls_policy"
  description            = "my description"
  location               = "global"
  labels                 = {
    foo = "bar"
  }
  allow_open = "false"
  server_certificate {
    certificate_provider_instance {
      plugin_instance = "google_cloud_private_spiffe"
    }
  }
  mtls_policy {
    client_validation_ca {
      certificate_provider_instance {
        plugin_instance = "google_cloud_private_spiffe" 
      }
    }
  }
}

resource "google_network_security_server_tls_policy" "td_with_server_cert_policy" {
  name                   = "td_with_server_cert_policy"
  description            = "my description"
  location               = "global"
  allow_open             = "false"
  server_certificate {
    grpc_endpoint {
        target_uri = "unix:mypath"
      }
  }
}

resource "google_network_security_server_tls_policy" "empty_description_policy" {
  name                   = "empty_description_policy"
  location               = "global"
  labels                 = {
    foo = "bar"
  }
  mtls_policy {
    client_validation_mode         = "REJECT_INVALID"
    client_validation_trust_config = "projects/307841421122/locations/global/trustConfigs/id-4adf7779-1e9f-4124-9438-652c80886074"
  }
}

resource "google_network_security_server_tls_policy" "empty_labels_policy" {
  name                   = "empty_labels_policy"
  description            = "my description"
  location               = "global"
  mtls_policy {
    client_validation_mode         = "REJECT_INVALID"
    client_validation_trust_config = "projects/307841421122/locations/global/trustConfigs/id-4adf7779-1e9f-4124-9438-652c80886074"
  }
}

resource "google_network_security_server_tls_policy" "regional_location_policy" {
  name                   = "regional_location_policy"
  description            = "my description"
  location               = "us-central1"
  labels                 = {
    foo = "bar"
  }
  mtls_policy {
    client_validation_mode         = "REJECT_INVALID"
    client_validation_trust_config = "projects/307841421122/locations/us-central1/trustConfigs/tsmx-20250609-tc1"
  }
}

resource "google_network_security_server_tls_policy" "lb_mtls_allow_invalid_cert_policy" {
  name                   = "lb_mtls_allow_invalid_cert_policy"
  labels                 = {
    foo = "bar"
  }
  description            = "my description"
  location               = "global"
  mtls_policy {
    client_validation_mode = "ALLOW_INVALID_OR_MISSING_CLIENT_CERT"
  }
}

resource "google_network_security_server_tls_policy" "td_allow_open_policy" {
  name                   = "td_allow_open_policy"
  description            = "my description"
  location               = "global"
  allow_open             = "true"
  server_certificate {
    grpc_endpoint {
        target_uri = "unix:mypath"
      }
  }
  mtls_policy {
    client_validation_ca {
      certificate_provider_instance {
        plugin_instance = "google_cloud_private_spiffe" 
      }
    }
  }
}

resource "google_network_security_server_tls_policy" "td_with_cert_provider_policy" {
  name                   = "td_with_cert_provider_policy"
  description            = "my description"
  location               = "global"
  server_certificate {
    certificate_provider_instance {
      plugin_instance = "google_cloud_private_spiffe"
    }
  }
}

resource "google_network_security_server_tls_policy" "td_mtls_client_validation_grpc_policy" {
  name                   = "td_mtls_client_validation_grpc_policy"
  description            = "my description"
  location               = "global"
  labels                 = {
    foo = "bar"
  }
  allow_open = "false"
  server_certificate {
    certificate_provider_instance {
      plugin_instance = "google_cloud_private_spiffe"
    }
  }
  mtls_policy {
    client_validation_ca {
      grpc_endpoint {
        target_uri = "unix:mypath"
      }
    }
  }
}
