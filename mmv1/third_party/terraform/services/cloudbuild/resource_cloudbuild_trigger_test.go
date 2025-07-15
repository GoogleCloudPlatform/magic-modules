package cloudbuild_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudBuildTrigger_basic(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudBuildTriggerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudBuildTrigger_basic(name),
			},
			{
				ResourceName:      "google_cloudbuild_trigger.build_trigger",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCloudBuildTrigger_updated(name),
			},
			{
				ResourceName:      "google_cloudbuild_trigger.build_trigger",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCloudBuildTrigger_available_secrets_config(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudBuildTriggerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudBuildTrigger_available_secrets_config(name),
			},
			{
				ResourceName:      "google_cloudbuild_trigger.build_trigger",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCloudBuildTrigger_available_secrets_config_update(name),
			},
			{
				ResourceName:      "google_cloudbuild_trigger.build_trigger",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCloudBuildTrigger_pubsub_config(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudBuildTriggerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudBuildTrigger_pubsub_config(name),
			},
			{
				ResourceName:      "google_cloudbuild_trigger.build_trigger",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCloudBuildTrigger_pubsub_config_update(name),
			},
			{
				ResourceName:      "google_cloudbuild_trigger.build_trigger",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCloudBuildTrigger_webhook_config(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudBuildTriggerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudBuildTrigger_webhook_config(name),
			},
			{
				ResourceName:      "google_cloudbuild_trigger.build_trigger",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCloudBuildTrigger_webhook_config_update(name),
			},
			{
				ResourceName:      "google_cloudbuild_trigger.build_trigger",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCloudBuildTrigger_customizeDiffTimeoutSum(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudBuildTriggerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccCloudBuildTrigger_customizeDiffTimeoutSum(name),
				ExpectError: regexp.MustCompile("cannot be greater than build timeout"),
			},
		},
	})
}

func TestAccCloudBuildTrigger_customizeDiffTimeoutFormat(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudBuildTriggerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccCloudBuildTrigger_customizeDiffTimeoutFormat(name),
				ExpectError: regexp.MustCompile("Error parsing build timeout"),
			},
		},
	})
}

