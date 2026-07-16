package colab_test

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	_ "github.com/hashicorp/terraform-provider-google/google/services/colab"
	_ "github.com/hashicorp/terraform-provider-google/google/services/storage"
)

func TestAccColabSchedule_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"location":           envvar.GetTestRegionFromEnv(),
		"project_id":         envvar.GetTestProjectFromEnv(),
		"service_account":    envvar.GetTestServiceAccountFromEnv(t),
		"end_time":           time.Date(time.Now().Year(), 12, 31, 0, 0, 0, 0, time.Now().Location()).AddDate(0, 0, 10).Format(time.RFC3339),
		"start_time":         time.Date(time.Now().Year(), 12, 31, 0, 0, 0, 0, time.Now().Location()).AddDate(0, 0, 1).Format(time.RFC3339),
		"random_suffix":      acctest.RandString(t, 10),
		"updated_start_time": time.Date(time.Now().Year(), 12, 31, 0, 0, 0, 0, time.Now().Location()).AddDate(0, 0, 2).Format(time.RFC3339),
		"updated_end_time":   time.Date(time.Now().Year(), 12, 31, 0, 0, 0, 0, time.Now().Location()).AddDate(0, 0, 5).Format(time.RFC3339),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckColabScheduleDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccColabSchedule_full(context),
			},
			{
				ResourceName:            "google_colab_schedule.schedule",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
			{
				Config: testAccColabSchedule_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_colab_schedule.schedule", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_colab_schedule.schedule",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
		},
	})
}

func TestAccColabSchedule_update_state(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"location":           envvar.GetTestRegionFromEnv(),
		"project_id":         envvar.GetTestProjectFromEnv(),
		"service_account":    envvar.GetTestServiceAccountFromEnv(t),
		"end_time":           time.Date(time.Now().Year(), 12, 31, 0, 0, 0, 0, time.Now().Location()).AddDate(0, 0, 10).Format(time.RFC3339),
		"start_time":         time.Date(time.Now().Year(), 12, 31, 0, 0, 0, 0, time.Now().Location()).AddDate(0, 0, 1).Format(time.RFC3339),
		"random_suffix":      acctest.RandString(t, 10),
		"updated_start_time": time.Date(time.Now().Year(), 12, 31, 0, 0, 0, 0, time.Now().Location()).AddDate(0, 0, 2).Format(time.RFC3339),
		"updated_end_time":   time.Date(time.Now().Year(), 12, 31, 0, 0, 0, 0, time.Now().Location()).AddDate(0, 0, 5).Format(time.RFC3339),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckColabScheduleDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccColabSchedule_full(context),
			},
			{
				ResourceName:            "google_colab_schedule.schedule",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "desired_state"},
			},
			{
				Config: testAccColabSchedule_active(context),
			},
			{
				ResourceName:            "google_colab_schedule.schedule",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "desired_state"},
			},
			{
				Config: testAccColabSchedule_paused(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_colab_schedule.schedule", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_colab_schedule.schedule",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "desired_state"},
			},
			{
				Config: testAccColabSchedule_active(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_colab_schedule.schedule", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_colab_schedule.schedule",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "desired_state"},
			},
			{
				Config: testAccColabSchedule_full(context),
			},
			{
				ResourceName:            "google_colab_schedule.schedule",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "desired_state"},
			},
			{
				Config: testAccColabSchedule_paused(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_colab_schedule.schedule", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_colab_schedule.schedule",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "desired_state"},
			},
			{
				Config: testAccColabSchedule_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_colab_schedule.schedule", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_colab_schedule.schedule",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "desired_state"},
			},
			{
				Config: testAccColabSchedule_active(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_colab_schedule.schedule", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_colab_schedule.schedule",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "desired_state"},
			},
		},
	})
}

func testAccColabSchedule_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_colab_runtime_template" "my_runtime_template" {
  name = "tf-test-runtime-template%{random_suffix}"
  display_name = "Runtime template"
  location = "us-central1"

  machine_spec {
    machine_type     = "e2-standard-4"
  }

  network_spec {
    enable_internet_access = true
  }
}

