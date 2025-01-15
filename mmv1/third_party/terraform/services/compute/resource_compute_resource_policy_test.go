package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/services/compute"
)

// Value returned from the API will always be in format "HH:MM", so we need the suppress only on new values
func TestHourlyFormatSuppressDiff(t *testing.T) {
	cases := map[string]struct {
		Old, New           string
		ExpectDiffSuppress bool
	}{
		"Same value": {
			Old:                "01:00",
			New:                "01:00",
			ExpectDiffSuppress: false,
		},
		"Same value but different format": {
			Old:                "01:00",
			New:                "1:00",
			ExpectDiffSuppress: true,
		},
		"Changed value": {
			Old:                "01:00",
			New:                "02:00",
			ExpectDiffSuppress: false,
		},
		"Changed value but different format": {
			Old:                "01:00",
			New:                "2:00",
			ExpectDiffSuppress: false,
		},
		"Check interference with unaffected values": {
			Old:                "11:00",
			New:                "22:00",
			ExpectDiffSuppress: false,
		},
	}

	for tn, tc := range cases {
		if compute.HourlyFormatSuppressDiff("", tc.Old, tc.New, nil) != tc.ExpectDiffSuppress {
			t.Errorf("bad: %s, %q => %q expect DiffSuppress to return %t", tn, tc.Old, tc.New, tc.ExpectDiffSuppress)
		}
	}
}

func TestAccComputeResourcePolicy_attached(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeResourcePolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeResourcePolicy_attached(acctest.RandString(t, 10)),
			},
			{
				ResourceName:      "google_compute_resource_policy.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeResourcePolicy_guestFlushEmptyValue(t *testing.T) {
	t.Parallel()

	context_1 := map[string]interface{}{
		"suffix":              fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10)),
		"snapshot_properties": ``,
	}

	context_2 := map[string]interface{}{
		"suffix": context_1["suffix"],
		"snapshot_properties": `
		snapshot_properties {
          labels = {
            foo = "bar"
          }
        }
		`,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeResourcePolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeResourcePolicy_guestFlushEmptyValue(context_1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("google_compute_resource_policy.foo", "guest_flush"),
				),
			},
			{
				ResourceName:      "google_compute_resource_policy.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeResourcePolicy_guestFlushEmptyValue(context_2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("google_compute_resource_policy.foo", "guest_flush"),
				),
			},
			{
				ResourceName:      "google_compute_resource_policy.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeResourcePolicy_attached(suffix string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance" "foobar" {
  name           = "tf-test-%s"
  machine_type   = "e2-medium"
  zone           = "us-central1-a"
  can_ip_forward = false
  tags           = ["foo", "bar"]

  //deletion_protection = false is implicit in this config due to default value

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = "default"
  }

  metadata = {
    foo            = "bar"
    baz            = "qux"
    startup-script = "echo Hello"
  }

  labels = {
    my_key       = "my_value"
    my_other_key = "my_other_value"
  }

  resource_policies = [google_compute_resource_policy.foo.self_link]
}

resource "google_compute_resource_policy" "foo" {
  name   = "tf-test-policy-%s"
  region = "us-central1"
  group_placement_policy {
    availability_domain_count = 2
  }
}

`, suffix, suffix)
}

func testAccComputeResourcePolicy_guestFlushEmptyValue(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_resource_policy" "foo" {
  name   = "tf-test-gce-policy%{suffix}"
  region = "us-central1"
  snapshot_schedule_policy {
	schedule {
	  daily_schedule {
		days_in_cycle = 1
		start_time    = "04:00"
	  }
	}
	%{snapshot_properties}
  }
}
`, context)
}