func TestAccCloudBuildTrigger_disable(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudBuildTriggerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudBuildTrigger_basic(name),
			},
			{
				ResourceName:      "google_cloudbuild_trigger.build_trigger",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCloudBuildTrigger_basicDisabled(name),
			},
			{
				ResourceName:      "google_cloudbuild_trigger.build_trigger",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCloudBuildTrigger_fullStep(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudBuildTriggerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudBuildTrigger_fullStep(),
			},
			{
				ResourceName:      "google_cloudbuild_trigger.build_trigger",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCloudBuildTrigger_basic_bitbucket(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudBuildTriggerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudBuildTrigger_basic_bitbucket(name),
			},
			{
				ResourceName:      "google_cloudbuild_trigger.build_trigger",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCloudBuildTrigger_basic(name string) string {
	return fmt.Sprintf(`
resource "google_cloudbuild_trigger" "build_trigger" {
  name        = "%s"
  description = "acceptance test build trigger"
  trigger_template {
    branch_name = "main"
    repo_name   = "some-repo"
  }
  build {
    tags   = ["team-a", "service-b"]
    timeout = "1800s"
    step {
      name = "gcr.io/cloud-builders/gsutil"
      args = ["cp", "gs://mybucket/remotefile.zip", "localfile.zip"]
      timeout = "300s"
    }
    step {
      name = "gcr.io/cloud-builders/go"
      args = ["build", "my_package"]
      env  = ["env1=two"]
      timeout = "300s"
    }
    step {
      name = "gcr.io/cloud-builders/docker"
      args = ["build", "-t", "gcr.io/$PROJECT_ID/$REPO_NAME:$COMMIT_SHA", "-f", "Dockerfile", "."]
      timeout = "300s"
    }
    artifacts {
      images = ["gcr.io/$PROJECT_ID/$REPO_NAME:$COMMIT_SHA"]
      objects {
        location = "gs://bucket/path/to/somewhere/"
        paths = ["path"]
      }
    }
    options {
      source_provenance_hash = ["MD5"]
      requested_verify_option = "VERIFIED"
      machine_type = "N1_HIGHCPU_8"
      disk_size_gb = 100
      substitution_option = "ALLOW_LOOSE"
      dynamic_substitutions = false
      log_streaming_option = "STREAM_OFF"
      worker_pool = "pool"
      logging = "LEGACY"
      env = ["ekey = evalue"]
      secret_env = ["secretenv = svalue"]
      volumes {
        name = "v1"
        path = "v1"
      }
    }
  }
}
`, name)
}

func testAccCloudBuildTrigger_basic_bitbucket(name string) string {
	return fmt.Sprintf(`
resource "google_cloudbuild_trigger" "build_trigger" {
  name        = "%s"
  description = "acceptance test build trigger on bitbucket"
  trigger_template {
    branch_name = "main"
    repo_name   = "some-repo"
  }
  git_file_source {
    path      = "cloudbuild.yaml"
    uri       = "https://bitbucket.org/myorg/myrepo"
    revision  = "refs/heads/develop"
    repo_type = "BITBUCKET_SERVER"
    bitbucket_server_config = "projects/123456789/locations/us-central1/bitbucketServerConfigs/myBitbucketConfig"
  }
}
`, name)
}

func testAccCloudBuildTrigger_basicDisabled(name string) string {
	return fmt.Sprintf(`
resource "google_cloudbuild_trigger" "build_trigger" {
  disabled    = true
  name        = "%s"
  description = "acceptance test build trigger"
  trigger_template {
    branch_name = "main"
    repo_name   = "some-repo"
  }
  build {
    images = ["gcr.io/$PROJECT_ID/$REPO_NAME:$COMMIT_SHA"]
    tags   = ["team-a", "service-b"]
    step {
      name = "gcr.io/cloud-builders/gsutil"
      args = ["cp", "gs://mybucket/remotefile.zip", "localfile.zip"]
    }
    step {
      name = "gcr.io/cloud-builders/go"
      args = ["build", "my_package"]
      env  = ["env1=two"]
    }
    step {
      name = "gcr.io/cloud-builders/docker"
      args = ["build", "-t", "gcr.io/$PROJECT_ID/$REPO_NAME:$COMMIT_SHA", "-f", "Dockerfile", "."]
    }
  }
}
`, name)
}

func testAccCloudBuildTrigger_fullStep() string {
	return fmt.Sprintf(`
resource "google_cloudbuild_trigger" "build_trigger" {
  description = "acceptance test build trigger"
  trigger_template {
    branch_name = "main"
    repo_name   = "some-repo"
	invert_regex = false
  }
  build {
    images = ["gcr.io/$PROJECT_ID/$REPO_NAME:$COMMIT_SHA"]
    tags   = ["team-a", "service-b"]
    step {
      name       = "gcr.io/cloud-builders/go"
      args       = ["build", "my_package"]
      env        = ["env1=two"]
      dir        = "directory"
      id         = "12345"
      secret_env = ["fooo"]
      timeout    = "100s"
      wait_for   = ["something"]
    }
  }
}
`)
}

func testAccCloudBuildTrigger_updated(name string) string {
	return fmt.Sprintf(`
resource "google_cloudbuild_trigger" "build_trigger" {
  description = "acceptance test build trigger updated"
  name        = "%s"
  trigger_template {
    branch_name = "main-updated"
    repo_name   = "some-repo-updated"
	invert_regex = true
  }
  build {
    images = ["gcr.io/$PROJECT_ID/$REPO_NAME:$SHORT_SHA"]
    tags   = ["team-a", "service-b", "updated"]
    timeout = "2100s"
    step {
      name = "gcr.io/cloud-builders/gsutil"
      args = ["cp", "gs://mybucket/remotefile.zip", "localfile-updated.zip"]
      timeout = "300s"
    }
    step {
      name = "gcr.io/cloud-builders/go"
      args = ["build", "my_package_updated"]
      timeout = "300s"
    }
    step {
      name = "gcr.io/cloud-builders/docker"
      args = ["build", "-t", "gcr.io/$PROJECT_ID/$REPO_NAME:$SHORT_SHA", "-f", "Dockerfile", "."]
      timeout = "300s"
    }
    step {
      name = "gcr.io/$PROJECT_ID/$REPO_NAME:$SHORT_SHA"
      args = ["test"]
      timeout = "300s"
    }
    logs_bucket = "gs://mybucket/logs"
    options {
      # this field is always enabled for triggered build and cannot be overridden in the build configuration file.
      dynamic_substitutions = true
    }
  }
}
  `, name)
}

func testAccCloudBuildTrigger_available_secrets_config(name string) string {
	return fmt.Sprintf(`
resource "google_cloudbuild_trigger" "build_trigger" {
  name        = "%s"
  description = "acceptance test build trigger"
  trigger_template {
    branch_name = "main"
    repo_name   = "some-repo"
  }
  build {
    tags   = ["team-a", "service-b"]
    timeout = "1800s"
    step {
      name = "gcr.io/cloud-builders/gsutil"
      args = ["cp", "gs://mybucket/remotefile.zip", "localfile.zip"]
      timeout = "300s"
    }
    available_secrets {
      secret_manager {
        env          = "MY_SECRET"
        version_name = "projects/myProject/secrets/mySecret/versions/latest"
      }
    }
  }
}
`, name)
}

func testAccCloudBuildTrigger_available_secrets_config_update(name string) string {
	return fmt.Sprintf(`
resource "google_cloudbuild_trigger" "build_trigger" {
  name        = "%s"
  description = "acceptance test build trigger updated"
  trigger_template {
    branch_name = "main"
    repo_name   = "some-repo"
  }
  build {
    tags   = ["team-a", "service-b"]
    timeout = "1800s"
    step {
      name = "gcr.io/cloud-builders/gsutil"
      args = ["cp", "gs://mybucket/remotefile.zip", "localfile.zip"]
      timeout = "300s"
    }
  }
}
`, name)
}

func testAccCloudBuildTrigger_pubsub_config(name string) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "build-trigger" {
  name = "%s"
}

resource "google_cloudbuild_trigger" "build_trigger" {
  name        = "%s"
  description = "acceptance test build trigger"
  pubsub_config {
    topic = "${google_pubsub_topic.build-trigger.id}"
  }
  build {
    tags   = ["team-a", "service-b"]
    timeout = "1800s"
    step {
      name = "gcr.io/cloud-builders/gsutil"
      args = ["cp", "gs://mybucket/remotefile.zip", "localfile.zip"]
      timeout = "300s"
    }
  }
  depends_on = [
    google_pubsub_topic.build-trigger
  ]
}
`, name, name)
}

func testAccCloudBuildTrigger_pubsub_config_update(name string) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "build-trigger" {
  name = "%s"
}

resource "google_cloudbuild_trigger" "build_trigger" {
  name        = "%s"
  description = "acceptance test build trigger updated"
  pubsub_config {
    topic = "${google_pubsub_topic.build-trigger.id}"
  }
  build {
    tags   = ["team-a", "service-b"]
    timeout = "1800s"
    step {
      name = "gcr.io/cloud-builders/gsutil"
      args = ["cp", "gs://mybucket/remotefile.zip", "localfile.zip"]
      timeout = "300s"
    }
  }
  depends_on = [
    google_pubsub_topic.build-trigger
  ]
}
`, name, name)
}

func testAccCloudBuildTrigger_webhook_config(name string) string {
	return fmt.Sprintf(`
resource "google_secret_manager_secret" "webhook_trigger_secret_key" {
  secret_id = "webhook_trigger-secret-key"

  replication {
    user_managed {
      replicas {
        location = "us-central1"
      }
    }
  }
}

resource "google_secret_manager_secret_version" "webhook_trigger_secret_key_data" {
  secret = google_secret_manager_secret.webhook_trigger_secret_key.id

  secret_data = "secretkeygoeshere"
}

data "google_project" "project" {}

data "google_iam_policy" "secret_accessor" {
  binding {
    role = "roles/secretmanager.secretAccessor"
    members = [
      "serviceAccount:service-${data.google_project.project.number}@gcp-sa-cloudbuild.iam.gserviceaccount.com",
    ]
  }
}

resource "google_secret_manager_secret_iam_policy" "policy" {
  project = google_secret_manager_secret.webhook_trigger_secret_key.project
  secret_id = google_secret_manager_secret.webhook_trigger_secret_key.secret_id
  policy_data = data.google_iam_policy.secret_accessor.policy_data
}

resource "google_cloudbuild_trigger" "build_trigger" {
  name        = "%s"

  webhook_config {
    secret = "${google_secret_manager_secret_version.webhook_trigger_secret_key_data.id}"
  }

  build {
    step {
      name = "ubuntu"
      args = [
        "-c", 
        <<EOT
          echo data
        EOT
      ]
      entrypoint = "bash"
    }
  }

  depends_on = [
    google_secret_manager_secret_version.webhook_trigger_secret_key_data,
    google_secret_manager_secret_iam_policy.policy
  ]
}
`, name)
}

func testAccCloudBuildTrigger_webhook_config_update(name string) string {
	return fmt.Sprintf(`
resource "google_secret_manager_secret" "webhook_trigger_secret_key" {
  secret_id = "webhook_trigger-secret-key"

  replication {
    user_managed {
      replicas {
        location = "us-central1"
      }
    }
  }
}

resource "google_secret_manager_secret_version" "webhook_trigger_secret_key_data" {
  secret = google_secret_manager_secret.webhook_trigger_secret_key.id

  secret_data = "secretkeygoeshere"
}

data "google_project" "project" {}

data "google_iam_policy" "secret_accessor" {
  binding {
    role = "roles/secretmanager.secretAccessor"
    members = [
      "serviceAccount:service-${data.google_project.project.number}@gcp-sa-cloudbuild.iam.gserviceaccount.com",
    ]
  }
}

resource "google_secret_manager_secret_iam_policy" "policy" {
  project = google_secret_manager_secret.webhook_trigger_secret_key.project
  secret_id = google_secret_manager_secret.webhook_trigger_secret_key.secret_id
  policy_data = data.google_iam_policy.secret_accessor.policy_data
}

resource "google_cloudbuild_trigger" "build_trigger" {
  name        = "%s"

  webhook_config {
    secret = "${google_secret_manager_secret_version.webhook_trigger_secret_key_data.id}"
  }

  build {
    step {
      name = "ubuntu"
      args = [
        "-c", 
        <<EOT
          echo data-updated
        EOT
      ]
      entrypoint = "bash"
    }
  }

  depends_on = [
    google_secret_manager_secret_version.webhook_trigger_secret_key_data,
    google_secret_manager_secret_iam_policy.policy
  ]
}
`, name)
}

func testAccCloudBuildTrigger_customizeDiffTimeoutSum(name string) string {
	return fmt.Sprintf(`
resource "google_cloudbuild_trigger" "build_trigger" {
  name = "%s"
  description = "acceptance test build trigger"
  trigger_template {
    branch_name = "main"
    repo_name   = "some-repo"
  }
  build {
    images = ["gcr.io/$PROJECT_ID/$REPO_NAME:$COMMIT_SHA"]
    tags = ["team-a", "service-b"]
    timeout = "900s"
    step {
      name = "gcr.io/cloud-builders/gsutil"
      args = ["cp", "gs://mybucket/remotefile.zip", "localfile.zip"]
      timeout = "500s"
    }
    step {
      name = "gcr.io/cloud-builders/go"
      args = ["build", "my_package"]
      env = ["env1=two"]
      timeout = "500s"
    }
    step {
      name = "gcr.io/cloud-builders/docker"
      args = ["build", "-t", "gcr.io/$PROJECT_ID/$REPO_NAME:$COMMIT_SHA", "-f", "Dockerfile", "."]
      timeout = "500s"
    }
  }
}
  `, name)
}

func testAccCloudBuildTrigger_customizeDiffTimeoutFormat(name string) string {
	return fmt.Sprintf(`
resource "google_cloudbuild_trigger" "build_trigger" {
  name = "%s"
  description = "acceptance test build trigger"
  trigger_template {
    branch_name = "main"
    repo_name   = "some-repo"
  }
  build {
    images = ["gcr.io/$PROJECT_ID/$REPO_NAME:$COMMIT_SHA"]
    tags = ["team-a", "service-b"]
    timeout = "1200"
    step {
      name = "gcr.io/cloud-builders/gsutil"
      args = ["cp", "gs://mybucket/remotefile.zip", "localfile.zip"]
      timeout = "500s"
    }
  }
}
`, name)
}

func TestAccCloudBuildTrigger_developerConnect_pullRequest(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	connectionName := fmt.Sprintf("tf-test-conn-%d", acctest.RandInt(t))
	repoLinkName := fmt.Sprintf("tf-test-repo-link-%d", acctest.RandInt(t))
	location := acctest.GetEnvDefaultOverride("LOCATION", "us-central1")

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudBuildTriggerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudBuildTrigger_developerConnect_pullRequestConfig(name, connectionName, repoLinkName, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_cloudbuild_trigger.build_trigger", "name", name),
					resource.TestCheckResourceAttrSet("google_cloudbuild_trigger.build_trigger", "developer_connect_event_config.0.git_repository_link"),
					resource.TestCheckResourceAttr("google_cloudbuild_trigger.build_trigger", "developer_connect_event_config.0.pull_request.0.branch", "main"),
					resource.TestCheckResourceAttr("google_cloudbuild_trigger.build_trigger", "developer_connect_event_config.0.pull_request.0.comment_control", "COMMENTS_ENABLED_FOR_EXTERNAL_CONTRIBUTORS_ONLY"),
					resource.TestCheckResourceAttr("google_cloudbuild_trigger.build_trigger", "developer_connect_event_config.0.pull_request.0.invert_regex", "true"),
				),
			},
			{
				ResourceName:      "google_cloudbuild_trigger.build_trigger",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCloudBuildTrigger_developerConnect_pullRequestConfigUpdated(name, connectionName, repoLinkName, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_cloudbuild_trigger.build_trigger", "name", name),
					resource.TestCheckResourceAttr("google_cloudbuild_trigger.build_trigger", "developer_connect_event_config.0.pull_request.0.branch", "feature/.*"),
					resource.TestCheckResourceAttr("google_cloudbuild_trigger.build_trigger", "developer_connect_event_config.0.pull_request.0.comment_control", "COMMENTS_DISABLED"),
				),
			},
			{
				ResourceName:      "google_cloudbuild_trigger.build_trigger",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCloudBuildTrigger_developerConnect_push(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	connectionName := fmt.Sprintf("tf-test-conn-%d", acctest.RandInt(t))
	repoLinkName := fmt.Sprintf("tf-test-repo-link-%d", acctest.RandInt(t))
	location := acctest.GetEnvDefaultOverride("LOCATION", "us-central1")

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudBuildTriggerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudBuildTrigger_developerConnect_pushConfig(name, connectionName, repoLinkName, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_cloudbuild_trigger.build_trigger", "name", name),
					resource.TestCheckResourceAttrSet("google_cloudbuild_trigger.build_trigger", "developer_connect_event_config.0.git_repository_link"),
					resource.TestCheckResourceAttr("google_cloudbuild_trigger.build_trigger", "developer_connect_event_config.0.push.0.branch", "release/v.*"),
					resource.TestCheckResourceAttr("google_cloudbuild_trigger.build_trigger", "developer_connect_event_config.0.push.0.invert_regex", "false"),
					resource.TestCheckNoResourceAttr("google_cloudbuild_trigger.build_trigger", "developer_connect_event_config.0.push.0.tag"),
				),
			},
			{
				ResourceName:      "google_cloudbuild_trigger.build_trigger",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCloudBuildTrigger_developerConnect_pushConfigUpdated(name, connectionName, repoLinkName, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_cloudbuild_trigger.build_trigger", "name", name),
					resource.TestCheckResourceAttr("google_cloudbuild_trigger.build_trigger", "developer_connect_event_config.0.push.0.tag", "v\\d+\\.\\d+\\.\\d+"),
					resource.TestCheckResourceAttr("google_cloudbuild_trigger.build_trigger", "developer_connect_event_config.0.push.0.invert_regex", "true"),
					resource.TestCheckNoResourceAttr("google_cloudbuild_trigger.build_trigger", "developer_connect_event_config.0.push.0.branch"),
				),
			},
			{
				ResourceName:      "google_cloudbuild_trigger.build_trigger",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCloudBuildTrigger_developerConnect_pullRequestConfig(name, connectionName, repoLinkName, location string) string {
	return fmt.Sprintf(`
resource "google_developer_connect_connection" "dev_connect_connection" {
  project = %[1]s
  location = "%[4]s"
  name = "%[2]s"
  github_config {
    app_installation_id = "1234567"
    authorizer_credential {
        oauth_token_secret_version = "projects/%[1]s/secrets/my-github-pat/versions/latest"
    }
  }
}

resource "google_developer_connect_git_repository_link" "dev_connect_repo_link" {
  parent = google_developer_connect_connection.dev_connect_connection.id
  name = "%[3]s"
  github {
    repo_uri = "https://github.com/GoogleCloudPlatform/my-test-repo"
  }
}

resource "google_cloudbuild_trigger" "build_trigger" {
  project = %[1]s
  location = "%[4]s"
  name = "%[5]s"
  description = "Developer Connect Pull Request Trigger"

  developer_connect_event_config {
    git_repository_link = google_developer_connect_git_repository_link.dev_connect_repo_link.id
    pull_request {
      branch          = "main"
      comment_control = "COMMENTS_ENABLED_FOR_EXTERNAL_CONTRIBUTORS_ONLY"
      invert_regex    = true
    }
  }

  build {
    step {
      name = "gcr.io/cloud-builders/gcloud"
      args = ["builds", "triggers", "describe", google_cloudbuild_trigger.build_trigger.name, "--location=%[4]s", "--format=json"]
    }
  }

  depends_on = [
    google_developer_connect_git_repository_link.dev_connect_repo_link
  ]
}
`, acctest.ProjectID(), connectionName, repoLinkName, location, name)
}

func testAccCloudBuildTrigger_developerConnect_pullRequestConfigUpdated(name, connectionName, repoLinkName, location string) string {
	return fmt.Sprintf(`
resource "google_developer_connect_connection" "dev_connect_connection" {
  project = %[1]s
  location = "%[4]s"
  name = "%[2]s"
  github_config {
    app_installation_id = "1234567"
    authorizer_credential {
        oauth_token_secret_version = "projects/%[1]s/secrets/my-github-pat/versions/latest"
    }
  }
}

resource "google_developer_connect_git_repository_link" "dev_connect_repo_link" {
  parent = google_developer_connect_connection.dev_connect_connection.id
  name = "%[3]s"
  github {
    repo_uri = "https://github.com/GoogleCloudPlatform/my-test-repo"
  }
}

resource "google_cloudbuild_trigger" "build_trigger" {
  project = %[1]s
  location = "%[4]s"
  name = "%[5]s"
  description = "Developer Connect Pull Request Trigger Updated"

  developer_connect_event_config {
    git_repository_link = google_developer_connect_git_repository_link.dev_connect_repo_link.id
    pull_request {
      branch          = "feature/.*"
      comment_control = "COMMENTS_DISABLED"
      invert_regex    = false
    }
  }

  build {
    step {
      name = "gcr.io/cloud-builders/gcloud"
      args = ["builds", "triggers", "describe", google_cloudbuild_trigger.build_trigger.name, "--location=%[4]s", "--format=json"]
    }
  }

  depends_on = [
    google_developer_connect_git_repository_link.dev_connect_repo_link
  ]
}
`, acctest.ProjectID(), connectionName, repoLinkName, location, name)
}

func testAccCloudBuildTrigger_developerConnect_pushConfig(name, connectionName, repoLinkName, location string) string {
	return fmt.Sprintf(`
resource "google_developer_connect_connection" "dev_connect_connection" {
  project = %[1]s
  location = "%[4]s"
  name = "%[2]s"
  github_config {
    app_installation_id = "1234567"
    authorizer_credential {
        oauth_token_secret_version = "projects/%[1]s/secrets/my-github-pat/versions/latest"
    }
  }
}

resource "google_developer_connect_git_repository_link" "dev_connect_repo_link" {
  parent = google_developer_connect_connection.dev_connect_connection.id
  name = "%[3]s"
  github {
    repo_uri = "https://github.com/GoogleCloudPlatform/my-test-repo"
  }
}

resource "google_cloudbuild_trigger" "build_trigger" {
  project = %[1]s
  location = "%[4]s"
  name = "%[5]s"
  description = "Developer Connect Push Trigger"

  developer_connect_event_config {
    git_repository_link = google_developer_connect_git_repository_link.dev_connect_repo_link.id
    push {
      branch       = "release/v.*"
      invert_regex = false
    }
  }

  build {
    step {
      name = "gcr.io/cloud-builders/gcloud"
      args = ["builds", "triggers", "describe", google_cloudbuild_trigger.build_trigger.name, "--location=%[4]s", "--format=json"]
    }
  }

  depends_on = [
    google_developer_connect_git_repository_link.dev_connect_repo_link
  ]
}
`, acctest.ProjectID(), connectionName, repoLinkName, location, name)
}

func testAccCloudBuildTrigger_developerConnect_pushConfigUpdated(name, connectionName, repoLinkName, location string) string {
	return fmt.Sprintf(`
resource "google_developer_connect_connection" "dev_connect_connection" {
  project = %[1]s
  location = "%[4]s"
  name = "%[2]s"
  github_config {
    app_installation_id = "1234567"
    authorizer_credential {
        oauth_token_secret_version = "projects/%[1]s/secrets/my-github-pat/versions/latest"
    }
  }
}

resource "google_developer_connect_git_repository_link" "dev_connect_repo_link" {
  parent = google_developer_connect_connection.dev_connect_connection.id
  name = "%[3]s"
  github {
    repo_uri = "https://github.com/GoogleCloudPlatform/my-test-repo"
  }
}

resource "google_cloudbuild_trigger" "build_trigger" {
  project = %[1]s
  location = "%[4]s"
  name = "%[5]s"
  description = "Developer Connect Push Trigger Updated"

  developer_connect_event_config {
    git_repository_link = google_developer_connect_git_repository_link.dev_connect_repo_link.id
    push {
      tag          = "v\\d+\\.\\d+\\.\\d+"
      invert_regex = true
    }
  }

  build {
    step {
      name = "gcr.io/cloud-builders/gcloud"
      args = ["builds", "triggers", "describe", google_cloudbuild_trigger.build_trigger.name, "--location=%[4]s", "--format=json"]
    }
  }

  depends_on = [
    google_developer_connect_git_repository_link.dev_connect_repo_link
  ]
}
`, acctest.ProjectID(), connectionName, repoLinkName, location, name)
}