resource "google_storage_bucket" "output_bucket" {
  name          = "tf_test_my_bucket%{random_suffix}"
  location      = "US"
  force_destroy = true
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_object" "notebook" {
  name   = "hello_world.ipynb"
  bucket = google_storage_bucket.output_bucket.name
  content = <<EOF
    {
      "cells": [
        {
          "cell_type": "code",
          "execution_count": null,
          "metadata": {},
          "outputs": [],
          "source": [
            "print(\"Hello, World!\")"
          ]
        }
      ],
      "metadata": {
        "kernelspec": {
          "display_name": "Python 3",
          "language": "python",
          "name": "python3"
        },
        "language_info": {
          "codemirror_mode": {
            "name": "ipython",
            "version": 3
          },
          "file_extension": ".py",
          "mimetype": "text/x-python",
          "name": "python",
          "nbconvert_exporter": "python",
          "pygments_lexer": "ipython3",
          "version": "3.8.5"
        }
      },
      "nbformat": 4,
      "nbformat_minor": 4
    }
    EOF
}

resource "google_colab_schedule" "schedule" {
  display_name = "tf-test-schedule%{random_suffix}"
  location = "%{location}"
  allow_queueing = true
  max_concurrent_run_count = 2
  cron = "TZ=America/Los_Angeles * * * * *"
  max_run_count = 5
  start_time = "%{start_time}"
  end_time = "%{end_time}"

  create_notebook_execution_job_request {
    notebook_execution_job {
      display_name = "Notebook execution"
      gcs_notebook_source {
        uri = "gs://${google_storage_bucket_object.notebook.bucket}/${google_storage_bucket_object.notebook.name}"
        generation = google_storage_bucket_object.notebook.generation
      }

      notebook_runtime_template_resource_name = "projects/${google_colab_runtime_template.my_runtime_template.project}/locations/${google_colab_runtime_template.my_runtime_template.location}/notebookRuntimeTemplates/${google_colab_runtime_template.my_runtime_template.name}"
      gcs_output_uri = "gs://${google_storage_bucket.output_bucket.name}"
      service_account = "%{service_account}"
      }
  }

  depends_on = [
    google_colab_runtime_template.my_runtime_template,
    google_storage_bucket.output_bucket,
  ]
}
`, context)
}

func testAccColabSchedule_paused(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_colab_runtime_template" "my_runtime_template" {
  name = "tf-test-runtime-template%{random_suffix}"
  display_name = "Runtime template"
  location = "us-central1"

  machine_spec {
    machine_type     = "e2-standard-4"
  }

  network_spec {
    enable_internet_access = true
  }
}

resource "google_storage_bucket" "output_bucket" {
  name          = "tf_test_my_bucket%{random_suffix}"
  location      = "US"
  force_destroy = true
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_object" "notebook" {
  name   = "hello_world.ipynb"
  bucket = google_storage_bucket.output_bucket.name
  content = <<EOF
    {
      "cells": [
        {
          "cell_type": "code",
          "execution_count": null,
          "metadata": {},
          "outputs": [],
          "source": [
            "print(\"Hello, World!\")"
          ]
        }
      ],
      "metadata": {
        "kernelspec": {
          "display_name": "Python 3",
          "language": "python",
          "name": "python3"
        },
        "language_info": {
          "codemirror_mode": {
            "name": "ipython",
            "version": 3
          },
          "file_extension": ".py",
          "mimetype": "text/x-python",
          "name": "python",
          "nbconvert_exporter": "python",
          "pygments_lexer": "ipython3",
          "version": "3.8.5"
        }
      },
      "nbformat": 4,
      "nbformat_minor": 4
    }
    EOF
}

resource "google_colab_schedule" "schedule" {
  display_name = "tf-test-schedule%{random_suffix}"
  location = "%{location}"
  allow_queueing = true
  max_concurrent_run_count = 2
  cron = "TZ=America/Los_Angeles * * * * *"
  max_run_count = 5
  start_time = "%{start_time}"
  end_time = "%{end_time}"

  desired_state = "PAUSED"

  create_notebook_execution_job_request {
    notebook_execution_job {
      display_name = "Notebook execution"
      gcs_notebook_source {
        uri = "gs://${google_storage_bucket_object.notebook.bucket}/${google_storage_bucket_object.notebook.name}"
        generation = google_storage_bucket_object.notebook.generation
      }

      notebook_runtime_template_resource_name = "projects/${google_colab_runtime_template.my_runtime_template.project}/locations/${google_colab_runtime_template.my_runtime_template.location}/notebookRuntimeTemplates/${google_colab_runtime_template.my_runtime_template.name}"
      gcs_output_uri = "gs://${google_storage_bucket.output_bucket.name}"
      service_account = "%{service_account}"
      }
  }

  depends_on = [
    google_colab_runtime_template.my_runtime_template,
    google_storage_bucket.output_bucket,
  ]
}
`, context)
}

