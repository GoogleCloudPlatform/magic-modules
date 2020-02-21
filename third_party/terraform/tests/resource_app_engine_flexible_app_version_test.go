package google

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccAppEngineFlexibleAppVersion_update(t *testing.T) {
	t.Parallel()

	resourceName := fmt.Sprintf("tf-test-ae-service-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppEngineFlexibleAppVersionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAppEngineFlexibleAppVersion_python(resourceName),
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

func testAccAppEngineFlexibleAppVersion_python(resourceName string) string {
	return fmt.Sprintf(`
resource "google_app_engine_flexible_app_version" "foo" {
  version_id = "v1"
  service    = "%s"
  runtime    = "python"

  runtime_api_version = "1"

  resources {
    cpu       = 1
    memory_gb = 0.5
    disk_gb   = 10
  }

  entrypoint {
    shell = "gunicorn -b :$PORT project_name.wsgi"
  }

  deployment {
    zip {
      source_url = "https://storage.googleapis.com/${google_storage_bucket.bucket.name}/${google_storage_bucket_object.object.name}"
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
    name       = "default"
    subnetwork = "default"
  }

  instance_class = "B1"

  manual_scaling {
    instances = 1
  }

  delete_service_on_destroy = true
}

resource "google_storage_bucket" "bucket" {
  name = "%s-bucket"
}

resource "google_storage_bucket_object" "object" {
  name   = "hello-world-django.zip"
  bucket = google_storage_bucket.bucket.name
  source = "./test-fixtures/appengine/hello-world-django.zip"
}`, resourceName, resourceName)
}
