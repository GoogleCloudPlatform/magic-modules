resource "google_access_context_manager_service_perimeter" "storage-perimeter" {
  parent = "accesspolicies/${google_access_context_manager_access_policy.access-policy.name}"
  name   = "accesspolicies/${google_access_context_manager_access_policy.access-policy.name}/serviceperimeters/storage-perimeter"
  title  = "Storage Perimeter"
  status {
    restricted_services = ["storage.googleapis.com"]
  }
  lifecycle {
    ignore_changes = [status[0].resources]
  }
}

# Allow for anyone
resource "google_access_context_manager_service_perimeter_ingress_policy" "ingress_policy" {
  perimeter = "${google_access_context_manager_service_perimeter.storage-perimeter.name}"

  ingress_from {
    identity_type = "ANY_IDENTITY"
    sources {
      access_level = "*"
    }
  }

  ingress_to {
    resources = ["*"]
    operations {
      service_name = "bigquery.googleapis.com"
      method_selectors {
        method = "*"
      }
    }
  }
  lifecycle {
    create_before_destroy = true
  }
}

# Allow just from a specific VPC
resource "google_access_context_manager_service_perimeter_ingress_policy" "restricted_network" {
  perimeter = "${google_access_context_manager_service_perimeter.storage-perimeter.name}"

  # Allow ingress from a specific VPC in another project
  ingress_from {
    sources {
      resource = "//compute.googleapis.com/projects/87654321/global/networks/network-in-another-project"
    }
  }

  ingress_to {
    resources = ["*"]
    operations {
      service_name = "bigquery.googleapis.com"
      method_selectors {
        method = "*"
      }
    }
  }
  lifecycle {
    create_before_destroy = true
  }
}

resource "google_access_context_manager_access_policy" "access-policy" {
  parent = "organizations/123456789"
  title  = "Storage Policy"
}
