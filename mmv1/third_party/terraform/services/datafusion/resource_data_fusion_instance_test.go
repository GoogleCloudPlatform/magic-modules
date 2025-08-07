package datafusion_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	datafusion "cloud.google.com/go/datafusion/apiv1"
	datafusionpb "cloud.google.com/go/datafusion/apiv1/datafusionpb"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataFusionInstance_update(t *testing.T) {
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataFusionInstance_basic(instanceName),
			},
			{
				ResourceName:            "google_data_fusion_instance.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccDataFusionInstance_updated(instanceName),
			},
			{
				ResourceName:            "google_data_fusion_instance.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func testAccDataFusionInstance_basic(instanceName string) string {
	return fmt.Sprintf(`
resource "google_data_fusion_instance" "foobar" {
  name   = "%s"
  region = "us-central1"
  type   = "BASIC"
  # See supported versions here https://cloud.google.com/data-fusion/docs/support/version-support-policy
  version = "6.10.0"
  # Mark for testing to avoid service networking connection usage that is not cleaned up
  options = {
  	prober_test_run = "true"
  }
  accelerators {
    accelerator_type = "CDC"
    state = "DISABLED"
  }
}
`, instanceName)
}

func testAccDataFusionInstance_updated(instanceName string) string {
	return fmt.Sprintf(`
resource "google_data_fusion_instance" "foobar" {
  name                          = "%s"
  region                        = "us-central1"
  type                          = "DEVELOPER"
  enable_stackdriver_monitoring = true
  enable_stackdriver_logging    = true

  labels = {
    label1 = "value1"
    label2 = "value2"
  }
  version = "6.10.1"

  accelerators {
    accelerator_type = "CCAI_INSIGHTS"
    state = "ENABLED"
  }
  # Mark for testing to avoid service networking connection usage that is not cleaned up
  options = {
  	prober_test_run = "true"
  }
}
`, instanceName)
}

func TestAccDataFusionInstanceEnterprise_update(t *testing.T) {
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataFusionInstanceEnterprise_basic(instanceName),
			},
			{
				ResourceName:            "google_data_fusion_instance.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccDataFusionInstanceEnterprise_updated(instanceName),
			},
			{
				ResourceName:            "google_data_fusion_instance.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func testAccDataFusionInstanceEnterprise_basic(instanceName string) string {
	return fmt.Sprintf(`
resource "google_data_fusion_instance" "foobar" {
  name   = "%s"
  region = "us-central1"
  type   = "ENTERPRISE"
  # Mark for testing to avoid service networking connection usage that is not cleaned up
  options = {
  	prober_test_run = "true"
  }
}
`, instanceName)
}

func testAccDataFusionInstanceEnterprise_updated(instanceName string) string {
	return fmt.Sprintf(`
resource "google_data_fusion_instance" "foobar" {
  name                          = "%s"
  region                        = "us-central1"
  type                          = "ENTERPRISE"
  enable_stackdriver_monitoring = true
  enable_stackdriver_logging    = true
  enable_rbac                   = true

  labels = {
    label1 = "value1"
    label2 = "value2"
  }
  # Mark for testing to avoid service networking connection usage that is not cleaned up
  options = {
  	prober_test_run = "true"
  }
}
`, instanceName)
}

// Corrected destroy check function
func testAccCheckDatafusionInstanceDestroyProducer(t *testing.T) func(*terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_data_fusion_instance" {
				continue
			}

			instanceName := rs.Primary.Attributes["name"]
			project := rs.Primary.Attributes["project"]
			location := rs.Primary.Attributes["region"]
			if location == "" {
				location = rs.Primary.Attributes["location"]
			}

			if location == "" {
				return fmt.Errorf("could not determine location for Data Fusion instance %s", instanceName)
			}

			ctx := context.Background()
			datafusionClient, err := datafusion.NewClient(ctx)
			if err != nil {
				return fmt.Errorf("failed to create datafusion client: %v", err)
			}
			defer datafusionClient.Close()

			req := &datafusionpb.GetInstanceRequest{
				Name: fmt.Sprintf("projects/%s/locations/%s/instances/%s", project, location, instanceName),
			}

			_, err = datafusionClient.GetInstance(ctx, req)
			if err == nil {
				return fmt.Errorf("Data Fusion instance %s still exists", instanceName)
			}
			if !strings.Contains(err.Error(), "NotFound") {
				return fmt.Errorf("error checking Data Fusion instance %s existence: %v", instanceName, err)
			}
		}
		return nil
	}
}

func TestAccDatafusionInstance_tags(t *testing.T) {
	t.Parallel()

	// The tag key and value bootstrap functions create organization-level tags.
	// We will use the short names here for the labels, which is a common
	// pattern in test environments. We must convert the names to be valid labels.
	tagKeyShortName := "test-tagkey-" + acctest.RandString(t, 6)
	tagValueShortName := "test-tagvalue-" + acctest.RandString(t, 6)

	testContext := map[string]interface{}{
		"tagKey":        tagKeyShortName,
		"tagValue":      tagValueShortName,
		"random_suffix": acctest.RandString(t, 10),
	}
	resourceName := "google_data_fusion_instance.test"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDatafusionInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDatafusionInstanceLabels(testContext),
				Check: resource.ComposeTestCheckFunc(
					checkDatafusionInstanceLabels(resourceName, tagKeyShortName, tagValueShortName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels"},
			},
		},
	})
}

// CORRECTED HCL template to use 'labels' instead of 'tags'
func testAccDatafusionInstanceLabels(testContext map[string]interface{}) string {
	return acctest.Nprintf(`
	resource "google_data_fusion_instance" "test" {
	  name = "tf-test-instance-%{random_suffix}"
	  type = "BASIC"
	  region = "us-central1"
	
	  labels = {
	    "%{tagKey}" = "%{tagValue}"
	  }
	}
	`, testContext)
}

// CORRECTED check function to verify labels on the instance
func checkDatafusionInstanceLabels(resourceName, expectedTagKey, expectedTagValue string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource not found: %s", resourceName)
		}
		project := rs.Primary.Attributes["project"]
		location := rs.Primary.Attributes["region"]
		instanceName := rs.Primary.Attributes["name"]

		ctx := context.Background()
		datafusionClient, err := datafusion.NewClient(ctx)
		if err != nil {
			return fmt.Errorf("failed to create datafusion client: %v", err)
		}
		defer datafusionClient.Close()

		req := &datafusionpb.GetInstanceRequest{
			Name: fmt.Sprintf("projects/%s/locations/%s/instances/%s", project, location, instanceName),
		}

		instance, err := datafusionClient.GetInstance(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to get datafusion instance '%s': %v", req.Name, err)
		}

		labels := instance.GetLabels()
		if labels == nil {
			return fmt.Errorf("expected labels not found on instance '%s'", req.Name)
		}

		if actualValue, ok := labels[expectedTagKey]; ok {
			if actualValue == expectedTagValue {
				return nil
			}
			return fmt.Errorf("label key '%s' found with incorrect value. Expected: %s, Got: %s", expectedTagKey, expectedTagValue, actualValue)
		}

		return fmt.Errorf("expected label key '%s' not found on instance '%s'", expectedTagKey, req.Name)
	}
}
