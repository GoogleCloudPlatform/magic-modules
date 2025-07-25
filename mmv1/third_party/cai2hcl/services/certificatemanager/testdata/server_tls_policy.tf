################################################################################
# Core Usecases
################################################################################

# Load Balancer mTLS
resource "google_network_security_server_tls_policy" "test1" {
  name                   = "test1"
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

# Traffic Director mTLS
resource "google_network_security_server_tls_policy" "test2" {
  name                   = "test2"
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

# Traffic Director with server certificate 
resource "google_network_security_server_tls_policy" "test3" {
  name                   = "test3"
  description            = "my description"
  location               = "global"
  allow_open             = "false"
  server_certificate {
    grpc_endpoint {
        target_uri = "unix:mypath"
      }
  }
}

################################################################################
# Enumerating other possible configurations
################################################################################

# Empty description
resource "google_network_security_server_tls_policy" "test4" {
  name                   = "test4"
  location               = "global"
  labels                 = {
    foo = "bar"
  }
  mtls_policy {
    client_validation_mode         = "REJECT_INVALID"
    client_validation_trust_config = "projects/307841421122/locations/global/trustConfigs/id-4adf7779-1e9f-4124-9438-652c80886074"
  }
}

# Empty labels
resource "google_network_security_server_tls_policy" "test5" {
  name                   = "test5"
  description            = "my description"
  location               = "global"
  mtls_policy {
    client_validation_mode         = "REJECT_INVALID"
    client_validation_trust_config = "projects/307841421122/locations/global/trustConfigs/id-4adf7779-1e9f-4124-9438-652c80886074"
  }
}

# Regional location
resource "google_network_security_server_tls_policy" "test6" {
  name                   = "test6"
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

# Load Balancer mTLS but allowing invalid or missing client certificates
resource "google_network_security_server_tls_policy" "test7" {
  name                   = "test7"
  labels                 = {
    foo = "bar"
  }
  description            = "my description"
  location               = "global"
  mtls_policy {
    client_validation_mode = "ALLOW_INVALID_OR_MISSING_CLIENT_CERT"
  }
}

# Traffic Director with allow_open true
resource "google_network_security_server_tls_policy" "test8" {
  name                   = "test8"
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

# Traffic Director with certificate provider instance
resource "google_network_security_server_tls_policy" "test9" {
  name                   = "test9"
  description            = "my description"
  location               = "global"
  server_certificate {
    certificate_provider_instance {
      plugin_instance = "google_cloud_private_spiffe"
    }
  }
}

# Traffic Director mTLS with ClientValidation CA: gRPC endpoint
resource "google_network_security_server_tls_policy" "test10" {
  name                   = "test10"
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