func testAccColabSchedule_active(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_colab_runtime_template" "my_runtime_template" {
  name = "tf-test-runtime-template%{random_suffix}"
  display_name = "Runtime template"
  location = "us-central1"

  machine_spec {
    machine_type     = "e2-standard-4"
  }

  network_spec {
    enable_internet_access = true
  }
}

resource "google_storage_bucket" "output_bucket" {
  name          = "tf_test_my_bucket%{random_suffix}"
  location      = "US"
  force_destroy = true
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_object" "notebook" {
  name   = "hello_world.ipynb"
  bucket = google_storage_bucket.output_bucket.name
  content = <<EOF
    {
      "cells": [
        {
          "cell_type": "code",
          "execution_count": null,
          "metadata": {},
          "outputs": [],
          "source": [
            "print(\"Hello, World!\")"
          ]
        }
      ],
      "metadata": {
        "kernelspec": {
          "display_name": "Python 3",
          "language": "python",
          "name": "python3"
        },
        "language_info": {
          "codemirror_mode": {
            "name": "ipython",
            "version": 3
          },
          "file_extension": ".py",
          "mimetype": "text/x-python",
          "name": "python",
          "nbconvert_exporter": "python",
          "pygments_lexer": "ipython3",
          "version": "3.8.5"
        }
      },
      "nbformat": 4,
      "nbformat_minor": 4
    }
    EOF
}

resource "google_colab_schedule" "schedule" {
  display_name = "tf-test-schedule%{random_suffix}"
  location = "%{location}"
  allow_queueing = true
  max_concurrent_run_count = 2
  cron = "TZ=America/Los_Angeles * * * * *"
  max_run_count = 5
  start_time = "%{start_time}"
  end_time = "%{end_time}"

  desired_state = "ACTIVE"

  create_notebook_execution_job_request {
    notebook_execution_job {
      display_name = "Notebook execution"
      gcs_notebook_source {
        uri = "gs://${google_storage_bucket_object.notebook.bucket}/${google_storage_bucket_object.notebook.name}"
        generation = google_storage_bucket_object.notebook.generation
      }

      notebook_runtime_template_resource_name = "projects/${google_colab_runtime_template.my_runtime_template.project}/locations/${google_colab_runtime_template.my_runtime_template.location}/notebookRuntimeTemplates/${google_colab_runtime_template.my_runtime_template.name}"
      gcs_output_uri = "gs://${google_storage_bucket.output_bucket.name}"
      service_account = "%{service_account}"
      }
  }

  depends_on = [
    google_colab_runtime_template.my_runtime_template,
    google_storage_bucket.output_bucket,
  ]
}
`, context)
}

func testAccColabSchedule_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_colab_runtime_template" "my_runtime_template" {
  name = "tf-test-runtime-template%{random_suffix}"
  display_name = "Runtime template"
  location = "us-central1"

  machine_spec {
    machine_type     = "e2-standard-4"
  }

  network_spec {
    enable_internet_access = true
  }
}

resource "google_storage_bucket" "output_bucket" {
  name          = "tf_test_my_bucket%{random_suffix}"
  location      = "US"
  force_destroy = true
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_object" "notebook" {
  name   = "hello_world.ipynb"
  bucket = google_storage_bucket.output_bucket.name
  content = <<EOF
    {
      "cells": [
        {
          "cell_type": "code",
          "execution_count": null,
          "metadata": {},
          "outputs": [],
          "source": [
            "print(\"Hello, World!\")"
          ]
        }
      ],
      "metadata": {
        "kernelspec": {
          "display_name": "Python 3",
          "language": "python",
          "name": "python3"
        },
        "language_info": {
          "codemirror_mode": {
            "name": "ipython",
            "version": 3
          },
          "file_extension": ".py",
          "mimetype": "text/x-python",
          "name": "python",
          "nbconvert_exporter": "python",
          "pygments_lexer": "ipython3",
          "version": "3.8.5"
        }
      },
      "nbformat": 4,
      "nbformat_minor": 4
    }
    EOF
}

resource "google_colab_schedule" "schedule" {
  display_name = "tf-test-schedule-updated%{random_suffix}"
  location = "%{location}"
  allow_queueing = false
  max_concurrent_run_count = 1
  cron = "TZ=America/Los_Angeles 0 * * * *"
  max_run_count = 3
  start_time = "%{updated_start_time}"
  end_time = "%{updated_end_time}"

  create_notebook_execution_job_request {
    notebook_execution_job {
      display_name = "Notebook execution"
      gcs_notebook_source {
        uri = "gs://${google_storage_bucket_object.notebook.bucket}/${google_storage_bucket_object.notebook.name}"
        generation = google_storage_bucket_object.notebook.generation
      }

      notebook_runtime_template_resource_name = "projects/${google_colab_runtime_template.my_runtime_template.project}/locations/${google_colab_runtime_template.my_runtime_template.location}/notebookRuntimeTemplates/${google_colab_runtime_template.my_runtime_template.name}"
      gcs_output_uri = "gs://${google_storage_bucket.output_bucket.name}"
      service_account = "%{service_account}"
      }
  }

  depends_on = [
    google_colab_runtime_template.my_runtime_template,
    google_storage_bucket.output_bucket,
  ]
}
`, context)
}

