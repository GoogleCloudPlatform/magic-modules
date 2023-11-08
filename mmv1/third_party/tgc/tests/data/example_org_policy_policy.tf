terraform {
  required_providers {
    google = {
      source = "hashicorp/google-beta"
      version = "~> {{.Provider.version}}"
    }
  }
}

provider "google" {
  {{if .Provider.credentials }}credentials = "{{.Provider.credentials}}"{{end}}
}

resource "google_org_policy_policy" "folder_policy" {
  name   = "folders/{{.FolderID}}/policies/samplePolicy"
  parent = "folders/{{.FolderID}}"
  spec {
    rules {
      deny_all = "TRUE"
    }
    inherit_from_parent = true
  }
}

resource "google_org_policy_policy" "organizationPolicy" {
  name   = "organizations/{{.OrgID}}/policies/gcp.detailedAuditLoggingMode"
  parent = "organizations/{{.OrgID}}"
  spec {
    reset = true
  }
}

resource "google_org_policy_policy" "project_policy" {
  name   = "projects/{{.Provider.project}}/policies/gcp.resourceLocations"
  parent = "projects/{{.Provider.project}}"
  spec {
    rules {
      condition {
        description = "A sample condition for the policy"
        expression  = "resource.matchLabels('labelKeys/123', 'labelValues/345')"
        location    = "sample-location.log"
        title       = "sample-condition"
      }


      values {
        allowed_values = ["projects/allowed-project1", "projects/allowed-project2"]
        denied_values  = ["projects/denied-project"]
      }
    }

    rules {
      allow_all = "TRUE"
    }

    inherit_from_parent = true
  }
}
