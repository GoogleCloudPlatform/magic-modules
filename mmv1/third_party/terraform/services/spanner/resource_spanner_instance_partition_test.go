package spanner_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"google.golang.org/api/spanner/v1"
)

func TestAccSpannerInstancePartition_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSpannerInstancePartitionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerInstancePartition_basic(context),
			},
			{
				ResourceName:      "google_spanner_instance_partition.partition",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSpannerInstancePartition_update(context),
			},
			{
				ResourceName:      "google_spanner_instance_partition.partition",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSpannerInstancePartition_processingUnits(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSpannerInstancePartitionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerInstancePartition_processingUnits(context),
			},
			{
				ResourceName:      "google_spanner_instance_partition.partition",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSpannerInstancePartition_processingUnitsUpdate(context),
			},
			{
				ResourceName:      "google_spanner_instance_partition.partition",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSpannerInstancePartition_autoscaling(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSpannerInstancePartitionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerInstancePartition_autoscaling(context),
			},
			{
				ResourceName:      "google_spanner_instance_partition.partition",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSpannerInstancePartition_autoscaling_update(context),
			},
			{
				ResourceName:      "google_spanner_instance_partition.partition",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSpannerInstancePartition_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_spanner_instance" "main" {
  name         = "tf-test-spanner-main-%{random_suffix}"
  config       = "nam6"
  display_name = "main-instance"
  num_nodes    = 1
}

resource "google_spanner_instance_partition" "partition" {
  name         = "tf-test-partition-%{random_suffix}"
  instance     = google_spanner_instance.main.name
  config       = "regional-us-central1"
  display_name = "test-spanner-partition"
  node_count   = 1
}
`, context)
}

func testAccSpannerInstancePartition_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_spanner_instance" "main" {
  name         = "tf-test-spanner-main-%{random_suffix}"
  config       = "nam6"
  display_name = "main-instance"
  num_nodes    = 1
}

resource "google_spanner_instance_partition" "partition" {
  name         = "tf-test-partition-%{random_suffix}"
  instance     = google_spanner_instance.main.name
  config       = "regional-us-central1"
  display_name = "updated-spanner-partition"
  node_count   = 2
}
`, context)
}

func testAccSpannerInstancePartition_processingUnits(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_spanner_instance" "main" {
  name             = "tf-test-spanner-main-%{random_suffix}"
  config           = "nam6"
  display_name     = "main-instance"
  processing_units = 1000
}

resource "google_spanner_instance_partition" "partition" {
  name             = "tf-test-partition-%{random_suffix}"
  instance         = google_spanner_instance.main.name
  config           = "regional-us-central1"
  display_name     = "test-spanner-partition"
  processing_units = 1000
}
`, context)
}

func testAccSpannerInstancePartition_processingUnitsUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_spanner_instance" "main" {
  name             = "tf-test-spanner-main-%{random_suffix}"
  config           = "nam6"
  display_name     = "main-instance"
  processing_units = 1000
}

resource "google_spanner_instance_partition" "partition" {
  name             = "tf-test-partition-%{random_suffix}"
  instance         = google_spanner_instance.main.name
  config           = "regional-us-central1"
  display_name     = "updated-spanner-partition"
  processing_units = 2000
}
`, context)
}

func testAccSpannerInstancePartition_autoscaling(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_spanner_instance" "main" {
  name         = "tf-test-spanner-main-%{random_suffix}"
  config       = "nam6"
  display_name = "main-instance"
  num_nodes    = 2
}

resource "google_spanner_instance_partition" "partition" {
  name         = "tf-test-partition-%{random_suffix}"
  instance     = google_spanner_instance.main.name
  config       = "regional-us-central1"
  display_name = "test-spanner-partition"
  autoscaling_config {
    autoscaling_limits {
	  min_nodes = 1
	  max_nodes = 2
    }
    autoscaling_targets {
      high_priority_cpu_utilization_percent = 65
      storage_utilization_percent           = 95
    }
  }
}
`, context)
}

func testAccSpannerInstancePartition_autoscaling_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_spanner_instance" "main" {
  name         = "tf-test-spanner-main-%{random_suffix}"
  config       = "nam6"
  display_name = "main-instance"
  num_nodes    = 3
}

resource "google_spanner_instance_partition" "partition" {
  name         = "tf-test-partition-%{random_suffix}"
  instance     = google_spanner_instance.main.name
  config       = "regional-us-central1"
  display_name = "updated-spanner-partition"
  autoscaling_config {
    autoscaling_limits {
	  min_nodes = 1
	  max_nodes = 3
    }
    autoscaling_targets {
      high_priority_cpu_utilization_percent = 75
      storage_utilization_percent           = 90
    }
  }
}
`, context)
}

func testAccCheckSpannerInstancePartitionDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_spanner_instance_partition" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)
			spannerService, err := spanner.NewService(context.Background(), option.WithHTTPClient(config.Client))
			if err != nil {
				return err
			}

			_, err = spannerService.Projects.Instances.InstancePartitions.Get(rs.Primary.ID).Do()
			if err != nil {
				if isGoogleApiErrorWithCode(err, 404) {
					return nil
				}
				return fmt.Errorf("Error retrieving instance partition %s: %s", rs.Primary.ID, err)
			}
			return fmt.Errorf("Instance partition %s still exists", rs.Primary.ID)
		}
		return nil
	}
}

func isGoogleApiErrorWithCode(err error, code int) bool {
	gerr, ok := err.(*googleapi.Error)
	return ok && gerr.Code == code
}
