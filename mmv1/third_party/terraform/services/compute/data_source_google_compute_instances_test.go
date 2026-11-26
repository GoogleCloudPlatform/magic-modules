package compute_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	_ "github.com/hashicorp/terraform-provider-google/google/services/compute"
)

func TestAccDataSourceGoogleComputeInstances_zone(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleComputeInstancesConfig(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_compute_instances.bar", "instances.#", "2"),
					resource.TestCheckResourceAttrSet("data.google_compute_instances.bar", "instances.0.name"),
					resource.TestCheckResourceAttrSet("data.google_compute_instances.bar", "instances.0.self_link"),
					resource.TestCheckResourceAttrSet("data.google_compute_instances.bar", "instances.0.machine_type"),
					resource.TestCheckResourceAttr("data.google_compute_instances.bar", "instances.0.zone", "us-central1-a"),
					resource.TestCheckResourceAttrSet("data.google_compute_instances.bar", "instances.1.name"),
				),
			},
		},
	})
}

func TestAccDataSourceGoogleComputeInstances_allZones(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleComputeInstancesConfig_allZones(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_compute_instances.bar", "instances.#", "2"),
					resource.TestCheckResourceAttrSet("data.google_compute_instances.bar", "instances.0.zone"),
					resource.TestCheckResourceAttrSet("data.google_compute_instances.bar", "instances.1.zone"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleComputeInstancesConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_instance" "foo" {
  name           = "tf-test-%{random_suffix}-foo"
  machine_type   = "e2-medium"
  zone           = "us-central1-a"

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-11"
    }
  }

  network_interface {
    network = "default"
  }
}

resource "google_compute_instance" "bar" {
  name           = "tf-test-%{random_suffix}-bar"
  machine_type   = "e2-medium"
  zone           = "us-central1-a"

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-11"
    }
  }

  network_interface {
    network = "default"
  }
}

data "google_compute_instances" "bar" {
  zone   = "us-central1-a"
  filter = "name:tf-test-%{random_suffix}-*"

  depends_on = [
    google_compute_instance.foo,
    google_compute_instance.bar,
  ]
}
`, context)
}

func testAccDataSourceGoogleComputeInstancesConfig_allZones(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_instance" "foo" {
  name           = "tf-test-%{random_suffix}-foo"
  machine_type   = "e2-medium"
  zone           = "us-central1-a"

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-11"
    }
  }

  network_interface {
    network = "default"
  }
}

resource "google_compute_instance" "bar" {
  name           = "tf-test-%{random_suffix}-bar"
  machine_type   = "e2-medium"
  zone           = "us-central1-b"

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-11"
    }
  }

  network_interface {
    network = "default"
  }
}

data "google_compute_instances" "bar" {
  filter = "name:tf-test-%{random_suffix}-*"

  depends_on = [
    google_compute_instance.foo,
    google_compute_instance.bar,
  ]
}
`, context)
}
