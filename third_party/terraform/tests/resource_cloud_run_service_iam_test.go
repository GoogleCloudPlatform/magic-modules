package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccCloudRunServiceIamBinding(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(10),
		"role":          "roles/viewer",
		"namespace":     getTestProjectFromEnv(),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudRunServiceIamBinding_basic(context),
			},
			{
				ResourceName:      "google_cloud_run_service_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/locations/%s/services/%s roles/viewer", getTestProjectFromEnv(), getTestRegionFromEnv(), fmt.Sprintf("tftest-cloudrun%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test Iam Binding update
				Config: testAccCloudRunServiceIamBinding_update(context),
			},
			{
				ResourceName:      "google_cloud_run_service_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/locations/%s/services/%s roles/viewer", getTestProjectFromEnv(), getTestRegionFromEnv(), fmt.Sprintf("tftest-cloudrun%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCloudRunServiceIamMember(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(10),
		"role":          "roles/viewer",
		"namespace":     getTestProjectFromEnv(),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Test Iam Member creation (no update for member, no need to test)
				Config: testAccCloudRunServiceIamMember_basic(context),
			},
			{
				ResourceName:      "google_cloud_run_service_iam_member.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/locations/%s/services/%s roles/viewer user:admin@hashicorptest.com", getTestProjectFromEnv(), getTestRegionFromEnv(), fmt.Sprintf("tftest-cloudrun%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCloudRunServiceIamPolicy(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(10),
		"role":          "roles/viewer",
		"namespace":     getTestProjectFromEnv(),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudRunServiceIamPolicy_basic(context),
			},
			{
				ResourceName:      "google_cloud_run_service_iam_policy.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/locations/%s/services/%s", getTestProjectFromEnv(), getTestRegionFromEnv(), fmt.Sprintf("tftest-cloudrun%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCloudRunServiceIamMember_basic(context map[string]interface{}) string {
	return Nprintf(`
resource "google_cloud_run_service" "default" {
  name     = "tftest-cloudrun%{random_suffix}"
  location = "us-central1"

  metadata {
    namespace = "%{namespace}"
  }

  template {
    spec {
      containers {
        image = "gcr.io/cloudrun/hello"
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }
}

resource "google_cloud_run_service_iam_member" "foo" {
	location = "${google_cloud_run_service.default.location}"
	project = "${google_cloud_run_service.default.project}"
	service = "${google_cloud_run_service.default.name}"
	role = "%{role}"
	member = "user:admin@hashicorptest.com"
}
`, context)
}

func testAccCloudRunServiceIamPolicy_basic(context map[string]interface{}) string {
	return Nprintf(`
resource "google_cloud_run_service" "default" {
  name     = "tftest-cloudrun%{random_suffix}"
  location = "us-central1"

  metadata {
    namespace = "%{namespace}"
  }

  template {
    spec {
      containers {
        image = "gcr.io/cloudrun/hello"
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }
}

data "google_iam_policy" "foo" {
	binding {
		role = "%{role}"
		members = ["user:admin@hashicorptest.com"]
	}
}

resource "google_cloud_run_service_iam_policy" "foo" {
	location = "${google_cloud_run_service.default.location}"
	project = "${google_cloud_run_service.default.project}"
	service = "${google_cloud_run_service.default.name}"
	policy_data = "${data.google_iam_policy.foo.policy_data}"
}
`, context)
}

func testAccCloudRunServiceIamBinding_basic(context map[string]interface{}) string {
	return Nprintf(`
resource "google_cloud_run_service" "default" {
  name     = "tftest-cloudrun%{random_suffix}"
  location = "us-central1"

  metadata {
    namespace = "%{namespace}"
  }

  template {
    spec {
      containers {
        image = "gcr.io/cloudrun/hello"
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }
}

resource "google_cloud_run_service_iam_binding" "foo" {
	location = "${google_cloud_run_service.default.location}"
	project = "${google_cloud_run_service.default.project}"
	service = "${google_cloud_run_service.default.name}"
	role = "%{role}"
	members = ["user:admin@hashicorptest.com"]
}
`, context)
}

func testAccCloudRunServiceIamBinding_update(context map[string]interface{}) string {
	return Nprintf(`
resource "google_cloud_run_service" "default" {
  name     = "tftest-cloudrun%{random_suffix}"
  location = "us-central1"

  metadata {
    namespace = "%{namespace}"
  }

  template {
    spec {
      containers {
        image = "gcr.io/cloudrun/hello"
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }
}

resource "google_cloud_run_service_iam_binding" "foo" {
	location = "${google_cloud_run_service.default.location}"
	project = "${google_cloud_run_service.default.project}"
	service = "${google_cloud_run_service.default.name}"
	role = "%{role}"
	members = ["user:admin@hashicorptest.com", "user:paddy@hashicorp.com"]
}
`, context)
}
