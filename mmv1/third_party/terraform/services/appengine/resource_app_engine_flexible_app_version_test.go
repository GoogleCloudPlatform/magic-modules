package appengine_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccAppEngineFlexibleAppVersion_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
		"random_suffix":   acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAppEngineFlexibleAppVersionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAppEngineFlexibleAppVersion_python(context),
			},
			{
				ResourceName:            "google_app_engine_flexible_app_version.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"env_variables", "deployment", "entrypoint", "service", "noop_on_destroy"},
			},
			{
				Config: testAccAppEngineFlexibleAppVersion_pythonUpdate(context),
			},
			{
				ResourceName:            "google_app_engine_flexible_app_version.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"env_variables", "deployment", "entrypoint", "service", "delete_service_on_destroy"},
			},
		},
	})
}

func testAccAppEngineFlexibleAppVersion_python(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "my_project" {
  name = "tf-test-appeng-flex%{random_suffix}"
  project_id = "tf-test-appeng-flex%{random_suffix}"
  org_id = "%{org_id}"
  billing_account = "%{billing_account}"
  deletion_policy = "NONE"
}

resource "google_project_service" "compute" {
  project = google_project.my_project.project_id
  service = "compute.googleapis.com"

  disable_dependent_services = false
}

resource "google_project_service" "appengineflex" {
  project = google_project.my_project.project_id
  service = "appengineflex.googleapis.com"

  disable_dependent_services = false
  depends_on = [google_project_service.compute]
}

resource "google_compute_network" "network" {
  project                 = google_project_service.compute.project
  name                    = "custom"
  auto_create_subnetworks = "false"
}

resource "google_compute_subnetwork" "subnetwork" {
  project                  = google_project_service.compute.project
  name                     = "custom"
  region                   = "us-central1"
  network                  = google_compute_network.network.id
  ip_cidr_range            = "10.0.0.0/16"
  private_ip_google_access = true
}

resource "google_app_engine_application" "app" {
  project     = google_project.my_project.project_id
  location_id = "us-central"
}

resource "google_project_iam_member" "gae_api" {
  project = google_project_service.appengineflex.project
  role    = "roles/compute.networkUser"
  member  = "serviceAccount:service-${google_project.my_project.number}@gae-api-prod.google.com.iam.gserviceaccount.com"
}

resource "google_app_engine_standard_app_version" "foo" {
  project    = google_project_iam_member.gae_api.project
  version_id = "v1"
  service    = "default"
  runtime    = "python38"

  entrypoint {
    shell = "gunicorn -b :$PORT main:app"
  }

  deployment {
    files {
      name = "main.py"
      source_url = "https://storage.googleapis.com/${google_storage_bucket.bucket.name}/${google_storage_bucket_object.main.name}"
    }

    files {
      name = "requirements.txt"
      source_url = "https://storage.googleapis.com/${google_storage_bucket.bucket.name}/${google_storage_bucket_object.requirements.name}"
    }
  }

  env_variables = {
    port = "8000"
  }

  noop_on_destroy = true
}

resource "google_app_engine_flexible_app_version" "foo" {
  project    = google_project_iam_member.gae_api.project
  version_id = "v1"
  service    = "custom"
  runtime    = "python"

  runtime_api_version = "1"

  resources {
    cpu       = 1
    memory_gb = 0.5
    disk_gb   = 10
  }

  entrypoint {
    shell = "gunicorn -b :$PORT main:app"
  }

  flexible_runtime_settings {
    operating_system = "ubuntu22"
    runtime_version = "3.11"
  }

  deployment {
    files {
      name = "main.py"
      source_url = "https://storage.googleapis.com/${google_storage_bucket.bucket.name}/${google_storage_bucket_object.main.name}"
    }

    files {
      name = "requirements.txt"
      source_url = "https://storage.googleapis.com/${google_storage_bucket.bucket.name}/${google_storage_bucket_object.requirements.name}"
    }

    files {
      name = "app.yaml"
      source_url = "https://storage.googleapis.com/${google_storage_bucket.bucket.name}/${google_storage_bucket_object.yaml.name}"
    }
  }

  liveness_check {
    path = "alive"
  }

  readiness_check {
    path = "ready"
  }

  env_variables = {
    port = "8000"
  }

  network {
    name             = google_compute_network.network.name
    subnetwork       = google_compute_subnetwork.subnetwork.name
    instance_ip_mode = "EXTERNAL"
  }

  instance_class = "B1"

  manual_scaling {
    instances = 1
  }

  noop_on_destroy = true

  depends_on = [google_app_engine_standard_app_version.foo]
}

resource "google_storage_bucket" "bucket" {
  project  = google_project.my_project.project_id
  name     = "tf-test-%{random_suffix}-flex-ae-bucket"
  location = "US"
}

resource "google_storage_bucket_object" "yaml" {
  name   = "app.yaml"
  bucket = google_storage_bucket.bucket.name
  source = "./test-fixtures/hello-world-flask/app.yaml"
}

resource "google_storage_bucket_object" "requirements" {
  name   = "requirements.txt"
  bucket = google_storage_bucket.bucket.name
  source = "./test-fixtures/hello-world-flask/requirements.txt"
}

resource "google_storage_bucket_object" "main" {
  name   = "main.py"
  bucket = google_storage_bucket.bucket.name
  source = "./test-fixtures/hello-world-flask/main.py"
}`, context)
}

func testAccAppEngineFlexibleAppVersion_pythonUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "my_project" {
  name = "tf-test-appeng-flex%{random_suffix}"
  project_id = "tf-test-appeng-flex%{random_suffix}"
  org_id = "%{org_id}"
  billing_account = "%{billing_account}"
  deletion_policy = "NONE"
}

resource "google_project_service" "compute" {
  project = google_project.my_project.project_id
  service = "compute.googleapis.com"

  disable_dependent_services = false
}

resource "google_project_service" "appengineflex" {
  project = google_project.my_project.project_id
  service = "appengineflex.googleapis.com"

  disable_dependent_services = false
  depends_on = [google_project_service.compute]
}

resource "google_compute_network" "network" {
  project                 = google_project_service.compute.project
  name                    = "custom"
  auto_create_subnetworks = "false"
}

resource "google_compute_subnetwork" "subnetwork" {
  project                  = google_project_service.compute.project
  name                     = "custom"
  region                   = "us-central1"
  network                  = google_compute_network.network.id
  ip_cidr_range            = "10.0.0.0/16"
  private_ip_google_access = true
}

resource "google_app_engine_application" "app" {
  project     = google_project.my_project.project_id
  location_id = "us-central"
}

resource "google_project_iam_member" "gae_api" {
  project = google_project_service.appengineflex.project
  role    = "roles/compute.networkUser"
  member  = "serviceAccount:service-${google_project.my_project.number}@gae-api-prod.google.com.iam.gserviceaccount.com"
}

resource "google_app_engine_standard_app_version" "foo" {
  project    = google_project_iam_member.gae_api.project
  version_id = "v1"
  service    = "default"
  runtime    = "python38"

  entrypoint {
    shell = "gunicorn -b :$PORT main:app"
  }

  deployment {
    files {
      name = "main.py"
      source_url = "https://storage.googleapis.com/${google_storage_bucket.bucket.name}/${google_storage_bucket_object.main.name}"
    }

    files {
      name = "requirements.txt"
      source_url = "https://storage.googleapis.com/${google_storage_bucket.bucket.name}/${google_storage_bucket_object.requirements.name}"
    }
  }

  env_variables = {
    port = "8000"
  }

  noop_on_destroy = true
}

resource "google_app_engine_flexible_app_version" "foo" {
  project    = google_project_iam_member.gae_api.project
  version_id = "v1"
  service    = "custom"
  runtime    = "python"

  runtime_api_version = "1"

  resources {
    cpu       = 1
    memory_gb = 1
    disk_gb   = 10
  }

  entrypoint {
    shell = "gunicorn -b :$PORT main:app"
  }

  flexible_runtime_settings {
    operating_system = "ubuntu22"
    runtime_version = "3.11"
  }

  deployment {
    files {
      name = "main.py"
      source_url = "https://storage.googleapis.com/${google_storage_bucket.bucket.name}/${google_storage_bucket_object.main.name}"
    }

    files {
      name = "requirements.txt"
      source_url = "https://storage.googleapis.com/${google_storage_bucket.bucket.name}/${google_storage_bucket_object.requirements.name}"
    }

    files {
      name = "app.yaml"
      source_url = "https://storage.googleapis.com/${google_storage_bucket.bucket.name}/${google_storage_bucket_object.yaml.name}"
    }
  }

  liveness_check {
    path = ""
  }

  readiness_check {
    path = ""
  }

  env_variables = {
    port = "8000"
  }

  network {
    name             = google_compute_network.network.name
    subnetwork       = google_compute_subnetwork.subnetwork.name
    instance_ip_mode = "INTERNAL"
  }

  instance_class = "B2"

  manual_scaling {
    instances = 2
  }

  delete_service_on_destroy = true
  
  depends_on = [google_app_engine_standard_app_version.foo]
}

resource "google_storage_bucket" "bucket" {
  project  = google_project.my_project.project_id
  name     = "tf-test-%{random_suffix}-flex-ae-bucket"
  location = "US"
}

resource "google_storage_bucket_object" "yaml" {
  name   = "app.yaml"
  bucket = google_storage_bucket.bucket.name
  source = "./test-fixtures/hello-world-flask/app.yaml"
}

resource "google_storage_bucket_object" "requirements" {
  name   = "requirements.txt"
  bucket = google_storage_bucket.bucket.name
  source = "./test-fixtures/hello-world-flask/requirements.txt"
}

resource "google_storage_bucket_object" "main" {
  name   = "main.py"
  bucket = google_storage_bucket.bucket.name
  source = "./test-fixtures/hello-world-flask/main.py"
}`, context)
}
