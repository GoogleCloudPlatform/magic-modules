package compute_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/go-version"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccComputeInstancePowerAction_basic(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	project := os.Getenv("GOOGLE_PROJECT")

	zone := os.Getenv("GOOGLE_ZONE")
	if zone == "" {
		zone = "us-central1-a"
	}

	rName := fmt.Sprintf("tf-test-power-action%s", acctest.RandString(t, 10))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.14.0"))),
		},
		CheckDestroy: testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstancePowerActionConfig(project, zone, rName, "stop"),
				Check: resource.ComposeTestCheckFunc(
					testAccWaitForStatus(ctx, t, project, zone, rName, "TERMINATED", 5*time.Minute),
				),
			},
			{
				Config: testAccComputeInstancePowerActionUpdateConfig(project, zone, rName, "start"),
				Check: resource.ComposeTestCheckFunc(
					testAccWaitForStatus(ctx, t, project, zone, rName, "RUNNING", 5*time.Minute),
				),
			},
			{
				Config: testAccComputeInstancePowerActionUpdateConfig(project, zone, rName, "restart"),
				Check: resource.ComposeTestCheckFunc(
					testAccWaitForStatus(ctx, t, project, zone, rName, "RUNNING", 5*time.Minute),
				),
			},
		},
	})
}

// Poll for instance status until expected or timeout
func testAccWaitForStatus(ctx context.Context, t *testing.T, project, zone, name, expected string, timeout time.Duration) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		cfg := acctest.GoogleProviderConfig(t)
		client := cfg.NewComputeClient(cfg.UserAgent)

		deadline := time.Now().Add(timeout)
		for {
			inst, err := client.Instances.Get(project, zone, name).Context(ctx).Do()
			if err != nil {
				return fmt.Errorf("Error fetching instance: %v", err)
			}
			if inst.Status == expected {
				return nil
			}
			if time.Now().After(deadline) {
				return fmt.Errorf("Timeout waiting for status %q, last seen %q", expected, inst.Status)
			}
			time.Sleep(10 * time.Second)
		}
	}
}

// Initial config with lifecycle trigger for after_create
func testAccComputeInstancePowerActionConfig(project, zone, rName, op string) string {
	return fmt.Sprintf(`
resource "google_compute_instance" "test" {
  name         = "%[3]s"
  machine_type = "e2-micro"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-11"
    }
  }

  network_interface {
    network = "default"
  }

  labels = {
    action-operation = "%[4]s"
  }

  lifecycle {
    action_trigger {
      events  = [after_create]
      actions = [action.google_compute_instance_power.test_action_create]
    }
  }
}

action "google_compute_instance_power" "test_action_create" {
  config {
    instance  = google_compute_instance.test.name
    project   = google_compute_instance.test.project
    zone      = "us-central1-a"
    operation = "stop"
  }
}
`, project, zone, rName, op)
}

// Update config with combined lifecycle triggers
func testAccComputeInstancePowerActionUpdateConfig(project, zone, rName, op string) string {
	return fmt.Sprintf(`
resource "google_compute_instance" "test" {
  name         = "%[3]s"
  machine_type = "e2-micro"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-11"
    }
  }

  network_interface {
    network = "default"
  }

  labels = {
    action-operation = "%[4]s"
  }

  lifecycle {
    action_trigger {
      events  = [after_create]
      actions = [action.google_compute_instance_power.test_action_create]
    }
    action_trigger {
      events  = [after_update]
      actions = [action.google_compute_instance_power.test_action]
    }
  }
}

action "google_compute_instance_power" "test_action_create" {
  config {
    instance  = google_compute_instance.test.name
    project   = google_compute_instance.test.project
    zone      = "%[2]s"
    operation = "stop"
  }
}

action "google_compute_instance_power" "test_action" {
  config {
    instance  = google_compute_instance.test.name
    project   = google_compute_instance.test.project
    zone      = "%[2]s"
    operation = "%[4]s"
  }
}
`, project, zone, rName, op)
}