func TestAccColabSchedule_allFields(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"location":        envvar.GetTestRegionFromEnv(),
		"project_id":      envvar.GetTestProjectFromEnv(),
		"service_account": envvar.GetTestServiceAccountFromEnv(t),
		"random_suffix":   acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccColabSchedule_allFields(context),
			},
		},
	})
}

func testAccColabSchedule_allFields(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network" "my_network" {
  name                    = "tf-test-network%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "my_subnetwork" {
  name          = "tf-test-subnetwork%{random_suffix}"
  network       = google_compute_network.my_network.id
  region        = "us-central1"
  ip_cidr_range = "10.0.1.0/24"
}

resource "google_colab_schedule" "all_fields_schedule" {
  display_name             = "tf-test-schedule%{random_suffix}"
  location                 = "%{location}"
  max_concurrent_run_count = 2
  cron                     = "TZ=America/Los_Angeles * * * * *"

  create_notebook_execution_job_request {
    notebook_execution_job_id = "job-id"
    parent                    = "projects/%{project_id}/locations/%{location}"
    notebook_execution_job {
      display_name   = "Notebook execution"
      gcs_output_uri = "gs://my-bucket"
      execution_user = "user@example.com"
      kernel_name    = "python3"

      direct_notebook_source {
        content = "ewogICJjZWxscyI6IFtdLAogICJtZXRhZGF0YSI6IHt9LAogICJuYmZvcm1hdCI6IDQsCiAgIm5iZm9ybWF0X21pbm9yIjogMAp9Cg=="
      }

      custom_environment_spec {
        machine_spec {
          machine_type       = "n1-standard-4"
          accelerator_type   = "NVIDIA_TESLA_T4"
          accelerator_count  = 1
          gpu_partition_size = "1g.10gb"
          tpu_topology       = "2x2"
          reservation_affinity {
            key                       = "reservation-key"
            reservation_affinity_type = "NO_RESERVATION"
            use_reservation_pool      = false
            values                    = ["reservation-val"]
          }
        }
        network_spec {
          enable_internet_access = true
          network                = google_compute_network.my_network.id
          subnetwork             = google_compute_subnetwork.my_subnetwork.id
        }
        persistent_disk_spec {
          disk_size_gb = "100"
          disk_type    = "pd-standard"
        }
      }

      encryption_spec {
        kms_key_name = "projects/%{project_id}/locations/%{location}/keyRings/my-keyring/cryptoKeys/my-key"
      }

      labels = {
        test = "value"
      }

      parameters = {
        param1 = "val1"
      }
    }
  }

  create_pipeline_job_request {
    pipeline_job_id = "pipeline-job-id"
    parent          = "projects/%{project_id}/locations/%{location}"
    pipeline_job {
      display_name       = "test-pipeline-job"
      network            = google_compute_network.my_network.id
      service_account    = "%{service_account}"
      template_uri       = "https://us-kfp.pkg.dev/proj/repo/template/v1"
      reserved_ip_ranges = ["vertex-ai-ip-range"]

      encryption_spec {
        kms_key_name = "projects/%{project_id}/locations/%{location}/keyRings/my-keyring/cryptoKeys/my-key"
      }

      psc_interface_config {
        network_attachment = "projects/%{project_id}/regions/us-central1/networkAttachments/my-attachment"
        dns_peering_configs {
          domain         = "my-internal-domain.corp."
          target_network = google_compute_network.my_network.id
          target_project = "%{project_id}"
        }
      }

      runtime_config {
        gcs_output_directory = "gs://my-bucket/pipeline_root"
        input_artifacts = {
          artifact1 = "val1"
        }
        parameter_values = {
          param1 = "val1"
        }
      }
    }
  }
}
`, context)
}

