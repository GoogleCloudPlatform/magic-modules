package datafusion_test

import (
	"fmt"
	"testing"
	"regexp"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels", "deletion_protection"},
			},
			{
				Config: testAccDataFusionInstance_updated(instanceName),
			},
			{
				ResourceName:            "google_data_fusion_instance.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels", "deletion_protection"},
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
  deletion_protection = true
  # See supported versions here https://cloud.google.com/data-fusion/docs/support/version-support-policy
  version = "6.9.1"
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
  deletion_protection = true

  labels = {
    label1 = "value1"
    label2 = "value2"
  }
  version = "6.9.2"

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
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels", "deletion_protection"},
			},
			{
				Config: testAccDataFusionInstanceEnterprise_updated(instanceName),
			},
			{
				ResourceName:            "google_data_fusion_instance.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels", "deletion_protection"},
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
  deletion_protection = true
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
  deletion_protection = true

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

func TestAccDataFusionInstanceVersion_dataFusionInstanceUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"version":       "6.9.1",
	}

	contextUpdate := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"version":       "6.9.2",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataFusionInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataFusionInstanceVersion_dataFusionInstanceUpdate(context),
			},
			{
				ResourceName:            "google_data_fusion_instance.basic_instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"region", "deletion_protection"},
			},
			{
				Config: testAccDataFusionInstanceVersion_dataFusionInstanceUpdate(contextUpdate),
			},
			{
				ResourceName:            "google_data_fusion_instance.basic_instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"region", "deletion_protection"},
			},
		},
	})
}

func testAccDataFusionInstanceVersion_dataFusionInstanceUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_fusion_instance" "basic_instance" {
  name   = "tf-test-my-instance%{random_suffix}"
  region = "us-central1"
  type   = "BASIC"
  deletion_protection = true
  # Mark for testing to avoid service networking connection usage that is not cleaned up
  options = {
    prober_test_run = "true"
  }
  version = "%{version}"
}
`, context)
}

func TestAccDatafusionInstance_DeletionProtection(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataFusionInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDatafusionInstanceDeletionProtection("us-central1"),
			},
			{
				ResourceName:            "google_data_fusion_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels", "deletion_protection"},
			},
			{
                                Config: testAccDatafusionInstanceDeletionProtection("us-west2"),
                                ExpectError: regexp.MustCompile("deletion_protection"),
                        },
		},
	})
}

func testAccDatafusionInstanceDeletionProtection(region string) string {

	return fmt.Sprintf(`
	resource "google_data_fusion_instance" "instance" {
        name   = "my-instance"
        region = "%s"
        type   = "BASIC"
	deletion_protection = true
}
`, region)
}
