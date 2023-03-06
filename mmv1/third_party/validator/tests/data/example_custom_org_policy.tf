terraform {
  required_providers {
    google = {
      source = "hashicorp/google"
      version = "~> {{.Provider.version}}"
    }
  }
}

provider "google" {
  {{if .Provider.credentials }}credentials = "{{.Provider.credentials}}"{{end}}
}

resource "google_org_policy_policy" "project_policy" {
  name   = "projects/{{.Provider.project}}/policies/gcp.resourceLocations"
  parent = "projects/{{.Provider.project}}"
  
  spec {
    rules {
      condition {
        description = "Description the policy"
        expression  = "resource.matchLabels('label1', 'label2')"
        location    = "EU"
        title       = "Title of the condition"
      }

      values {
        allowed_values = ["projects/123","projects/456"]
        denied_values  = ["projects/789"]
      }
    }

    rules {
      allow_all = "TRUE"
    }
  }
}
